---
title: "Repair UI"
linkTitle: "Repair UI"
weight: 1
date: 2020-11-13
description: |
  Accessing the Cassandra Reaper, repair interface
---

# TODO Rewrite Intro
# TODO Remove Traefik references

Follow these steps to configure and install `Traefik Ingress` custom resources
for accessing your K8ssandra cluster's repair interface (provided by Cassandra
Reaper).

## Tools

* Web Browser

## Prerequisites

1. Kubernetes cluster with the following elements deployed:
   * [K8ssandra Operators]({{< ref "getting-started#install-k8ssandra" >}})
   * [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}})
   * Kubernetes Ingress Controller
1. DNS name where the repair service should be listening.

   _Note_ if you do not have a DNS name available, consider using a service like
   [xip.io](http://xip.io) to generate a domain name based on the ingress IP
   address. For local Kind clusters this may look like `repair.127.0.0.1.xip.io`
   which would return the address `127.0.0.1` during DNS lookup.

## Access Repair Interface

![Reaper UI](reaper-ui.png)

Now that Traefik is configured you may now access the web interface by visiting
the domain name provided within the `values.yaml` file. Traefik receives the
HTTP request then performs the following actions:

* Extract the HTTP `Host` header 
* Match the `Host` against the rules specified in our `IngressRoutes`
* Proxies the request to the upstream Kubernetes Service.

## What can I do in Reaper?

For details about the tasks you can perform in Reaper, see these topics in the
Cassandra Reaper documentation:

* [Check a cluster's health](http://cassandra-reaper.io/docs/usage/health/)
* [Run a cluster repair](http://cassandra-reaper.io/docs/usage/single/)
* [Schedule a cluster repair](http://cassandra-reaper.io/docs/usage/schedule/)
* [Monitor Cassandra diagnostic events](http://cassandra-reaper.io/docs/usage/cassandra-diagnostics/)
