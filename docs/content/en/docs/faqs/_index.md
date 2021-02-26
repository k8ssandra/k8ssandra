---
title: "K8ssandra FAQs"
linkTitle: "FAQs"
weight: 2
description: Frequently asked questions about K8ssandra.
---

If you're new to K8ssandra, this FAQ is for you! Whether you're viewing this page in [GitHub](https://github.com/k8ssandra/k8ssandra/blob/main/docs/content/en/docs/faqs/_index.md) or on the [Web](https://k8ssandra.io/docs/faqs/), you can also propose new or modified FAQs. For this open-source project, contributions are welcome from the community of K8ssandra users.

### What is K8ssandra?

K8ssandra is an open-source project that anyone in the community can use, improve, and enjoy. K8ssandra is a cloud native distribution of Apache Cassandra® that runs on Kubernetes. Accompanying Cassandra is a suite of tools to ease and automate operational tasks. This includes metrics, data anti-entropy services, and backup/restore tools. As part of K8ssandra’s installation process, all of these components are installed and wired together, freeing your teams from having to perform the tedious plumbing of components.

### Ok - how should I pronounce "K8ssandra"?

Any way you want! But think of it this way:  "Kate" + "Sandra".

### What does K8ssandra include?

At a pure component level, K8ssandra integrates and packages together:

* Apache Cassandra
* Stargate, the open source data gateway
* Kubernetes Operator for Apache Cassandra (cass-operator)
* Reaper, an anti-entropy repair feature for Apache Cassandra (reaper-operator)
* Medusa for backup and restore (medusa-operator)
* Metrics Collector, with Prometheus integration, and visualization via preconfigured Grafana dashboards
* Templates for connections into your Kubernetes environment via Ingress solutions

An illustration always helps:

![K8ssandra components](k8ssandra-components2.png)

In addition to the set of components, it's important to emphasize that K8ssandra is really a collection of experience from the community of Cassandra + Kubernetes users, packaged and ready for everyone to use freely. 

### How do I get started and install K8ssandra?

It's easy! There are several options, but we recommend using [Helm](https://helm.sh/docs/intro/install/) commands. 

```bash
helm repo add k8ssandra https://helm.k8ssandra.io/
helm repo add traefik https://helm.traefik.io/traefik
helm repo update
helm install k8ssandra k8ssandra/k8ssandra
```

For more, see [Getting Started]({{< ref "/docs/getting-started/" >}}).

### When I install K8ssandra, I see a some warning messages, is that a problem?

When installing K8ssandra on newer versions of Kubernetes (v1.19+), some warnings may be visible on the command line
related to deprecated API usage.  This is currently a known issue and will not impact the provisioning of the cluster.

```bash
W0128 11:24:54.792095  27657 warnings.go:70] apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
```

For more information, check out issue [#267](https://github.com/k8ssandra/k8ssandra/issues/267).

### What does K8ssandra install?

The chart installs Kubernetes Operator for Apache Cassandra (cass-operator), Prometheus Operator, Reaper (repairs), Medusa (backup/restore), the Grafana Operator, (optional) Stargate, and launches instances.

After those installs, and all the pods are in a Ready state, from `kubectl get pods` you'll see output similar to:

```bash
NAME                                                   READY   STATUS      RESTARTS   AGE
demo-cass-operator-65cc657-fq6bc                       1/1     Running     0          10m
demo-dc1-default-sts-0                                 3/3     Running     0          10m
demo-dc1-stargate-bb47877d5-54sdt                      1/1     Running     0          10m
demo-grafana-7f84d96d47-xd79s                          2/2     Running     0          10m
demo-kube-prometheus-stack-operator-76b984f9f4-pp745   1/1     Running     0          10m
demo-medusa-operator-6888946787-qwzsx                  1/1     Running     2          10m
demo-reaper-k8ssandra-656f5b77cc-nqfzv                 1/1     Running     0          10m
demo-reaper-k8ssandra-schema-88cpx                     0/1     Completed   0          10m
demo-reaper-operator-5b8c4c66b8-8cf86                  1/1     Running     2          10m
prometheus-demo-kube-prometheus-stack-prometheus-0     2/2     Running     1          10m
```

### Does K8ssandra have to be installed in a particular namespace?

The chart can be installed to any namespace. The following example demonstrates this:

```bash
helm install demo k8ssandra/k8ssandra -n k8ssandra --create-namespace
```

### Can I install multiple releases of K8ssandra?

Some of the objects installed by the K8ssandra chart are currently configured to be cluster-scoped; consequently, you should only install those components once. This should be fixed before version 1.0 to allow multiple installations. Other parts can be installed multiple times to allow creating multiple Cassandra clusters in a single k8s cluster.

### What components does K8ssandra install?

K8ssandra deploys the following components, some components are optional, and depending on the configuration, may not be deployed:

* [Apache Cassandra](https://cassandra.apache.org/) (version deployed is dependent upon configuration)
  * 3.11.7
  * 3.11.8
  * 3.11.9
  * 3.11.10 (default)
* [Management API for Apache Cassandra](https://github.com/datastax/management-api-for-apache-cassandra)
  * 0.1.19
* [Metric Collector for Apache Cassandra (MCAC)](https://github.com/datastax/metric-collector-for-apache-cassandra)
  * 0.1.9
* [Prometheus](https://prometheus.io/)
  * 2.22.1
* [Grafana](https://grafana.com/)
  * 7.3.5
* [Medusa for Apache Cassandra](https://github.com/thelastpickle/cassandra-medusa)
  * 0.9.0
* [Reaper for Apache Cassandra](http://cassandra-reaper.io/)
  * 2.2.1
* [Stargate](https://stargate.io/)
  * 1.0.9

{{% alert title="Note" color="primary" %}}
Throughout these docs, examples are shown to deploy [Traefik](https://traefik.io/) as a means to provide external access to the k8ssandra cluster.  It is deployed separately from K8ssandra, and as such, the version deployed will vary.*
{{% /alert %}}

### What is Stargate?

[Stargate](https://stargate.io/) is an open source data gateway that abstracts away many Apache Cassandra specific concepts, providing access to the database through various API options.  It helps to remove barriers of entry for developers new to Apache Cassandra by providing REST, GraphQL, and schemaless JSON document based APIs in addition to traditional CQL access.

### What is cass-operator?

Kubernetes Operator for Apache Cassandra -- [cass-operator](https://github.com/datastax/cass-operator) -- is the most critical element bridging Kubernetes and Cassandra. The community has been focusing much of its attention on operators over the past two years, as the appropriate starting place. If there is magic happening, it’s all in the operator. The cass-operator serves as the translation layer between the control plane of Kubernetes and actual operation done by the Cassandra cluster. Recently, the Apache Cassandra project agreed on gathering around a single operator: cass-operator. Some great contributions from Orange with CassKop will be merged with the DataStax operator and a final version will be merged into the Apache project. This is the best example of actual production knowledge finding its way into code. Community members contributing to cass-operator are running large amounts of Cassandra in Kubernetes every day. 

### What is Reaper for Apache Cassandra?

Reaper for Apache Cassandra (Reaper) is a tool that helps manage the critical maintenance task of anti-entropy **repair** in a Cassandra cluster. Originally created by Spotify, later adopted and maintained by The Last Pickle, and one of the features installed by K8ssandra. If you were to sit a group of Cassandra DBAs down to talk about what they do, chances are they would talk a lot about running repairs. It’s an important operation because it keeps data consistent despite inevitable issues that happen like node failures and network partitions. In K8ssandra, Reaper runs it for you automatically! And because this is built for SREs, you can expect a good set of pre-built metrics to verify everything is working great. See the [Reaper Web Interface for Cassandra repairs]({{< ref "/docs/topics/repair/" >}}).

### What is Medusa for Apache Cassandra?

Medusa for Apache Cassandra (Medusa) provides backup/restore functionality for Cassandra data; this project also originated at Spotify. Medusa not only helps coordinate backup &amp; restore tasks, it manages the placement of the data at rest. The initial implementation allows backup sets to be stored and retrieved on cloud object storage (such as AWS S3 buckets) with more options on the way. K8ssandra offers this [backup and restore]({{< ref "/docs/topics/restore-a-backup/" >}}) feature to help you recover Cassandra data when inevitable real-world issues occur.

### How can I access Kubernetes resources from outside the environment?

K8ssandra provides [preconfigured]({{< ref "/docs/topics/ingress/traefik/" >}}) Ingress integrations, such as Traefik, which is a modern reverse proxy and load balancer that makes deploying microservices easy. Traefik integrates with your existing infrastructure components and configures itself automatically and dynamically. Traefik handles advanced ingress deployments including mTLS of TCP with SNI and UDP. Operators define rules for routing traffic to downstream systems through Kubernetes Ingress objects or more specific Custom Resource Definitions. K8ssandra supports deploying `IngressRoute objects` as part of a deployment to expose metrics, repair, and Cassandra interfaces. For more, start in the [Traefik]({{< ref "/docs/topics/ingress/traefik/" >}}) topic.

### How can I monitor the health of my Kubernetes + Cassandra cluster?

Configure Traefik to expose the K8ssandra monitoring interfaces. See [Monitoring]({{< ref "/docs/topics/ingress/traefik/monitoring/" >}}) for the steps to enable the Traefik Ingress. Then see [Monitoring UI]({{< ref "/docs/topics/monitoring/" >}}) for details about how to access the preconfigured Grafana dashboards that K8ssandra provides. After completing the prerequisites, for example in your local environment, you can open http://grafana.localhost:8080/ in your browser.

### What is the login for the Grafana dashboards?

The default configured Grafana username is `admin`, and the password is `secret`. See the topic about managing [Grafana credentials]({{< ref "/docs/topics/monitoring/_index.md#grafana-credentials" >}}).

### What kind of provisioning tasks can I perform with K8ssandra?

Among the tasks are to dynamically scale up or down the size of your cluster. See [Provisioning a cluster]({{< ref "/docs/topics/provision-a-cluster/" >}}).

### How can I backup and restore my Cassandra data?

Backup and restore Cassandra data to/from a supported storage object, such as an Amazon S3 bucket or Google Cloud Storage. See [Backup and restore Cassandra]({{< ref "/docs/topics/restore-a-backup/" >}}).

### How do I schedule and orchestrate repairs of my Cassandra data?

Periodically run anti-entropy operations to repair your Cassandra data. A general recommendation is once every 7-10 days. With the Reaper UI, you can schedule repairs, run repairs, and check the cluster's health. See [Reaper for Apache Cassandra repairs]({{< ref "/docs/topics/repair/" >}}).

### Can you illustrate the steps and sample commands I'll use with K8ssandra?

Yes - here are the steps and commands in a single graphic:

![K8ssandra steps](k8ssandra-steps2.png)

For command-line and UI details, see K8ssandra [tasks]({{< ref "/docs/topics/" >}}).

### How can I contribute to the K8ssandra docs?

See the documentation [guidelines]({{< ref "/docs/contribution-guidelines/" >}}) topic.

## Next

Read the documentation [tasks]({{< ref "/docs/topics" >}}) and actively participate in the [community](https://k8ssandra.io/community/) of K8ssandra users.
