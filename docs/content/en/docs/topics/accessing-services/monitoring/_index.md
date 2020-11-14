---
title: "Monitoring UI"
linkTitle: "Monitoring UI"
weight: 1
date: 2020-11-13
description: |
  Follow these simple steps to access the Prometheus and Grafana monitoring interfaces.
---

## Tools

* Web Browser

## Prerequisites

1. Kubernetes cluster with the following elements deployed:
   * [K8ssandra Operators]({{< ref "getting-started#install-k8ssandra" >}})
   * [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}})
   * [Ingress Controller]({{< ref "ingress" >}})
1. DNS name for the Grafana service, referred to as _GRAFANA DOMAIN_ below.
1. DNS name for the Prometheus service, referred to as _PROMETHEUS DOMAIN_
   below.

## Access Grafana Interface

![Grafana UI](grafana-dashboard.png)

With the prerequisites satisfied the repair GUI should be available at the
following address:

**http://GRAFANA_DOMAIN/**

### What can I do in Grafana?

* Cluster health
* Traffic metrics

## Access Prometheus Interface

![Prometheus UI](prometheus-dashboard.png)

Prometheus is available at the following address:

**http://PROMETHEUS_DOMAIN/**

### What can I do in Prometheus?

* Validate serves being scraped
* Confirm metrics collection

## Next

Access the [Repair Web interface]({{ ref "repair" }}).
