---
title: "Backup and restore Cassandra data"
linkTitle: "Backup/restore"
no_list: true
weight: 4
description: Use Medusa to backup and restore Apache Cassandra® data in Kubernetes.
---

Medusa is a Cassandra backup and restore tool. It's packaged with K8ssandra Operator and supports a variety of backends. 

## Supported object storage types for backups

Supported in K8ssandra Operator's Medusa:

* Amazon S3  
* Google Cloud Storage (GCS)
* MinIO 
* IBM Cloud Object Storage
* OVHCloud Object Storage
* Riak S2
* Dell EMC ECS
* CEPH Object Gateway
* Azure Storage

## Backup and restore steps

For detailed walk-throughs of Medusa backup and restore operations, see:

* Backup and restore Cassandra with S3-compatible [MinIO]({{< relref "/tasks/backup-restore/minio/" >}}).
* Backup and restore Cassandra with [Amazon S3]({{< relref "/tasks/backup-restore/amazon-s3/" >}}).
* Backup and restore Cassandra with [GCS]({{< relref "/tasks/backup-restore/gcs/" >}}).
* Backup and restore Cassandra with [Azure]({{< relref "/tasks/backup-restore/azure/" >}}).

## Next steps

See the following Custom Resource Definition (CRD) reference topics:

* [Medusa CRD]({{< relref "/reference/crd/medusa" >}})
* [CassandraBackup CRD]({{< relref "/reference/crd/cassandrabackup" >}})
* [CassansdraRestore CRD]({{< relref "/reference/crd/cassandrarestore" >}})
