---
title: "Cassandra"
linkTitle: "Cassandra"
weight: 1
description: >
  Apache Cassandra, the core of K8ssandra.
---
 Apache Cassandra is an open-source, NoSQL database built from the ground up with geographically distributed and fault tolerant data repliction as a base concept. Given the ephemeral nature of containers it is a logical fit as _the_ cloud-native data plane for Kubernetes. 
 
 ## Operations with cass-operator
 
 K8ssandra leverages the DataStax Kubernetes Operator for Apache Cassandra, [cass-operator](https://github.com/datastax/cass-operator), as the operator for managing the Cassandra component. It handles provisioning nodes, scaling operations, and automated container failure remediation.
 
 ## Logical Datacenters
 
Apache Cassandra clusters are composed of one or more logical datacenters. Datacenters are usually aligned to cloud regions, but may reside within the same geography as other datacenters for workload isolation purposes (OLTP vs OLAP).  

## Logical Racks (StatefulSets)
A single logical datacenter is composed of multiple logical racks (think about racks in datacenters). Cassandra ensures that data is replicated across rack boundaries such that the loss of a single rack does not affect data availability. With K8ssandra, logical Cassandra racks are mapped to Kubernetes Stateful Sets. Thus a datacenter with three logical racks will be composed or three Stateful Sets. Stateful Sets allow for reliable and consistent identity and storage between instances of containers running.

## Nodes (Pods)

