---
title: "K8ssandra documentation"
linkTitle: "Documentation"
no_list: true
weight: 20
menu:
  main:
    weight: 20
description: |
  Welcome to K8ssandra!
---

The K8ssandra Documentation covers everything from architecture and
configuration references to guided outcome-based tasks. Check out the sections
at left based on your needs and be sure to leave us a <a class="github-button"
href="https://github.com/k8ssandra/k8ssandra" data-icon="octicon-star"
aria-label="Star k8ssandra/k8ssandra on GitHub">star</a> on Github!

{{% alert title="Tip" color="primary" %}}
If you're impatient, jump right in with our **[Quick start]({{< ref "getting-started" >}})**!
{{% /alert %}}

## What is K8ssandra?

K8ssandra is a cloud native distribution of Apache Cassandra meant to run on Kubernetes. Accompanying Cassandra is a suite of tools to ease and automate operational tasks. This includes metrics, data anti-entropy services, and backup tooling. As part of K8ssandra's installation process all of these components are installed and wired together freeing your teams from having to perform the tedious plumbing of components.

Cassandra may be deployed in a number of environments. This includes on bare metal hosts, virtual machines, and within container platforms. Each deployment type has its pros and cons, but in all cases it is **_essential_** that automation be leveraged to ensure that all node are configured homogenously and without failure.

K8ssandra focuses on deploying Cassandra within Kubernetes. Kubernetes was chosen as it allows for the consumption of a common, versioned, set of APIs and tooling across multiple cloud platforms and environments.

## Why do I want K8ssandra?

Apache Cassandra is _the_ NoSQL database for applications that require resilience and scalability. Unfortunately this comes with the same burdens as other distributed systems. There are multiple nodes replicating data all the time. Understanding the health of these systems requires advanced tooling and knowledge of the constituent parts. Users could spend time investigating and building out solutions to ensure operational stability of their Cassandra clusters. K8ssandra looks to provide those integrations from the start in a simple easy to deploy package. 

## What is K8ssandra good for?

K8ssandra is a great fit for operators looking for easy to install and manage Cassandra clusters. Even if your environment currently does not run Cassandra on Kubernetes we believe that simple installation and upkeep will win you over. Consider some of the integrations listed below:

* All Cassandra containers are preinstalled with
  * [Metrics Collector for Apache Cassandra](https://github.com/datastax/metric-collector-for-apache-cassandra)
  * [Management API for Apache Cassandra](https://github.com/datastax/management-api-for-apache-cassandra)
* Prometheus Operator `ServiceMonitor` custom resources complete with metric relabelling.
* Grafana Operator `Dashboard` custom resources configured with metrics exposed by Prometheus
* Reaper Operator custom resources connected to the cluster.

## What is K8ssandra *not yet* good for?

Right now K8ssandra is deployed as an entire stack. It currently assumes your deployment uses the entire stack. Trading out certain components for others is not supported. As part of the [Roadmap]({{< ref "roadmap" >}}) we would like to see this change to support a la carte composition of components.

## Where should I go next?

Depending on your needs, see the following:

* [Quick start]({{< ref "getting-started" >}}): Get started with K8ssandra!
* [Tasks]({{< ref "topics" >}}): Need to get something done? Check out the Tasks section for a helpful collection of outcome-based solutions.
* [Reference]({{< ref "reference" >}}): Explore the K8ssandra configuration interface and options available.
* [Architecture]({{< ref "architecture" >}}): Dig in to each operational component of the K8ssandra stack and see how it communicates with the others.
