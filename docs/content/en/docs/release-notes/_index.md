
---
title: "K8ssandra release notes"
linkTitle: "Release notes"
weight: 1
description: Release notes for the open-source K8ssandra community project.
---

K8ssandra provides a production-ready platform for running Apache Cassandra&reg; on Kubernetes. This includes automation for operational tasks such as repairs, backup and restores, and monitoring. 

Also deployed is Stargate, an open source data gateway that lets you interact programmatically with your Kubernetes-hosted Cassandra resources via a well-defined API. 

The **K8ssandra 1.1.0** release implements a number of key changes, enhancements, and bug fixes, as noted below. For the complete list, see [CHANGELOG-110.md](https://github.com/k8ssandra/k8ssandra/blob/main/CHANGELOG.md).  

**Release date:** 09-April-2021.

## Prerequisites

* A Kubernetes v1.15 or later environment - local or via a supported cloud provider
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

* Apache Cassandra  the deployed version depends on the configured `serverImage` setting:
  * 3.11.7
  * 3.11.8
  * 3.11.9
  * 3.11.10 (default)
* DataStax Kubernetes Operator for Apache Cassandra (cass-operator) 1.6.0
* Management API for Apache Cassandra 0.1.19
* Stargate 1.0.11
* Metric Collector for Apache Cassandra (MCAC) 0.1.9
* Prometheus Operator 2.22.1
* Grafana Operator 7.3.5
* Medusa Operator for Apache Cassandra 0.2.0
* Reaper Operator for Apache Cassandra 2.2.2


## K8ssandra 1.1.0 changes

 *(DRAFT comment: in CHANGELOG-110.md, can we link to the issue numbers (as done below), and can each issue on GH then link to relevant commit IDs?)*

* Remove Jolokia integration [#533](https://github.com/k8ssandra/k8ssandra/issues/533)
* Upgrade to medusa-operator 0.2.0 [#630](https://github.com/k8ssandra/k8ssandra/issues/630)
* Mount Cassandra pod labels in volume [#613](https://github.com/k8ssandra/k8ssandra/issues/613)
* Shut down cluster by default with in-place restores [#611](https://github.com/k8ssandra/k8ssandra/issues/611)

## K8ssandra 1.1.0 enhancements

* Add option to disable Cassandra logging sidecar [#576](https://github.com/k8ssandra/k8ssandra/issues/576)
* Upgrade Reaper to 2.2.2 and Medusa to 0.9.1 [#530](https://github.com/k8ssandra/k8ssandra/issues/530)
* Split dashboards into separate configmaps [#504](https://github.com/k8ssandra/k8ssandra/issues/504)
* Upgrade Stargate to 1.0.11, and add a `preStop` lifecycle hook to improve behavior when reducing the number of Stargate replicas in the presence of live traffic [#436](https://github.com/k8ssandra/k8ssandra/issues/436)
* Add automation for stable and next release streams [#419](https://github.com/k8ssandra/k8ssandra/issues/419)
* Add support for `additionalSeeds` in the CassandraDatacenter [#547](https://github.com/k8ssandra/k8ssandra/issues/547)
* Documentation updates and additions:
  * Developer documentation [#239](https://github.com/k8ssandra/k8ssandra/issues/239)
  * Sample values.yaml explanations and examples [#510](https://github.com/k8ssandra/k8ssandra/issues/510)
  * S3-compliant MinIO buckets for Medusa backup and restore operations, and related edits for the separate Amazon S3 topic [#556](https://github.com/k8ssandra/k8ssandra/issues/556)
  * Migrating existing Cassandra to K8ssandra [#377](https://github.com/k8ssandra/k8ssandra/issues/377) 
  * Underlying considerations and impacts of scaling nodes up or down in Kubernetes clusters [#501](https://github.com/k8ssandra/k8ssandra/issues/501) 

## K8ssandra 1.1.0 bug fixes

* Cassandra config clobbering when enabling Medusa [#475](https://github.com/k8ssandra/k8ssandra/issues/475)
* cqlsh commands show warnings [#396](https://github.com/k8ssandra/k8ssandra/issues/396)
* Fix issue with scripts not being checked out before attempting to run them [#516](https://github.com/k8ssandra/k8ssandra/issues/516)
* Remove GitHub Actions for pre-releasing off of main [#517](https://github.com/k8ssandra/k8ssandra/issues/517)
* Fix Cassandra config clobbering when enabling Medusa [#475](https://github.com/k8ssandra/k8ssandra/issues/475)
* Create cass-operator webhook secret [#590](https://github.com/k8ssandra/k8ssandra/issues/590)
* helm uninstall can leave CassandraDatacenter behind [#623](https://github.com/k8ssandra/k8ssandra/issues/623)
* Fix indentation error in backup-restore-values.yaml [#602](https://github.com/k8ssandra/k8ssandra/issues/602)

## Contributions
​
Your feedback and contributions are welcome! To contribute, fork the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra) and submit Pull Requests (PRs) for review.

To submit documentation comments or edits, see [Contribution guidelines]({{< ref "/docs/contribution-guidelines/" >}}).

## Next

Read the K8ssandra [FAQs]({{< ref "/docs/faqs/" >}}) - for starters, how to pronounce "K8ssandra."

