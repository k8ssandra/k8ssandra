---
title: "K8ssandra operator"
linkTitle: "K8ssandra operator"
weight: 1
description: Introducing the K8ssandra operator.
---


K8ssandra operator can be installed to a local Kubernetes (such as kind, K3D, minikube) or to a cloud provider's Kubernetes platform: Amazon AKS, Digital Ocean, Google Cloud GKE, Microsoft Azure. 

## Architecture

K8ssandra Operator consists of two primary components:

* A control plane
* A data plane

 The control plane creates and manages objects that exist only in the API server. The control plane does not deploy or manage pods.

{{% alert title="Note" color="success" %}}
The control plane can be installed in only one cluster; that is, in the control plane cluster.
{{% /alert %}}

The data plane can be installed on any number of clusters. The control plane cluster can also function as the data plane.

The data plane deploys and manages pods. Moreover, the data plane may interact directly with the managed applications. For example, the operator may call the `management-api` to create keyspaces in Cassandra.

In each cluster, the deployed and managed pods can include [Stargate]({{< relref "/components/stargate/" >}}) and [cass-operator]({{< relref "/components/cass-operator/" >}}).   

Here's how the K8ssandra operator components fit together:

![How the K8ssandra operator fit together](k8ssandra-operator-architecture.png)

## Multi-cluster requirements

There must exist routable pod IPs between Kubernetes clusters; however this requirement may be relaxed in the future.

If you are running in a cloud provider, you can get routable IPs by installing the Kubernetes clusters in the same VPC.

If you run multiple kind clusters locally, you will have routable pod IPs, assuming that they run on the same Docker network, which is normally the case. The K8ssandra project leverage this setup for our multi-cluster end-to-end tests.

## Next steps

* Learn how to install K8ssandra operator in multiple clusters.
* Also see the topics covering other [components]({{< relref "/components/" >}}) deployed by K8ssandra. 
* For information on using the deployed K8ssandra components, see the [Tasks]({{< relref "/tasks/" >}}) topics.
