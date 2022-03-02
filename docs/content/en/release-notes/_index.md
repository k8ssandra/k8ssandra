---
title: "K8ssandra project release notes"
linkTitle: "Release notes"
weight: 2
description: "Release notes for the open-source K8ssandra project."
---

The open-source K8ssandra project provides a production-ready platform for running Apache Cassandra&reg; on Kubernetes. This functionality includes automation for operational tasks such as repairs, backup and restores, and monitoring. Also deployed is Stargate, an open source data gateway that lets you interact programmatically with your Kubernetes-hosted Cassandra resources via a well-defined API. 

## Latest releases

* K8ssandra Operator 1.0.0, released 17-Feb-2022
* K8ssandra 1.4.1, released 02-December-2021

## New &amp; noteworthy

**K8ssandra Operator** is our latest implementation. It provides a cloud-native distribution of Cassandra that runs on Kubernetes. Significantly, K8ssandra Operator provides a new `K8ssandraCluster` custom resource that enables support for single- or **multi-cluster, multi-region** deployments of Cassandra and related services. It's all part of the overall K8ssandra project, but you'll need to deploy with K8ssandra Operator to use the latest multi-cluster/region features. For details, start in the K8ssandra Operator [install topics](https://docs-staging-v2.k8ssandra.io/install/). You'll find there topics for single- and multi-cluster instructions that use Helm or Kustomize tools.

For more, see the [K8ssandra Operator architecture](https://docs-staging-v2.k8ssandra.io/components/k8ssandra-operator/architecture/).

## GitHub repos

* K8ssandra Operator: https://github.com/k8ssandra/k8ssandra-operator
* K8ssandra: https://github.com/k8ssandra
* Cass Operator: https://github.com/k8ssandra/cass-operator

## Prerequisites

* A Kubernetes environment from v1.17 (minimum supported) up to v1.22 (current tested upper bound) - local or via a supported cloud provider
* [Helm](https://helm.sh/) v3.5.x or later. Recommendation: avoid Helm 3.7.0 due to a known CVE and subsequent regression. See issue [1103](https://github.com/k8ssandra/k8ssandra/issues/1103) in GitHub.
* Additional prereqs are listed in the K8ssandra Operator [install](https://docs-staging-v2.k8ssandra.io/install/) topic.

## Supported Kubernetes environments

* Open-source [kubernetes.io](https://kubernetes.io)
* [Amazon Elastic Kubernetes Service](https://aws.amazon.com/eks/) (EKS)
* [DigitalOcean Kubernetes](https://www.digitalocean.com/products/kubernetes/) (DOKS)
* [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) (GKE)
* [Microsoft Azure Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/) (AKS)
* [MiniKube](https://minikube.sigs.k8s.io/docs/)
* [Kind](https://kind.sigs.k8s.io/)
* [K3D](https://k3d.io/)

## Deployed components

### K8ssandra Operator 1.0.x deployments

K8ssandra Operator deploys and manages the following. 

* [Apache Cassandra](https://cassandra.apache.org/) - the deployed version depends on the configured setting:
  * 4.0.1 (default)
  * 3.11.11
  * 3.11.10
  * 3.11.9
  * 3.11.8
  * 3.11.7
* [cass-operator](https://github.com/k8ssandra/cass-operator) v1.10.0.
* Reaper vN.n.n
* Medusa vN.n.n
* [Stargate](https://github.com/stargate/stargate) vN.n.n
* Others TBS here ... 

### K8ssandra 1.4.x deployments

The K8ssandra helm chart deploys the following components. Some are optional, and depending on the configuration, may not be deployed:

* [Apache Cassandra](https://cassandra.apache.org/) - the deployed version depends on the configured `cassandra.version` setting:
  * 4.0.1 (default)
  * 3.11.11
  * 3.11.10
  * 3.11.9
  * 3.11.8
  * 3.11.7
* [cass-operator](https://github.com/k8ssandra/cass-operator) v1.8.0
* Management API for Apache Cassandra ([MAAC](https://github.com/datastax/management-api-for-apache-cassandra) v0.1.33
* [Stargate](https://github.com/stargate/stargate) Stargate v1.0.40
* Metric Collector for Apache Cassandra ([MCAC](https://github.com/datastax/metric-collector-for-apache-cassandra) v0.2.0
* kube-prometheus-stack v12.11.3 [chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
* Medusa for Apache Cassandra v0.11.13
* medusa-operator v0.4.0
* Reaper for Apache Cassandra v3.0.0
* reaper-operator v2.3.0

## K8ssandra Operator 1.0.x revisions

For the latest K8ssandra Operator changes, features, enhancements, and bug fixes, refer to the [CHANGELOG](https://github.com/k8ssandra/k8ssandra-operator/blob/main/CHANGELOG/CHANGELOG-1.0.md).


## K8ssandra 1.4.x revisions

For the latest K8ssandra changes, features, enhancements, and bug fixes, refer to the [CHANGELOG](https://github.com/k8ssandra/k8ssandra/blob/main/CHANGELOG-1.4.md).

## Contributions

Your feedback and contributions are welcome! To contribute, fork the [K8ssandra Operator repo]() and the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra), and submit Pull Requests (PRs) for review.

To submit documentation comments or edits, see [Contribution guidelines]({{< relref "contribute" >}}).

## Next steps

Read the [FAQs]({{< relref "faqs" >}}) - for starters, how to pronounce "K8ssandra." 
