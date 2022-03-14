---
title: "Configure Traefik ingress"
linkTitle: "Traefik ingress"
toc_hide: true
simple_list: true
description: "Provide access to your Apache Cassandra® database and utilities using a Kubernetes ingress."
---

{{< tbs >}}

External connectivity can be tricky with Kubernetes. There are many solutions in this space, ranging from some that focus on pure HTTP workloads, to others that allow for custom TCP and UDP routing. 

Kubernetes provides an [`Ingress`](https://kubernetes.io/docs/concepts/services-networking/ingress/) resource type, but does **not** provide an [`Ingress Controller`](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/). Instead Kubernetes operators must decide on the Ingress controller they choose to run, whether it is provided by a cloud vendor, a software vendor, or an open-source project. 

Transport Layer Security (TLS) to Traefik Ingress integrates with your existing infrastructure components and configures itself automatically and dynamically. Traefik handles advanced Ingress deployments including Mutual TLS (mTLS) of TCP with Server Name Indication (SNI) and User Datagram Protocol (UDP). Operators define rules for routing traffic to downstream systems through Kubernetes `Ingress` objects or more specific `Custom Resource Definitions`. K8ssandra supports deploying `IngressRoute` objects as part of a deployment to expose metrics, repair, and Cassandra interfaces. See the topics listed below.

{{% alert title="Note" color="warning" %}}
The provided Traefik ingress solutions are not recommended for production environments. As an alternative, consider port forwarding. It's another way to provide external access to resources that have been deployed by K8ssandra in your Kubernetes environment. The `kubectl port-forward` command does not require an Ingress/Traefik to work. See:
* Developers, see [Set up port forwarding]({{< relref "/quickstarts/developer/#set-up-port-forwarding" >}}).  
* Site reliability engineers, see [Configure port forwarding]({{< relref "/quickstarts/site-reliability-engineer/#port-forwarding" >}}).
{{% /alert %}}

Traefik ingress topics:
