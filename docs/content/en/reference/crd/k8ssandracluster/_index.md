---
title: "K8ssandraCluster CRD"
linkTitle: "K8ssandraCluster CRD"
no_list: true
toc_hide: true
simple_list: false
weight: 6
description: >
  K8ssandraCluster Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [CassandraClusterTemplate](#cassandraclustertemplate)
* [CassandraDatacenterTemplate](#cassandradatacentertemplate)
* [EmbeddedObjectMeta](#embeddedobjectmeta)
* [K8ssandraCluster](#k8ssandracluster)
* [K8ssandraClusterCondition](#k8ssandraclustercondition)
* [K8ssandraClusterList](#k8ssandraclusterlist)
* [K8ssandraClusterSpec](#k8ssandraclusterspec)
* [K8ssandraClusterStatus](#k8ssandraclusterstatus)
* [K8ssandraStatus](#k8ssandrastatus)

#### CassandraClusterTemplate



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| superuserSecretRef | The reference to the superuser secret to use for Cassandra. If unspecified, a default secret will be generated with a random password; the generated secret name will be \"<cluster_name>-superuser\" where <cluster_name> is the K8ssandraCluster CRD name. | corev1.LocalObjectReference | false |
| serverImage | ServerImage is the image for the cassandra container. Note that this should be a management-api image. If left empty the operator will choose a default image based on ServerVersion. | string | false |
| serverVersion | ServerVersion is the Cassandra version. | string | false |
| jmxInitContainerImage | The image to use in each Cassandra pod for the (short-lived) init container that enables JMX remote authentication on Cassandra pods. This is only useful when authentication is enabled in the cluster. The default is \"busybox:1.34.1\". | *images.Image | false |
| resources | Resources is the cpu and memory resources for the cassandra container. | *corev1.ResourceRequirements | false |
| config | CassandraConfig is configuration settings that are applied to cassandra.yaml and jvm-options for 3.11.x or jvm-server-options for 4.x. | *CassandraConfig | false |
| storageConfig | StorageConfig is the persistent storage requirements for each Cassandra pod. This includes everything under /var/lib/cassandra, namely the commit log and data directories. | *cassdcapi.StorageConfig | false |
| networking | Networking enables host networking and configures a NodePort ports. | *cassdcapi.NetworkingConfig | false |
| racks | Racks is a list of named racks. Note that racks are used to create node affinity. // | []cassdcapi.Rack | false |
| datacenters | Datacenters a list of the DCs in the cluster. | [][CassandraDatacenterTemplate](#cassandradatacentertemplate) | false |
| telemetry | Telemetry defines the desired state for telemetry resources in this K8ssandraCluster. If telemetry configurations are defined, telemetry resources will be deployed to integrate with a user-provided monitoring solution (at present, only support for Prometheus is available). | *telemetryapi.TelemetrySpec | false |
| mgmtAPIHeap | MgmtAPIHeap defines the amount of memory devoted to the management api heap. | *resource.Quantity | false |
| additionalSeeds | AdditionalSeeds specifies Cassandra node IPs for an existing datacenter. This is primarily intended for migrations from an existing Cassandra cluster that is not managed by k8ssandra-operator. Note that this property should NOT be used to set seeds for a DC that is or will be managed by k8ssandra-operator. k8ssandra-operator already manages seeds for DCs that it manages. If you have DNS set up such that you can resolve hostnames for the remote Cassandra cluster, then you can specify hostnames here; otherwise, use IP addresses. | []string | false |
| softPodAntiAffinity | SoftPodAntiAffinity sets whether multiple Cassandra instances can be scheduled on the same node. This should normally be false to ensure cluster resilience but may be set true for test/dev scenarios to minimise the number of nodes required. | *bool | false |
| tolerations | Tolerations applied to every Cassandra pod. | []corev1.Toleration | false |
| serverEncryptionStores | Internode encryption stores which are used by Cassandra and Stargate. | *encryption.Stores | false |
| clientEncryptionStores | Client encryption stores which are used by Cassandra and Reaper. | *encryption.Stores | false |

[Back to Custom Resources](#custom-resources)

#### CassandraDatacenterTemplate



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [EmbeddedObjectMeta](#embeddedobjectmeta) | false |
| k8sContext |  | string | false |
| serverImage |  | string | false |
| size | Size is the number Cassandra pods to deploy in this datacenter. This number does not include Stargate instances. | int32 | true |
| stopped | Stopped means that the datacenter will be stopped. Use this for maintenance or for cost saving. A stopped CassandraDatacenter will have no running server pods, like using \"stop\" with  traditional System V init scripts. Other Kubernetes resources will be left intact, and volumes will re-attach when the CassandraDatacenter workload is resumed. | bool | false |
| serverVersion | ServerVersion is the Cassandra version. | string | false |
| jmxInitContainerImage | The image to use in each Cassandra pod for the (short-lived) init container that enables JMX remote authentication on Cassandra pods. This is only useful when authentication is enabled in the cluster. The default is \"busybox:1.34.1\". | *images.Image | false |
| config | CassandraConfig is configuration settings that are applied to cassandra.yaml and jvm-options for 3.11.x or jvm-server-options for 4.x. | *CassandraConfig | false |
| resources | Resources is the cpu and memory resources for the cassandra container. | *corev1.ResourceRequirements | false |
| racks |  | []cassdcapi.Rack | false |
| networking | Networking enables host networking and configures a NodePort ports. | *cassdcapi.NetworkingConfig | false |
| storageConfig | StorageConfig is the persistent storage requirements for each Cassandra pod. This includes everything under /var/lib/cassandra, namely the commit log and data directories. | *cassdcapi.StorageConfig | false |
| stargate | Stargate defines the desired deployment characteristics for Stargate in this datacenter. Leave nil to skip deploying Stargate in this datacenter. | *stargateapi.StargateDatacenterTemplate | false |
| mgmtAPIHeap | MgmtAPIHeap defines the amount of memory devoted to the management api heap. | *resource.Quantity | false |
| telemetry | Telemetry defines the desired state for telemetry resources in this datacenter. If telemetry configurations are defined, telemetry resources will be deployed to integrate with a user-provided monitoring solution (at present, only support for Prometheus is available). | *telemetryapi.TelemetrySpec | false |
| softPodAntiAffinity | SoftPodAntiAffinity sets whether multiple Cassandra instances can be scheduled on the same node. This should normally be false to ensure cluster resilience but may be set true for test/dev scenarios to minimise the number of nodes required. | *bool | false |
| tolerations | Tolerations applied to every Cassandra pod. | []corev1.Toleration | false |

[Back to Custom Resources](#custom-resources)

#### EmbeddedObjectMeta



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| namespace |  | string | false |
| name |  | string | true |
| labels |  | map[string]string | false |
| annotations |  | map[string]string | false |

[Back to Custom Resources](#custom-resources)

#### K8ssandraCluster

K8ssandraCluster is the Schema for the k8ssandraclusters API. The K8ssandraCluster CRD name is also the name of the Cassandra cluster (which corresponds to cluster_name in cassandra.yaml).

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec |  | [K8ssandraClusterSpec](#k8ssandraclusterspec) | false |
| status |  | [K8ssandraClusterStatus](#k8ssandraclusterstatus) | false |

[Back to Custom Resources](#custom-resources)

#### K8ssandraClusterCondition



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type |  | K8ssandraClusterConditionType | true |
| status |  | corev1.ConditionStatus | true |
| lastTransitionTime | LastTransitionTime is the last time the condition transited from one status to another. | *metav1.Time | false |

[Back to Custom Resources](#custom-resources)

#### K8ssandraClusterList

K8ssandraClusterList contains a list of K8ssandraCluster

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][K8ssandraCluster](#k8ssandracluster) | true |

[Back to Custom Resources](#custom-resources)

#### K8ssandraClusterSpec

K8ssandraClusterSpec defines the desired state of K8ssandraCluster

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| auth | Whether to enable authentication in this cluster. The default is true; it is highly recommended to always leave authentication turned on. When enabled, authentication will be enforced not only on Cassandra nodes, but also on Reaper, Medusa and Stargate nodes, if any. | *bool | false |
| cassandra | Cassandra is a specification of the Cassandra cluster. This includes everything from the number of datacenters, the k8s cluster where each DC should be deployed, node affinity (via racks), individual C* node settings, JVM settings, and more. | *[CassandraClusterTemplate](#cassandraclustertemplate) | false |
| stargate | Stargate defines the desired deployment characteristics for Stargate in this K8ssandraCluster. If this is non-nil, Stargate will be deployed on every Cassandra datacenter in this K8ssandraCluster. | *stargateapi.StargateClusterTemplate | false |
| reaper | Reaper defines the desired deployment characteristics for Reaper in this K8ssandraCluster. If this is non-nil, Reaper will be deployed on every Cassandra datacenter in this K8ssandraCluster. | *reaperapi.ReaperClusterTemplate | false |
| medusa | Medusa defines the desired deployment characteristics for Medusa in this K8ssandraCluster. If this is non-nil, Medusa will be deployed in every Cassandra pod in this K8ssandraCluster. | *medusaapi.MedusaClusterTemplate | false |
| externalDatacenters | During a migration the operator should alter keyspaces replication settings including the following external DCs. This avoids removing replicas from datacenters which are outside of the operator scope (not referenced in the CR). Replication settings changes will only apply to system_* keyspaces as well as reaper_db and data_endpoint_auth (Stargate). | []string | false |

[Back to Custom Resources](#custom-resources)

#### K8ssandraClusterStatus

K8ssandraClusterStatus defines the observed state of K8ssandraCluster

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| conditions |  | [][K8ssandraClusterCondition](#k8ssandraclustercondition) | false |
| datacenters | Datacenters maps the CassandraDatacenter name to a K8ssandraStatus. The naming is a bit confusing but the mapping makes sense because we have a CassandraDatacenter and then define other components like Stargate and Reaper relative to it. I wanted to inline the field but when I do it won't serialize. | map[string][K8ssandraStatus](#k8ssandrastatus) | false |

[Back to Custom Resources](#custom-resources)

#### K8ssandraStatus

K8ssandraStatus defines the observed of a k8ssandra instance

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| decommissionProgress |  | DecommissionProgress | false |
| cassandra |  | *cassdcapi.CassandraDatacenterStatus | false |
| stargate |  | *stargateapi.StargateStatus | false |
| reaper |  | *reaperapi.ReaperStatus | false |

[Back to Custom Resources](#custom-resources)
