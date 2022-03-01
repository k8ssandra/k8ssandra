---
title: "Backup and restore Cassandra data"
linkTitle: "Backup/restore"
no_list: true
weight: 4
description: Use Medusa to backup and restore Apache Cassandra® data in Kubernetes.
---

Medusa is a Cassandra backup and restore tool. It's packaged with K8ssandra Operator and supports a variety of backends. 

These instructions use a local `minio` bucket as an example.

## Supported object storage types for backups

Supported in K8ssandra Operator's Medusa:

* local (`host: minio.minio.svc.cluster.local` in this topic)
* s3
* s3_compatible
* s3_rgw
* azure_blobs
* google_storage

# Deploying Medusa

You can deploy Medusa on all Cassandra datacenters in the cluster through the addition of settings in the `K8ssandraCluster` definition. Example:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    ...
    ...
  medusa:
    containerImage:
      registry: docker.io
      repository: k8ssandra
      tag: 0.11.3
    storageProperties:
      # Can be either of local, google_storage, azure_blobs, s3, s3_compatible, s3_rgw or ibm_storage 
      storageProvider: s3_compatible
      # Name of the secret containing the credentials file to access the backup storage backend
      storageSecretRef:
        name: medusa-bucket-key
      # Name of the storage bucket
      bucketName: k8ssandra-medusa
      # Prefix for this cluster in the storage bucket directory structure, used for multitenancy
      prefix: test
      # Host to connect to the storage backend (Omitted for GCS, S3, Azure and local).
      host: minio.minio.svc.cluster.local
      # Port to connect to the storage backend (Omitted for GCS, S3, Azure and local).
      port: 9000
      # Region of the storage bucket
      # region: us-east-1
      
      # Whether or not to use SSL to connect to the storage backend
      secure: false 
      
      # Maximum backup age that the purge process should observe.
      # 0 equals unlimited
      # maxBackupAge: 0

      # Maximum number of backups to keep (used by the purge process).
      # 0 equals unlimited
      # maxBackupCount: 0

      # AWS Profile to use for authentication.
      # apiProfile: 
      # transferMaxBandwidth: 50MB/s

      # Number of concurrent uploads.
      # Helps maximizing the speed of uploads but puts more pressure on the network.
      # Defaults to 1.
      # concurrentTransfers: 1
      
      # File size in bytes over which cloud specific cli tools are used for transfer.
      # Defaults to 100 MB.
      # multiPartUploadThreshold: 104857600
      
      # Age after which orphan sstables can be deleted from the storage backend.
      # Protects from race conditions between purge and ongoing backups.
      # Defaults to 10 days.
      # backupGracePeriodInDays: 10
      
      # Pod storage settings to use for local storage (testing only)
      # podStorage:
      #   storageClassName: standard
      #   accessModes:
      #     - ReadWriteOnce
      #   size: 100Mi
```

The definition above requires a `medusa-bucket-key` to be created in the target namespace before the `K8ssandraCluster` object gets created. Use the following format for this secret: 

```yaml
apiVersion: v1
kind: Secret
metadata:
 name: medusa-bucket-key
type: Opaque
stringData:
 # Note that this currently has to be set to credentials!
 credentials: |-
   [default]
   aws_access_key_id = minio_key
   aws_secret_access_key = minio_secret
```

The file should always specify `credentials` as shown in the example above; in that section, provide the expected format and credential values that are expected by Medusa for the chosen storage backend. For more, refer to the [Medusa documentation](https://github.com/thelastpickle/cassandra-medusa/blob/master/docs/Installation.md) to know which file format should used for each supported storage backend.

A successful deployment should inject a new init container named `medusa-restore` and a new container named `medusa` in the Cassandra STS pods.  

# Creating a Backup

To perform a backup of a Cassandra datacenter, create the following custom resource in the namespace where K8ssandra was deployed:

```yaml
apiVersion: medusa.k8ssandra.io/v1alpha1
kind: CassandraBackup
metadata:
  name: medusa-backup1
spec:
  cassandraDatacenter: dc1
  name: medusa-backup1
```

The `metadata.name` value can match the `spec.name` value for convenience, but it is not mandatory. The latter will be used to identify the backup in the storage backend, the former being the name of the `CassandraBackup` custom resource in Kubernetes.

## Checking Backup Completion

K8ssandra Operator will detect the `CassandraBackup` object creation and trigger a backup asynchronously.

To monitor the backup completion, check if the `finishTime` value isn't empty in the CassandraBackup object status. Example:

```sh
% kubectl get cassandrabackup/medusa-backup1 -o yaml

kind: CassandraBackup
metadata:
    name: medusa-backup1
spec:
  backupType: differential
  cassandraDatacenter: dc1
  name: medusa-backup1
status:
  ...
  ...
  finishTime: "2022-01-06T16:34:35Z"
  finished:
  - demo-dc1-default-sts-0
  - demo-dc1-default-sts-1
  - demo-dc1-default-sts-2
  startTime: "2022-01-06T16:34:30Z"

```

All pods having completed the backup will be in the `finished` list.

# Restoring a Backup

To restore an existing backup for a Cassandra datacenter, create the following custom resource in the namespace where K8ssandra was deployed. Example:

```yaml
apiVersion: medusa.k8ssandra.io/v1alpha1
kind: CassandraRestore
metadata:
  name: restore-backup1
  namespace: k8ssandra-operator
spec:
  cassandraDatacenter: 
    name: dc1
    clusterName: demo
  backup: medusa-backup1
  inPlace: true
  shutdown: true
```

The `spec.backup` value should match the CassandraBackup `spec.name` value.  
Once the K8ssandra Operator detects on the `CassandraRestore` object creation, it will control the shutdown of all Cassandra pods, and the `medusa-restore` container will perform the actual data restore upon pod restart.

## Checking Restore Completion

To monitor the restore completion, check if the `finishTime` value isn't empty in the `CassandraRestore` object status. Example:

```sh
% kubectl get cassandrarestore/restore-backup1 -o yaml

apiVersion: medusa.k8ssandra.io/v1alpha1
kind: CassandraRestore
metadata:
  name: restore-backup1
spec:
  backup: medusa-backup1
  cassandraDatacenter:
    clusterName: demo
    name: dc1
  inPlace: true
  shutdown: true
status:
  datacenterStopped: "2022-01-06T16:45:09Z"
  finishTime: "2022-01-06T16:48:23Z"
  restoreKey: ec5b35c1-f2fe-4465-a74f-e29aa1d467ff
  startTime: "2022-01-06T16:44:53Z"
```

## Next steps

See the following Custom Resource Definition (CRD) reference topics:

* [Medusa CRD]({{< relref "/reference/crd/medusa" >}})
* [CassandraBackup CRD]({{< relref "/reference/crd/cassandrabackup" >}})
* [CassansdraRestore CRD]({{< relref "/reference/crd/cassandrarestore" >}})
