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
   * [Traefik]({{< ref "traefik" >}})
   * [K8ssandra Operators]({{< ref "getting-started#install-k8ssandra" >}})
   * [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}})

   See the [Configuring Kind]({{< ref "kind-deployment" >}}) for an example of
   how to set up a local installation.
1. DNS name for the Grafana service
1. DNS name for the Prometheus service

## Access Grafana Interface

![Grafana UI](grafana-dashboard.png)

Now that Traefik is configured you may now access the web interface by visiting
the domain name provided within the `values.yaml` file. Traefik receives the
HTTP request then performs the following actions:

* Extract the HTTP `Host` header 
* Match the `Host` against the rules specified in our `IngressRoutes`
* Proxies the request to the upstream Kubernetes Service.

## Access Prometheus Interface

![Prometheus UI](grafana-dashboard.png)

Now that Traefik is configured you may now access the web interface by visiting
the domain name provided within the `values.yaml` file. Traefik receives the
HTTP request then performs the following actions:

* Extract the HTTP `Host` header 
* Match the `Host` against the rules specified in our `IngressRoutes`
* Proxies the request to the upstream Kubernetes Service.

## Next

Access the [Repair Web interface](docs/topics/access-repair-interface/).
