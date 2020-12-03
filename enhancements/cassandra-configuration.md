# Overview
Cassandra has a lot of configuration settings that are exposed in different places like,

* cassandra.yaml
* jvm.options
* cassandra-env.sh
* logback.xml

The purpose of this document is determine what configuration settings to expose in k8ssandra 
and how. Doing so requires some understanding of how things are done in cass-operator.

**Note:** This doc currently does not take into consideration any differences between 3.11 
and 4.0.

# Cass Operator

## Topology
Cass Operator utilizes Cassandra racks to establish anti-affinity for pods. Let's look at an 
example CassandraDatacenter:

```yaml

apiVersion: cassandra.datastax.com/v1beta1
kind: CassandraDatacenter
metadata:
  name: dc1
spec:
  clusterName: multi-rack
  serverType: cassandra
  serverVersion: 3.11.7
  managementApiAuth:
    insecure: {}
  size: 9
  racks:
  - name: rack-1
    zone: us-east1-b
  - name: rack-2
    zone: us-east1-c
  - name: rack-3
    zone: us-east1-d
  ...
```

**Note:** Zones will soon be deprecated in favor of free form labels.

We are specifying a 9 node cluster with 3 racks. cass-operator will evenly distribute C* 
nodes across the racks.

Each rack includes a zone property which corresponds to an availability zone. cass-operator 
will use affinity and anti-affinity to ensure that pods from each rack are scheduled into 
the corresponding availability zone and also to ensure that C* pods are not co-located on 
k8s worker nodes.

Cass Operator always configures the cluster to use `GossipingPropertyFileSnitch`. I am not 
sure what happens if I try to specify a different value for `endpoint_snitch`, but 
regardless, this is an example of a setting that we should not expose in k8ssandra.

Now let's look at how to create a multi-DC cluster (multi-region is outside the scope of 
this doc).

Here is the first DC:

```yaml
apiVersion: cassandra.datastax.com/v1beta1
kind: CassandraDatacenter
metadata:
  name: dc1
  namespace: dev
spec:
  clusterName: multi-dc
  ...
```
Here is the second DC:

```yaml
apiVersion: cassandra.datastax.com/v1beta1
kind: CassandraDatacenter
metadata:
  name: dc2
  namespace: dev
spec:
  clusterName: multi-dc
  ...
```

We create two CassandraDatacenters that have the same clusterName and same namespace.

If we want common settings across the two DCs, e.g., auth, heap settings, etc., we need
to specify those settings for each CassandraDatacenter.

## Configuration
Now let's look at how we specify some settings from cassandra.yaml and jvm.options:

```yaml
apiVersion: cassandra.datastax.com/v1beta1
kind: CassandraDatacenter
metadata:
  name: dc1
# ...
spec:
# ...
  # Everything under the config key is passed to the 
  # cass-config-builder init container that runs on pod creation, and 
  # then marshalled into config files.
  config:
    cassandra-yaml:
      num_tokens: 8
      memtable_allocation_type: heap_buffers
      file_cache_size_in_mb: 1000
      
      authenticator: org.apache.cassandra.auth.PasswordAuthenticator
      authorizer: org.apache.cassandra.auth.CassandraAuthorizer
      role_manager: org.apache.cassandra.auth.CassandraRoleManager
      
      roles_validity_in_ms: 2000
      permissions_validity_in_ms: 2000
      credentials_validity_in_ms: 2000
      
      # How long before a node logs slow queries. Select queries that 
      # take longer than this timeout to execute, will generate an
      # aggregated log message, so that slow queries can be identified. 
      # Set this value to zero to disable slow query logging. 
      slow_query_log_timeout_in_ms: 500

    jvm-options:
      # Set the database to use 14 GB of Java heap
      initial_heap_size: "14G"
      max_heap_size: "14G"

      additional-jvm-opts:
        # As the database comes up for the first time, set system keyspaces to RF=3
        - "-Dcassandra.system_distributed_replication_dc_names=dc1"
        - "-Dcassandra.system_distributed_replication_per_dc=3"
```

I tried to include a number of properties for illustration. There are some however that 
deserve further discussion.

First, I want to point out that there are number of properties that we do not want to expose 
including the following:

* storage_port
* ssl_storage_port
* native_transport_port
* seed_provider
* listen_address
* native_transport_port
* rpc_address
* data_files_directories
* commitlog_directories

There are a number of properties around authentication/authorization. we may not want to 
directly expose properties like `authenticator` and `authorizer`. The auth subsystem is 
pluggable which is great. I can use a custom `IAuthenticator` by specifying the class name 
for `authenticator` and including it in Cassandra's classpath. Doing the latter is a 
non-trivial exercise in Kubernetes.

My example also included some auth cache settings:

* roles_validity_in_ms
* permissions_validity_in_ms
* credentials_validity_in_ms

These are not very good defaults. They can easily contribute to latency issues when auth is 
enabled. We probably want to expose these but we also should override the defaults if auth 
is enabled.

Now let's talk about the slow query log which is enabled via the 
`slow_query_log_timeout_in_ms` property. The slow query log can be really helpful for 
debugging latency and performance issues. If this is enabled, we might want to consider 
adding a logging sidecar for the debug log to make those logs more easily accessible. 

Similarly, if the user wants to enable CDC we should make it easy to configure another 
volume for it.

For `jvm.options`, the big thing that immediately comes to mind is GC settings. 
Depending on what collectors you use, e.g., ParNew/CMS, G1, etc., there are a number of 
properties that you may want to configure. It should be enough for the user to simply 
specify the collector to use and we should then provide sensible defaults. The user should 
still have the ability to set specific GC properties.

# K8ssandra
## Topology
```yaml
k8ssandra:
  clusterName: multi-rack-multi-dc
  datacenters:
  - name: dc1
    racks:
    - name: rack-1
      zone: us-east1-a
    - name: rack-2
      zone: us-east1-b
    - name: rack-3
      zone: us-east1-c
  - name: dc2:
    racks:
    - name: rack-1
      zone: us-west1-a
    - name: rack-2
      zone: us-west1-b
    - name: rack-3
      zone: us-west1-c
```
The idea is to provide a unified view of the k8ssandra cluster where we specify the DCs and 
racks in one place. It will be implementation details left to k8ssandra and cass-operator to 
ensure that the correct operational steps are carried out.

## Configuration
```yaml
k8ssandra:
  clusterName: multi-rack-multi-dc

  configuration:
    # For 3.11, we should consider using the even token distribution
    # algorithm for using small values of num_tokens. I believe the
    # operational steps are simplified a bit in 4.0.
    num_tokens: 8
    auth:
      enabled: true
    file_cache_size_in_mb: 1000
    
    jvm: 
      heapSize: 8G
      gc:
        collector: G1
        # Specify settings specific to G1
          
 
  datacenters:
  - name: dc1
    racks:
    - name: rack-1
      zone: us-east1-a
    - name: rack-2
      zone: us-east1-b
    - name: rack-3
      zone: us-east1-c
  - name: dc2:
    racks:
    - name: rack-1
      zone: us-west1-a
    - name: rack-2
      zone: us-west1-b
    - name: rack-3
      zone: us-west1-c
```
The first thing to point out is that configuration properties are specified outside of the 
DCs. This is in contrast to how it is done in a CassandraDatacenter. 

Doing it this way will make it possible and easier to have common settings across DCs. 

We should also allow settings to be specified at the DC-level. This will accomodate
different workloads per DC, e.g., OLAP vs OLTP.

We might also want to expose rack-level or node-level settings to faciliate things like
debugging and canary rollouts. 

Notice that for auth I have just exposed an enabled property, but we will set a number of 
properties. 

We may also want to expose properties for configuring roles and permissions.

**TODO** Add a section on how we might configure CDC or slow query logging, two things which 
may require volume configuration.