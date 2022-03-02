---
title: "K8ssandra glossary"
linkTitle: "Glossary"
weight: 7
description: Helpful definitions of common Kubernetes terms along with the K8ssandra or Apache Cassandra context.
---

### AKS
The [Azure Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/) from Microsoft. One of the "top 3" major cloud providers supported by K8ssandra, along with EKS from Amazon, and GKE from Google. AKS offers serverless Kubernetes, an integrated continuous integration and continuous delivery (CI/CD) experience, and enterprise-grade security and governance.  

### anti-entropy
The process of comparing the data of all replicas and updating each replica to the newest version. Cassandra has two phases to the process: Build a Merkle tree for each replica. Compare the Merkle trees to discover differences. K8ssandra provides [Reaper](https://docs-v2.k8ssandra.io/components/reaper/) as one of its deployed components, enabling you to perform Cassandra [repair](https://docs-v2.k8ssandra.io/tasks/repair/) operations. 

### Astra DB
A [CNDB]({{< relref "#cndb" >}}) product from DataStax that gives you the ability to develop and deploy data-driven applications with a cloud-native service, without the hassles of database and infrastructure administration. By automating tuning and configuration, [Astra](https://astra.datastax.com/) radically simplifies database and streaming operations. See the Astra DB [documentation](https://docs.datastax.com/en/astra-streaming/docs/).

### Astra Streaming
A cloud native messaging and event streaming platform powered by Apache Pulsar. You can quickly create Pulsar instances, manage clusters, scale across cloud regions, and manage Pulsar resources such as topics, connectors, functions, and subscriptions. See the Astra Streaming [documentation](https://docs.datastax.com/en/astra-streaming/docs/).

### charts
[Helm charts](https://helm.sh/) are a YAML-based packaging format to create, version, share, and publish software in Kubernetes. A Helm chart is a collection of templates and settings that describe a set of Kubernetes resources. For details about each Helm chart provided by K8ssandra 1.4.x, see the [Helm chart](https://docs-v1.k8ssandra.io/reference/helm-charts/) reference topics.


### CNDB
An acronym for Cloud Native DataBase, which refers to a database that is created and managed in a cloud environment. DataStax Astra DB is a CNDB, as well as an Apache Cassandra instance that's deployed to a Kubernetes cloud provider (such as GKE, EKS, AKS, DOKS).

### container
An image that is a ready-to-run software package with everything needed to run an application: code, runtime, system tools, system libraries, and settings.

### CQL
Cassandra Query Language is a set of DDL and DML statements designed for communicating with Apache Cassandra databases. CQL offers a model close to SQL in the sense that data is put in tables containing rows of columns. 

### CQLSH
A command-line shell (Cassandra Query Language Shell) for interacting with Cassandra through CQL. CQLSH is included with every Cassandra package, and can be found in the `bin/` directory alongside the cassandra executable. CQLSH utilizes the Python native protocol driver, and connects to the single node specified on the command line.

### EKS
Amazon [Elastic Kubernetes Service](https://aws.amazon.com/eks/) is one of the "top 3" major cloud providers supported by K8ssandra, along with GKE from Google, and AKS from Microsoft. EKS allows you to start, run, and scale Kubernetes applications in the AWS cloud or on-premises. 

### GKE
[Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) is one of the "top 3" major cloud providers supported by K8ssandra, along with EKS from Amazon, and AKS from Microsoft. GKE includes a set of UI-based tools that are part of the Google Cloud Console for GCP environments. 

### gossip
In Cassandra, a protocol to discover location and state information about the other nodes participating in the cluster. Gossip is a peer-to-peer communication protocol in which nodes periodically exchange state information about themselves and about other known nodes.

### Grafana
A multi-platform open source analytics and interactive visualization web application. It provides charts, graphs, and alerts for the web when connected to supported data sources. K8ssandra provides preconfigured Grafana dashboards that visualize Cassandra, cluster, OS and node metrics that are captured at runtime by [Prometheus]({{< relref "#prometheus" >}}) (also provided by K8ssandra). See the [Metrics Collector](https://docs-v2.k8ssandra.io/components/metrics-collector/) component, and the monitor [task](https://docs-v2.k8ssandra.io/tasks/monitor/).

### Helm
Commonly used [tool](https://helm.sh/) that helps you manage Kubernetes applications. K8ssandra works with Helm v3. It includes a command-line tool, a standard for chart definitions, and a repository for use in Kubernetes.

### Helm chart
Used to define, install, and upgrade Kubernetes applications. See the chart [summary](https://helm.sh/docs/topics/charts/) on the Helm site. Refer to the K8ssandra [reference](https://docs-v1.k8ssandra.io/reference/helm-charts/) topics for details about the Helm charts deployed by K8ssandra. Also see the single K8ssandra Operator [Helm chart](https://docs-v2.k8ssandra.io/reference/helm-chart/k8ssandra-operator/) reference topic. 

### Helm repository
The place where charts are collected and shared for Kubernetes packages. For example, you can use `helm repo add k8ssandra https://helm.k8ssandra.io/stable`, and `helm repo update`, to stay current with the latest software. A subsequent command will deploy a `K8ssandraCluster`, using K8ssandra Operator, with a command such as: 

`helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace`

### keyspace
The top-level database object that controls the replication for the object it contains at each datacenter in the cluster. Keyspaces contain tables, materialized views, user-defined types, functions,  and aggregates. 

### kubectl
A command-line tool that allows you to run commands against Kubernetes clusters. You can use `kubectl` ([Kubernetes control](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)) to deploy applications, inspect and manage cluster resources, and view logs.

### Kubernetes
A portable, extensible, [open source platform](https://kubernetes.io/) for managing containerized workloads and services, that facilitates both declarative configuration and automation. It has a large, rapidly growing ecosystem. The name Kubernetes originates from Greek, meaning helmsman or pilot.

### K8ssandraCluster

A custom resource that covers all the bases necessary for installing a production-ready, multi-cluster deployment using K8ssandra Operator. See this DataStax [Tech blog](https://medium.com/building-the-open-data-stack/deploying-to-multiple-kubernetes-clusters-with-the-k8ssandra-operator-f7562bee1841) to learn how to specify your remote clusters with the `K8ssandraCluster`, its deployment architecture, and what's coming next in our continued development of the K8ssandra Operator. We built the new K8ssandra Operator to simplify deploying multiple Apache Cassandra data centers in different regions and across multiple Kubernetes (K8s) clusters. Now, it’s easier than ever to run Apache Cassandra® across multiple K8s clusters in multiple regions with the `K8ssandraCluster`. 

### K8ssandra Operator
Our latest implementation that provides single- or **multi-cluster, multi-region** support in Kubernetes. It's all part of the overall K8ssandra project, but you'll need to deploy with K8ssandra Operator to use the latest multi-cluster/region features. K8ssandra Operator is a cloud native distribution of Apache Cassandra® that runs on Kubernetes. The K8ssandra Operator GitHub repo is:

https://github.com/k8ssandra/k8ssandra-operator

Accompanying Cassandra is a suite of tools to ease and automate operational tasks. This includes metrics, data anti-entropy services, and backup/restore tools. As part of the install process, by using K8ssandra Operator, all of these components are installed and wired together, freeing your teams from having to perform the tedious plumbing of components.

### K8ssandra (1.4.x)
The initial project implementation, an open source, production-ready platform for running Apache Cassandra® on Kubernetes. [K8ssandra](https://k8ssandra.io) includes automation for operational tasks such as [repairs](https://docs-v2.k8ssandra.io/tasks/repair), [backup/restore](https://docs-v2.k8ssandra.io/reference/tasks/backup-restore), and [monitoring](https://docs-v2.k8ssandra.io/tasks/monitor).

### Medusa
An open source backup and restore tool for Cassandra data, deployed by K8ssandra for Kubernetes environments. For more, see [Medusa component](https://docs-v2.k8ssandra.io/components/medusa/) and [backup and restore tasks](https://docs-v2.k8ssandra.io/tasks/backup-restore/). 

### minikube
A tool that lets you run Kubernetes locally. [Minikube](https://minikube.sigs.k8s.io/docs/) runs a single-node Kubernetes cluster on your personal computer (including Windows, macOS and Linux PCs) so that you can try out Kubernetes, or for daily development work.

### MinIO
An Amazon S3-compatible server-side software storage stack. MinIO is one of the local or cloud-based storage objects ("buckets") supported by K8ssandra's Medusa backup/restore operations. For more, see [Backup and restore with MinIO](https://docs-v2.k8ssandra.io/tasks/backup-restore/minio/) buckets.

### namespace
A way to provide a scope for names. Names of resources need to be unique within a namespace, but not across namespaces. Note that namespaces cannot be nested inside one another and each Kubernetes resource can only be in one namespace.

### NetworkTopologyStrategy
In Cassandra, a data replication strategy that places replicas in the same CassandraDatacenter by walking the ring clockwise until reaching the first node in another rack. See also [SimpleStrategy]({{< relref "#simplestrategy" >}}).

### nodetool
A Cassandra [command-line interface](https://cassandra.apache.org/doc/latest/tools/nodetool/nodetool.html) for monitoring a cluster and performing routine database operations. It is typically run from an operational node, and includes commands such as `nodetool repair`. For repair operations in Kubernetes, we recommend an alternative: [Reaper](https://docs-v2.k8ssandra.io/components/reaper/), which is deployed by K8ssandra Operator. Also see the Cassandra [repair tasks](https://docs-v2.k8ssandra.io/tasks/repair/).

### pod
Represents a single instance of a running process in your cluster. Pods contain one or more containers, such as Docker containers. When a Pod runs multiple containers, the containers are managed as a single entity and share the Pod's resources. here, just checking format options. For Cassandra and DataStax Enterprise users, a "node" in a cluster is the equivalent of a pod.

### port forwarding
An application of network address translation that redirects a communication request from one address and port number combination to another, while the packets are traversing a network gateway, such as a router or firewall. For information about using port forwarding with K8ssandra Operator deployments, see:
* Developers, see [Set up port forwarding](https://docs-v2.k8ssandra.io/quickstarts/developer/#set-up-port-forwarding).  
* Site reliability engineers, see [Configure port forwarding](https://docs-v2.k8ssandra.io/quickstarts/site-reliability-engineer/#port-forwarding).

### Prometheus
An open source tool deployed by K8ssandra and used for event monitoring and alerting. [Prometheus](https://prometheus.io) records real-time metrics in a time series database built using a HTTP pull model, with flexible queries and real-time alerting. K8ssandra Operator provides preconfigured [Grafana]({{< relref "#grafana" >}}) dashboards that display the cluster, OS, and node metrics collected by Prometheus in your Kubernetes environment.  

### rack
In the context of a CassandraDatacenter topology, a rack is a logical grouping of Cassandra nodes within the ring. Cassandra uses racks so that it can ensure replicas are distributed among different logical groupings. The number of racks should equal the replication factor (RF) of your application keyspaces. Cassandra ensures that replicas are spread across racks, versus having multiple replicas within the same rack. For example, let’s say you are using RF = 3 with a 9-node cluster and 3 racks (and 3 nodes per rack). There will be one replica of the dataset spread across each rack. See the rack-related properties in the K8ssandra v1.4.x Helm Chart [reference](https://docs-v1.k8ssandra.io/reference/helm-charts/k8ssandra/), or the single K8ssandra Operator Helm chart [reference](https://docs-v2.k8ssandra.io/reference/helm-chart/k8ssandra-operator).

### Reaper
An open source tool deployed by K8ssandra that lets you schedule and orchestrate repairs of Apache Cassandra clusters. Reaper improves the existing Cassandra `nodetool repair` process by:
* Splitting repair jobs into smaller tunable segments.
* Handling back-pressure through monitoring running repairs and pending compactions.
* Adding ability to pause or cancel repairs and track progress precisely.
For details, see the K8ssandra documentation topics covering the [Reaper component](https://docs-v2.k8ssandra.io/components/reaper/) and [repair tasks](https://docs-v2.k8ssandra.io/tasks/repair/).

### repair
In the context of Cassandra data, anti-entropy is the process of comparing the data of all replicas, and updating each replica to the newest version. Cassandra has two phases to the process: Build a Merkle tree for each replica, and then compare the Merkle trees to discover differences. K8ssandra deploys Reaper to your Kubernetes environment. See [Repair Cassandra with Reaper](https://docs-v2.k8ssandra.io/tasks/repair/).

### schemaless
A database in which there is no formal or rigid schema. The work to provide attributes to the data is performed in client apps, rather than by RDBMS-style DDL definitions at database creation time.

### secret
A way to store and manage sensitive information, such as passwords, OAuth tokens, and ssh keys. Storing confidential information in a secret is safer and more flexible than putting it verbatim in a Pod definition or in a container image. For more, see [K8ssandra security](https://docs-v2.k8ssandra.io/tasks/secure/).

### seeds
In Cassandra, a seed node is used to bootstrap the [gossip]({{< relref "#gossip" >}}) process for new nodes joining a cluster. To learn the topology of the ring, a joining node contacts one of the nodes in the `-seeds` list in `cassandra.yaml`. The first time you bring up a node in a new cluster, only one node is the seed node.  

### serverless
A cloud computing execution model in which the cloud provider allocates machine resources on demand, taking care of the servers on behalf of their customers. Serverless computing does not hold resources in volatile memory; computing is rather done in short bursts with the results persisted to storage. When an app is not in use, there are no computing resources allocated to the app. Pricing is based on the actual amount of resources consumed by an application.

### service
In Kubernetes, a service describes a set of pods that perform the same task.

### SimpleSnitch
In Cassandra, the default [snitch]({{< relref "#networktopologystrategy" >}}) type. Used only for single-datacenter deployments. It does not recognize datacenter or rack information and can be used only for single-datacenter deployments or single-zone in public clouds. It treats strategy order as proximity, which can improve cache locality when disabling read repair.

### SimpleStrategy
In Cassandra, a data replication strategy that places the first replica on a node determined by the partitioner. This strategy specifies how many replicas you want in each CassandraDatacenter. See also [NetworkTopologyStrategy]({{< relref "#networktopologystrategy" >}}).

### Site Reliability Engineer (SRE)
See [SRE]({{< relref "#sre" >}}).

### snitch
In Cassandra, the mapping from the IP addresses of nodes to physical and virtual locations, such as racks and data centers. There are several types of snitches. The type of snitch affects the request routing mechanism. See also [SimpleSnitch]({{< relref "#simplesnitch" >}}).

### SRE 
An acronym for Site Reliability Engineer. A computing professional who applies aspects of software engineering to infrastructure and operations problems. The main goal of an SRE is to create scalable and highly reliable software systems. SRE is a more recent term for a discipline that was often called operations. 

### Stargate
An open source data gateway that sits between your app and your databases. Stargate brings together an API platform and data request coordination code into one OSS project. See https://stargate.io. 

### StatefulSet
The workload API object used to manage stateful applications. Manages the deployment and scaling of a set of [pods]({{< relref "#pod" >}}), and provides guarantees about the ordering and uniqueness of these pods. Like a deployment, a StatefulSet manages pods that are based on an identical container spec. After installing K8ssandra, in the output of subsequent commands like `kubectl get pods`, notice the naming convention of using `-sts-` in the DC1 StatefulSet pod name: `demo-dc1-default-sts-0`. That important pod deployed by K8ssandra Operator includes the container that manages each Cassandra node. 

```bash
kubectl get pods
```
**Output:**
```bash
NAME                                                    READY   STATUS    RESTARTS   AGE
demo-dc1-default-stargate-deployment-547df5877d-bvnz2   1/1     Running   0          66m
demo-dc1-default-sts-0                                  2/2     Running   0          80m
demo-dc1-default-sts-1                                  2/2     Running   0          80m
demo-dc1-default-sts-2                                  2/2     Running   0          80m
k8ssandra-operator-7cfd7977cb-wxww5                     1/1     Running   0          97m
k8ssandra-operator-cass-operator-794f65d9f4-s697p       1/1     Running   0          97m
```

Given the output above, the command previously entered was, for example:

```bash
helm install demo k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace`
```

...where the `clusterName` command parameter was `demo`, and the default configured CassandraDatacenter value was `dc1`.

### table
In a database such as Cassandra, a collection of ordered (by name) columns fetched by row. A row consists of columns and have a primary key. The first part of the key is a column name. Subsequent parts of a compound key are other column names that define the order of columns in the table.

### Traefik
An HTTP reverse proxy and load balancer that makes deploying microservices easier. Traefik (pronounced "Traffic") integrates with your existing infrastructure components and configures itself automatically and dynamically. The K8ssandra GitHub code and documentation include Traefik ingress configuration examples. See the [Traefik ingress](https://docs-v2.k8ssandra.io/tasks/connect/ingress/) topic.

## Next steps

* [FAQs](https://docs.k8ssandra.io/faqs/): If you're new to the K8ssandra project, these FAQs are for you. 
* [Install](https://docs-v2.k8ssandra.io/install/local/): K8ssandra Operator install steps for local development or production-ready cloud platforms in single- or multi-cluster Kubernetes.
* [Quickstarts](https://docs-v2.k8ssandra.io/quickstarts/): Post-install K8ssandra topics for developers or Site Reliability Engineers.
* [Components](https://docs-v2.k8ssandra.io/components/): Dig in to each deployed component of the K8ssandra Operator stack and see how it communicates with the others.
* [Tasks](https://docs-v2.k8ssandra.io/tasks/): Check out the Tasks topics for a helpful collection of outcome-based solutions.
* [Reference](https://docs-v2.k8ssandra.io/reference/): Explore the Custom Resource Definitions (CRDs) reference topics for use with K8ssandra Operator deployments.

