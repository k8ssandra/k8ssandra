---
title: "K8ssandra tasks"
linkTitle: "Tasks"
simple_list: true
weight: 5
description: Topics to help you get things done, **now**.
---

**Tip:** These topics are specific to K8ssandra 1.4.x. Consider exploring our most recent (and recommended) implementation: **K8ssandra Operator**. It includes a `K8ssandraCluster` custom resource and supports single- and **multi-cluster** Cassandra deployments in Kubernetes, for High Availability (HA) capabilities. See the [K8ssandra Operator documentatation](https://docs-staging-v2.k8ssandra.io). 

The 1.4.x instructions here build on the [Quickstart]({{< relref "/quickstarts/" >}}) steps for developers or SREs, and on the cloud-provider specific [Install]({{< relref "/install/" >}}) topics. If you haven't already installed K8ssandra 1.4.x and its deployed components using default or custom settings for your preferred cloud provider, see those topics.

Accessing resources from outside of the K8ssandra cluster requires tooling to make the internal resources available at external connection points. This may be accomplished through a number of means, many explored here. Some of the tasks assume a Kubernetes Ingress Controller has been installed and configured. If you need help with that setup, head over to the [Traefik ingress]({{< relref "/tasks/connect/ingress" >}}) section and explore these topics.

