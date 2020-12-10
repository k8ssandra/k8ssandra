---
title: "Stargate"
linkTitle: "Stargate"
weight: 1
description: |
  Accessing the K8ssandra Stargate interfaces.
---

[Stargate](https://stargate.io/) is an open-source framework providing common
API interfaces for backend databases. With K8ssandra, Stargate may be deployed
in front of the Apache Cassandra cluster providing CQL, REST, and GraphQL API
endpoints. These endpoints may be scaled horizontally independently of data
layer scaling needs. This guide covers accessing the various API endpoints
provided by Stargate.

## Tools

* HTTP client (cURL, Postman, etc.)
* GraphQL client
* CQL client (application, cqlsh, etc.)

## Prerequisites

1. Kubernetes cluster with the following elements deployed:
   * [K8ssandra Operators]({{< ref "getting-started#install-k8ssandra" >}}) Helm
     Chart
   * [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}}) Helm
     Chart
   * [Ingress Controller]({{< ref "ingress" >}})
1. DNS name configured for the Stargate interface, referred to as _STARGATE
   DOMAIN_ below.
1. Port number for the Stargate CQL interface

## Access REST interface

### TODO 

## Access GraphQL interface

### TODO

## Access CQL interface

### TODO
