---
title: "Configure Ingress"
linkTitle: "Configure Ingress"
weight: 7
description: |
  Routing from the outside in.
---

External connectivity can be tricky with Kubernetes. There are many solutions in
this space ranging from some which focus on pure HTTP workloads to others which
allow for custom TCP and UDP routing. Kubernetes provides an
[`Ingress`](https://kubernetes.io/docs/concepts/services-networking/ingress/)
resource type, but does **not** provide an [`Ingress
Controller`](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/).
Instead Kubernetes operators must decide on the Ingress controller they choose
to run whether it is provided by a cloud vendor, software vendor, or open-source
project. Below you will find a collection of ingress software solutions and the
appropriate integration steps to leverage them with K8ssandra.
