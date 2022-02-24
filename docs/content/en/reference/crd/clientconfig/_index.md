---
title: "ClientConfig CRD"
linkTitle: "ClientConfig CRD"
simple_list: false
weight: 6
description: >
  ClientConfig Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [ClientConfig](#clientconfig)
* [ClientConfigList](#clientconfiglist)
* [ClientConfigSpec](#clientconfigspec)

#### ClientConfig

ClientConfig is the Schema for the kubeconfigs API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec |  | [ClientConfigSpec](#clientconfigspec) | false |

[Back to Custom Resources](#custom-resources)

#### ClientConfigList

ClientConfigList contains a list of KubeConfig

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][ClientConfig](#clientconfig) | true |

[Back to Custom Resources](#custom-resources)

#### ClientConfigSpec

ClientConfigSpec defines the desired state of KubeConfig

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| contextName | ContextName allows to override the object name for context-name. If not set, the ClientConfig.Name is used as context name | string | false |
| kubeConfigSecret | KubeConfigSecret should reference an existing secret; the actual configuration will be read from this secret's \"kubeconfig\" key. | corev1.LocalObjectReference | false |

[Back to Custom Resources](#custom-resources)
