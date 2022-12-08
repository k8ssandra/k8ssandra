---
title: "Single-cluster install with kustomize"
linkTitle: "Single-cluster/kustomize"
no_list: false
weight: 3
description: "Quickstart with Kustomize to install K8ssandraCluster in single-cluster Kubernetes."
---

This topic shows how you can use Kustomize to declaratively install and configure the `K8ssandraCluster` custom resource in **single-cluster** local Kubernetes. 

## Prerequisites

If you haven't already, see the install [prerequisites]({{< relref "install/local/" >}}).

## Introduction

You can install K8ssandra Operator with [Kustomize](https://kustomize.io/), which takes 
a declarative approach to configuring and deploying resources, whereas Helm takes more of 
an imperative approach.

Kustomize is integrated directly into `kubectl`. For example, `kubectl apply -k` essentially runs `kustomize build` over the specified directory followed by `kubectl apply`. See this [topic](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/) for details on the integration of Kustomize into `kubectl`.

K8ssandra Operator uses some features of Kustomize that are only available in `kubectl` v1.23 or later. For this reason the following examples use `kustomize build <dir> | kubectl apply -f -`.


## Single-cluster local Kubernetes
Let's look at a single-cluster install to demonstrate that while K8ssandra 
Operator is designed for multi-cluster use, it can be used in a single cluster without 
any extra configuration.

### Create kind cluster
Run `setup-kind-multicluster.sh` as follows:

```sh
./setup-kind-multicluster.sh --kind-worker-nodes 4
```

### Install Cert Manager
We need to first install Cert Manager because it is a dependency of cass-operator:

```console
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```

### Install K8ssandra Operator
Install with (replace `X.X.X` with [the release](https://github.com/k8ssandra/k8ssandra-operator/releases) you which to install):

```console
kustomize build "github.com/k8ssandra/k8ssandra-operator/config/deployments/control-plane?ref=vX.X.X" | kubectl apply --server-side -f -

Verify that the following CRDs are installed:

* `cassandrabackups.medusa.k8ssandra.io`
* `cassandrarestores.medusa.k8ssandra.io`
* `certificaterequests.cert-manager.io`
* `certificates.cert-manager.io`
* `challenges.acme.cert-manager.io`
* `clientconfigs.config.k8ssandra.io`
* `clusterissuers.cert-manager.io`
* `issuers.cert-manager.io`
* `k8ssandraclusters.k8ssandra.io`
* `orders.acme.cert-manager.io`
* `reapers.reaper.k8ssandra.io`
* `replicatedsecrets.replication.k8ssandra.io`
* `stargates.stargate.k8ssandra.io`


Check that there are two Deployments. The output should look similar to this:

```console
kubectl -n k8ssandra-operator get deployment
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
cass-operator-controller-manager   1/1     1            1           77s
k8ssandra-operator                 1/1     1            1           77s
```

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `true`:

```console
kubectl -n k8ssandra-operator get deployment k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

### Deploy a K8ssandraCluster
Now we will deploy a K8ssandraCluster that consists of a 3-node Cassandra cluster and a 
Stargate node.

```sh
cat <<EOF | kubectl -n k8ssandra-operator apply -f -
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    serverVersion: "4.0.3"
    datacenters:
      - metadata:
          name: dc1
        size: 3
        storageConfig:
          cassandraDataVolumeClaimSpec:
            storageClassName: standard
            accessModes:
              - ReadWriteOnce
            resources:
              requests:
                storage: 5Gi
        config:
          jvmOptions:
            heapSize: 512M
        stargate:
          size: 1
          heapSize: 256M
EOF
```

Confirm that the resource has been created:

```console
kubectl -n k8ssandra-operator get k8ssandraclusters
```

```console
NAME   AGE
demo   45s
```

```console
kubectl -n k8ssandra-operator describe k8ssandracluster demo
```

```console
Name:         demo
Namespace:    k8ssandra-operator
Labels:       <none>
Annotations:  <none>
API Version:  k8ssandra.io/v1alpha1
Kind:         K8ssandraCluster
...
Status:
  Datacenters:
    dc1:
      Cassandra:
        Cassandra Operator Progress:  Updating
        Node Statuses:
Events:  <none>
```

Monitor the status of the deployment, eventually resulting in all the resources being in 
the `Ready` state:

```console
kubectl -n k8ssandra-operator describe K8ssandraCluster demo
```

```console
Name:         demo
Namespace:    k8ssandra-operator
Labels:       <none>
Annotations:  <none>
API Version:  k8ssandra.io/v1alpha1
Kind:         K8ssandraCluster
...
Status:
  Datacenters:
    dc1:
      Cassandra:
        Cassandra Operator Progress:  Ready
      ...
      Stargate:
        Available Replicas:  1
        Conditions:
          Last Transition Time:  2021-09-28T03:32:07Z
          Status:                True
          Type:                  Ready
        Deployment Refs:
          demo-dc1-default-stargate-deployment
        Progress:              Running
        Ready Replicas:        1
        Ready Replicas Ratio:  1/1
        Replicas:              1
        Service Ref:           demo-dc1-stargate-service
        Updated Replicas:      1
Events:                        <none>
```

## Next steps

* See other [local install]({{< relref "install/local/" >}}) options, including K8ssandra Operator in multi-cluster Kubernetes.
* Also, dig into the K8ssandra Operator [components]({{< relref "components" >}}).

