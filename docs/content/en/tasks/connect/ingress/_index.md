---
title: "Configure Ingress"
linkTitle: "Ingress"
description: "Expose access to your Apache Cassandra® database and utilities for monitoring and repair using a Kubernetes ingress"
---

External connectivity can be tricky with Kubernetes. There are many solutions in this space ranging from some which focus on pure HTTP workloads to others which allow for custom TCP and UDP routing. Kubernetes provides an [`Ingress`](https://kubernetes.io/docs/concepts/services-networking/ingress/) resource type, but does **not** provide an [`Ingress Controller`](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/). Instead Kubernetes operators must decide on the Ingress controller they choose to run whether it is provided by a cloud vendor, software vendor, or open-source project. Below you will find a collection of ingress software solutions and the appropriate integration steps to leverage them with K8ssandra.

{{% alert title="Tip" color="success" %}}
As an alternative to configuring an Ingress, consider port forwarding. It's another way to provide external access to  resources that have been deployed by K8ssandra in your Kubernetes environment. Those resources could include Prometheus metrics, pre-configured Grafana dashboards, and the Reaper web interface for repairs of Cassandra&reg; data. The `kubectl port-forward` command does not require an Ingress/Traefik to work. 

* Developers, see [Set up port forwarding]({{< relref "/quickstarts/developer/#set-up-port-forwarding" >}}).  
* Site reliability engineers, see [Configure port forwarding]({{< relref "/quickstarts/sre/#port-forwarding" >}}).
{{% /alert %}}
