---
title: "Kubernetes Operator for Apache Cassandra (cass-operator)"
linkTitle: "Cass Operator"
weight: 2
description: K8ssandra deploys Kubernetes Operator for Apache Cassandra&reg; (cass-operator) to support management tasks in Kubernetes.
---

Kubernetes Operator for Apache Cassandra&reg;, also known as Cass Operator (cass-operator), is deployed by K8ssandra as part of its Helm chart install. If you haven't already, see the [quickstart]({{< relref "/quickstarts/" >}}) and [install]({{< relref "/install" >}}) topics.

## Introduction

Cass Operator automates the process of deploying and managing Cassandra in a Kubernetes cluster. Cass Operator distills the user-supplied information down to the number of nodes and cluster name to manage the lifecycle of individual Kubernetes resources. Additional options are available, but for starters, that's essentially all you'll need to specify. Now the process of managing the distributed Cassandra or DSE data platform is turnkey and much easier, which means your team is free to focus on the application layer and its functionality.

Let's start by looking at containers and the emergence of Kubernetes as the premier platform for application orchestration.

### Optimizing data management in containers with Kubernetes

Containers are a popular technology used to accelerate today's application development. Thanks to prevalent container platforms like Docker, you can package applications efficiently compared with virtual machines. With containers, apps and all of their dependencies are packaged together into a minimal deployable image. As a developer, you can use containers to move applications between environments and guarantee that your apps behave as expected. These goals led to the creation of container orchestration platforms. The leader in this space is Kubernetes.

Highlighting just a few of the advantages:

* Kubernetes accepts definitions for services and handles the assignment of containers to servers and connecting them together.
* Kubernetes dynamically tracks the health of the running containers. If a container goes down, Kubernetes handles restarting it, and can schedule its container replacement on other hardware.
* By using Kubernetes to orchestrate containers, you can rapidly build microservice-powered applications and ensure they run as designed across any Kubernetes platform.

### Cassandra managed by K8ssandra cass-operator in Kubernetes clusters

Cassandra substantially simplify development. All nodes are equal, and each node is capable of handling read and write requests with no single point of failure. Data is automatically replicated between failure zones to prevent the loss of a single container taking down your application. With simple configuration options in Cass Operator, Cassandra databases can rapidly take advantage of Kubernetes orchestration and are well suited for the container-first approach in your enterprise.

## More...

Content TBS...

## Next

See the other [components]({{< relref "/components/" >}}) deployed by K8ssandra. For information on using the deployed components, see the [Tasks]({{< relref "/tasks/" >}}) topics.
