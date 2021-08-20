---
title: "K8ssandra release notes"
linkTitle: "Release notes"
weight: 2
description: Release notes for the open-source K8ssandra community project.
---

K8ssandra provides a production-ready platform for running Apache Cassandra&reg; on Kubernetes. This includes automation for operational tasks such as repairs, backup and restores, and monitoring. 

Also deployed is Stargate, an open source data gateway that lets you interact programmatically with your Kubernetes-hosted Cassandra resources via a well-defined API. 

{{% alert title="Note" color="success" %}}
**K8ssandra 1.3.0** implements a number of changes, enhancements, and bug fixes. This topic summarizes the key revisions in 1.3.0, provides links to the associated issues in our GitHub repo, and contains important upgrade considerations.

**Reminder**: We've migrated the cass-operator GitHub repo from https://github.com/datastax/cass-operator to https://github.com/k8ssandra/cass-operator. Refer to the new repo for the latest Cass Operator developments.
{{% /alert %}}

**K8ssandra 1.3.0 release date:** 27-July-2021.

## Prerequisites

* A Kubernetes environment, from v1.17 (minimum supported), up to v1.21.1 (current tested upper bound) - local or via a supported cloud provider
* [Helm](https://helm.sh/) v3

## Supported Kubernetes environments

* Open-source [kubernetes.io](https://kubernetes.io)
* [Amazon Elastic Kubernetes Service](https://aws.amazon.com/eks/) (EKS)
* [DigitalOcean Kubernetes](https://www.digitalocean.com/products/kubernetes/) (DOKS)
* [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) (GKE)
* [Microsoft Azure Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/) (AKS)
* [MiniKube](https://minikube.sigs.k8s.io/docs/)
* [Kind](https://kind.sigs.k8s.io/)
* [K3D](https://k3d.io/)

## K8ssandra deployed components

The K8ssandra helm chart deploys the following components. Some are optional, and depending on the configuration, may not be deployed:

* [Apache Cassandra](https://cassandra.apache.org/) - the deployed version depends on the configured `cassandra.version` setting:
  * 4.0.0 (default)
  * 3.11.10
  * 3.11.9
  * 3.11.8
  * 3.11.7
* DataStax Kubernetes Operator for Apache Cassandra ([cass-operator](https://github.com/k8ssandra/cass-operator)) 1.7.1
* Management API for Apache Cassandra ([MAAC](https://github.com/datastax/management-api-for-apache-cassandra)) 0.1.27
* [Stargate](https://github.com/stargate/stargate) 1.0.29
* Metric Collector for Apache Cassandra ([MCAC](https://github.com/datastax/metric-collector-for-apache-cassandra)) 0.2.0
* kube-prometheus-stack 12.11.3 [chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
* Medusa for Apache Cassandra 0.11.0
* medusa-operator 0.3.3
* Reaper for Apache Cassandra 2.2.5
* reaper-operator 2.3.0

{{% alert title="Tip" color="success" %}}
Operators are software extensions to Kubernetes that make use of custom resources to manage applications and their components. Thus, for example, "Reaper Operator" deploys and configures Reaper. "Reaper" itself manages the actual Cassandra repair operations. Similarly, "Prometheus Operator" deploys and configures Prometheus. "Prometheus" itself manages the actual collection of relevant OS / Cassandra metrics. "Medusa Operator" configures and orchestrates the backup and restore operations. "Medusa" itself runs the container that performs backups of Cassandra data. 
{{% /alert %}}

## Upgrade notices

{{% alert title="Important!" color="warning" %}}
Before upgrading to K8ssandra 1.3.0, be sure to read the sections below.
{{% /alert %}}

### Upgrading to K8ssandra 1.3.0 and Cassandra 4.0

When you upgrade from a prior K8ssandra release to 1.3.0, the default `cassandra.version` is now `4.0.0`. A changed `num_tokens` default may prevent Cassandra 4.0 from starting. This section explains how to avoid the issue. 

First, some background information. The default value for `num_tokens` in Cassandra 3.11 is 256. The default in Cassandra 4.0 is 16. All nodes are assigned the same number of tokens based on the `num_tokens` setting. More tokens means that the ring is divided up into smaller ranges. A larger number of tokens means that each individual range will be smaller.

The issue: [Upgrading from Cassandra 3.11 to 4.0 fails if num_tokens is not set](https://github.com/k8ssandra/k8ssandra/issues/1029). As noted for this scenario, the Cassandra logs will include the following error, and Cassandra will not start:

```
org.apache.cassandra.exceptions.ConfigurationException: Cannot change the number of tokens from 256 to 16
```

To avoid the issue, first check the `num_tokens` settings in your K8ssandra-deployed Cassandra 3.11 database. Let's walk through an example.

#### Initial installation

This example starts by installing K8ssandra 1.3.0 with Cassandra 3.11.10 and the latter's default `num_tokens` (256). 

```bash
helm install k8ssandra -f k8ssandra.values.yml k8ssandra/k8ssandra
```

k8ssandra.values.yml:

```yaml
cassandra:
  version: "3.11.10"
  datacenters:
    - 
      name: dc1
      size: 1
```

**Resulting deployment:**

```bash
helm list
```
**Output**:
```
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART           APP VERSION
k8ssandra       default         1               2021-08-17 18:58:03.831212 -0400 EDT    deployed        k8ssandra-1.3.0 
```

After about 10 minutes, verify the Cassandra Operator status:

```bash
kubectl describe cassdc dc1 | grep "Cassandra Operator Progress:"
```
**Output**:
```
  Cassandra Operator Progress:  Ready
```

Verify the server version:

```bash
kubectl describe cassdc dc1 | grep "Server Version:" 
```
**Output**:
```
  Server Version:  3.11.10
```

Verify the `num_tokens`, in this case defined by default value:

```bash
kubectl describe cassdc dc1 | grep "num_tokens:"
```
**Output**:
```
      num_tokens:                         256
```

#### Upgrades

Upgrade to Cassandra 4.0.0 and explicitly sets the `num_tokens` to match the previous default (256).

```bash
helm upgrade k8ssandra -f k8ssandra-upgrade.values.yml k8ssandra/k8ssandra
```

k8ssandra-upgrade.values.yml:

```yaml
cassandra:
  version: "4.0.0"
  datacenters:
    - 
      name: dc1
      size: 1
      num_tokens: 256
```

**Resulting deployment:**

After about 10 minutes, verify the Cassandra Operator status:

```bash
kubectl describe cassdc dc1 | grep "Cassandra Operator Progress:"                              
```
**Output**:
```
  Cassandra Operator Progress:  Ready
```

Check the pods:

```bash
kubectl get pods
```
**Output**:
```
NAME                                                READY   STATUS      RESTARTS   AGE
k8ssandra-cass-operator-5bbf6584b9-pgldj            1/1     Running     0          9m23s
k8ssandra-crd-upgrader-job-k8ssandra-skrjh          0/1     Completed   0          2m6s
k8ssandra-dc1-default-sts-0                         2/2     Running     0          106s
k8ssandra-grafana-679b4bbd74-ngrz2                  2/2     Running     0          9m23s
k8ssandra-kube-prometheus-operator-85695ffb-qhqc7   1/1     Running     0          9m23s
k8ssandra-reaper-694cf49b96-jgnvg                   0/1     Running     1          6m12s
k8ssandra-reaper-operator-5c6cc55869-tdhjw          1/1     Running     0          9m23s
prometheus-k8ssandra-kube-prometheus-prometheus-0   2/2     Running     1          9m14s
```

Verify the server version:

```bash
kubectl describe cassdc dc1 | grep "Server Version:" 
```
**Output**:
```
  Server Version:  4.0.0
```

Verify `num_tokens`:

```bash
kubectl describe cassdc dc1 | grep "num_tokens:"
```
**Output**:
```
      num_tokens:                         256
```

Because the explicitly set `num_tokens` value in 4.0 matches the previous `num_tokens` default in 3.11 (256), Cassandra starts. 

### Upgrading from K8ssandra 1.0.0 to 1.3.0 

Upgrading directly from K8ssandra 1.0.0 to 1.3.0 causes a StatefulSet update (due to [#533](https://github.com/k8ssandra/k8ssandra/issues/533) and [#613](https://github.com/k8ssandra/k8ssandra/issues/613)). A StatefulSet update has the effect of a rolling restart. Because of [#411](https://github.com/k8ssandra/k8ssandra/issues/411) this could require you to perform a manual restart of all Stargate nodes after the Cassandra cluster is back online. This behavior also impacts in-place restore operations of Medusa backups [#611](https://github.com/k8ssandra/k8ssandra/issues/611). To manually restart Stargate nodes:

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


## K8ssandra 1.3.0 revisions

Release date: 27-July-2021

The following sections summarize and link to key revisions in K8ssandra 1.3.0. For the latest, refer to the [CHANGELOG](https://github.com/k8ssandra/k8ssandra/blob/main/CHANGELOG-1.3.md).  

### Changes

* Support for the General Availability (GA) official release of [Apache Cassandra 4.0.0](https://cassandra.apache.org/doc/latest/cassandra/new/index.html). 
* Upgrade to reaper-operator 0.3.3 and Reaper 2.3.0.
* Upgrade from Stargate 1.0.18 to 1.0.29.
* Upgrade from Medusa 0.10.1 to 0.11.0.
* Upgrade from Reaper 2.2.2 to 2.2.5.
* Integrate Fossa component/license scanning, [#812](https://github.com/k8ssandra/k8ssandra/issues/812).
* Upgrade medusa-operator to v0.3.3, [#905](https://github.com/k8ssandra/k8ssandra/issues/905).

### New features

* Upgrade the Management API from 0.1.26 to 0.1.27 to provide support for Cassandra 4.0.0 (GA), and make Cassandra 4.0.0 the default release, [#949](https://github.com/k8ssandra/k8ssandra/issues/949).
* Make affinity configurable for Stargate, [#617](https://github.com/k8ssandra/k8ssandra/issues/617).
* Make affinity configurable for Reaper, [#847](https://github.com/k8ssandra/k8ssandra/issues/847).
* Experimental support for custom init containers, [#952](https://github.com/k8ssandra/k8ssandra/issues/952).

### Enhancements

* Allow configuring the namespace of service monitors, [#844](https://github.com/k8ssandra/k8ssandra/issues/844).
* Detect IEC formatted c* heap.size and heap.newGenSize; return error identifying issue, [#29](https://github.com/k8ssandra/k8ssandra/issues/29). 
Also see: Add validation check for Cassandra heap size properties, [#701](https://github.com/k8ssandra/k8ssandra/issues/701).
* Add support for private registries, [#420](https://github.com/k8ssandra/k8ssandra/issues/420).
* Add support for Medusa backups on Azure, [#685](https://github.com/k8ssandra/k8ssandra/issues/685).

### Bug fixes

* Fix property name in scaling docs, [#853](https://github.com/k8ssandra/k8ssandra/issues/853).
* Hot replace disallowed characters in generated secret names, [#870](https://github.com/k8ssandra/k8ssandra/issues/870).
* Stargate metrics don't show up in the dashboards, [#412](https://github.com/k8ssandra/k8ssandra/issues/412).

### Doc updates

* See the new topics that cover:
  *  [Backup and restore with Azure Storage]({{< relref "tasks/backup-restore/azure/" >}})
  *  [Private registries]({{< relref "tasks/manage/private-registries/" >}})
* [The topics that walk through installing K8ssandra]({{< relref "install" >}}) on AKS, EKS, and GKE include settings and guidelines from the [performance benchmark blog](https://k8ssandra.io/blog/articles/k8ssandra-performance-benchmarks-on-cloud-managed-kubernetes/), which compares throughput and latency between:
  * The baseline performance of a Cassandra cluster running on AWS EC2 instances -- a common setup for enterprises operating Cassandra clusters
  * The performance of K8ssandra running on AKS, EKS, GKE
* [The reference topics]({{< relref "reference/helm-charts" >}}) for the K8ssandra deployed Helm charts have been updated with the latest descriptions.

## K8ssandra 1.2.0 revisions

Release date: 02-June-2021

The following sections briefly summarize and link to key developments in K8ssandra 1.2.0. For the latest list, refer to the [CHANGELOG](https://github.com/k8ssandra/k8ssandra/blob/main/CHANGELOG-1.2.md).  

### Changes

* Upgrade to Cassandra 4.0-rc1 [#726](https://github.com/k8ssandra/k8ssandra/issues/726).

### New features

* Make tolerations configurable [#673](https://github.com/k8ssandra/k8ssandra/issues/673), [#698](https://github.com/k8ssandra/k8ssandra/issues/698). 
{{% alert title="Tips" color="success" %}} 
In Kubernetes, **node affinity** is a property of Pods that attracts them to a set of nodes, as a preference or a hard requirement. **Taints** are the opposite. They allow a node to repel a set of pods. **Tolerations** are applied to pods, and allow (but do not require) the pods to schedule onto nodes with matching taints. Taints and tolerations work together to ensure that pods are not scheduled onto inappropriate nodes. One or more taints are applied to a node. Once applied, the node should not accept any pods that do not tolerate the taints. For more, see [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) in the Kubernetes.io documentation. 
{{% /alert %}} 

### Enhancements

* Add the ability to attach additional Persistent Volumes (PVs) for Medusa backups [#560](https://github.com/k8ssandra/k8ssandra/issues/560).
* Reduce initial delay of Stargate readiness probe [#654](https://github.com/k8ssandra/k8ssandra/issues/654).
* Update cass-operator to v1.7.0 [#693](https://github.com/k8ssandra/k8ssandra/issues/693).
* Make `allocate_tokens_for_replication_factor` configurable [#741](https://github.com/k8ssandra/k8ssandra/pull/741).

### Bug fixes

* Token allocations are random when using 4.0 and lead to collisions [#732](https://github.com/k8ssandra/k8ssandra/issues/732). Related enhancement: Make `allocate_tokens_for_replication_factor` configurable [#741](https://github.com/k8ssandra/k8ssandra/pull/741).
* Upgrade to Medusa 0.10.1 fixing failed backups after a restore [#678](https://github.com/k8ssandra/k8ssandra/issues/678).

### Doc updates

* [Reference topics]({{< relref "reference/helm-charts" >}}) for the K8ssandra deployed Helm charts have been updated with the latest descriptions.


## K8ssandra 1.1.0 revisions

Each of the following sections present a **subset** of key developments in K8ssandra 1.1.0. For the complete list, see the [CHANGELOG](https://github.com/k8ssandra/k8ssandra/blob/main/CHANGELOG-1.1.md).
 
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

* S3-compliant MinIO buckets for Medusa backup and restore operations, and related edits for the separate Amazon S3 topic [#556](https://github.com/k8ssandra/k8ssandra/issues/556). For the updates, start in [Backup and restore Cassandra data]({{< relref "backup-restore" >}}).
* Migrating existing Cassandra to K8ssandra [#377](https://github.com/k8ssandra/k8ssandra/issues/377). See [Migrating a Cassandra cluster to K8ssandra]({{< relref "migrate" >}}).
* Underlying considerations for scaling nodes up/down [#501](https://github.com/k8ssandra/k8ssandra/issues/501). See [Scale your Cassandra cluster in K8ssandra]({{< relref "scale" >}}).

## Contributions
​
Your feedback and contributions are welcome! To contribute, fork the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra) and submit Pull Requests (PRs) for review.

To submit documentation comments or edits, see [Contribution guidelines]({{< relref "contribute" >}}).

## Next steps

Read the K8ssandra [FAQs]({{< relref "faqs" >}}) - for starters, how to pronounce "K8ssandra." 

If you're impatient, jump right in with the K8ssandra [install]({{< relref "install" >}}) steps for these platforms:

* [Local]({{< relref "install/local" >}})
* Amazon Elastic Kubernetes Service ([EKS]({{< relref "install/eks" >}}))
* DigitalOcean Kubernetes ([DOKS]({{< relref "install/doks" >}}))
* Google Kubernetes Engine ([GKE]({{< relref "install/gke" >}}))
* Microsoft Azure Kubernetes Service ([AKS]({{< relref "install/aks" >}}))
