---
title: "Backup and restore Cassandra data"
linkTitle: "Backup/restore"
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

Added in K8ssandra 1.3.0:

* Support in K8ssandra's Medusa for Azure Storage

## Backup and restore steps

For detailed walk-throughs of Medusa backup and restore operations, see:

* Backup and restore Cassandra with S3-compatible [MinIO]({{< relref "/tasks/backup-restore/minio/" >}}).
* Backup and restore Cassandra with [Amazon S3]({{< relref "/tasks/backup-restore/amazon-s3/" >}}).
* Backup and restore Cassandra with [GCS]({{< relref "/tasks/backup-restore/gcs/" >}}).
* Backup and restore Cassandra with [Azure]({{< relref "/tasks/backup-restore/azure/" >}}).

## Next steps

Also see the following reference topics:

* [Medusa Operator Helm Chart]({{< relref "/reference/helm-charts/medusa-operator" >}})
* [Backup Helm Chart]({{< relref "/reference/helm-charts/backup" >}})
* [Restore Helm Chart]({{< relref "/reference/helm-charts/restore" >}})
