---
title: "ReplicatedSecret CRD"
linkTitle: "ReplicatedSecret CRD"
no_list: true
toc_hide: true
simple_list: false
weight: 6
description: >
  ReplicatedSecret Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [ReplicatedSecret](#replicatedsecret)
* [ReplicatedSecretList](#replicatedsecretlist)
* [ReplicatedSecretSpec](#replicatedsecretspec)
* [ReplicatedSecretStatus](#replicatedsecretstatus)
* [ReplicationCondition](#replicationcondition)
* [ReplicationTarget](#replicationtarget)

#### ReplicatedSecret

ReplicatedSecret is the Schema for the replicatedsecrets API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec |  | [ReplicatedSecretSpec](#replicatedsecretspec) | false |
| status |  | [ReplicatedSecretStatus](#replicatedsecretstatus) | false |

[Back to Custom Resources](#custom-resources)

#### ReplicatedSecretList

ReplicatedSecretList contains a list of ReplicatedSecret

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][ReplicatedSecret](#replicatedsecret) | true |

[Back to Custom Resources](#custom-resources)

#### ReplicatedSecretSpec

ReplicatedSecretSpec defines the desired state of ReplicatedSecret

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| selector | Selector defines which secrets are replicated. If left empty, all the secrets are replicated | *metav1.LabelSelector | false |
| replicationTargets | TargetContexts indicates the target clusters to which the secrets are replicated to. If empty, no clusters are targeted | [][ReplicationTarget](#replicationtarget) | false |

[Back to Custom Resources](#custom-resources)

#### ReplicatedSecretStatus

ReplicatedSecretStatus defines the observed state of ReplicatedSecret

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| conditions |  | [][ReplicationCondition](#replicationcondition) | false |

[Back to Custom Resources](#custom-resources)

#### ReplicationCondition



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| cluster | Cluster | string | true |
| type | Type of condition | ReplicationConditionType | true |
| status | Status of the replication to target cluster | corev1.ConditionStatus | true |
| lastTransitionTime | LastTransitionTime is the last time the condition transited from one status to another. | *metav1.Time | false |

[Back to Custom Resources](#custom-resources)

#### ReplicationTarget



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| namespace | Namespace to replicate the data to in the target cluster. If left empty, current namespace is used. | string | false |
| k8sContextName | K8sContextName defines the target cluster name as set in the ClientConfig. If left empty, current cluster is assumed | string | false |

[Back to Custom Resources](#custom-resources)
