---
title: "K8ssandra FAQs"
linkTitle: "FAQs"
weight: 1
description: Frequently asked questions about the K8ssandra project and K8ssandra Operator.
---

If you're new to K8ssandra and the K8ssandra Operator, this FAQ is for you! 

### What is K8ssandra?

K8ssandra is an open-source project that anyone in the community can use, improve, and enjoy. We've added a unified K8ssandra Operator in March 2022. While you may continue using the initial K8ssandra 1.4.x features (single cluster only), we recommend that you consider moving up to the new **K8ssandra Operator**. For that reason, most of the FAQs here will refer to K8ssandra Operator.

### Ok - how should I pronounce "K8ssandra"?

Any way you want! But think of it this way:  "Kate" + "Sandra".

### What is K8ssandra Operator?

K8ssandra Operator is our latest implementation that provides single- or **multi-cluster, multi-region** support in Kubernetes. It's all part of the overall K8ssandra project, but you'll need to deploy with K8ssandra Operator to use the latest multi-cluster/region features. K8ssandra Operator is a cloud native distribution of Apache Cassandra® that runs on Kubernetes. The K8ssandra Operator GitHub repo is:

https://github.com/k8ssandra/k8ssandra-operator

Accompanying Cassandra is a suite of tools to ease and automate operational tasks. This includes metrics, data anti-entropy services, and backup/restore tools. As part of the install process, by using K8ssandra Operator, all of these components are installed and wired together, freeing your teams from having to perform the tedious plumbing of components.

### What is `K8ssandraCluster`?

The `K8ssandraCluster` is a new custom resource that covers all the bases necessary for installing a production-ready, multi-cluster deployment using K8ssandra Operator. Head over to the DataStax Tech blog to learn more about how to specify your remote clusters with the K8ssandraCluster, its deployment architecture, and what’s coming next in our continued development of the K8ssandra operator.

We built the new K8ssandra Operator to simplify deploying multiple Apache Cassandra data centers in different regions and across multiple Kubernetes (K8s) clusters. Now, it’s easier than ever to run Apache Cassandra® across multiple K8s clusters in multiple regions with the `K8ssandraCluster`. 

### What does K8ssandra Operator include?

At a pure component level, K8ssandra Operator integrates and packages together:

* Apache Cassandra
* Stargate, the open source data gateway
* Cass Operator, also known as `cass-operator`
* Reaper for Apache Cassandra anti-entropy data repair feature
* Medusa for Apache Cassandra backup and restore
* Observability service for integration with metrics and visualization tools
* Templates for connections into your Kubernetes environment via Ingress solutions; or, instructions for using `port-forwarding`, as an alternative approach)

For the full list of deployed components and latest versions, see the [Release notes]({{< relref "release-notes" >}}). 

In addition to the set of components, it's important to emphasize that the K8ssandra project is really a collection of experience from the community of Cassandra + Kubernetes users, packaged and ready for everyone to use freely. 

### How do I get started and install a `K8ssandraCluster` using K8ssandra Operator?

Start with the steps for a [local dev install](https://docs-staging-v2.k8ssandra.io/install/local/). Related install topics cover single- and multi-cluster, multi-region configurations. 

Then proceed to the post-installation quickstart steps for [developers](https://docs-staging-v2.k8ssandra.io/quickstarts/developer/) or [Site Reliability Engineers (SREs)](https://docs-staging-v2.k8ssandra.io/quickstarts/site-reliability-engineer/).

### How do I get started and install K8ssandra 1.4.x?

For single-cluster, single-region only environments, start in the [install](https://docs-staging-v1.k8ssandra.io/install/) topics.

### When I install using K8ssandra or K8ssandra Operator, I see some warning messages. Is that a problem?

When installing with K8ssandra or K8ssandra Operator on newer versions of Kubernetes (v1.19+), some warnings may be visible on the helm command line related to deprecated API usage. This is a known issue and will not impact the provisioning of the cluster.

```bash
W0128 11:24:54.792095  27657 warnings.go:70] apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
```

For more information, see issue [#267](https://github.com/k8ssandra/k8ssandra/issues/267) in the K8ssandra GitHub repo.

### Does K8ssandra Opertor have to be installed in a particular namespace?

You can install a `K8ssandraCluster` using K8ssandra Operator as namespace-scoped or cluste-scoped. The default is namespace-scoped. For details, start in the [local install](https://docs-staging-v2.k8ssandra.io/install/local/) topic. Examples:

Namespace-scoped:

```bash
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
```

Cluster-scoped:

```bash
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
     --set global.clusterScoped=true
```

### Can I install multiple releases of K8ssandra Operator?

You can install multiple releases of K8ssandra Operator in a Kubernetes environment, provided:

* You install one release per namespace
* The release names are unique across the entire Kubernetes cluster
* You install each `K8ssandraCluster` resource as cluster-scoped, one per cluster:

```bash
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
     --set global.clusterScoped=true
```

### How do I install K8ssandra Operator using the K8 included with Docker Desktop for Mac?

When installing K8ssandra Operator in the K8 instance of Docker Desktop for Mac, you may encounter the following error:

```bash
Error: failed pre-install: timed out waiting for the condition
```

To solve the issue, add the `--set cassandra.cassandraLibDirVolume.storageClass=hostpath` option. 

Namespace-scoped example:

```bash
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
     --set cassandra.cassandraLibDirVolume.storageClass=hostpath
```

Cluster-scoped example:

```bash
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
     --set cassandra.cassandraLibDirVolume.storageClass=hostpath --set global.clusterScoped=true
```

### What is Stargate?

[Stargate](https://stargate.io/) is an open source data gateway that abstracts away many Apache Cassandra specific concepts, providing access to the database through various API options.  It helps to remove barriers of entry for developers new to Apache Cassandra by providing REST, GraphQL, and schemaless JSON document based APIs in addition to traditional CQL access.

### What is cass-operator?

Kubernetes Operator for Apache Cassandra, also known as Cass Operator or cass-operator, is the most critical element bridging Kubernetes and Cassandra. The community has been focusing much of its attention on operators over the past two years, as the appropriate starting place. If there is magic happening, it’s all in the operator. The cass-operator serves as the translation layer between the control plane of Kubernetes and actual operation done by the Cassandra cluster. Recently, the Apache Cassandra project agreed on gathering around a single operator: cass-operator. Some great contributions from Orange with CassKop will be merged with the DataStax operator and a final version will be merged into the Apache project. This is the best example of actual production knowledge finding its way into code. Community members contributing to cass-operator are running large amounts of Cassandra in Kubernetes every day. 

### What is Reaper for Apache Cassandra?

Reaper for Apache Cassandra is a tool that helps manage the critical maintenance task of anti-entropy **repairs** in a Cassandra cluster. Originally created by Spotify, later adopted and maintained by The Last Pickle, and one of the features installed by K8ssandra. If you were to sit a group of Cassandra DBAs down to talk about what they do, chances are they would talk a lot about running repairs. It’s an important operation because it keeps data consistent despite inevitable issues that happen like node failures and network partitions. In K8ssandra, Reaper runs it for you automatically! And because this is built for SREs, you can expect a good set of pre-built metrics to verify everything is working great. 

For more, see the [Reaper repair](https://docs-staging-v2.k8ssandra.io/tasks/repair/) topic in the K8ssandra Operator docs.

### What is Medusa for Apache Cassandra?

Medusa for Apache Cassandra provides backup/restore functionality for Cassandra data. This project also originated at Spotify. Medusa not only helps coordinate backup &amp; restore tasks, it manages the placement of the data at rest. The initial implementation allows backup sets to be stored and retrieved on cloud object storage (such as AWS S3 buckets) with more options on the way. 

For more, see the [backup/restore](https://docs-staging-v2.k8ssandra.io/tasks/backup-restore/) topic in the K8ssandra Operator docs.

### How can I access Kubernetes resources from outside the environment?

K8ssandra provides [preconfigured]({{< relref "/tasks/connect/ingress/" >}}) Ingress integrations, such as Traefik, which is a modern reverse proxy and load balancer that makes deploying microservices easy. Traefik integrates with your existing infrastructure components and configures itself automatically and dynamically. Traefik handles advanced ingress deployments including mTLS of TCP with SNI and UDP. Operators define rules for routing traffic to downstream systems through Kubernetes Ingress objects or more specific Custom Resource Definitions. K8ssandra supports deploying `IngressRoute objects` as part of a deployment to expose metrics, repair, and Cassandra interfaces. For more, see the [Connect](https://docs-staging-v2.k8ssandra.io/tasks/connect/) topic in the K8ssandra Operator docs. 

### How can I monitor the health of my Kubernetes + Cassandra cluster?

You can configure Traefik to expose the K8ssandra monitoring interfaces, or use port-forwarding. See:

* [Monitoring with Traefik](https://docs-staging-v2.k8ssandra.io/tasks/connect/ingress/monitoring/)
* [Configure port-forwarding](https://docs-staging-v1.k8ssandra.io/quickstarts/site-reliability-engineer/#port-forwarding)

After completing the prerequisites, for example in your local environment, you can open <http://grafana.localhost:8080>. 

### What is the login for the Grafana dashboards?

The default configured Grafana username is `admin`, and the password is `secret`. See the topic about managing [Grafana credentials](https://docs-staging-v2.k8ssandra.io/tasks/monitor#grafana-credentials).

### What kind of provisioning tasks can I perform with K8ssandra?

Among the tasks are to dynamically scale up or down the size of your cluster. See the [scaling task](https://docs-staging-v2.k8ssandra.io/tasks/scale/").

### How can I backup and restore my Cassandra data?

Backup and restore Cassandra data to/from a supported storage object, such as an Amazon S3 bucket or Google Cloud Storage. See [Backup and restore Cassandra](https://docs-staging-v2.k8ssandra.io/tasks/backup-restore/).

### How do I schedule and orchestrate repairs of my Cassandra data?

Periodically run anti-entropy operations to repair your Cassandra data. A general recommendation is once every 7-10 days. With the Reaper UI, you can schedule repairs, run repairs, and check the cluster's health. See [Reaper for Apache Cassandra repairs]({{< relref "/tasks/repair" >}}).

For command-line and UI details, see the [tasks](https://docs-staging-v2.k8ssandra.io/tasks/) topics.

### How can I contribute to the K8ssandra Operator docs?

See the code and documentation [Contributing guidelines]({{< relref "contribute" >}}).

## Next steps

* [Install](https://docs-staging-v2.k8ssandra.io/install/local/): Learn how to deploy a `K8ssandraCluster` custom resource, using K8ssansdra Operator, in your local dev Kubernetes. 
* [Components](https://docs-staging-v2.k8ssandra.io/components/): Dig in to each deployed component of the K8ssandra stack and see how it communicates with the others.
* [Tasks](https://docs-staging-v2.k8ssandra.io/tasks/): Need to get something done? Check out the Tasks topics for a helpful collection of outcome-based solutions.
* [Reference](https://docs-staging-v2.k8ssandra.io/reference/): Explore the K8ssandra Operator reference topics on the provided Cassandra Custom Resource Definitions (CRDs).

We encourage you to actively participate in the [K8ssandra community](https://k8ssandra.io/community/).
