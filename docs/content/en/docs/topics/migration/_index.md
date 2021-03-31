---
title: "Migrating a Cassandra cluster to K8ssandra"
linkTitle: "Migrating a cluster to K8ssandra"
no_list: true
weight: 1
description: How to migrate an existing Cassandra cluster that's running in Kubernetes to K8ssandra.
---

The strategy to perform this migration to K8ssandra focuses on a datacenter migration.

## The environment

It's assumed that the Cassandra cluster is running in the same Kubernetes cluster in which K8ssandra will run. The Cassandra cluster may have been installed with another operator (such as CassKop), a Helm chart (such as bitnami/cassandra), or by directly creating the YAML manfiests for the StatefulSets and any other objects that were created.

{{% alert title="Tip" color="success" %}}
See <https://thelastpickle.com/blog/2019/02/26/data-center-switch.html> for a thorough guide on how to migrate to a new datacenter.
{{% /alert %}}

## Check the replication strategies of keyspaces

Confirm that the user-defined keyspaces and each of the `system_auth`, `system_distributed`, and `system_traces` are using `NetworkTopologyStrategy`.

{{% alert title="Tip" color="success" %}}
The `system_*` keyspaces include `system_auth`, `system_schema`, `system_distributed`, and `system_traces`.
{{% /alert %}}

It is generally recommended to use a replicator factor of `3`. If your keyspaces currently use `SimpleStrategy` and assuming you have at least 3 Cassandra nodes, then you would change the replication factor for `system_auth` as follows:

```cql
ALTER KEYSPACE system_auth WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': 3}
```

`dc1` is the name of the datacenter as defined in `cassandra-rackdc.properties`.

Changing the recommendation may result in a topology change; that is, a token ownership change. 

{{% alert title="Recommendation" color="success" %}}
If you change the replication strategy of a keyspace, run a full, cluster-wide repair on it. See the [nodetool repair](https://docs.datastax.com/en/cassandra-oss/3.x/cassandra/tools/toolsRepair.html) or the [Cassandra Reaper]({{< ref "/docs/topics/repair/" >}}) topic for information on repair operations. 
{{% /alert %}}

## Check the endpoint snitch

Make sure that `GossipingPropertyFileSnitch` is used, and not `SimpleSnitch`.

{{% alert title="Recommendation" color="success" %}}
If you change the snitch, run a full, cluster-wide repair. Again, refer to the `nodetool repair` or the Reaper documentation. 
{{% /alert %}}

## Client changes

Make sure client services are using a `LoadBalancingPolicy` that will route requests to the local datacenter. With version 4.11 of the Java driver you might configure this setting as follows:

```yaml
datastax-java-driver.basic.load-balancing-policy {
  local-datacenter = dc1
}
```

Also, make sure your clients are using `LOCAL_*` and not `QUORUM` consistency levels. You can check this setting in your `cassandra.yaml` configuration file. 

## Install K8ssandra

Before installing K8ssandra, make a note of the IP addresses of the seed nodes of the current datacenter. You will use the following to configure the `additionalSeeds` property in the K8ssandra chart. Here is an example of a chart properties overrides file. Assume that it is named `k8ssandra-values.yaml`:

```yaml
cassandra:
  clusterName: cassandra
  datacenters:
  - name: dc2
    size: 3
  additionalSeeds:
  - <dc1-seed>
  - <dc1-seed>
  - <dc2-seed>
```

Run `nodetool status <keyspace-name>` to verify that the nodes in the new datacenter are gossiping with the old datacenter. For the keyspace argument, you can specify a user-defined one or `system_auth`.  Here is some example output:

```bash
nodetool status system_auth
```

**Output:**

```bash
Datacenter: datacenter1
=======================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load        Tokens       Owns (effective)  Host ID                               Rack
UN  10.40.4.16  236 KiB     256          100.0%            f659967d-07b8-49b8-9ca8-bd02e2a58911  rack1
UN  10.40.5.2   318.03 KiB  256          100.0%            4b58ef5a-5578-4126-b1b6-4fb9a2d2cd40  rack1
UN  10.40.2.2   341.91 KiB  256          100.0%            380848f0-297b-47fd-9509-56d70835f410  rack1
Datacenter: dc1
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load        Tokens       Owns (effective)  Host ID                               Rack
UN  10.40.3.42  71.04 KiB   256          0.0%              0380fb66-2697-4a90-80a3-cf6f4b1f3476  default
UN  10.40.4.19  97.22 KiB   256          0.0%              1fd067c8-6eb2-4fd6-bd88-f6fbb72e8ede  default
UN  10.40.5.4   97.21 KiB   256          0.0%              2a986937-856b-4758-9026-146dc4620eb4  default
```

You should expect to see 0.0% ownership for the nodes in the new datacenter. This is because we are not yet replicating to the new datacenter.

{{% alert title="Tip" color="success" %}}
If you run `nodetool status` without the keyspace argument, you may find that nodes in the new datacenter report something greater than 0.0% for the token ownership. This happens because, when Stargate is installed, a keyspace named `data_endpoint_auth` is created.
{{% /alert %}}

## Update replication strategy

For each keyspace that uses `NetworkTopologyStrategy`, update the replication strategy to include the new datacenter as follows:

```cql
ALTER KEYSPACE system_auth WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': 3, 'dc2': 3}
```

Run the same `nodetool status` command that you ran previously. You should see that the token ownership has changed for the nodes in the new datacenter.

At this point nodes in the new datacenter will start receiving writes.

## Stream data to the new nodes

Next you need to stream data from the original datacenter to the new datacenter to get the latter caught up with past writes. This can be done with the `nodetool rebuild` command. It needs to be run on each node in the new datacenter as follows:

```bash
nodetool rebuild <old-datacenter-name>
```

## Stop sending traffic to the old datacenter

To stop sending traffic to the old datacenter, there are two steps:

1. Route clients to the new datacenter
2. Update replication strategies

### Route clients to the new datacenter

Update the `LoadBalancingPolicy` of client services to route requests to the new datacenter. Here is an example for v4.11 of the Java driver where the new datacenter is named dc2:

```yaml
datastax-java-driver.basic.load-balancing-policy {
  local-datacenter = dc2
}
```

{{% alert title="Recommendation" color="success" %}}
Verify that there are no client connections to nodes in the old datacenter.
{{% /alert %}}

### Update replication strategies

Next we need to stop replicating data to nodes in the old datacenter. For each keyspace that was previously updated, we need to update it again. This time we will remove the old datacenter. 

Here is an example that specifies the `system_auth keyspace`:

```cql
ALTER KEYSPACE system_auth WITH replication = {'class': 'NetworkTopologyStrategy', 'dc2': 3}
```

## Remove the old datacenter

Now we are ready to do some cleanup and remove the old datacenter and associated resources that are no longer in use. There are three steps:

1. Decommission nodes
2. Remove old seed nodes
3. Delete old datacenter resources from Kubernetes

### Decommission nodes

Run the following on each node in the old datacenter:

```bash
nodetool decommission
```

{{% alert title="Note" color="success" %}}
Decommissioning should be done serially, one node at a time.
{{% /alert %}}

### Remove the old seed nodes

Remove the seeds nodes from the `additionalSeeds` chart property. You can simply remove the `additionaSeeds` property from your chart overrides file; refer to the [Install K8ssandra](#Install-K8ssandra) topic above. 

You need to run a `helm upgrade` for the changes to take effect. Assuming that the edited chart properties are stored in a file named `k8ssandra-values.yaml`:

```bash
helm upgrade <release-name> k8ssandra/k8ssandra -f k8ssandra-values.yaml
```

The command above will trigger a statefulset update, which in effect is a cluster-wide (with respect to Cassandra) rolling restart.

### Delete the old datacenter resources from Kubernetes

The steps necessary here will vary depending on how you installed and managed the old datacenter. You want to make sure that the StatefulSets, PersistentVolumeClaims, PersistentVolumes, and any other objects created in association with the old datacenter are deleted in the Kubernetes cluster.

## Next

Learn how to [monitor Cassandra]({{< ref "/docs/topics/monitoring/" >}}) in a Kubernetes cluster.
