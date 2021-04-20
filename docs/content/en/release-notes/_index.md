
---
title: "K8ssandra release notes"
linkTitle: "Release notes"
weight: 1
description: Release notes for the open-source K8ssandra community project.
---

K8ssandra provides a production-ready platform for running Apache Cassandra&reg; on Kubernetes. This includes automation for operational tasks such as repairs, backup and restores, and monitoring. 

Also deployed is Stargate, an open source data gateway that lets you interact programmatically with your Kubernetes-hosted Cassandra resources via a well-defined API. 

{{% alert title="Note" color="success" %}}
The **K8ssandra 1.1.0** release implements a number of changes, enhancements, and bug fixes. This Release Notes topic lists a subset of the key updates. For the complete list, see the [CHANGELOG](https://github.com/k8ssandra/k8ssandra/blob/main/CHANGELOG-1.1.md).  
{{% /alert %}}

**Release date:** 09-April-2021.

## Prerequisites

* A Kubernetes v1.16 or later environment - local or via a supported cloud provider
* [Helm](https://helm.sh/) v3

## Supported Kubernetes environments

* Open-source [kubernetes.io](https://kubernetes.io)
* [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) (GKE)
* [Microsoft Azure Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/) (AKS)
* [Amazon Elastic Kubernetes Service](https://aws.amazon.com/eks/) (EKS)
* [MiniKube](https://minikube.sigs.k8s.io/docs/)
* [Kind](https://kind.sigs.k8s.io/)
* [K3D](https://k3d.io/)

## K8ssandra deployed components

The K8ssandra helm chart deploys the following components. Some are optional, and depending on the configuration, may not be deployed:

* [Apache Cassandra](https://cassandra.apache.org/) - the deployed version depends on the configured `cassandra.version` setting:
  * 3.11.7
  * 3.11.8
  * 3.11.9
  * 3.11.10 (default)
* DataStax Kubernetes Operator for Apache Cassandra ([cass-operator](https://github.com/datastax/cass-operator)) 1.6.0
* Management API for Apache Cassandra ([MAAC](https://github.com/datastax/management-api-for-apache-cassandra)) 0.1.24
* [Stargate](https://github.com/stargate/stargate) 1.0.18
* Metric Collector for Apache Cassandra ([MCAC](https://github.com/datastax/metric-collector-for-apache-cassandra)) 0.2.0
* kube-prometheus-stack 12.11.3 [chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
* Medusa for Apache Cassandra 0.10.0
* Reaper for Apache Cassandra 2.2.2

## Upgrade notice

{{% alert title="Important!" color="warning" %}}
Upgrading from K8ssandra 1.0.0 to 1.1.0 causes a StatefulSet update (due to [#533](https://github.com/k8ssandra/k8ssandra/issues/533) and [#613](https://github.com/k8ssandra/k8ssandra/issues/613)). A StatefulSet update has the effect of a rolling restart. Because of [#411](https://github.com/k8ssandra/k8ssandra/issues/411) this could require you to perform a manual restart of all Stargate nodes after the Cassandra cluster is back online. This behavior also impacts in-place restore operations of Medusa backups [#611](https://github.com/k8ssandra/k8ssandra/issues/611). 

To manually restart Stargate nodes:

1. Get the Deployment object in your Kubernetes environment:
   ```bash
   kubectl get deployment | grep stargate
   ```
2. Scale down with this command:
   ```bash
   kubectl scale deployment <stargate-deployment> --replicas 0
   ```
3. Run this next command and wait until it reports 0/0 ready replicas. This should happen within a couple seconds.
   ```bash
   kubectl get deployment <stargate-deployment>
   ```
4. Now scale up with:
   ```bash
    kubectl scale deployment <stargate-deployment> --replicas 1
    ```
{{% /alert %}}

## K8ssandra 1.1.0 revisions

Each of the following sections present a **subset** of key devlopments in K8ssandra 1.1.0. For the complete list, see the [CHANGELOG](https://github.com/k8ssandra/k8ssandra/blob/main/CHANGELOG-1.1.md).  

### Changes

* Shut down cluster by default with in-place restores of Medusa backups [#611](https://github.com/k8ssandra/k8ssandra/issues/611). 
  {{% alert title="Important!" color="warning" %}} 
  The new default behavior for in-place restores, which now shut down the whole cluster by default, will require Stargate nodes to be restarted after the Cassandra cluster is back online. See the [Upgrade notice]({{< relref "#upgrade-notice" >}}) above.
  {{% /alert %}} 
* Update Management API image locations [#637](https://github.com/k8ssandra/k8ssandra/issues/533).

### Enhancements

* Add option to disable Cassandra logging sidecar [#576](https://github.com/k8ssandra/k8ssandra/issues/576).
* Developer documentation [#239](https://github.com/k8ssandra/k8ssandra/issues/239).
* Add support for `additionalSeeds` in the CassandraDatacenter [#547](https://github.com/k8ssandra/k8ssandra/issues/547).

### Bug fixes

* Scale up fails if a restore was performed in the past [#616](https://github.com/k8ssandra/k8ssandra/issues/616).

### Doc updates

* S3-compliant MinIO buckets for Medusa backup and restore operations, and related edits for the separate Amazon S3 topic [#556](https://github.com/k8ssandra/k8ssandra/issues/556). For the updates, start in [Backup and restore Cassandra data]({{< ref "/topics/restore-a-backup/" >}}).
* Migrating existing Cassandra to K8ssandra [#377](https://github.com/k8ssandra/k8ssandra/issues/377). See [Migrating a Cassandra cluster to K8ssandra]({{< ref "/topics/migration/" >}}).
* Underlying considerations for scaling nodes up/down [#501](https://github.com/k8ssandra/k8ssandra/issues/501). See [Scale your Cassandra cluster in K8ssandra]({{< ref "/topics/provision-a-cluster/" >}}).

## Contributions
​
Your feedback and contributions are welcome! To contribute, fork the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra) and submit Pull Requests (PRs) for review.

To submit documentation comments or edits, see [Contribution guidelines]({{< ref "/contribution-guidelines/" >}}).

## Next

Read the K8ssandra [FAQs]({{< ref "/faqs/" >}}) - for starters, how to pronounce "K8ssandra." 

If you're impatient, jump right in with our **[Quick start]({{< ref "getting-started" >}})**.
