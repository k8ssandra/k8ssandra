---
title: "Traefik"
linkTitle: "Traefik"
weight: 1
description: |
  Traefik is a modern reverse proxy and load balancer that makes deploying microservices easy.
---

It integrates with your existing infrastructure components and configures itself automatically and dynamically. Traefik handles advanced ingress deployments including mTLS of TCP with SNI and UDP. Operators define rules for routing traffic to downstream systems through Kubernetes `Ingress` objects or more specific `Custom Resource Definitions`. K8ssandra supports deploying `IngressRoute` objects as part of a deployment to expose metrics, repair, and Apache Cassandra® interfaces.

{{% alert title="Warning" color="warning" %}}
The provided Traefik examples, such as the ones in the [Minikube deployment]({{< ref "/docs/topics/ingress/traefik/minikube-deployment/" >}}) topic, are not recommended for production deployments. 
{{% /alert %}}
