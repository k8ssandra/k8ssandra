---
date: 2021-02-25
title: "K8ssandra 1.0 Stable Release and What's Next"
linkTitle: "K8ssandra 1.0 Release"
description: >
  The K8ssandra 1.0 release is a significant milestone for the project, reflecting months of work and contributions from the community. This release represents a production-ready version of the project and sets the stage for continued development.
author: Chris Bradford ([@bradfordcp](https://twitter.com/bradfordcp))

---

A project's 1.0 release is a significant milestone. In some cases it is the solidification of an API, a reflection of its suitability for production workloads, or occasionally because the developers like round numbers. For K8ssandra, the 1.0 release represents a version that is ready for widespread adoption. Configuration interfaces have been solidified, defaults defined and tested, and documentation developed to guide the various user personas through deployment and operation. Beyond the code bits and Helm charts, the project has grown significantly. We've seen an evolution of how features are developed, validated, and deployed. Automation has risen to drive some of these processes and ease the development workflow.

The K8ssandra team has regularly shipped updates to the Helm repository, each one reflecting new functionality, bug fixes, and enhancements. Part of a 1.0 release is a confidence in stability and measured release of major, minor, and patch level changes. K8ssandra 1.0 also brings online a pair of release streams: Stable and Next. The Stable stream contains releases which have undergone additional testing and validation prior to being published. Alternatively, users looking for the latest in-flight functionality may leverage the Next stream. This new stream exposes chart updates as they are merged into the main branch. Whether you are looking to take the latest feature for a spin, or validate a bugfix before wider release, it's available without having to clone the repo and build from source.

Since the initial release at Kubecon NA 2020, there have been plenty of incremental version updates. We were at v0.55.0 at the time I wrote this post! Let's take a look at how much ground has been covered over that time.

## Stats
* 259 commits on the main branch
* 380 Issues opened
* 267 Issues closed
* 308 files changed
* 27 contributors
* 2 workshops

## Features
* Installation with a single chart
* Dedicated charts for sub-projects
* Stargate integration
* Authentication enabled by default
* Non-root containers
* Migration to kube-prometheus-stack
* Comprehensive Getting Started and Tasks documentation
* Integration with native K8s Ingress resources for HTTP workloads
* Additional C* configuration points (topology, heap, storage class)
* Namespace-scoped resources

## Components

### K8ssandra Projects
* K8ssandra chart - v1.0.0
* Cass Operator - v1.5.2
* Medusa Operator - v0.1.0
* Reaper Operator - v0.1.1

### External Projects
* Cassandra - 3.11.6, 3.11.7, 3.11.8, 3.11.9, 3.11.10, 4.0-beta4
* Stargate - v1.0.9
* Medusa - v0.9.0
* Reaper - v2.2.1
* kube-prometheus-stack chart - v12.11.3
* Prometheus - v0.44.0
* Grafana - v6.1

## Contributors
A huge THANK YOU to everyone who contributed to the project and made 1.0 happen. Whether it's a typo correction, YAML indentation tweak, or operator go code, every bit shared makes this project better. A huge thank you to the core team: 

* [@adejanovski](https://github.com/adejanovski)
* [@burmanm](https://github.com/burmanm)
* [@jakerobb](https://github.com/jakerobb)
* [@jdonenine](https://github.com/jdonenine)
* [@jeffbanks](https://github.com/jeffbanks)
* [@jeffreyscarpenter](https://github.com/jeffreyscarpenter)
* [@johnsmartco](https://github.com/johnsmartco)
* [@johnwfrancis](https://github.com/johnwfrancis)
* [@jsanda](https://github.com/jsanda)

And new contributors:
* [@florissmit10](https://github.com/florissmit10)
* [@idleyoungman](https://github.com/idleyoungman)
* [@michaelsembwever](https://github.com/michaelsembwever)
* [@Miles-Garnsey](https://github.com/Miles-Garnsey)
* [@mproch](https://github.com/mproch)
* [@parham-pythian](https://github.com/parham-pythian)
* [@stanislawwos](https://github.com/stanislawwos)
* [@tlasica](https://github.com/tlasica)

Check out ways to join the [community](/community) as we move forward towards the next release. Looking for a guided experience with K8ssandra? Join us at the [workshop](https://www.datastax.com/workshops/142078180663) next week as we walk through installation and exploration of this cloud-native data platform.
