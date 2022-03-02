---
title: "Medusa for Cassandra backup and restore"
linkTitle: "Medusa"
weight: 4
description: K8ssandra deploys Medusa to support backup and restore of Apache Cassandra&reg; tables.
---

Medusa for Apache Cassandra&reg; is deployed by K8ssandra Operator as part of its K8ssandraCluster install. 

If you haven't already installed K8ssandra Operator, see the [install](https://docs-v2.k8ssandra.io/install) topics.

## Introduction

Even with the heightened availability of Apache Cassandra® a proper backup schedule and testing of restore procedures is good practice in case catastrophe strikes. With distributed systems backups can be tricky, there's the timing of the snapshot process on all nodes, correlation of data files to remote storage, and eventual restore.

K8ssandra provides Helm charts for taking backups or triggering the restoration of data. This is accomplished via the [Medusa for Apache Cassandra](https://github.com/thelastpickle/cassandra-medusa) project from The Last Pickle and Spotify.

## Supported storage objects

K8ssandra Operator's Medusa supports:

* Google Cloud Storage (GCS)

* Amazon S3  

* All S3-compatible implementations, which include:

  * MinIO 
  * IBM Cloud Object Storage
  * OVHCloud Object Storage
  * Riak S2
  * Dell EMC ECS
  * CEPH Object Gateway
  * Others - this list is not exhaustive

## Next steps

* [Backup and restore Cassandra data](https://docs-v2.k8ssandra.io/tasks/backup-restore/)
* Reference topics:
  * [Medusa Custom Resource Definition (CRD)](https://docs-v2.k8ssandra.io/reference/crd/medusa/)
  * [CassandraBackup CRD](https://docs-v2.k8ssandra.io/reference/crd/cassandrabackup/)
  * [CassansdraRestore CRD](https://docs-v2.k8ssandra.io/reference/crd/cassandrarestore/)
* Additional [components](https://docs-v2.k8ssandra.io/components/)
* [Tasks](https://docs-v2.k8ssandra.io/tasks)
