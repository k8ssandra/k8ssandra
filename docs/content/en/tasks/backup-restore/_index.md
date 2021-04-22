---
title: "Backup and restore Apache Cassandra data"
linkTitle: "Backup and restore Cassandra"
no_list: true
weight: 4
description: Use Medusa to backup and restore Apache Cassandra® data in Kubernetes.
---

Medusa is a Cassandra backup and restore tool. It's packaged with K8ssandra and supports a variety of backends. 

## Supported object storage types for backups

Supported in K8ssandra's Medusa since the 1.0.0 release:

* Amazon S3  

* Google Cloud Storage (GCS)

Added in K8ssandra 1.1.0:

* Support in K8ssandra's Medusa for all S3-compatible implementations, which include:

  * MinIO 
  * IBM Cloud Object Storage
  * OVHCloud Object Storage
  * Riak S2
  * Dell EMC ECS
  * CEPH Object Gateway
  * Others - this list is not exhaustive

## Backup and restore steps

For detailed walk-throughs of Medusa backup and restore operations, see:

* Backup and restore Cassandra with S3-compatible [MinIO]({{< relref "/tasks/backup-restore/s3-compatible/" >}}).

* Backup and restore Cassandra with [Amazon S3]({{< relref "/tasks/backup-restore/s3-compatible/amazon-s3.md" >}}).

* Backup and restore Cassandra with [GCS]({{< relref "/tasks/backup-restore/gcs" >}}).

For information about GCS, see the [Google Cloud Storage documentation](https://cloud.google.com/storage).
