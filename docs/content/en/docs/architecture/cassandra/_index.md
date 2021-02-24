---
title: "Cassandra"
linkTitle: "Cassandra"
weight: 1
description: >
  Apache Cassandra, the core of K8ssandra.
---

 Apache Cassandra is an open-source, NoSQL database built from the foundation of geographically distributed and fault tolerant data replication. Given the ephemeral nature of containers Cassandra is a logical fit as _the_ cloud-native data plane for Kubernetes. 
 
## Operations with cass-operator
 
K8ssandra delegates core Cassandra management to the DataStax Kubernetes Operator for Apache Cassandra, [cass-operator](https://github.com/datastax/cass-operator). cass-operator handles the provisioning of datacenters, scaling operations, rolling restarts, upgrades, and container failure remediation. 
 
## Anatomy of a Cassandra Cluster

 Cassandra clusters are separated into a topology of logical datacenters, racks and nodes. We will cover each level of the topology along with its associated Kubernetes.

### Logical Datacenters (Namespaces?)

Apache Cassandra clusters are composed of one or more logical datacenters. Datacenters are usually aligned to cloud regions or geographical areas, but may reside within the same geography as other datacenters for workload isolation purposes.

![Single DC, Cassandra Cluster on Kubernetes](cassandra-bootstrap-5.png)

_1x Datacenter, 3x Rack, 6x node Cassandra Cluster_

Here we have a single Cassandra datacenter occupying a cloud region. In this deployment there are three failure domains, or logical racks where six nodes are deployed.

### Logical Racks (StatefulSets)

Each logical datacenter is composed of multiple logical racks (named such as they previously represented physical racks in datacenters). Cassandra ensures that data is replicated across rack boundaries such that the loss of a single rack does not effect data availability. With K8ssandra, logical Cassandra racks are mapped to Kubernetes Stateful Sets. Thus a datacenter with three logical racks will be composed or three Stateful Sets. Stateful Sets allow for reliable and consistent identity and storage between instances of containers running.

![Single Rack / Stateful Set](cassandra-rack.png)

If the replication factor in use matches the number of racks being deployed across then each rack contains a single copy of the data. It is important to note that while an entire rack may be taken down and still support operations at local quorum sizing _must_ take into account the additional query load on each of the remaining racks should one become unavailable.

### Nodes (Pods)

The smallest unit within the topology of a Cassandra cluster is a single node. A Cassandra node is represented by a JVM process. It _is_ possible to run multiple instances or nodes of Cassandra per physical host, but care should be that there are enough fault domains to keep multiple record copies off the same host. 

![Cassandra Pod](cassandra-pod.png)

In Kubernetes, each Cassandra pod is composed of a number of containers. The first container run in any Cassandra pod is the `cass-config-builder`. It handles rendering out configurations on a per pod basis with input from the `CassandraDatacenter` custom resource. Next the `cassandra` container is started, but it doesn't begin with the Cassandra JVM. Instead the [Management API for Apache Cassandra](https://github.com/datastax/management-api-for-apache-cassandra) is started first. This boots a REST API for lifecycle and operations tasks to be requested by `cass-operator`. For instance all nodes in the cluster may be scheduled and start their management APIs before the operator starts triggering the bootstrap for nodes. Finally an _optional_ third container is started whose sole purpose is to `tail` logs from the management API process. Given the management API has a separate lifetime (and even binary) than Cassandra the log streams are separated for easier debugging.

## Next

Check out [Monitoring with Prometheus and Grafana]({{< ref "/docs/architecture/monitoring" >}})
