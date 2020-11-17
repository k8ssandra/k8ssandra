---
title: "Cassandra"
linkTitle: "Cassandra"
weight: 1
date: 2020-11-14
description: |
  Accessing the K8ssandra Apache Cassandra interfaces.
---

## Tools

* Cassandra enabled application
* CQLSH command-line tool

## Prerequisites

1. Kubernetes cluster with the following elements deployed:
   * [K8ssandra Operators]({{< ref "getting-started#install-k8ssandra" >}})
   * [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}})
   * [Ingress Controller]({{< ref "ingress" >}})
1. DNS name _**and port**_ for the non-TLS Cassandra service
1. _Optional_ DNS name for the TLS Cassandra service
1. _Optional_ CA certificate
1. _Optional_ TLS client certificate and key

## Access Cassandra without TLS

### TODO

## Access Cassandra with TLS

### TODO

## Access Cassandra with mTLS

### TODO
