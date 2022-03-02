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
description: "K8ssandra documentation: architecture, configuration, guided tasks."
type: docs
---

The K8ssandra 1.4.x documentation covers everything from install details, deployed components, configuration references, and guided outcome-based tasks. Check out the sections at left based on your needs and be sure to leave us a <a class="github-button" href="https://github.com/k8ssandra/k8ssandra" data-icon="octicon-star" aria-label="Star k8ssandra/k8ssandra on GitHub">star</a> on Github!

**Tip:** Starting with the March 2022 introduction of K8ssandra Operator, the K8ssandra Operator custom resource, and its support for single- or **multi-cluster** deployments, the documentation site has been organized into three areas. 

* [docs.k8ssandra.io](https://docs.k8ssandra.io/) provides topics that are of common interest to users of K8ssandra Operator and K8ssandra, such as FAQs, Release Notes, Components, and a Glossary.
* [docs-v1.k8ssandra.io](https://docs-v1.k8ssandra.io/) (**you are here**) provides topics that are specific to K8ssandra 1.4.x users (the initial project releases).
* [docs-v2.k8ssandra.io](https://docs-v2.k8ssandra.io/) provides topics that are specific to the more recent (and recommended) K8ssandra Operator software, including single- or multi-cluster installs.

From each page’s top banner, use the **Versions** menu to navigate back to the Common, v1, or v2 documentation Home.

![K8ssandra documentation Versions](k8ssandra-doc-versions.png)

## What is K8ssandra?

K8ssandra is a cloud native distribution of Apache Cassandra® (Cassandra) meant to run on Kubernetes. Accompanying Cassandra is a suite of tools to ease and automate operational tasks. This includes metrics, data anti-entropy services, and backup tooling. As part of K8ssandra's installation process all of these components are installed and wired together freeing your teams from having to perform the tedious plumbing of components.

Cassandra may be deployed in a number of environments. This includes on bare metal hosts, virtual machines, and within container platforms. Each deployment type has its pros and cons, but in all cases it is **_essential_** that automation be leveraged to ensure that all node are configured homogeneously and without failure.

K8ssandra focuses on deploying Cassandra within Kubernetes. Kubernetes was chosen as it allows for the consumption of a common, versioned, set of APIs and tooling across multiple cloud platforms and environments.

If you're impatient, jump right in with our K8ssandra 1.4.x **[install]({{< relref "install" >}})** topics!

## Why do I want K8ssandra?

Apache Cassandra is _the_ NoSQL database for applications that require resilience and scalability. Unfortunately this comes with the same burdens as other distributed systems. There are multiple nodes replicating data all the time. Understanding the health of these systems requires advanced tooling and knowledge of the constituent parts. Users could spend time investigating and building out solutions to ensure operational stability of their Cassandra clusters. K8ssandra looks to provide those integrations from the start in a simple easy to deploy package.

## What is K8ssandra good for?

K8ssandra is a great fit for operators looking for easy to install and manage Cassandra clusters. Even if your environment currently does not run Cassandra on Kubernetes we believe that simple installation and upkeep will win you over. Consider some of the integrations listed below:

* All Cassandra containers are preinstalled with:
  * [Metrics Collector for Apache Cassandra](https://github.com/datastax/metric-collector-for-apache-cassandra)
  * [Management API for Apache Cassandra](https://github.com/datastax/management-api-for-apache-cassandra)
* Cass Operator to support management tasks in Kubernetes. 
* Prometheus Operator `ServiceMonitor` custom resources complete with metric labeling.
* Grafana Operator `Dashboard` custom resources configured with metrics exposed by Prometheus.
* Reaper to repair Cassandra data.
* Medusa for backup/restore operations.

## Next steps

* [FAQs](https://docs.k8ssandra.io/faqs/): If you're new to K8ssandra, these FAQs are for you. 
* [Install]({{< relref "install" >}}): K8ssandra install steps for local development or production-ready cloud platforms.
* [Quickstarts]({{< relref "quickstarts" >}}): Post-install K8ssandra topics for developers or Site Reliability Engineers.
* [Components]({{< relref "components" >}}): Dig in to each deployed component of the K8ssandra stack and see how it communicates with the others.
* [Tasks]({{< relref "tasks" >}}): Need to get something done? Check out the Tasks topics for a helpful collection of outcome-based solutions.
* [Reference]({{< relref "reference" >}}): Explore the K8ssandra configuration interface (Helm charts), the available options, and a Glossary.

We encourage you to actively participate in the [K8ssandra community](https://k8ssandra.io/community/).
