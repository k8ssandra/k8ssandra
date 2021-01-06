---
title: "K8ssandra FAQs"
linkTitle: "K8ssandra FAQs"
weight: 1
description: Frequently asked questions about K8ssandra.
---

If you're new to K8ssandra, this FAQ is for you! Whether you're viewing this page in [GitHub](https://github.com/k8ssandra/k8ssandra/blob/main/docs/content/en/docs/faqs/_index.md) or on the [Web](https://k8ssandra.io/docs/faqs/), you can also propose new or modified FAQs. For this open-source project, contributions are welcome from the community of K8ssandra users. 

### What is K8ssandra?

K8ssandra is an open-source project that anyone in the community can use, improve, and enjoy. K8ssandra is a cloud native distribution of Apache Cassandra&reg; that runs on Kubernetes. Accompanying Cassandra is a suite of tools to ease and automate operational tasks. This includes metrics, data anti-entropy services, and backup/restore tools. As part of K8ssandra’s installation process, all of these components are installed and wired together, freeing your teams from having to perform the tedious plumbing of components.

### Ok - how should I pronounce "K8ssandra"?

Any way you want! But think of it this way:  "Kate" + "Sandra".

### What does K8ssandra include?

At a pure component level, K8ssandra integrates and packages together:

* Apache Cassandra 3.11.7
* Kubernetes Operator for Apache Cassandra (cass-operator)
* Reaper, also known as the Repair Web Interface
* Medusa for backup and restore
* Metrics Collector, with Prometheus integration, and visualization via preconfigured Grafana dashboards
* Stargate, the open source data gateway
* Templates for connections into your Kubernetes environment via Ingress solutions such as Traefik

An illustration always helps:

![K8ssandra components](k8ssandra-components.png)

In addition to the set of components, it's important to emphasize that K8ssandra is really a collection of experience from the community of Cassandra + Kubernetes users, packaged and ready for everyone to use freely. 

### How do I get started and install K8ssandra?

It's easy! There are several options, but we recommend using [Helm](https://helm.sh/docs/intro/install/) commands. 

```
helm repo add k8ssandra https://helm.k8ssandra.io/
helm repo add traefik https://helm.traefik.io/traefik
helm repo update
helm install k8ssandra-tools k8ssandra/k8ssandra
helm install k8ssandra-cluster-a k8ssandra/k8ssandra-cluster  
```

For more, see [Getting Started]({{< ref "/docs/getting-started/" >}}).

### What exactly do k8ssandra and k8ssandra-cluster install?

Referring to the helm commands from the prior FAQ:

* `k8ssandra` installs Kubernetes Operator for Apache Cassandra (cass-operator) and the Prometheus Operator.
* `k8ssandra-cluster` installs an instance of the stack: reaper (repairs), medusa (backup/restores), the Grafana Operator, (optional) Stargate, and instances.

After those installs, and all the pods are in a Ready state, from `kubectl get pods` you'll see output similar to:

```
NAME                                                              READY   STATUS      RESTARTS   AGE
cass-operator-65956c4f6d-f25nl                                    1/1     Running     0          10m
grafana-deployment-8467d8bc9d-czsg5                               1/1     Running     0          6m23s
k8ssandra-cluster-a-grafana-operator-k8ssandra-5bcb746b8d-4nlhz   1/1     Running     0          6m20s
k8ssandra-cluster-a-reaper-k8ssandra-6cf5b87b8f-vxrwj             1/1     Running     6          6m20s
k8ssandra-cluster-a-reaper-k8ssandra-schema-pjmv8                 0/1     Completed   5          6m20s
k8ssandra-cluster-a-reaper-operator-k8ssandra-55dc486998-f4r46    1/1     Running     2          6m20s
k8ssandra-dc1-default-sts-0                                       2/2     Running     0          10m
k8ssandra-tools-kube-prome-operator-6d57f758dd-7zd92              1/1     Running     0          10m
prometheus-k8ssandra-cluster-a-prometheus-k8ssandra-0             2/2     Running     1          10m
```

### Do k8ssandra and k8ssandra-cluster have to be installed in a particular namespace?

Both charts can be installed in any namespace. Furthermore, you can install them in separate namespaces. The following example demonstrates this:

```
# Install k8ssandra-tool in the k8ssandra namespace
$ helm install k8ssandra-tools k8ssandra/k8ssandra -n k8ssandra --create-namespace

# Install k8ssandra in the k8ssandra-dev namespace
$ helm install dev-cluster k8ssandra/k8ssandra-cluster -n k8ssandra-dev --create-namespace
```

### Can I install multiple releases of k8ssandra?

The objects installed by the k8ssandra chart are all currently configured to be cluster-scoped; consequently, you should only install it once.

### Can I install multiple releases of k8ssandra-cluster?

Yes, you can install multiple releases of k8ssandra-cluster. Do to this [issue](https://github.com/integr8ly/grafana-operator/issues/306) with grafana-operator, each release should be installed in a separate namespace.

### What is cass-operator?

Kubernetes Operator for Apache Cassandra -- [cass-operator](https://github.com/datastax/cass-operator) -- is the most critical element bridging Kubernetes and Cassandra. The community has been focusing much of its attention on operators over the past two years, as the appropriate starting place. If there is magic happening, it’s all in the operator. The cass-operator serves as the translation layer between the control plane of Kubernetes and actual operation done by the Cassandra cluster. Recently, the Apache Cassandra project agreed on gathering around a single operator: cass-operator. Some great contributions from Orange with CassKop will be merged with the DataStax operator and a final version will be merged into the Apache project. This is the best example of actual production knowledge finding its way into code. Community members contributing to cass-operator are running large amounts of Cassandra in Kubernetes every day. 

### What version of Cassandra does K8ssandra install and manage via the cass-operator?

It's Apache Cassandra 3.11.7 currently.

### What is Reaper?

Reaper is a tool that helps manage the critical maintenance task of anti-entropy **repair** in a Cassandra cluster. We also refer to Reaper as the [Repair Web Interface]({{< ref "/docs/topics/accessing-services/repair/" >}}). Originally created by Spotify, later adopted and maintained by The Last Pickle. If you were to sit a group of Cassandra DBAs down to talk about what they do, chances are they would talk a lot about running repairs. It’s an important operation because it keeps data consistent despite inevitable issues that happen like node failures and network partitions. In K8ssandra, Reaper runs it for you automatically! And because this is built for SREs, you can expect a good set of pre-built metrics to verify everything is working great. 

### What is Medusa?

Medusa provides backup/restore functionality for Cassandra data; this project also originated at Spotify. Medusa not only helps coordinate backup &amp; restore tasks, it manages the placement of the data at rest. The initial implementation allows backup sets to be stored and retrieved on cloud object storage (such as AWS S3 buckets) with more options on the way. K8ssandra offers this [backup and restore]({{< ref "/docs/topics/restore-a-backup/" >}}) feature to help you recover Cassandra data when inevitable real-world issues occur.

### What is Stargate?

[Stargate](https://stargate.io/) is an open source data gateway that abstracts away many Apache Cassandra specific concepts, providing access to the database through various API options.  It helps to remove barriers of entry for developers new to Apache Cassandra by providing REST, GraphQL, and schemaless JSON document based APIs in addition to traditional CQL access.

### How can I access Kubernetes resources from outside the environment?

K8ssandra provides [preconfigured]({{< ref "/docs/topics/ingress/traefik/" >}}) Traefik Ingress integrations. Traefik is a modern reverse proxy and load balancer that makes deploying microservices easy.  Traefik integrates with your existing infrastructure components and configures itself automatically and dynamically. Traefik handles advanced ingress deployments including mTLS of TCP with SNI and UDP. Operators define rules for routing traffic to downstream systems through Kubernetes Ingress objects or more specific Custom Resource Definitions. K8ssandra supports deploying `IngressRoute objects` as part of a deployment to expose metrics, repair, and Cassandra interfaces. For more, start in the [Traefik]({{< ref "/docs/topics/ingress/traefik/" >}}) topic.

### How can I monitor the health of my Kubernetes + Cassandra cluster?

Configure Traefik to expose the K8ssandra monitoring interfaces. See [Monitoring]({{< ref "/docs/topics/ingress/traefik/monitoring/" >}}) for the steps to enable the Traefik Ingress. Then see [Monitoring UI]({{< ref "/docs/topics/accessing-services/monitoring/" >}}) for details about how to access the preconfigured Grafana dashboards that K8ssandra provides. After completing the prerequisites, for example in your local environment, you can open http://grafana.localhost:8080/ in your browser. 

### What is the login for the Grafana dashboards?

The default configured Grafana username is `admin`, and the password is `secret`. See the topic about managing [Grafana credentials]({{< ref "/docs/topics/accessing-services/monitoring/_index.md#grafana-credentials" >}}).

### What kind of provisioning tasks can I perform with K8ssandra?

Among the tasks are dynamically scaling up or down the size of your cluster. See [Provisioning a cluster]({{< ref "/docs/topics/provision-a-cluster/" >}}). 

### Can you illustrate the steps and sample commands I'll use with K8ssandra?

Yes - here are the steps and commands in a single graphic:

<p><img style="height: 100%; width: 100%; object-fit: contain" src="https://k8ssandra.io/docs/faqs/k8ssandra-steps.png"></p>

### How can I contribute to the K8ssandra docs?

See the documentation [guidelines]({{< ref "/docs/contribution-guidelines/" >}}) topic. 

## Next

Read the documentation [topics]({{< ref "/docs/topics" >}}) and actively participate in the [community](https://k8ssandra.io/community/) of K8ssandra users.
