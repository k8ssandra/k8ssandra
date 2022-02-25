---
title: "Reaper CRD"
linkTitle: "Reaper CRD"
simple_list: false
weight: 6
description: >
  Reaper Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [AutoScheduling](#autoscheduling)
* [CassandraDatacenterRef](#cassandradatacenterref)
* [Reaper](#reaper)
* [ReaperClusterTemplate](#reaperclustertemplate)
* [ReaperCondition](#reapercondition)
* [ReaperList](#reaperlist)
* [ReaperSpec](#reaperspec)
* [ReaperStatus](#reaperstatus)
* [ReaperTemplate](#reapertemplate)

#### AutoScheduling

AutoScheduling includes options to configure the auto scheduling of repairs for new clusters.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| enabled |  | bool | false |
| repairType | RepairType is the type of repair to create: - REGULAR creates a regular repair (non-adaptive and non-incremental); - ADAPTIVE creates an adaptive repair; adaptive repairs are most suited for Cassandra 3. - INCREMENTAL creates an incremental repair; incremental repairs should only be used with Cassandra 4+. - AUTO chooses between ADAPTIVE and INCREMENTAL depending on the Cassandra server version; ADAPTIVE for Cassandra 3 and INCREMENTAL for Cassandra 4+. | string | false |
| percentUnrepairedThreshold | PercentUnrepairedThreshold is the percentage of unrepaired data over which an incremental repair should be started. Only relevant when using repair type INCREMENTAL. | int | false |
| initialDelayPeriod | InitialDelay is the amount of delay time before the schedule period starts. Must be a valid ISO-8601 duration string. The default is \"PT15S\" (15 seconds). | string | false |
| periodBetweenPolls | PeriodBetweenPolls is the interval time to wait before checking whether to start a repair task. Must be a valid ISO-8601 duration string. The default is \"PT10M\" (10 minutes). | string | false |
| timeBeforeFirstSchedule | TimeBeforeFirstSchedule is the grace period before the first repair in the schedule is started. Must be a valid ISO-8601 duration string. The default is \"PT5M\" (5 minutes). | string | false |
| scheduleSpreadPeriod | ScheduleSpreadPeriod is the time spacing between each of the repair schedules that is to be carried out. Must be a valid ISO-8601 duration string. The default is \"PT6H\" (6 hours). | string | false |
| excludedClusters | ExcludedClusters are the clusters that are to be excluded from the repair schedule. | []string | false |
| excludedKeyspaces | ExcludedKeyspaces are the keyspaces that are to be excluded from the repair schedule. | []string | false |

[Back to Custom Resources](#custom-resources)

#### CassandraDatacenterRef

CassandraDatacenterRef references the target Cassandra DC that Reaper should manage.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | The datacenter name. | string | true |
| namespace | The datacenter namespace. If empty, the datacenter will be assumed to reside in the same namespace as the Reaper instance. | string | false |

[Back to Custom Resources](#custom-resources)

#### Reaper

Reaper is the Schema for the reapers API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec |  | [ReaperSpec](#reaperspec) | false |
| status |  | [ReaperStatus](#reaperstatus) | false |

[Back to Custom Resources](#custom-resources)

#### ReaperClusterTemplate



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| deploymentMode |  | string | false |

[Back to Custom Resources](#custom-resources)

#### ReaperCondition



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type |  | ReaperConditionType | true |
| status |  | corev1.ConditionStatus | true |
| lastTransitionTime | LastTransitionTime is the last time the condition transited from one status to another. | *metav1.Time | false |

[Back to Custom Resources](#custom-resources)

#### ReaperList

ReaperList contains a list of Reaper

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][Reaper](#reaper) | true |

[Back to Custom Resources](#custom-resources)

#### ReaperSpec

ReaperSpec defines the desired state of Reaper

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| datacenterRef | DatacenterRef is the reference of a CassandraDatacenter resource that this Reaper instance should manage. It will also be used as the backend for persisting Reaper's state. Reaper must be able to access the JMX port (7199 by default) and the CQL port (9042 by default) on this DC. | [CassandraDatacenterRef](#cassandradatacenterref) | true |
| datacenterAvailability | DatacenterAvailability indicates to Reaper its deployment in relation to the target datacenter's network. For single-DC clusters, the default (ALL) is fine. For multi-DC clusters, it is recommended to use EACH, provided that there is one Reaper instance managing each DC in the cluster; otherwise, if one single Reaper instance is going to manage more than one DC in the cluster, use ALL. See https://cassandra-reaper.io/docs/usage/multi_dc/. | string | false |
| clientEncryptionStores | Client encryption stores which are used by Cassandra and Reaper. | *encryption.Stores | false |
| skipSchemaMigration | Whether to skip schema migration. Schema migration is done in an init container on every Reaper deployment and can slow down Reaper's startup time. Besides, schema migration requires reading data at QUORUM. It can be skipped if you know that the schema is already up-to-date, or if you know upfront that QUORUM cannot be achieved (for example, because a DC is down). | bool | false |

[Back to Custom Resources](#custom-resources)

#### ReaperStatus

ReaperStatus defines the observed state of Reaper

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| progress | Progress is the progress of this Reaper object. | ReaperProgress | false |
| conditions |  | [][ReaperCondition](#reapercondition) | false |

[Back to Custom Resources](#custom-resources)

#### ReaperTemplate



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| keyspace | The keyspace to use to store Reaper's state. Will default to \"reaper_db\" if unspecified. Will be created if it does not exist, and if this Reaper resource is managed by K8ssandra. | string | false |
| cassandraUserSecretRef | Defines the username and password that Reaper will use to authenticate CQL connections to Cassandra clusters. These credentials will be automatically turned into CQL roles by cass-operator when bootstrapping the datacenter, then passed to the Reaper instance, so that it can authenticate against nodes in the datacenter using CQL. If CQL authentication is not required, leave this field empty. The secret must be in the same namespace as Reaper itself and must contain two keys: \"username\" and \"password\". | corev1.LocalObjectReference | false |
| jmxUserSecretRef | Defines the username and password that Reaper will use to authenticate JMX connections to Cassandra clusters. These credentials will be automatically passed to each Cassandra node in the datacenter, as well as to the Reaper instance, so that the latter can authenticate against the former. If JMX authentication is not required, leave this field empty. The secret must be in the same namespace as Reaper itself and must contain two keys: \"username\" and \"password\". | corev1.LocalObjectReference | false |
| uiUserSecretRef | Defines the secret which contains the username and password for the Reaper UI and REST API authentication. | corev1.LocalObjectReference | false |
| containerImage | The image to use for the Reaper pod main container. The default is \"thelastpickle/cassandra-reaper:3.1.1\". | *images.Image | false |
| initContainerImage | The image to use for the Reaper pod init container (that performs schema migrations). The default is \"thelastpickle/cassandra-reaper:3.1.1\". | *images.Image | false |
| ServiceAccountName |  | string | false |
| autoScheduling | Auto scheduling properties. When you enable the auto-schedule feature, Reaper dynamically schedules repairs for all non-system keyspaces in a cluster. A cluster's keyspaces are monitored and any modifications (additions or removals) are detected. When a new keyspace is created, a new repair schedule is created automatically for that keyspace. Conversely, when a keyspace is removed, the corresponding repair schedule is deleted. | [AutoScheduling](#autoscheduling) | false |
| livenessProbe | LivenessProbe sets the Reaper liveness probe. Leave nil to use defaults. | *corev1.Probe | false |
| readinessProbe | ReadinessProbe sets the Reaper readiness probe. Leave nil to use defaults. | *corev1.Probe | false |
| affinity | Affinity applied to the Reaper pods. | *corev1.Affinity | false |
| tolerations | Tolerations applied to the Reaper pods. | []corev1.Toleration | false |
| podSecurityContext | PodSecurityContext contains a pod-level SecurityContext to apply to Reaper pods. | *corev1.PodSecurityContext | false |
| securityContext | SecurityContext applied to the Reaper main container. | *corev1.SecurityContext | false |
| initContainerSecurityContext | InitContainerSecurityContext is the SecurityContext applied to the Reaper init container, used to perform schema migrations. | *corev1.SecurityContext | false |

[Back to Custom Resources](#custom-resources)
