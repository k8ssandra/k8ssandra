---
title: "Install K8ssandra on DOKS"
linkTitle: "DigitalOcean DOKS"
weight: 2
description: >
  Pointer to K8ssandra install on DigitalOcean DOKS topic.
---

[DigitalOcean Kubernetes](https://www.digitalocean.com/products/kubernetes/) or "DOKS" is a managed Kubernetes environment on DigitalOcean. DOKS is a fully managed experience; it handles the management/upgrading of the Kubernetes cluster master as well as autoscaling of "nodes" through "node pools."

## Available topics

* For details about the initial K8ssandra (1.4.x) implementation, see this [DigitalOcean DOKS](https://docs-staging-v1.k8ssandra.io/install/doks/) install topic.

* For details about the new K8ssandra Operator, see this [local install](https://docs-staging-v2.k8ssandra.io/install/local/) topic. It showcases the K8ssandra Operator support for single- and multi-cluster environments. Examples use local **kind** Kubernetes along with helper scripts invoked by `make` commands, as well as `helm` and `kustomize` options. Topics covering installs that use K8ssandra Operator for Google GKE, Amazon EKS, and Azure AKS may be added at a later time. 
