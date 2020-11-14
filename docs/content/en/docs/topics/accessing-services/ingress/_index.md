---
title: "Ingress"
linkTitle: "Ingress"
weight: 1
date: 2020-11-07
description: |
  
---

External connectivity, meaning connectivity for applications running outside of
your Kubernetes cluster can be tricky. There are many solutions in this space
ranging from which focus on pure HTTP workloads to others which allow for custom
TCP and UDP routing. Kubernetes provides an
[`Ingress`](https://kubernetes.io/docs/concepts/services-networking/ingress/)
resource type, but does **not** provide an [`Ingress
Controller`](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/).
Instead Kubernetes operators must decide on the Ingress controller they choose
to run whether it is provided by a cloud vendor, software vendor, or open-source
project. Below you will find a collection of ingress software solutions and the
appropriate integration steps to leverage them with K8ssandra.