---
title: "Install K8ssandra on AKS"
linkTitle: "Azure AKS"
toc_hide: true
weight: 2
description: >
  Pointer to K8ssandra install on Azure AKS topic
---

[Azure Elastic Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/) or "AKS" is a managed Kubernetes service that makes it easy for you to run Kubernetes on Azure. AKS offers serverless Kubernetes, an integrated continuous integration and continuous delivery (CI/CD) experience, and enterprise-grade security and governance.

## Available topics

* For details about the initial K8ssandra (1.4.x) implementation, see this [Azure AKS](https://docs-staging-v1.k8ssandra.io/install/aks) install topic.

* For details about the new K8ssandra Operator, see this [local install](https://docs-staging-v2.k8ssandra.io/install/local/) topic. It showcases the K8ssandra Operator support for single- and multi-cluster environments. Examples use local **kind** Kubernetes along with helper scripts invoked by `make` commands, as well as `helm` and `kustomize` options. Topics covering installs that use K8ssandra Operator for Google GKE, Amazon EKS, and Azure AKS may be added at a later time. 
