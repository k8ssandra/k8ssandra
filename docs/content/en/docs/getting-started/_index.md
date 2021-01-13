---
title: "Getting Started"
linkTitle: "Getting Started"
weight: 2
description: |
  Kick the tires and take it for a spin!
---

Welcome to K8ssandra! This guide gets you up and running with a single-node Apache Cassandra&reg; cluster on Kubernetes. If you are interested in a more detailed component walkthroughs check out the [topics]({{< ref "topics">}}) section.

## Prerequisites

In your local environment the following tools are required for provisioning a K8ssandra cluster.

* [Helm v3+](https://helm.sh/docs/intro/install/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

As K8ssandra deploys on a Kubernetes cluster one must be available to target for installation. This may be a local version running on your development machine, on-premises self-hosted environment, or managed cloud offering. To that end the cluster must be up and available to your `kubectl` command.

```console
# Validate cluster connectivity
kubectl cluster-info
```

If you do not have a Kubernetes cluster available consider one of the following local versions that run within Docker or a virtual machine.

* [K3D](https://k3d.io/)
* [Kind](https://kind.sigs.k8s.io/)
* [OpenShift CodeReady Containers](https://developers.redhat.com/products/codeready-containers/overview)

## Configure Helm Repository

K8ssandra is delivered as a collection of Helm Charts. In order to leverage these charts we have provided a k8ssandra Helm Repository for easy installation. 

Also add the Traefik Ingress repo - you'll need its resources to access services from outside the Kubernetes cluster.

```console
helm repo add k8ssandra https://helm.k8ssandra.io/
helm repo add traefik https://helm.traefik.io/traefik
helm repo update
```

Alternatively, you may download the individual charts directly from the project's [releases](https://github.com/k8ssandra/k8ssandra/releases) page.

## Enabling Stargate

Stargate deployment is *disabled* by default. To enable it, modify the values provided to the helm install/upgrade of the
`k8ssandra` chart.

This can be done by adding to the values.yaml file provided to helm:

```
stargate:
  enabled: true
```

or set via the command line when installing via helm:

```
--set stargate.enabled=true
```

## Install K8ssandra

From a packaging perspective, K8ssandra is composed of a number of helm charts. There's separate charts handling the installation of operators, custom resources as well as
provisioning the cluster instances.

```console
helm install k8ssandra-tools k8ssandra/k8ssandra
```

In later steps, you can upgrade your k8ssandra via `helm upgrade` commands, for example to access services from outside Kubernetes via a Traefik Ingress controller.
