---
title: "Repair UI"
linkTitle: "Repair UI"
weight: 1
date: 2020-11-13
description: |
  Follow these simple steps to access the Reaper repair interface.
---

## Tools

* Web Browser

## Prerequisites

1. Kubernetes cluster with the following elements deployed:
   * [K8ssandra Operators]({{< ref "getting-started#install-k8ssandra" >}}) Helm
     Chart
   * [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}}) Helm
     Chart
   * [Ingress Controller]({{< ref "ingress" >}})
1. DNS name configured for the repair interface, referred to as _REPAIR DOMAIN_
   below.

## Access Repair Interface

![Reaper UI](reaper-ui.png)

With the prerequisites satisfied the repair GUI should be available at the
following address:

http://REPAIR_DOMAIN/webui

## What can I do in Reaper?

For details about the tasks you can perform in Reaper, see these topics in the
Cassandra Reaper documentation:

* [Check a cluster's health](http://cassandra-reaper.io/docs/usage/health/)
* [Run a cluster repair](http://cassandra-reaper.io/docs/usage/single/)
* [Schedule a cluster repair](http://cassandra-reaper.io/docs/usage/schedule/)
* [Monitor Cassandra diagnostic
  events](http://cassandra-reaper.io/docs/usage/cassandra-diagnostics/)
