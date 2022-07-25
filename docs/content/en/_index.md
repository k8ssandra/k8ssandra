---
title: "K8ssandra Documentation"
linkTitle: "Docs"
no_list: true
weight: 20
menu:
  main:
    weight: 20
  footer:
    weight: 60
description: "K8ssandra documentation: architecture, configuration, guided tasks"
type: docs
---

This documentation covers everything from install details, deployed components, configuration references, and guided outcome-based tasks. 

Be sure to leave us a <a class="github-button" href="https://github.com/k8ssandra/k8ssandra" data-icon="octicon-star" aria-label="Star k8ssandra/k8ssandra on GitHub">star</a> on Github!

## Features for single- and multi-cluster Kubernetes environments

| K8ssandra Operator: enhanced capabilities | Initial K8ssandra project|
| ----------- | ----------- |
| K8ssandra Operator is our most recent offering. In a **unified operator**, K8ssandra Operator provides an entirely new, solidified set of features for Kubernetes + Cassandra deployments. The features include robust management (cass-operator), API integration (Stargate), anti-entropy repairs (Reaper), and backup/restore (Medusa). Important enhancements include **multi-cluster** and **multi-region** support, which enables greater scalability and availability for enterprise apps and data. Single cluster/region deployments are also supported with K8ssandra Operator.| K8ssandra v1.4.x is our project's initial implementation. It continues to provide a set of separate Helm charts that you can use to configure and deploy Apache Cassandra&reg; into a single-cluster, single-region Kubernetes environment. |
| For enhanced capabilities, we recommend that you explore K8ssandra Operator [local install]({{< relref "install/local" >}}) topic, which focuses on single- or multi-cluster deployments on local dev  **kind** Kubernetes clusters, using the provided `make` commands, `helm`, or `kustomize`. | Start in the K8ssandra v1.4.x [install](https://docs-v1.k8ssandra.io/install/local/) topics, which include the steps for single-cluster installs on local or cloud-provider Kubernetes platforms. |

If you're using K8ssandra v1.4.x, you may continue to do so. Or consider stepping up to the project's latest implementation with K8ssandra Operator.

## Compatibility matrix

| Kubernetes                  | **v1.17** | **v1.18** | **v1.19** | **v1.20** | **v1.21** | **v1.22** | **v1.23** | **v1.24** |
|-----------------------------|:---------:|:---------:|:---------:|:---------:|:---------:|:---------:|:---------:|:---------:|
| **K8ssandra v1.5**          |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |
| **K8ssandra-operator v1.0** |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |     ✅     |
| **K8ssandra-operator v1.1** |           |           |           |           |     ✅     |     ✅     |     ✅     |     ✅     |
| **K8ssandra-operator v1.2** |           |           |           |           |     ✅     |     ✅     |     ✅     |     ✅     |

## Next steps

We encourage you to actively participate in the [K8ssandra community](https://k8ssandra.io/community/).
