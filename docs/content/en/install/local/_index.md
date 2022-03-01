---
title: "Install prerequisites"
linkTitle: "Install prerequisites"
no_list: true
weight: 1
description: "Install K8ssandraCluster custom resource for Apache Cassandra&reg; in local Kubernetes."
---

This topic identifies the K8ssandra Operator install prerequisites. It then directs you to related topics that describe single- or multi-cluster installs and configuration options for the `K8ssandraCluster` custom resource and Cassandra deployments. 

## Prerequisites

Make sure you have the following installed before going through the related install topics. 

* [kind](#kind)
* [kubectx](#kubectx)
* [yq (YAML processor)](#yq)
* [gnu-getopt](#gnu)
* [kubectl](https://kubernetes.io/docs/tasks/tools/)
 and [helm v3+](https://helm.sh/docs/intro/install/) on your preferred OS. 

### **kind**

The local install examples use [kind](https://kind.sigs.k8s.io/) clusters. If you have not already, install kind.

By default, kind clusters run on the same Docker network, which means we will have routable pod IPs across clusters.

### **kubectx**

[kubectx](https://github.com/ahmetb/kubectx) is a really handy tool when you are dealing with multiple clusters.  

### **yq**

[yq](https://github.com/mikefarah/yq#install) is lightweight and portable command-line YAML processor.

### **gnu-getopt**

[gnu-getopt](https://formulae.brew.sh/formula/gnu-getopt) is a command-line option parsing utility. 

To make sure that the command line is using the intended version on your local machine, add in your shell profile:

```bash
export PATH="/usr/local/opt/gnu-getopt/bin:$PATH"
```

In our testing on Linux, we used `gnu-getopt` version 2.37.3. The default downloaded version of `gnu-getopt` on macOS might cause issues.

### setup-kind-multicluster.sh (FYI) 

Note that the `make NUM_CLUSTERS=<number> create-kind-multicluster` command, which is shown in subsequent install topics, is a reference to a `Makefile` target within the k8ssandra-operator repo. The `Makefile` is [here](https://github.com/k8ssandra/k8ssandra-operator/blob/main/Makefile). The command invokes the [setup-kind-multicluster.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/setup-kind-multicluster.sh) script. It's used extensively during development and testing. Not only does it configure and create kind clusters, it also generates `kubeconfig` files for each cluster.

**Tip:** kind generates a `kubeconfig` with the IP address of the API server set to `localhost` because the cluster is intended for local development. We need a `kubeconfig` with the IP address set to the internal address of the API server. The `setup-kind-mulitcluster.sh` script takes care of this requirement for you.  

### create-clientconfig.sh (FYI) 

[create-clientconfig.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/create-clientconfig.sh) is in the k8ssandra-operator repo. This script is used to configure access to remote clusters, as described in subsequent topics. 

## Next steps

After confirming you have the prerequisite software, proceed to the detailed steps for single- or multi-clusters, using your preferred tools:

### Quickstarts with helm

* [Quickstart **single-cluster** install]({{< relref "single-cluster-helm/" >}}) of K8ssandra Operator with `helm`.

* [Quickstart **multi-cluster** install]({{< relref "multi-cluster-helm/" >}}) of K8ssandra Operator with `helm`. 

### Quickstarts with Kustomize

* [Quickstart **single-cluster** install]({{< relref "single-cluster-kustomize/" >}}) of K8ssandra Operator with `kustomize`.

* [Quickstart **multi-cluster** install]({{< relref "multi-cluster-kustomize/" >}}) of K8ssandra Operator with `kustomize`.









