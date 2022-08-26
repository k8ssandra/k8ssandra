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

* local
* s3
* s3_compatible
* s3_rgw
* azure_blobs
* google_storage

## Deploying Medusa

You can deploy Medusa on all Cassandra datacenters in the cluster through the addition of the `medusa` section in the `K8ssandraCluster` definition. Example:

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

The definition above requires a secret named `medusa-bucket-key` to be created in the target namespace before the `K8ssandraCluster` object gets created. Use the following format for this secret: 

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

A successful deployment should inject a new init container named `medusa-restore` and a new container named `medusa` in the Cassandra StatefulSet pods.  

## Creating a Backup

To perform a backup of a Cassandra datacenter, create the following custom resource in the same namespace and Kubernetes cluster as the CassandraDatacenter resource, `cassandradatacenter/dc1` in this case :

```yaml
apiVersion: medusa.k8ssandra.io/v1alpha1
kind: MedusaBackupJob
metadata:
  name: medusa-backup1
spec:
  cassandraDatacenter: dc1
```

### Checking Backup Completion

K8ssandra Operator will detect the `MedusaBackupJob` object creation and trigger a backup asynchronously.

To monitor the backup completion, check if the `finishTime` is set in the `MedusaBackupJob` object status. Example:

```sh
% kubectl get medusabackupjob/medusa-backup1 -o yaml

kind: MedusaBackupJob
metadata:
    name: medusa-backup1
spec:
  cassandraDatacenter: dc1
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
At the end of the backup operation, a `MedusaBackup` custom resource will be created with the same name as the `MedusaBackupJob` object. It materializes the backup locally on the Kubernetes cluster.
For a restore to be possible, a `MedusaBackup` object must exist.


## Creating a Backup Schedule

K8ssandra-operator v1.2 introduced a new `MedusaBackupSchedule` CRD to manage backup schedules using a [cron expression](https://crontab.guru/):  

```yaml
apiVersion: medusa.k8ssandra.io/v1alpha1
kind: MedusaBackupSchedule
metadata:
  name: medusa-backup-schedule
  namespace: k8ssandra-operator
spec:
  backupSpec:
    backupType: differential
    cassandraDatacenter: dc1
  cronSchedule: 30 1 * * *
  disabled: false
```

This resource must be created in the same Kubernetes cluster and namespace as the `CassandraDatacenter` resource referenced in the spec, here `cassandradatacenter/dc1`.  
The above definition would trigger a differential backup of `dc1` every day at 1:30 AM. The status of the backup schedule will be updated with the last execution and next execution times:

```yaml
...
status:
  lastExecution: "2022-07-26T01:30:00Z"
  nextSchedule: "2022-07-27T01:30:00Z"
...
```

The `MedusaBackupJob` and `MedusaBackup` objects will be created with the name of the `MedusaBackupSchedule` object as prefix and a timestamp as suffix, for example: `medusa-backup-schedule-1658626200`.

## Restoring a Backup

To restore an existing backup for a Cassandra datacenter, create the following custom resource in the same namespace as the referenced CassandraDatacenter resource, `cassandradatacenter/dc1` in this case :

```yaml
apiVersion: medusa.k8ssandra.io/v1alpha1
kind: MedusaRestoreJob
metadata:
  name: restore-backup1
  namespace: k8ssandra-operator
spec:
  cassandraDatacenter: dc1
  backup: medusa-backup1
```

The `spec.backup` value should match the `MedusaBackup` `metadata.name` value.  
Once the K8ssandra Operator detects on the `MedusaRestoreJob` object creation, it will orchestrate the shutdown of all Cassandra pods, and the `medusa-restore` container will perform the actual data restore upon pod restart.

### Checking Restore Completion

To monitor the restore completion, check if the `finishTime` value isn't empty in the `MedusaRestoreJob` object status. Example:

```yaml
% kubectl get cassandrarestore/restore-backup1 -o yaml

apiVersion: medusa.k8ssandra.io/v1alpha1
kind: MedusaRestoreJob
metadata:
  name: restore-backup1
spec:
  backup: medusa-backup1
  cassandraDatacenter: dc1
status:
  datacenterStopped: "2022-01-06T16:45:09Z"
  finishTime: "2022-01-06T16:48:23Z"
  restoreKey: ec5b35c1-f2fe-4465-a74f-e29aa1d467ff
  restorePrepared: true
  startTime: "2022-01-06T16:44:53Z"
```

## Synchronizing MedusaBackup objects with a Medusa storage backend (S3, GCS, etc.)

In order to restore a backup taken on a different Cassandra cluster, a synchronization task must be executed to create the corresponding `MedusaBackup` objects locally.  
This can be achieved by creating a `MedusaTask` custom resource in the Kubernetes cluster and namespace where the referenced `CassandraDatacenter` was deployed, using a `sync` operation:

```yaml
apiVersion: medusa.k8ssandra.io/v1alpha1
kind: MedusaTask
metadata:
  name: sync-backups-1
  namespace: k8ssandra-operator
spec:
  cassandraDatacenter: dc1
  operation: sync
```

Such backups can come from Apache Cassandra clusters running outside of Kubernetes as well as clusters running in Kubernetes, as long as they were created using Medusa.  
**Warning:** backups created with K8ssandra-operator v1.0 and K8ssandra up to v1.5 are not suitable for remote restores due to pod name resolving issues in these versions.
  
Reconciliation will be triggered by the `MedusaTask` object creation, executing the following operations:

- Backups will be listed in the remote storage system
- Backups missing locally will be created
- Backups missing in the remote storage system will be deleted locally

Upon completion, the `MedusaTask` object status will be updated with the finish time and the name of the pod which was used to communicate with the storage backend:

```yaml
...
status:
  finishTime: '2022-07-26T08:15:55Z'
  finished:
    - podName: demo-dc2-default-sts-0
  startTime: '2022-07-26T08:15:54Z'
...
```

## Purging backups

Medusa has two settings to control the retention of backups: `max_backup_age` and `max_backup_count`.
These settings are used by the `medusa purge` operation to determine which backups to delete.
In order to trigger a purge, a `MedusaTask` custom resource should be created in the same Kubernetes cluster and namespace where the referenced `CassandraDatacenter` was created, using the `purge` operation:

```yaml
apiVersion: medusa.k8ssandra.io/v1alpha1
kind: MedusaTask
metadata:
  name: purge-backups-1
  namespace: k8ssandra-operator
spec:
  cassandraDatacenter: dc1
  operation: purge
```

The purge operation will be scheduled on all Cassandra pods in the datacenter and apply Medusa's purge rules.
Once the purge finishes, the `MedusaTask` object status will be updated with the finish time and the purge stats of each Cassandra pod:

```yaml
status:
  finishTime: '2022-07-26T08:42:33Z'
  finished:
    - nbBackupsPurged: 3
      nbObjectsPurged: 814
      podName: demo-dc2-default-sts-1
      totalObjectsWithinGcGrace: 542
      totalPurgedSize: 10770961
    - nbBackupsPurged: 3
      nbObjectsPurged: 852
      podName: demo-dc2-default-sts-2
      totalObjectsWithinGcGrace: 520
      totalPurgedSize: 10787447
    - nbBackupsPurged: 3
      nbObjectsPurged: 808
      podName: demo-dc2-default-sts-0
      totalObjectsWithinGcGrace: 444
      totalPurgedSize: 10903221
  startTime: '2022-07-26T08:37:48Z'
```

A `sync` task will then be generated by the `purge` task to delete the purged backups from the local Kubernetes storage.

In order to schedule purge tasks, a `CronJob` resource using the following template can be used:

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: k8ssandra-medusa-backup
  namespace: k8ssandra-operator
spec:
  schedule: "0 0 * * *"
  jobTemplate:
    spec:
      template:
        metadata:
          name: k8ssandra-medusa-backup
        spec:
          serviceAccountName: medusa-backup
          containers:
          - name: medusa-backup-cronjob
            image: bitnami/kubectl:1.17.3
            imagePullPolicy: IfNotPresent
            command:
             - 'bin/bash'
             - '-c'
             - 'printf "apiVersion: medusa.k8ssandra.io/v1alpha1\nkind: MedusaTask\nmetadata:\n  name: purge-backups-timestamp\n  namespace: k8ssandra-operator\nspec:\n  cassandraDatacenter: dc1\n  operation: purge" | sed "s/timestamp/$(date +%Y%m%d%H%M%S)/g" | kubectl apply -f -'
          restartPolicy: OnFailure
```

The above `CronJob` will be scheduled to run every day at midnight, and trigger a `MedusaTask` object creation to purge the backups.

## Deprecation Notice

The `CassandraBackup` and `CassandraRestore` CRDs are deprecated in K8ssandra-operator v1.1 and will be removed in a future version of the operator.

## Next steps

See the [Custom Resource Definition (CRD) reference]({{< relref "/reference/crd" >}}) topics.
