---
title: "Install K8ssandra on EKS"
linkTitle: "Amazon EKS"
toc_hide: true
weight: 2
description: >
  Pointer to K8ssandra install on Amazon EKS topic.
---

Amazon [Elastic Kubernetes Service](https://aws.amazon.com/eks/features/) or "EKS" is a managed Kubernetes service that makes it easy for you to run Kubernetes on AWS and on-premises. EKS is certified Kubernetes conformant, so existing applications that run on upstream Kubernetes are compatible with EKS. AWS automatically manages the availability and scalability of the Kubernetes control plane nodes responsible scheduling containers, managing the availability of applications, storing cluster data, and other key tasks.

## Available topics

* For details about the initial K8ssandra (1.4.x) implementation, see the [Amazon EKS](https://docs-v1.k8ssandra.io/install/eks/) install topic.

* For details about the new K8ssandra Operator, see this [local install](https://docs-v2.k8ssandra.io/install/local/) topic. It showcases the K8ssandra Operator support for single- and multi-cluster environments. Examples use local **kind** Kubernetes along with helper scripts invoked by `make` commands, as well as `helm` and `kustomize` options. Topics covering installs that use K8ssandra Operator for Google GKE, Amazon EKS, and Azure AKS may be added at a later time. 
