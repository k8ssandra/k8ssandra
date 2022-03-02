---
title: "Stargate CRD"
linkTitle: "Stargate CRD"
no_list: true
toc_hide: true
simple_list: false
weight: 6
description: >
  Stargate Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [CassandraEncryption](#cassandraencryption)
* [Stargate](#stargate)
* [StargateClusterTemplate](#stargateclustertemplate)
* [StargateCondition](#stargatecondition)
* [StargateDatacenterTemplate](#stargatedatacentertemplate)
* [StargateList](#stargatelist)
* [StargateRackTemplate](#stargateracktemplate)
* [StargateSpec](#stargatespec)
* [StargateStatus](#stargatestatus)
* [StargateTemplate](#stargatetemplate)

#### CassandraEncryption

Still it is required to pass the encryption stores secrets to the Stargate pods, so that they can be mounted as volumes.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| clientEncryptionStores | Client encryption stores which are used by Cassandra and Reaper. | *encryption.Stores | false |
| serverEncryptionStores | Internode encryption stores which are used by Cassandra and Stargate. | *encryption.Stores | false |

[Back to Custom Resources](#custom-resources)

#### Stargate

Stargate is the Schema for the stargates API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec | Specification of the desired behavior of this Stargate resource. | [StargateSpec](#stargatespec) | false |
| status | Most recently observed status of this Stargate resource. | [StargateStatus](#stargatestatus) | false |

[Back to Custom Resources](#custom-resources)

#### StargateClusterTemplate

StargateClusterTemplate defines global rules to apply to all Stargate pods in all datacenters in the cluster. These rules will be merged with rules defined at datacenter level in a StargateDatacenterTemplate; dc-level rules have precedence over cluster-level ones.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| size | Size is the number of Stargate instances to deploy in each datacenter. They will be spread evenly across racks. | int32 | true |

[Back to Custom Resources](#custom-resources)

#### StargateCondition



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type |  | StargateConditionType | true |
| status |  | corev1.ConditionStatus | true |
| lastTransitionTime | LastTransitionTime is the last time the condition transited from one status to another. | *metav1.Time | false |

[Back to Custom Resources](#custom-resources)

#### StargateDatacenterTemplate

StargateDatacenterTemplate defines rules to apply to all Stargate pods in a given datacenter. These rules will be merged with rules defined at rack level in a StargateRackTemplate; rack-level rules have precedence over datacenter-level ones.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| racks | Racks allow customizing Stargate characteristics for specific racks in the datacenter. | [][StargateRackTemplate](#stargateracktemplate) | false |

[Back to Custom Resources](#custom-resources)

#### StargateList

StargateList contains a list of Stargate

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][Stargate](#stargate) | true |

[Back to Custom Resources](#custom-resources)

#### StargateRackTemplate

StargateRackTemplate defines custom rules for Stargate pods in a given rack. These rules will be merged with rules defined at datacenter level in a StargateDatacenterTemplate; rack-level rules have precedence over datacenter-level ones.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Name is the rack name. It must correspond to an existing rack name in the CassandraDatacenter resource where Stargate is being deployed, otherwise it will be ignored. | string | true |

[Back to Custom Resources](#custom-resources)

#### StargateSpec

StargateSpec defines the desired state of a Stargate resource.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| datacenterRef | DatacenterRef is the namespace-local reference of a CassandraDatacenter resource where Stargate should be deployed. | corev1.LocalObjectReference | true |
| auth | Whether to enable authentication for Stargate. The default is true; it is highly recommended to always leave authentication turned on, not only on Stargate nodes, but also on data nodes as well. Note that Stargate REST APIs are currently only accessible if authentication is enabled, and if the authenticator in use in the whole cluster is PasswordAuthenticator. The usage of any other authenticator will cause the REST API to become inaccessible, see https://github.com/stargate/stargate/issues/792 for more. Stargate CQL API however remains accessible even if authentication is disabled in the cluster, or when a custom authenticator is being used. | *bool | false |
| cassandraEncryption |  | *[CassandraEncryption](#cassandraencryption) | false |

[Back to Custom Resources](#custom-resources)

#### StargateStatus

StargateStatus defines the observed state of a Stargate resource.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| progress | Progress is the progress of this Stargate object. | StargateProgress | false |
| conditions |  | [][StargateCondition](#stargatecondition) | false |
| deploymentRefs | DeploymentRefs is the names of the Deployment objects that were created for this Stargate object. | []string | false |
| serviceRef | ServiceRef is the name of the Service object that was created for this Stargate object. | *string | false |
| readyReplicasRatio | ReadyReplicasRatio is a \"X/Y\" string representing the ratio between ReadyReplicas and Replicas in the Stargate deployment. | *string | false |
| replicas | Total number of non-terminated pods targeted by the Stargate deployment (their labels match the selector). Will be zero if the deployment has not been created yet. | int32 | true |
| readyReplicas | ReadyReplicas is the total number of ready pods targeted by the Stargate deployment. Will be zero if the deployment has not been created yet. | int32 | true |
| updatedReplicas | UpdatedReplicas is the total number of non-terminated pods targeted by the Stargate deployment that have the desired template spec. Will be zero if the deployment has not been created yet. | int32 | true |
| availableReplicas | Total number of available pods targeted by the Stargate deployment. Will be zero if the deployment has not been created yet. | int32 | true |

[Back to Custom Resources](#custom-resources)

#### StargateTemplate

StargateTemplate defines a template for deploying Stargate.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| containerImage | ContainerImage is the image characteristics to use for Stargate containers. Leave nil to use a default image. | *images.Image | false |
| serviceAccount | ServiceAccount is the service account name to use for Stargate pods. | *string | false |
| resources | Resources is the Kubernetes resource requests and limits to apply, per Stargate pod. Leave nil to use defaults. | *corev1.ResourceRequirements | false |
| heapSize | HeapSize sets the JVM heap size to use for Stargate. If no Resources are specified, this value will also be used to set a default memory request and limit for the Stargate pods: these will be set to HeapSize x2 and x4, respectively. | *resource.Quantity | false |
| livenessProbe | LivenessProbe sets the Stargate liveness probe. Leave nil to use defaults. | *corev1.Probe | false |
| readinessProbe | ReadinessProbe sets the Stargate readiness probe. Leave nil to use defaults. | *corev1.Probe | false |
| nodeSelector | NodeSelector is an optional map of label keys and values to restrict the scheduling of Stargate nodes to workers with matching labels. Leave nil to let the controller reuse the same node selectors used for data pods in this datacenter, if any. See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector | map[string]string | false |
| tolerations | Tolerations are tolerations to apply to the Stargate pods. Leave nil to let the controller reuse the same tolerations used for data pods in this datacenter, if any. See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/ | []corev1.Toleration | false |
| affinity | Affinity is the affinity to apply to all the Stargate pods. Leave nil to let the controller reuse the same affinity rules used for data pods in this datacenter, if any. See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity | *corev1.Affinity | false |
| allowStargateOnDataNodes | AllowStargateOnDataNodes allows Stargate pods to be scheduled on a worker node already hosting data pods for this datacenter. The default is false, which means that Stargate pods will be scheduled on separate worker nodes. Note: if the datacenter pods have HostNetwork:true, then the Stargate pods will inherit of it, in which case it is possible that Stargate nodes won't be allowed to sit on data nodes even if this property is set to true, because of port conflicts on the same IP address. | bool | false |
| cassandraConfigMapRef | CassandraConfigMapRef is a reference to a ConfigMap that holds Cassandra configuration. The map should have a key named cassandra_yaml. | *corev1.LocalObjectReference | false |
| telemetry | Telemetry defines the desired telemetry integrations to deploy targeting the Stargate pods for all DCs in this cluster (unless overriden by DC specific settings) | *telemetryapi.TelemetrySpec | false |

[Back to Custom Resources](#custom-resources)
