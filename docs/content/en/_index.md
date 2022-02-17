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

The K8ssandra documentation covers everything from install details, deployed components, configuration references, and guided outcome-based tasks. 

Be sure to leave us a <a class="github-button" href="https://github.com/k8ssandra/k8ssandra" data-icon="octicon-star" aria-label="Star k8ssandra/k8ssandra on GitHub">star</a> on Github!

## Features for single- and multi-cluster Kubernetes environments

| K8ssandra v1.4.x      | K8ssandra Operator v1.0.0 |
| ----------- | ----------- |
| K8ssandra v1.4.x is our initial implementation. It provides a set of separate Helm charts you can use to configure and deploy Apache Cassandra&reg; into a single-cluster, single-region Kubernetes environment. | K8ssandra Operator v1.0.0 is our most recent offering. It combines API, management, &amp; observability features under the control of a unified operator. Important enhancements include **multi-cluster** and **multi-region** support for Cassandra deployments in Kubernetes, which enables greater scalability and availability. Single cluster/region deployments are also supported with K8ssandra Operator.|
| Start in the K8ssandra v1.4.x [install](https://docs-staging-v1.k8ssandra.io/install/local/) topics, which include the steps for single-cluster installs on local or cloud-provider Kubernetes platforms. | Start in the K8ssandra Operator v1.0.0 [local install](https://docs-staging-v2.k8ssandra.io/install/local/) topic, which focuses on multi-cluster or single-cluster installs on local dev  **kind** Kubernetes, using helper scripts, `helm`, or `kustomize`.

If you're using K8ssandra v1.4.x, you may continue to do so. Or consider stepping up to the project's latest implementation with K8ssandra Operator v1.0.0.

## Next steps

* [FAQs]({{< relref "faqs" >}}): If you're new to K8ssandra, these FAQs are for you. 
* [Install]({{< relref "install" >}}): K8ssandra install steps for local development or production-ready cloud platforms.
* [Quickstarts]({{< relref "quickstarts" >}}): Post-install K8ssandra topics for developers or Site Reliability Engineers.
* [Components]({{< relref "components" >}}): Dig in to each deployed component of the K8ssandra stack and see how it communicates with the others.
* [Tasks]({{< relref "tasks" >}}): Need to get something done? Check out the Tasks topics for a helpful collection of outcome-based solutions.
* [Reference]({{< relref "reference" >}}): Explore the K8ssandra configuration interface (Helm charts), the available options, and a Glossary.

We encourage you to actively participate in the [K8ssandra community](https://k8ssandra.io/community/).
