---
title: "Multi-cluster install with kustomize"
linkTitle: "Multi-cluster/kustomize"
no_list: false
weight: 4
description: "Quickstart with Kustomize to install K8ssandraCluster in multi-cluster Kubernetes."
---

This topic shows how you can use Kustomize to declaratively install and configure the `K8ssandraCluster` custom resource in **multi-cluster** local Kubernetes. 

## Prerequisites

If you haven't already, see the install [prerequisites]({{< relref "install/local/" >}}).

## Introduction

You can install K8ssandra Operator with [Kustomize](https://kustomize.io/), which takes 
a declarative approach to configuring and deploying resources, whereas Helm takes more of 
an imperative approach.

Kustomize is integrated directly into `kubectl`. For example, `kubectl apply -k` essentially runs `kustomize build` over the specified directory followed by `kubectl apply`. See this [topic](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/) for details on the integration of Kustomize into `kubectl`.

K8ssandra Operator uses some features of Kustomize that are only available in `kubectl` v1.23 or later. For this reason the following examples use `kustomize build <dir> | kubectl apply -f -`.


## Multi-cluster local Kubernetes

If you previously created a [single cluster]({{< relref "install/local/" >}}) with `setup-kind-multicluster.sh`, you will need to delete it in order to perform the multi-cluster setup. The script currently does not support 
adding clusters to an existing setup (see [#128](https://github.com/k8ssandra/k8ssandra-operator/issues/128)).

In this multi-cluster topic, we will create two kind clusters with four worker nodes per clusters. Remember that 
K8ssandra Operator requires clusters to have routable pod IPs. kind clusters by default 
will run on the same Docker network, which means that they will have routable IPs.

### Create kind clusters
Run `setup-kind-multicluster.sh` as follows:

```sh
scripts/setup-kind-multicluster.sh --clusters 2 --kind-worker-nodes 4
```

When creating a cluster, kind generates a kubeconfig with the address of the API server 
set to localhost. We need a kubeconfig that has the API server address set to its 
internal ip address. `setup-kind-multi-cluster.sh` takes care of this for us. Generated 
files are written into a `build` directory.

Run `kubectl config get-contexts` without any arguments and verify that you see the following contexts 
listed in the output:

```
          kind-k8ssandra-0                                            kind-k8ssandra-0                                            kind-k8ssandra-0                                               
          kind-k8ssandra-1                                            kind-k8ssandra-1                                            kind-k8ssandra-1 
```

### Install Cert Manager
Set the active context to `kind-k8ssandra-0`:

```console
kubectl config use-context kind-k8ssandra-0
```

Install Cert Manager:

```console
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```

Set the active context to `kind-k8ssandra-1`:

```console
kubectl config use-context kind-k8ssandra-1
```

Install Cert Manager:

```console
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```

### Install the control plane
We will install the control plane in `kind-k8ssandra-0`. Make sure your active context is 
configured correctly:

```console
kubectl config use-context kind-k8ssandra-0
```
Now install the operator (replace `X.X.X` with [the release](https://github.com/k8ssandra/k8ssandra-operator/releases) you which to install):

```console
kustomize build "github.com/k8ssandra/k8ssandra-operator/config/deployments/control-plane?ref=vX.X.X" | kubectl apply --server-side -f -
```

This installs the operator in the `k8ssandra-operator` namespace.

Verify that the following CRDs are installed:

* `cassandrabackups.medusa.k8ssandra.io`
* `cassandradatacenters.cassandra.datastax.com`
* `cassandrarestores.medusa.k8ssandra.io`
* `cassandratasks.control.k8ssandra.io`
* `clientconfigs.config.k8ssandra.io`
* `k8ssandraclusters.k8ssandra.io`
* `medusabackupjobs.medusa.k8ssandra.io`
* `medusabackups.medusa.k8ssandra.io`
* `medusabackupschedules.medusa.k8ssandra.io`
* `medusarestorejobs.medusa.k8ssandra.io`
* `medusatasks.medusa.k8ssandra.io`
* `reapers.reaper.k8ssandra.io`
* `replicatedsecrets.replication.k8ssandra.io`
* `stargates.stargate.k8ssandra.io`

Check that there are two Deployments. The output should look similar to this:

```console
kubectl get deployment -n k8ssandra-operator
```

```console
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
cass-operator-controller-manager   1/1     1            1           2m
k8ssandra-operator                 1/1     1            1           2m
```

The operator looks for an environment variable named `K8SSANDRA_CONTROL_PLANE`. When set 
to `true` the control plane is enabled. It is enabled by default.

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `true`:

```sh
kubectl -n k8ssandra-operator get deployment k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

### Install the data plane
Now we will install the data plane in `kind-k8ssandra-1`. Switch the active context:

```
kubectl config use-context kind-k8ssandra-1
```

Now install the operator (using the same release as the control plane deployment):

```console
kustomize build "github.com/k8ssandra/k8ssandra-operator/config/deployments/data-plane?ref=vX.X.X" | kubectl apply --server-side -f -
```

This installs the operator in the `k8ssandra-operator` namespace.

Verify that the following CRDs are installed:

* `cassandrabackups.medusa.k8ssandra.io`
* `cassandradatacenters.cassandra.datastax.com`
* `cassandrarestores.medusa.k8ssandra.io`
* `cassandratasks.control.k8ssandra.io`
* `clientconfigs.config.k8ssandra.io`
* `k8ssandraclusters.k8ssandra.io`
* `medusabackupjobs.medusa.k8ssandra.io`
* `medusabackups.medusa.k8ssandra.io`
* `medusabackupschedules.medusa.k8ssandra.io`
* `medusarestorejobs.medusa.k8ssandra.io`
* `medusatasks.medusa.k8ssandra.io`
* `reapers.reaper.k8ssandra.io`
* `replicatedsecrets.replication.k8ssandra.io`
* `stargates.stargate.k8ssandra.io`

Check that there are two Deployments. The output should look similar to this:

```console
kubectl -n k8ssandra-operator get deployment
```
```console
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
cass-operator-controller-manager   1/1     1            1           2m
k8ssandra-operator                 1/1     1            1           2m```
```

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `false`:

```console
kubectl -n k8ssandra-operator get deployment k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

### Create a ClientConfig
Now we need to create a `ClientConfig` for the `kind-k8ssandra-1` cluster. We will use the 
`create-clientconfig.sh` script which can be found
[here](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/create-clientconfig.sh).

Here is a summary of what the script does:

* Get the k8ssandra-operator service account from the data plane cluster
* Extract the service account token 
* Extract the CA cert
* Create a kubeonfig using the token and cert
* Create a secret for the kubeconfig in the control plane cluster
* Create a ClientConfig in the control plane cluster that references the secret

Create a `ClientConfig` in the `kind-k8ssandra-0` cluster using the service account 
token and CA cert from `kind-k8ssandra-1`:

```sh
scripts/create-clientconfig.sh --namespace k8ssandra-operator \
    --src-kubeconfig ./build/kind-kubeconfig \
    --dest-kubeconfig ./build/kind-kubeconfig \
    --src-context kind-k8ssandra-1 \
    --dest-context kind-k8ssandra-0 \
    --output-dir clientconfig
```
The script stores all of the artifacts that it generates in a directory which is specified with the `--output-dir` option. If not specified, a temp directory is created.

### Restart the control plane

There is a controller in the operator that watches for ClientConfig changes. When it detects a change (create/update/delete), it automatically restarts the operator.

**Note:** See https://github.com/k8ssandra/k8ssandra-operator/issues/178 for details on
why it is necessary to restart the control plane operator.

## Deploy a K8ssandraCluster
Now we will create a `K8ssandraCluster` custom resource that consists of a Cassandra cluster with 2 DCs and 3 
nodes per DC, and a Stargate node per DC.

```sh
cat <<EOF | kubectl -n k8ssandra-operator apply -f -
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    serverVersion: "4.0.3"
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
    networking:
      hostNetwork: true    
    datacenters:
      - metadata:
          name: dc1
        size: 3
        stargate:
          size: 1
          heapSize: 256M
      - metadata:
          name: dc2
        k8sContext: kind-k8ssandra-1
        size: 3
        stargate:
          size: 1
          heapSize: 256M 
EOF
```

## Next steps

* See other [local install]({{< relref "install/local/" >}}) options.
* Also, dig into the K8ssandra Operator [components]({{< relref "components" >}}).
