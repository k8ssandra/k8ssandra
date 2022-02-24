---
title: "CassandraRestore CRD"
linkTitle: "CassandraRestore CRD"
simple_list: false
weight: 6
description: >
  CassandraRestore Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [CassandraDatacenterConfig](#cassandradatacenterconfig)
* [CassandraRestore](#cassandrarestore)
* [CassandraRestoreList](#cassandrarestorelist)
* [CassandraRestoreSpec](#cassandrarestorespec)
* [CassandraRestoreStatus](#cassandrarestorestatus)

#### CassandraDatacenterConfig



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | The name to give the new, restored CassandraDatacenter | string | true |
| clusterName | The name to give the C* cluster. | string | true |

[Back to Custom Resources](#custom-resources)

#### CassandraRestore

CassandraRestore is the Schema for the cassandrarestores API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec |  | [CassandraRestoreSpec](#cassandrarestorespec) | false |
| status |  | [CassandraRestoreStatus](#cassandrarestorestatus) | false |

[Back to Custom Resources](#custom-resources)

#### CassandraRestoreList

CassandraRestoreList contains a list of CassandraRestore

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][CassandraRestore](#cassandrarestore) | true |

[Back to Custom Resources](#custom-resources)

#### CassandraRestoreSpec

CassandraRestoreSpec defines the desired state of CassandraRestore

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| backup | The name of the CassandraBackup to restore | string | true |
| inPlace | When true the restore will be performed on the source cluster from which the backup was taken. There will be a rolling restart of the source cluster. | bool | true |
| shutdown | When set to true, the cluster is shutdown before the restore is applied. This is necessary process if there are schema changes between the backup and current schema. Recommended. | bool | true |
| cassandraDatacenter |  | [CassandraDatacenterConfig](#cassandradatacenterconfig) | true |

[Back to Custom Resources](#custom-resources)

#### CassandraRestoreStatus

CassandraRestoreStatus defines the observed state of CassandraRestore

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| restoreKey | A unique key that identifies the restore operation. | string | true |
| startTime |  | metav1.Time | false |
| finishTime |  | metav1.Time | false |
| datacenterStopped |  | metav1.Time | false |
| inProgress |  | []string | false |
| finished |  | []string | false |
| failed |  | []string | false |

[Back to Custom Resources](#custom-resources)
