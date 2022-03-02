---
title: "Medusa CRD"
linkTitle: "Medusa CRD"
no_list: true
toc_hide: true
simple_list: false
weight: 6
description: >
  Medusa Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [MedusaClusterTemplate](#medusaclustertemplate)
* [PodStorageSettings](#podstoragesettings)
* [Storage](#storage)

#### MedusaClusterTemplate



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| containerImage | MedusaContainerImage is the image characteristics to use for Medusa containers. Leave nil to use a default image. | *images.Image | false |
| securityContext | SecurityContext applied to the Medusa containers. | *corev1.SecurityContext | false |
| cassandraUserSecretRef | Defines the username and password that Medusa will use to authenticate CQL connections to Cassandra clusters. These credentials will be automatically turned into CQL roles by cass-operator when bootstrapping the datacenter, then passed to the Medusa instances, so that it can authenticate against nodes in the datacenter using CQL. The secret must be in the same namespace as Cassandra and must contain two keys: \"username\" and \"password\". | corev1.LocalObjectReference | false |
| storageProperties | Provides all storage backend related properties for backups. | [Storage](#storage) | false |

[Back to Custom Resources](#custom-resources)

#### PodStorageSettings



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| storageClassName | Storage class name to use for the pod's storage. | string | false |
| size | Size of the pod's storage in bytes. Defaults to 10 GB. | resource.Quantity | false |
| accessModes | Pod local storage access modes | []corev1.PersistentVolumeAccessMode | false |

[Back to Custom Resources](#custom-resources)

#### Storage



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| storageProvider | The storage backend to use for the backups. | string | false |
| storageSecretRef | Kubernetes Secret that stores the key file for the storage provider's API. If using 'local' storage, this value is ignored. | corev1.LocalObjectReference | false |
| bucketName | The name of the bucket to use for the backups. | string | false |
| prefix | Name of the top level folder in the backup bucket. If empty, the cluster name will be used. | string | false |
| maxBackupAge | Maximum backup age that the purge process should observe. | int | false |
| maxBackupCount | Maximum number of backups to keep (used by the purge process). Default is unlimited. | int | false |
| apiProfile | AWS Profile to use for authentication. | string | false |
| transferMaxBandwidth | Max upload bandwidth in MB/s. Defaults to 50 MB/s. | string | false |
| concurrentTransfers | Number of concurrent uploads. Helps maximizing the speed of uploads but puts more pressure on the network. Defaults to 1. | int | false |
| multiPartUploadThreshold | File size over which cloud specific cli tools are used for transfer. Defaults to 100 MB. | int | false |
| host | Host to connect to for the storage backend. | string | false |
| region | Region of the storage bucket. Defaults to \"default\". | string | false |
| port | Port to connect to for the storage backend. | int | false |
| secure | Whether to use SSL for the storage backend. | bool | false |
| backupGracePeriodInDays | Age after which orphan sstables can be deleted from the storage backend. Protects from race conditions between purge and ongoing backups. Defaults to 10 days. | int | false |
| podStorage | Pod storage settings for the local storage provider | *[PodStorageSettings](#podstoragesettings) | false |

[Back to Custom Resources](#custom-resources)
