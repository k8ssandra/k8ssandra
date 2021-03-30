---
title: "Backup and restore Apache Cassandra data"
linkTitle: "Backup and restore Cassandra"
no_list: true
weight: 4
description: Use Medusa to backup and restore Apache Cassandra® data in Kubernetes.
---

Medusa is a Cassandra backup and restore tool. It's packaged with K8ssandra and supports a variety of backends. 

## Supported object storage types for backups

Supported in K8ssandra since the 1.0.0 release:

* Amazon S3  

* Google Cloud Storage (GCS)

Added in K8ssandra 1.1.0:

* MinIO, a popular and S3-compatible object storage suite

## Backup and restore steps

See the following topics in our K8ssandra documentation:

* [Backup and restore Cassandra with Amazon S3]({{< relref "amazon-s3" >}})

* [Backup and restore Cassandra with MinIO]({{< relref "minio" >}})

For information about GCS, see the [Google Cloud Storage documentation](https://cloud.google.com/storage).
