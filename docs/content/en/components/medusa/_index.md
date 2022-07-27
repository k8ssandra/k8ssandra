---
title: "Medusa for Cassandra backup and restore"
linkTitle: "Medusa"
weight: 4
description: K8ssandra Operator deploys Medusa to support backup and restore operations of Apache Cassandra&reg; tables.
---

Medusa for Apache Cassandra&reg; is deployed by a K8ssandra Operator install, based on the [Medusa]({{< relref "/reference/crd" >}}) Custom Resource Definition (CRD).

If you haven't already installed a K8ssandraCluster using K8ssandra Operator, see the [local install]({{< relref "/install/local" >}}) topic.

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

## Next steps

* Backup and restore Cassandra with S3-compatible [MinIO]({{< relref "/tasks/backup-restore/minio/" >}}).
* Backup and restore Cassandra with [Amazon S3]({{< relref "/tasks/backup-restore/amazon-s3/" >}}).
* Backup and restore Cassandra with [Google Cloud Storage]({{< relref "/tasks/backup-restore/gcs/" >}}).
* For information about using a superuser and secrets with Medusa authentication, see [Medusa security]({{< relref "/tasks/secure/#medusa-security" >}}).
* For reference details, see the K8ssandra Operator [Custom Resource Definitions (CRDs)]({{< relref "reference" >}}).
* Also see the topics covering other [components]({{< relref "/components/" >}}) deployed by K8ssandra Operator. 
* For information on using other deployed components, see the [Tasks]({{< relref "/tasks/" >}}) topics.
