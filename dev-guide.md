# K8ssandra Developer Guide

Ready to write some code?  Jump straight to the one of our other guides:

* [K8ssandra Chart Development Quick Start](dev-quick-start.md)
* [Contributing to Cass Operator](https://github.com/k8ssandra/cass-operator/blob/master/README.md#contributing)

## Overview

K8ssandra is a series of projects orchestrated to deliver a cloud-native Apache Cassandra deployment on Kubernetes.  To learn more about the project in general jump over to the project's main [README](README.md).

This guide is aimed at introducing developers to the various technologies and sub-projects involved in K8ssanrdra as well as some of the processes that the project uses to plan and prioritize our work.

## Learning Resources

To develop within the K8ssandra ecosystem there are some high-level topics that you'll want to gain experience with:

* Kubernetes
* Kind
* Helm
* Kubernetes Operator Development
* Go
* Apache Cassandra
* Medusa for Apache Cassandra
* Reaper for Apache Cassandra
* Git/GitHub

There are a lot of great learning resources out there related to all of these topics.  Here is a small collection of resources that our team has found useful over time.

(If you have any you'd like to add here, we'd love to have the contribution!)

* [Kubernetes Documentation](https://kubernetes.io/docs/home/)
* [Kind Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/)
* [Helm Documentation](https://helm.sh/docs/)
* [Kubebuilder Book](https://book.kubebuilder.io/)
* [Building Operators with Go](https://sdk.operatorframework.io/docs/building-operators/golang/)
* [Tutorial: Get Started with Go](https://golang.org/doc/tutorial/getting-started)
* [Apache Cassandra Operations in Kubernetes](https://www.datastax.com/learn/apache-cassandra-operations-in-kubernetes)
* [DataStax Developers YouTube Channel](https://www.youtube.com/c/DataStaxDevs/videos)
* [Medusa for Apache Cassandra Documentation](https://github.com/thelastpickle/cassandra-medusa#documentation)
* [Reaper for Apache Cassandra Documentation](http://cassandra-reaper.io/docs/)

## Running K8ssandra

The best way to get into developing parts of K8ssandra is by running it and learning about how it works.  We have a ton of great resources for just that in the project [Get Started Guide](https://k8ssandra.io/get-started/) and [Documentation](https://docs.k8ssandra.io/).

One of the best ways to run K8ssandra locally is using [kind](https://kind.sigs.k8s.io/).

If you're using Docker Desktop it's important to properly tune the resources allocated to Docker.  Check out this [blog](https://k8ssandra.io/blog/articles/requirements-for-running-k8ssandra-for-development/) for some background on configuring a development type environment.

## Development Tools

Depending on the project in which you're working a different set of tools might be used.

### Source Control

[Git](https://git-scm.com/)

[GitHub CLI](https://cli.github.com/)

### Runtime Envrionment

[Docker Desktop](https://www.docker.com/products/docker-desktop)

[Kind](https://kind.sigs.k8s.io/)

[kubectl](https://kubernetes.io/docs/tasks/tools/)

### Language Support

[Go](https://golang.org/doc/install)

[Python](https://www.python.org/downloads/)

[Visual Studio Code](https://code.visualstudio.com/)

[GoLand](https://www.jetbrains.com/go/)

## Project Planning & Management

Within the core K8ssandra team we use a combination of tools to plan and track our roadmap and development process.

The project roadmap is maintained in a GitHub project which can be found [here](https://github.com/orgs/k8ssandra/projects/6).  This roadmap will give a high-level idea of where the project is planning to go in the coming months.

That higher-level roadmap ultimately feeds into the shorter term planning of the core K8ssandra engineering team.

The direct work to be done is captured in GitHub issues - those issues get spread amongst the various projects under the K8ssandra umbrella, described below.

We work on a 2-week sprint cadence - going through the typical ceremonies you'd likely expect from a scrum project: grooming, planning, review, retrospective, etc.

The core development team uses Jira to organize and manage work.  Jira issues are synced automatically with GitHub issues, so that both systems contain the same content and discussion on all issues.

The Jira project used is also publicly available (although you will need an Atlassian account) and can be found [here](https://k8ssandra.atlassian.net/secure/RapidBoard.jspa?rapidView=5&projectKey=K8SSAND).

## Projects & Code Repositories

There are a number of projects involved in the K8ssandra ecosystem.

The majority of the projects developed and contributed to through K8ssandra can be found within the [k8ssandra GitHub Organization](https://github.com/k8ssandra).

### k8ssandra/k8ssandra

`k8ssandra/k8ssandra` represents the high-level umbrella project for the ecosystem.  This is primarily a packaging of Helm charts and testing capabilities.  This is the project from where K8ssandra is assembled and deployed.

[GitHub Repository](https://github.com/k8ssandra/k8ssandra)

### k8ssandra/k8ssandra-operator

`k8ssandra/k8ssandra-operator` represents the next generation of the K8ssandra implementation.  This project implements a Kubernetes operator that is responsible for managing the full deployment of K8ssandra across multiple clusters.

[GitHub Repository](https://github.com/k8ssandra/k8ssandra-operator)

### k8ssandra/cass-operator

`k8ssandra/cass-operator` is the Kubernetes operator responsible for managing the deployment of Apache Cassandra within a K8ssandra cluster.  This project was originally developed under `datastax/cass-operator` and migrated to the K8ssandra organization.

[Github Repository](https://github.com/k8ssandra/cass-operator)

[Docker Hub Repository](https://hub.docker.com/r/k8ssandra/cass-operator)

### k8ssandra/management-api-for-apache-cassandra

`k8ssandra/management-api-for-apache-cassandra` is a sidecar service layer that attempts to build a well supported set of operational actions on Apache Cassandra nodes that can be administered centrally.  This is the layer through which the Apache Cassandra cluster nodes are managed.

[GitHub Repository](https://github.com/k8ssandra/management-api-for-apache-cassandra)

[Docker Hub Repository](https://hub.docker.com/r/k8ssandra/cass-management-api)

### k8ssandra/medusa-operator

`k8ssandra/medusa-operator` is the Kubernetes operator responsible for managing backup and restore capabilities for Apache Cassandra using Medusa.

[GitHub Repository](https://github.com/k8ssandra/medusa-operator)

[Docker Hub Repository](https://hub.docker.com/r/k8ssandra/medusa-operator)

### k8ssandra/reaper-operator

`k8ssandra/reaper-operator` is the Kubernetes operator responsible for managing repair capabilities for Apache Cassandra using Reaper.

[GitHub Repository](https://github.com/k8ssandra/reaper-operator)

[Docker Hub Repository](https://hub.docker.com/r/k8ssandra/reaper-operator)

### thelastpickle/medusa

`thelastpickle/medusa` is the tool that K8ssandra has chosen to manage backup and restore capabilities for the Apache Cassandra within the stack.  Medusa itself has a vibrant project community that the K8ssandra team regularly contributes to - both to help improve Medusa in general and also to provide features within K8ssandra.

[GitHub Respository](https://github.com/thelastpickle/cassandra-medusa)

### thelastpickle/reaper

`thelastpickle/reaper` is the tool that K8ssandra has chosen to manage repair capabilities for Apache Cassandra within the stack.  Like Medusa, Reaper itself has a large project community and has a large basis of usage outside of K8ssandra.  The K8ssandra team also regularly contributes to Reaper.

[Documentation Site](http://cassandra-reaper.io/)

[GitHub Repository](https://github.com/thelastpickle/cassandra-reaper)

## K8ssandra Chart Development Quick Start

Ready to get started working on the helm charts that support the overall K8ssandra deployment?  Head over to the [K8ssandra Chart Development Quick Start](dev-quick-start.md) to learn more about setting up to develop, test, and contribute to the K8ssandra charts.

## Contributing to Cass Operator

Looking to contribute to Cass Operator?  Check out the [Contributing](https://github.com/k8ssandra/cass-operator/blob/master/README.md#contributing) guide to learn a bit more about how to develop and test Cass Operator.
