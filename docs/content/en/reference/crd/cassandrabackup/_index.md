---
title: "CassandraBackup CRD"
linkTitle: "CassandraBackup CRD"
simple_list: false
weight: 6
description: >
  CassandraBackup Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources

* [CassandraBackup](#cassandrabackup)
* [CassandraBackupList](#cassandrabackuplist)
* [CassandraBackupSpec](#cassandrabackupspec)
* [CassandraBackupStatus](#cassandrabackupstatus)
* [CassandraDatacenterTemplateSpec](#cassandradatacentertemplatespec)

#### CassandraBackup

CassandraBackup is the Schema for the cassandrabackups API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec |  | [CassandraBackupSpec](#cassandrabackupspec) | false |
| status |  | [CassandraBackupStatus](#cassandrabackupstatus) | false |

[Back to Custom Resources](#custom-resources)

#### CassandraBackupList

CassandraBackupList contains a list of CassandraBackup

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][CassandraBackup](#cassandrabackup) | true |

[Back to Custom Resources](#custom-resources)

#### CassandraBackupSpec

CassandraBackupSpec defines the desired state of CassandraBackup

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | The name of the backup. | string | false |
| cassandraDatacenter | The name of the CassandraDatacenter to back up | string | true |
| backupType | The type of the backup: \"full\" or \"differential\" | BackupType | false |

[Back to Custom Resources](#custom-resources)

#### CassandraBackupStatus

CassandraBackupStatus defines the observed state of CassandraBackup

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| cassdcTemplateSpec |  | *[CassandraDatacenterTemplateSpec](#cassandradatacentertemplatespec) | false |
| startTime |  | metav1.Time | false |
| finishTime |  | metav1.Time | false |
| inProgress |  | []string | false |
| finished |  | []string | false |
| failed |  | []string | false |

[Back to Custom Resources](#custom-resources)

#### CassandraDatacenterTemplateSpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata | Standard object metadata | metav1.ObjectMeta | false |
| spec |  | cassdcapi.CassandraDatacenterSpec | true |

[Back to Custom Resources](#custom-resources)
