---
title: "K8ssandra glossary"
linkTitle: "Glossary"
weight: 6
description: Definition of frequently used Kubernetes terms plus the K8ssandra or Apache Cassandra context.
---

AKS
: The Azure Kubernetes Service from Microsoft. One of the "top 3" major cloud providers supported by K8ssandra, along with EKS from Amazon, and GKE from Google. AKS offers serverless Kubernetes, an integrated continuous integration and continuous delivery (CI/CD) experience, and enterprise-grade security and governance.  

anti-entropy
: The process of comparing the data of all replicas and updating each replica to the newest version. Cassandra has two phases to the process: Build a Merkle tree for each replica. Compare the Merkle trees to discover differences. K8ssandra provides Reaper as one of its deployed components, enabling you to perform repair operations. 

Astra
: A CNDB product from DataStax that gives you the ability to develop and deploy data-driven applications with a cloud-native service, without the hassles of database and infrastructure administration. By automating tuning and configuration, Astra radically simplifies database and streaming operations.

charts
: Helm charts are a YAML-based packaging format to create, version, share, and publish software in Kubernetes. A Helm chart is a collection of templates and settings that describe a set of Kubernetes resources. For details about each Helm chart provided by K8ssandra, start on this [topic]({{< relref "/reference" >}}).

CNDB
: An acronym for Cloud Native DataBase, which refers to a database that is created and managed in a cloud environment. DataStax Astra is a CNDB, as well as Apache Cassandra that's deployed to a Kubernetes cloud provider (such as GKE, EKS, AKS) by K8ssandra. 

CQL
: Cassandra Query Language is a set of DDL and DML statements designed for communicating with Apache Cassandra databases. CQL offers a model close to SQL in the sense that data is put in tables containing rows of columns. 

CQLSH
: A command-line shell (Cassandra Query Language Shell) for interacting with Cassandra through CQL. CQLSH is included with every Cassandra package, and can be found in the `bin/` directory alongside the cassandra executable. CQLSH utilizes the Python native protocol driver, and connects to the single node specified on the command line.

container
: An image that is a ready-to-run software package with everything needed to run an application: the code and any runtime it requires, application and system libraries, and default settings.

EKS
: Amazon Elastic Kubernetes Service is one of the "top 3" major cloud providers supported by K8ssandra. EKS allows you to start, run, and scale Kubernetes applications in the AWS cloud or on-premises. 

GKE
: Google Kubernetes Engine is one of the "top 3" major cloud providers supported by K8ssandra, along with EKS from Amazon, and AKS from Microsoft. GKE includes a set of UI-based tools that are part of the Google Cloud Console for GCP environments. 

gossip
: (definition coming next)

Grafana
: (definition coming next)

helm
: A command-line tool that helps you manage Kubernetes applications. Helm Charts help you define, install, and upgrade Kubernetes application.

keyspace
: (definition coming next)

kubectl
: (definition coming next)

Kubernetes
: A portable, extensible, open-source platform for managing containerized workloads and services, that facilitates both declarative configuration and automation. It has a large, rapidly growing ecosystem. The name Kubernetes originates from Greek, meaning helmsman or pilot.

K8ssandra
: An open-source, production-ready platform for running Apache Cassandra® on Kubernetes. This includes automation for operational tasks such as repairs, backups, and monitoring.

LoadBalancingStrategy
: definition here, just checking format options.

Medusa
: (definition coming next)

minikube
: A tool that lets you run Kubernetes locally. Minikube runs a single-node Kubernetes cluster on your personal computer (including Windows, macOS and Linux PCs) so that you can try out Kubernetes, or for daily development work.

MinIO
: An Amazon S3-compatible server-side software storage stack. MinIO is one of the local or cloud-based storage objects ("buckets") supported by K8ssandra's Medusa backup/restore operations. For more, see [Backup and restore Cassandra data]({{< relref "/tasks/backup-restore/minio/" >}}).

namespace
: A way to provide a scope for names. Names of resources need to be unique within a namespace, but not across namespaces. Namespaces cannot be nested inside one another and each Kubernetes resource can only be in one namespace.

NetworkTopologyStrategy
: (definition coming next)

nodetool
: (definition coming next)

pod
: Represents a single instance of a running process in your cluster. Pods contain one or more containers, such as Docker containers. When a Pod runs multiple containers, the containers are managed as a single entity and share the Pod's resources. here, just checking format options. For Cassandra and DataStax Enterprise users, a "node" in a cluster is the equivalent of a pod.

port forwarding
: An application of network address translation that redirects a communication request from one address and port number combination to another, while the packets are traversing a network gateway, such as a router or firewall. For information about using port forwarding with K8ssandra deployments, see:
* [Set up port forwarding]({{< relref "/quickstarts/developer/#set-up-port-forwarding" >}}).  
* Site reliability engineers, see [Configure port forwarding]({{< relref "/quickstarts/sre/#port-forwarding" >}}).

Prometheus
: (definition coming next)

Reaper
: (definition coming next)

repairs
: In the context of Cassandra data, anti-entropy is the process of comparing the data of all replicas, and updating each replica to the newest version. Cassandra has two phases to the process: Build a Merkle tree for each replica, and then compare the Merkle trees to discover differences.

schemaless
: A database in which there is no formal or rigid schema. The work to provide attributes to the data is performed in client apps, rather than by RDBMS-style DDL definitions at database creation time.

secret
: A way to store and manage sensitive information, such as passwords, OAuth tokens, and ssh keys. Storing confidential information in a secret is safer and more flexible than putting it verbatim in a Pod definition or in a container image. For more, see [K8ssandra security]({{< relref "/tasks/secure" >}}).

seeds
: In Cassandra, a seed node is used to bootstrap the gossip process for new nodes joining a cluster. To learn the topology of the ring, a joining node contacts one of the nodes in the `-seeds` list in `cassandra.yaml`. The first time you bring up a node in a new cluster, only one node is the seed node. In Kubernetes, 

serverless
: A cloud computing execution model in which the cloud provider allocates machine resources on demand, taking care of the servers on behalf of their customers. Serverless computing does not hold resources in volatile memory; computing is rather done in short bursts with the results persisted to storage. When an app is not in use, there are no computing resources allocated to the app. Pricing is based on the actual amount of resources consumed by an application.

service
: In Kubernetes, a service describes a set of pods that perform the same task.

SimpleSnitch
: (definition coming next)

SimpleStrategy
: (definition coming next)

StatefulSet
: The workload API object used to manage stateful applications. Manages the deployment and scaling of a set of Pods, and provides guarantees about the ordering and uniqueness of these Pods. Like a deployment, a StatefulSet manages Pods that are based on an identical container spec.

Stargate
: An open source data gateway that sits between your app and your databases. Stargate brings together an API platform and data request coordination code into one OSS project. See https://stargate.io. 

table
: (definition coming next)

Traefik
: An HTTP reverse proxy and load balancer that makes deploying microservices easier. Traefik integrates with your existing infrastructure components and configures itself automatically and dynamically. The K8ssandra GitHub code and documentation include Traefik configuration examples. See the [Ingress]({{< relref "/tasks/connect/ingress" >}}) topics.
