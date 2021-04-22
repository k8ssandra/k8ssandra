---
title: "Google Kubernetes Engine"
linkTitle: "Google Kubernetes Engine"
weight: 1
description: >
  Complete production ready environment of K8ssandra on Google Kubernetes Engine (GKE).
---

[Google Cloud Platform](https://cloud.google.com/) (GCP) provides a fully managed Kubernetes service known as [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) or GKE. With this service users pay for the infrastructure used to power the cluster along with a small fee per node for management.

## Prerequisites

1. Google Cloud Platform (GCP) account
1. Quotas permitting the [Resources](#resources)
1. Helm
1. Terraform _optional_

## Provision Resources

This guide will cover provisioning and installing the following resources.

* 1x _Regional_ GKE cluster with instances spread across multiple Availability Zones.
* 2x Node Pools
  * Cassandra Node Pool
    * 3x n2-highmem-8 Instances 
      * 8 vCPUs
      * 64 GB RAM
  * Ops Node Pool
    * 3x n2-standard-4 Instances
      * 4 vCPUs
      * 16 GB RAM
* x Load Balancers
  * x Backend services
* x 2TB PD-SSD Volumes
...

### Terraform

### Manual

## Retrieve `kubeconfig`

```console
gcloud ....
```

## Install K8ssandra

```yaml
values.yaml file
```

## Google Cloud Platform Customizations

### Backups

Detail usage of Google Cloud Storage (GCS)

### Ingress

Detail GKE Ingress
