---
title: "Access the Repair Web Interface (Reaper)"
linkTitle: "Access the Repair Web Interface (Reaper)"
weight: 1
date: 2020-11-11
description: Use Reaper to repair Cassandra in Kubernetes
---

## Tools

K8ssandra comes preconfigured with the Repair Web Interface, which is also known as Cassandra&reg; Reaper. Use this helpful open-source tool to schedule and orchestrate repairs of Apache Cassandra clusters in your Kubernetes environment.

For full details on how to use Reaper, see its [documentation](http://cassandra-reaper.io/docs/).

## Prerequisites

All prerequisites are met by K8ssandra. It's installed, configured, and ready to go.

## Steps

Reaper includes a community-driven web interface that can be accessed at:

http://$REAPER_HOST:8080/webui/index.html

From there, use Reaper to perform these tasks:

<!--- Point to existing topics vs repeat info here? --> 

* [Check a cluster's health](http://cassandra-reaper.io/docs/usage/health/)
* [Run a cluster repair](http://cassandra-reaper.io/docs/usage/single/)
* [Schedule a cluster repair](http://cassandra-reaper.io/docs/usage/schedule/)
* [Monitor Cassandra diagnostic events](http://cassandra-reaper.io/docs/usage/cassandra-diagnostics/)

## Next

Learn about [external Cassandra connectivity](docs/topics/external-connectivity/) techniques in Kubernetes.
