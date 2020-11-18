---
title: "K8ssandra Roadmap"
linkTitle: "Roadmap"
weight: 7
description: K8ssandra roadmap ideas for community consideration.
---

K8ssandra today is deployed as an entire stack. This open-source technology
currently assumes your deployment uses the entire stack. Trading out certain
components for others is not supported at this time. As part of the roadmap, one
goal is to support a la carte composition of components.

The following additional ideas are not yet in priority order. 

* Preconfigured alerts for metrics
  * Support sending to a configured single email address
  * Support sending to multiple addresses
  * Support sending to multiple addresses based on fired alert
* Annotations in metrics system signaling operations performed on the cluster
  * Restarts
  * Upgrades
  * Backups
  * Node lifecycle calls
* Extended authentication options for metrics system
* Network policies to isolate all components as appropriate
* Support for monitoring repair process via metrics system
* Support for monitoring backups process via metrics system
* Centralized logging
  * ELK
  * Loki
* Distributed tracing via Jaeger
* Load/stress/perf testing with nosqlbench - guidance on how to use this tool:
  * To find optimal load testing values (number of threads / inflight)
  * To determine node requirements
* Documentation enhancements
  * Istio best practices guide
  * Linkerd best practices guide
* Migrations
  * Connecting to existing on-prem datacenters / clusters
* Spark connection
* Kafka connection
* Data loading
* Service Broker
* Serverless / Faas
  * Kubeless
  * Knative
* Dynamics secrets with Vault
  * Roles via Cassandra plugin
  * Rotating TLS certificates for clients, nodes, ingress, etc.
