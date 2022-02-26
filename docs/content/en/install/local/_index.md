---
title: "Install K8ssandra on your local K8s"
linkTitle: "Local"
toc_hide: true
no_list: true
weight: 1
description: "Pointers to K8ssandra Operator and K8ssandra install topics for local dev Kubernetes."
---

## Available topics

You can deploy Apache Cassandra&reg; in a local Kubernetes cluster using environments such as kind, minikube, and K3D. 

* For details about the new K8ssandra Operator, see this [local install](https://docs-staging-v2.k8ssandra.io/install/local/) topic. It showcases the K8ssandra Operator support for single- and multi-cluster environments. Examples use local **kind** Kubernetes along with helper scripts invoked by `make` commands, as well as `helm` and `kustomize` options.  

* For details about using the initial K8ssandra (1.4.x) implementation, see this [local installs](https://docs-staging-v1.k8ssandra.io/install/local) topic. It includes single-cluster only support, with local **minikube** Kubernetes and `helm` examples. 
