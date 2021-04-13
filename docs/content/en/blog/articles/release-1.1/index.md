---
date: 2021-04-12
title: "K8ssandra 1.1 Release"
linkTitle: "K8ssandra 1.1 Release"
description: >
  The K8ssandra 1.1.0 release builds on the rock-solid 1.0 foundation with improved storage and a number of organizational updates and project maturity milestones.
author: Chris Bradford ([@bradfordcp](https://twitter.com/bradfordcp)), Jeff Carpenter ([@jscarp](https://twitter.com/jscarp))
---

The K8ssandra 1.1.0 release was published on Friday, 04/09/2021. This release of K8ssandra builds on the rock-solid 1.0 foundation with improved storage and a number of organizational updates and project maturity milestones.

## Installing K8ssandra 1.1

To create a new K8ssandra deployment using the latest release, follow the instructions on the [Getting Started]({{< ref "getting-started" >}}) page. 

To update an existing K8ssandra deployment to the 1.1 release, see the instructions in the [Upgrading K8ssandra]({{< ref "upgrade" >}}) page. You can check the version of your current installation by executing the command: `helm show chart k8ssandra/k8ssandra` and searching for the line that begins with `version:`.

## MinIO Storage

As part of improvements to support S3-compatible storage, we have completed support for [MinIO](https://min.io/) as a backend storage option. K8ssandra and Medusa maintainer Alexander Dejanovski ([@alexanderDeja](https://twitter.com/alexanderDeja)) has put together an [awesome walkthrough]({{< ref "minio-backup" >}}) showing how to backup and restore with MinIO.

## Repository Organization

With this release, the K8ssandra project has also seen a number of improvements to repository organization. Ongoing development of the [datastax/cass-operator](https://github.com/datastax/cass-operator) project has been moved to [k8ssandra/cass-operator](https://github.com/k8ssandra/cass-operator). The original repository has not been removed to avoid breaking Go modules for existing projects.

Additionally, the Management API for Apache Cassandra has been migrated to the K8ssandra organization. This change also includes updates to Docker Hub. Instead of leveraging multiple Docker repositories (currently one for each major/minor/patch group) these have been merged at [k8ssandra/cass-management-api](https://github.com/k8ssandra/management-api-for-apache-cassandra). Note we provide both short and long tag formats: for example, the long form [3.11.9-v0.1.24](https://hub.docker.com/layers/k8ssandra/cass-management-api/3.11.9-v0.1.24/images/sha256-8d8241c7fa194ceb0b9b321f29ec46aa13eec5ba11c8f3ec94eb00fe9a812ad2?context=explore)and short form [3.11.9](https://hub.docker.com/layers/k8ssandra/cass-management-api/3.11.9/images/sha256-4857fb1701d46fa7a481a96133dbe38db6f79b035af89e6dbfde97247e408a73?context=explore). The latter provides the same image hash, but a shorter `latest` syntax for users looking to leverage the most recent version of the Management API.

## Testing improvements

We’ve made a number of improvements since the 1.0 release in how K8ssandra releases are tested. This includes not only [expanded unit tests](https://github.com/k8ssandra/k8ssandra/tree/main/tests/unit), but also a suite of [automated integration tests](https://github.com/k8ssandra/k8ssandra/tree/main/tests/integration). These changes will help us release more frequently, with greater confidence. 

## What’s next

As we shared recently on the [K8ssandra blog]({{< ref "roadmap-update" >}}), the roadmap is maintained in [GitHub](https://github.com/orgs/k8ssandra/projects/6), and if you take a look you’ll be able to see we’re hard at work on making it easier to run K8ssandra on Google Kubernetes Engine (GKE), Amazon Elastic Kubernetes Service (EKS), and Azure Kubernetes Service (AKS), as well as numerous other fixes and improvements.

We’d love to get your feedback on this new release and the K8ssandra project in general. As always, we encourage new [issues](https://github.com/k8ssandra/k8ssandra/issues) and [pull requests](https://github.com/k8ssandra/k8ssandra/pulls) on the main K8ssandra repo or any of the included projects. 

## Don't Miss: Kubecon EU Workshop

K8ssandra will be featured in a [workshop at Kubecon Europe](http://dtsx.io/k8ssandraATkubeconeu) on Tuesday, May 4, which will include hands-on exercises that will guide you through deploying K8ssandra, generating load against the database, and using key features from the 1.0 and 1.1 releases. Hope to see you there!
