---
title: "Install K8ssandra on GKE"
linkTitle: "Google GKE"
toc_hide: true
weight: 2
description: >
  Pointer to K8ssandra install on Google GKE topic.
---

[Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) or "GKE" is a managed Kubernetes environment on the [Google Cloud Platform](https://cloud.google.com/) (GCP). GKE is a fully managed experience; it handles the management/upgrading of the Kubernetes cluster master as well as autoscaling of "nodes" through "node pool" templates.

Through GKE, your Kubernetes deployments will have first-class support for GCP IAM identities, built-in configuration of high-availability and secured clusters, as well as native access to GCP's networking features such as load balancers.

## Available topics

* For details about the initial K8ssandra (1.4.x) implementation, see this [Google GKE](https://docs-staging-v1.k8ssandra.io/install/gke/) install topic.

* For details about the new K8ssandra Operator, see this [local install](https://docs-staging-v2.k8ssandra.io/install/local/) topic. It showcases the K8ssandra Operator support for single- and multi-cluster environments. Examples use local **kind** Kubernetes along with helper scripts invoked by `make` commands, as well as `helm` and `kustomize` options. Topics covering installs that use K8ssandra Operator for Google GKE, Amazon EKS, and Azure AKS may be added at a later time. 
