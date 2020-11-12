---
title: "External Cassandra Connectivity"
linkTitle: "External Cassandra Connectivity"
weight: 1
date: 2020-11-07
description: K8ssandra provides Kong Ingress and Traefik Controller Ingress for external connectivity
---

## Overview

When applications run within a Kubernetes cluster, you need a way to access those services from outside the cluster. This topic describes two solutions that are built into K8ssandra -- Kong Ingress and Traefik Ingress Controller - along with the motivation for each. The following approaches assume that the Cassandra cluster is already up and reported as running.

For background information, see the Ingress samples in this [GitHub repo](https://github.com/datastax/cass-operator/tree/master/docs/ingress). However, note that the `helm install` step with `k8ssandra-tools` already set up these Ingress services and configurations for your Kubernetes environment. 

## What is Ingress?

[Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) is a feature that forwards requests to services running within a Kubernetes cluster based on rules. These rules may include specifying the protocol, port, or even path. They may provide additional functionality like termination of SSL / TLS traffic, load balancing across a number of protocols, and name-based virtual hosting.

Behind the Ingress Kubernetes type is an Ingress Controller. There are a number of controllers available with varying features to service the defined ingress rules. Think of Ingress as an interface for routing and an Ingress Controller as the implementation of that interface. In this way, any number of Ingress Controllers may be used based on the workload requirements.

Ingress Controllers function at Layer 4 & 7 of the OSI model. When the Ingress specification was created, it focused specifically on HTTP / HTTPS workloads. From the documentation: "An Ingress does not expose arbitrary ports or protocols. Exposing services other than HTTP and HTTPS to the internet typically uses a service of type *service-name*`.Type=NodePort` or *service-name*`.Type=LoadBalancer`."

Cassandra workloads don't use HTTP as a protocol, but rather a specific TCP protocol. Ingress Controllers that we want to leverage require support for TCP load balancing. This approach provides routing semantics similar to those of LoadBalancer Services.

If the Ingress Controller also supports SSL termination with Server Name Indication (SNI), then secure access is possible from outside the cluster while keeping Token Aware routing support. Additionally, you should consider whether the chosen Ingress Controller supports client SSL certificates allowing for Mutual TLS to restrict access from unauthorized clients.

## Kong as an Ingress

[Kong](https://konghq.com/kong/) is open-source API gateway. Built for multi-cloud and hybrid, Kong is optimized for microservices and distributed architectures. The k8ssandra-tools install provided Kong as an Ingress for a Kubernetes cluster.

* Kong Simple Load Balancing - When leveraging a single endpoint Ingress / Load Balancer, you lose the ability to provide token-aware routing without the use of SNI. (See the [mTLS with SNI](https://github.com/datastax/cass-operator/blob/master/docs/ingress/kong/mtls-sni) guide). 

    **WARNING:** This approach does not interact with the traffic at all. All traffic is sent over cleartext without any form of authentication of the server or client. Note that each Cassandra cluster running behind the Ingress will require its own endpoint / port. Without a way to detect the pod that we want to connect with, it is the only way to differentiate requests.

* SNI Ingress using Kong - SNI provides hints to the Ingress (via TLS extensions) where the traffic should be routed from the proxy. In this case, the Ingress uses  the `hostId` as the endpoint.

* mTLS with SNI Ingress using Kong - With mTLS, not only does the client authenticate the server, but the server **also** authenticates the client. This allows for bi-directional authentication and prevents a bad actor from connecting to your cluster without the appropriate certificate.

    **Note:** mTLS is only available with Kong Enterprise.

## Traefik as an Ingress Controller

[Traefik](https://traefik.io/traefik-enterprise/) is an edge router into your Kubernetes services. As noted in the Traefik documentation, it is like a front door to your platform. Traefik intercepts and routes every incoming request, using all the logic and every rule that determine which services handle which requests based on the path, the host, headers, and soon. 

The advantage of using Traefik as an Ingress Controller: when a new service is deployed, Traefik detects it immediately and updates the routing rules. And, when you  remove a service from your Kubernetes infrastructure, the route is removed as well. This means you do not need to create and synchronize configuration files that are cluttered with IP addresses or other routing rules.

When you installed k8ssandra-tools, Traefix was provided as an Ingress Controller.  K8ssandra services automatically take advantage of the routing provided by Traefik. 












