---
date: 2021-03-23
title: "The search for a cloud-native database"
linkTitle: "The search for a cloud-native database"
description: >
  Cloud-native applications need cloud-native databases to achieve the next stage of maturity and scale, but what is a cloud-native database? Let’s try to define what it means and figure out how it relates to Kubernetes.
author: Cedrick Lunven ([@clunven](https://twitter.com/clunven)), Jeff Carpenter ([@jscarp](https://twitter.com/jscarp))

---

The concept of “cloud-native” has come to stand for a collection of best practices for application logic and infrastructure, including databases. However, many of the databases supporting our applications have been around for decades, before the cloud or cloud-native was a thing. The data gravity associated with these legacy solutions has limited our ability to move applications and workloads. As we move to the cloud, how do we evolve our data storage approach? Do we need a cloud-native database? What would it even mean for a database to be cloud-native? Let’s take a look at these questions.

## What is Cloud-Native?

It’s helpful to start by defining terms. In unpacking  “cloud-native”, let’s start with the word “native”. For individuals, the word may evoke thoughts of your first language, or your country or origin - things that feel natural to you. Or in nature itself, we might consider the native habitats inhabited by wildlife, and how each species is adapted to its environment. We can use this as a basis to understand the meaning of cloud-native.

Here’s how the Cloud Native Computing Foundation (CNCF) [defines the term](https://github.com/cncf/toc/blob/main/DEFINITION.md):

> “Cloud native technologies empower organizations to build and run scalable applications in modern, dynamic environments such as public, private, and hybrid clouds: Containers, service meshes, microservices, immutable infrastructure, and declarative APIs exemplify this approach.
>
> These techniques enable loosely coupled systems that are resilient, manageable, and observable. Combined with robust automation, they allow engineers to make high-impact changes frequently and predictably with minimal toil.”

This is a rich definition, but it can be a challenge to use this to define what a cloud-native database is, as evidenced by the **Database** section of the CNCF Landscape Map:. 

{{< imgproc cncf-landscape-db Resize "800x" >}}
Databases are just a small portion of a crowded cloud computing landscape
{{< /imgproc >}}

Look closely, and you’ll notice a wide range of offerings: both traditional relational databases and NoSQL databases, supporting a variety of different data models including key/value, document, and graph. You’ll also find technologies that layer clustering, querying or schema management capabilities on top of existing databases. And this doesn’t even consider related categories in the CNCF landscape such as **Streaming and Messaging** for data movement, or **Cloud Native Storage** for persistence. 

Which of these databases are cloud-native? Only those that are designed for the cloud, should we include those that can be adapted to work in the cloud? Bill Wilder provides an interesting perspective in his 2012 book, “Cloud Architecture Patterns”, defining “cloud-native” as:

> “Any application that was architected to take full advantage of cloud platforms”

By this definition, cloud-native databases are those that have been architected to take full advantage of underlying cloud infrastructure. Obvious? Maybe. Contentious? Probably...

## Why should I care if my database is cloud-native?

Or to ask a different way, what are the advantages of a cloud-native database? Consider the two main factors driving the popularity of the cloud: cost and time-to-market.

*   **Cost** - the ability to pay-as-you-go has been vital in increasing cloud adoption. (But that doesn’t mean that cloud is cheap or that cost management is always straightforward.)
*   **Time-to-market** - the ability to quickly spin up infrastructure to prototype, develop, test, and deliver new applications and features. (But that doesn’t mean that cloud development and operations are easy.)

These goals apply to your database selection, just as they do to any other part of your stack.

## What are the characteristics of a cloud-native database?

Now we can revisit the CNCF definition and extract characteristics of a cloud-native database that will help achieve our cost and time-to-market goals:

*   **Scalability** - the system must be able to add capacity dynamically to absorb additional workload
*   **Elasticity** - it must also be able to scale back down, so that you only pay for the resources you need 
*   **Resiliency** - the system must survive failures without losing your data
*   **Observability** -  tracking your activity, but also health checking and handling failovers
*   **Automation** - implementing operations tasks as repeatable logic to reduce the possibility of error. This characteristic is the most difficult to achieve, but is essential to achieve a high delivery tempo at scale

Cloud-native databases are designed to embody these characteristics, which distinguish them from “cloud-ready” databases, that is, those that can be deployed to the cloud with some adaptation. 

## What’s a good example of a cloud-native database?

Let’s test this definition of a cloud-native database by applying it to Apache Cassandra™ as an example. While the term “cloud-native” was not yet widespread when Cassandra was developed, it bears many of the same architectural influences, since it was inspired by public cloud infrastructure such as Amazon’s Dynamo Paper and Google’s BigTable. Because of this lineage, Cassandra embodies the principles outlined above:

*   Cassandra demonstrates horizontal **scalability** through adding nodes, and can be scaled down **elastically** to free resources outside of peak load periods  
*   By default, Cassandra is an AP system, that is, it prioritizes availability and partition tolerance over consistency, as described in the CAP theorem. Cassandra’s built in replication, shared-nothing architecture and self-healing features help guarantee **resiliency**.
*   Cassandra nodes expose logging, metrics, and query tracing, which enable **observability**
*   **Automation** is the most challenging aspect for Cassandra, as typical for databases. 

While automating the initial deployment of a Cassandra cluster is a relatively simple task, other tasks such as scaling up and down or upgrading can be time-consuming and difficult to automate. After all, even single-node database operations can be challenging, as many a DBA can testify. Fortunately, the K8ssandra project provides best practices for deploying Cassandra on Kubernetes, including major strides forward in automating “day 2” operations. 

## Does a cloud-native database have to run on Kubernetes?

Speaking of Kubernetes… When we talk about databases in the cloud, we’re really talking about stateful workloads requiring some kind of storage. But in the cloud world, stateful is painful. Data gravity is a real challenge - data may be hard to move due to regulations and laws, and the cost can get quite expensive. This results in a premium on keeping applications close to their data. 

The challenges only increase when we begin deploying containerized applications using Kubernetes, since it was not originally designed for stateful workloads. There’s an emerging push toward deploying databases to run on Kubernetes as well, in order to maximize development and operational efficiencies by running the entire stack on a single platform. What additional requirements does Kubernetes put on a cloud-native database?

### Containerization

First, the database must run in containers. This may sound obvious, but some work is required. Storage must be externalized, the memory and other computing resources must be tuned appropriately, and the application logs and metrics must be made available to infrastructure for monitoring and log aggregation.

### Storage

Next, we need to map the database’s storage needs onto Kubernetes constructs. At a minimum, each database node will make a persistent volume claim that Kubernetes can use to allocate a storage volume with appropriate capacity and I/O characteristics. Databases are typically deployed using Kubernetes Stateful Sets, which help manage the mapping of storage volumes to pods and maintain consistent, predictable, identity.

### Automated Operations

Finally, we need tooling to manage and automate database operations, including installation and maintenance. This is typically implemented via the Kubernetes operator pattern. Operators are basically control loops that observe the state of Kubernetes resources and take actions to help achieve a desired state. In this way they are similar to Kubernetes built-in controllers, but with the key difference that they understand domain-specific state and thus help Kubernetes make better decisions. 

For example, the K8ssandra project uses [cass-operator](https://github.com/datastax/cass-operator), which defines a Kubernetes custom resource (CRD) called “CassandraDatacenter” to describe the desired state of each top-level failure domain of a Cassandra cluster. This provides a level of abstraction higher than dealing with Stateful Sets or individual pods.

Kubernetes database operators typically help to answer questions like:

*   What happens during failovers? (pods, disks, networks)
*   What happens when you scale out? (pod rescheduling)
*   How are backups performed?
*   How do we effectively detect and prevent failure?
*   How is software upgraded? (rolling restarts)

## Conclusion and what’s next

A cloud-native database is one that is designed with cloud-native principles in mind, including scalability, elasticity, resiliency, observability, and automation. As we’ve seen with Cassandra, automation is often the final milestone to be achieved, but running databases in Kubernetes can actually help us progress toward this goal of automation. 

What’s next in the maturation of cloud-native databases? We’d love to hear your input as we continue to invent the future of this technology together.    

_This blog is based on [Cedrick’s presentation](https://pheedloop.com/blueprintldn/virtual/?page=sessions&section=SES7OFUNTNTXCVKQH) “Databases in the Cloud-Native Era” from BluePrint London, March 11, 2021 (registration required)._

