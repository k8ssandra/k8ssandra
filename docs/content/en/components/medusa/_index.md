---
title: "Medusa for Cassandra repair operations"
linkTitle: "Medusa"
weight: 4
description: K8ssandra deploys Medusa to support backup and restore of Apache Cassandra&reg; tables.
---

Medusa for Apache Cassandra&reg; is deployed by K8ssandra as part of its Helm chart install. 

If you haven't already installed K8ssandra, see the [quickstarts]({{< relref "/quickstarts/" >}}) and [install]({{< relref "/install" >}}) topics.

## Introduction

Even with the heightened availability of Apache Cassandra® a proper backup schedule and testing of restore procedures is good practice in case catastrophe strikes. With distributed systems backups can be tricky, there's the timing of the snapshot process on all nodes, correlation of data files to remote storage, and eventual restore.

K8ssandra provides Helm charts for taking backups or triggering the restoration of data. This is accomplished via the [Medusa for Apache Cassandra](https://github.com/thelastpickle/cassandra-medusa) project from The Last Pickle and Spotify.

## Supported storage objects

K8ssandra's Medusa supports:

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

## Related topics

See these topics in the K8ssandra documentation:

* Backup and restore Cassandra with S3-compatible [MinIO]({{< relref "/tasks/backup-restore/minio/" >}}).
* Backup and restore Cassandra with [Amazon S3]({{< relref "/tasks/backup-restore/amazon-s3/" >}}).
* Backup and restore Cassandra with [GCS]({{< relref "/tasks/backup-restore/gcs/" >}}).

## Next

See the other [components]({{< relref "/components/" >}}) deployed by K8ssandra. For information on using the deployed components, see the [Tasks]({{< relref "/tasks/" >}}) topics.
