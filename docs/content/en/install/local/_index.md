---
title: "Install K8ssandra Operator on local K8s"
linkTitle: "Local"
no_list: true
weight: 1
description: "Details to install K8ssandra Operator on a local Kubernetes **kind** development environment."
---

This topic explains how to install K8ssandra Operator in a local dev **kind** Kubernetes (K8s) environment. The configuration results in a Apache Cassandra&reg; database deployment in a **multi-cluster, multi-region** K8s environment. Included in the deployment are additional services, such as Stargate (API), Reaper (anti-entropy data repairs), and Medusa (backup/restore). Also shown in this topic are K8ssandra Operator install steps for a single-datacenter `kind` K8s cluster.

{{% alert title="Tip" color="success" %}}
Follow-up topics cover the post-install steps and role-based considerations for [developers]({{< relref "/quickstarts/developer">}}) or [Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer">}}) (SREs).
{{% /alert %}}

## Introduction

In this topic for installs of K8ssandra Operator on your local development environment, using **kind**, we'll cover:

* [Prerequisites]({{< relref "#prerequisites" >}}): Required supporting software
* [Quick Start]({{< relref "#quick-start" >}}): Quick start install including helper scripts, with single-cluster and multi-cluster examples
* [Helm]({{< relref "#helm" >}}): Install via a single K8ssandra Operator Helm chart, with single-cluster and multi-cluster examples
* [Kustomize]({{< relref "#kustomize" >}}): Install via Kustomize - a declarative approach, with single-cluster and multi-cluster examples
* [Next steps]({{< relref "#next-steps" >}}): Quick start info including helper scripts

## Prerequisites
Make sure you have the following installed before going through the rest of the guide. 

* [kind](#kind)
* [kubectx](#kubectx)
* [yq (YAML processor)](#yq)
* [setup-kind-multicluster.sh](#setup-kind-multiclustersh)
* [create-clientconfig.sh](#create-clientconfigsh)

### **kind**

The examples in this topic use [kind](https://kind.sigs.k8s.io/) clusters. Install it now if you have not already done so.

By default, kind clusters run on the same Docker network, which means we will have routable pod IPs across clusters.

**Note:**  Issues creating multiple kind clusters have been observed on various versions of Docker Desktop for macOS.  These issues seem to be resolved with the 4.5.0 release of Docker Desktop.  Please be sure to upgrade Docker Desktop if you plan to deploy using kind. Other options for local dev K8s environments include [minikube](https://minikube.sigs.k8s.io/docs/start/) or [K3D](https://k3d.io/v5.3.0/). 

### **kubectx**

[kubectx](https://github.com/ahmetb/kubectx) is a really handy tool when you are dealing with multiple clusters. The examples will use it so go ahead and install it now.

### **yq**

[yq](https://github.com/mikefarah/yq#install) is lightweight and portable command-line YAML processor.

### setup-kind-multicluster.sh

[setup-kind-multicluster.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/setup-kind-multicluster.sh) lives in the k8ssandra-operator repo. It is used extensively during development and testing. Not only does it configure and create kind clusters, it also generates kubeconfig files for each cluster.

**Note:** kind generates a kubeconfig with the IP address of the API server set to 
localhost since the cluster is intended for local development. We need a kubeconfig with the IP address set to the internal address of the api server. `setup-kind-mulitcluster.sh` takes care of this for us.

### create-clientconfig.sh

[create-clientconfig.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/create-clientconfig.sh) lives in the k8ssandra-operator repo. It is used to configure access to remote clusters. 

## Quick Start

If you're interested in getting running as quickly as possible, there's a number of helper scripts that can be used to greatly reduce the steps to deploy a local K8ssandra cluster via kind for testing purposes.

Two base `make` commands are provided that deploy a basic kind-based Kubernetes cluster(s).  These commands encapsulate the more detailed step-by-step installation instructions otherwise captured in this document.

Each of these commands will do the following:

* Create the kind-based cluster(s)

Across the cluster:

* Install cert-manager in it's own namespace
* Install cass-operator in the `k8ssandra-operator` namespace
* Build the k8ssandra-operator from source, load the image into the kind nodes, and 
  install it in the `k8ssandra-operator` namespace
* Install relevant CRDs

At completion, the cluster is now ready to accept a `K8ssandraCluster` deployment.

**Note:** if a k8ssandra-0 and/or k8ssandra-1 kind cluster already exists, running `make 
single-up` or `make multi-up` will delete and recreate them.

**Note:** These steps will attempt to start a local Docker registry instance to be used by the kind cluster(s), if you are already running one locally it will need to be stopped before following these procedures.

### Single Cluster

Deploy a single kind based Kubernetes cluster.

```sh
make single-up
```

Once cluster should be available:

```sh
kubectx
```

```sh
kind-k8ssandra-0
```

The cluster should consist of the following nodes:

```sh
NAME                        STATUS   ROLES                  AGE     VERSION
k8ssandra-0-control-plane   Ready    control-plane,master   3m24s   v1.21.2
k8ssandra-0-worker          Ready    <none>                 2m53s   v1.21.2
k8ssandra-0-worker2         Ready    <none>                 3m5s    v1.21.2
k8ssandra-0-worker3         Ready    <none>                 2m53s   v1.21.2
k8ssandra-0-worker4         Ready    <none>                 2m53s   v1.21.2
```

Once the Kubernetes cluster is ready, deploy a `K8ssandraCluster` like:

```sh
cat <<EOF | kubectl -n k8ssandra-operator apply -f -
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    serverVersion: "4.0.1"
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

Monitor the status of the deployment, eventually resulting in all the resources being in the `Ready` state:

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

### Multi-Cluster

Deploy two kind based Kubernetes clusters with:

```console
make multi-up
```

Two clusters should be available:

```console
kubectx
```

```console
kind-k8ssandra-0
kind-k8ssandra-1
```

Each cluster should consist of the following nodes:

kind-k8ssandra-0:

```console
NAME                        STATUS   ROLES                  AGE     VERSION
k8ssandra-0-control-plane   Ready    control-plane,master   9m20s   v1.21.2
k8ssandra-0-worker          Ready    <none>                 8m49s   v1.21.2
k8ssandra-0-worker2         Ready    <none>                 8m49s   v1.21.2
k8ssandra-0-worker3         Ready    <none>                 8m48s   v1.21.2
k8ssandra-0-worker4         Ready    <none>                 8m49s   v1.21.2
```

kind-k8ssandra-1

```console
NAME                        STATUS   ROLES                  AGE     VERSION
k8ssandra-1-control-plane   Ready    control-plane,master   9m51s   v1.21.2
k8ssandra-1-worker          Ready    <none>                 9m32s   v1.21.2
k8ssandra-1-worker2         Ready    <none>                 9m20s   v1.21.2
k8ssandra-1-worker3         Ready    <none>                 9m32s   v1.21.2
k8ssandra-1-worker4         Ready    <none>                 9m20s   v1.21.2
```

You're now ready to deploy a `K8ssandraCluster`.

Set your context to the control-plane cluster (`kind-k8ssandra-0`):

```console
kubectx kind-k8ssandra-0
```

```console
Switched to context "kind-k8ssandra-0".
```

Deploy the `K8ssandraCluster` resource:

```sh
cat <<EOF | kubectl -n k8ssandra-operator apply -f -
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    serverVersion: "4.0.1"
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

Confirm that the resource has been created:

```console
kubectl -n k8ssandra-operator get k8ssandraclusters
```

```console
NAME   AGE
demo   45s
```

```console
kubectl describe -n k8ssandra-operator K8ssandraCluster demo
```

```sh
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
          Last Transition Time:  2021-09-27T17:52:41Z
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
    dc2:
      Cassandra:
        Cassandra Operator Progress:  Ready
      ...
      Stargate:
        Available Replicas:  1
        Conditions:
          Last Transition Time:  2021-09-27T17:53:40Z
          Status:                True
          Type:                  Ready
        Deployment Refs:
          demo-dc2-default-stargate-deployment
        Progress:              Running
        Ready Replicas:        1
        Ready Replicas Ratio:  1/1
        Replicas:              1
        Service Ref:           demo-dc2-stargate-service
        Updated Replicas:      1
Events:  <none>
```
## Helm
You need to have [Helm v3+](https://helm.sh/docs/intro/install/) installed.

Configure the K8ssandra Helm repository:

```console
helm repo add k8ssandra https://helm.k8ssandra.io/stable
```

Update your Helm repository cache:

```console
helm repo update
```

Verify that you see the `k8ssandra-operator` chart:

```console
helm search repo k8ssandra-operator
```

```console
NAME                                CHART VERSION   APP VERSION DESCRIPTION
k8ssandra/k8ssandra-operator        0.32.0          1.0.0       Kubernetes operator which handles the provision...
```

### Single Cluster 
We will first look at a single cluster install to demonstrate that while K8ssandra 
Operator is designed for multi-cluster use, it can be used in a single cluster without 
any extra configuration.

#### Create kind cluster
Run `setup-kind-multicluster.sh` as follows:

```sh
./setup-kind-multicluster.sh --kind-worker-nodes 4
```

#### Install K8ssandra Operator
Install the Helm chart:

```console
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
```
This `helm install` command does the following:

* Create the `k8ssandra-operator` namespace if necessary
* Install Cass Operator in the `k8ssandra-operator` namespace
* Install K8ssandra Operator in the `k8ssandra-operator` namespace

This does not currently install Cert Manager. Cass Operator requires Cert Manager when 
its webhook is enabled. This installs with the webhook disabled.

Verify that the Helm release is installed:

```console
helm ls -n k8ssandra-operator
```

```console
NAME                NAMESPACE           REVISION    UPDATED                                 STATUS      CHART                       APP VERSION
k8ssandra-operator  k8ssandra-operator  1           2021-09-30 16:28:08.722822 -0400 EDT    deployed    k8ssandra-operator-0.32.0   1.0.0
```

Verify that the following CRDs are installed:

* `cassandradatacenters.cassandra.datastax.com`
* `clientconfigs.config.k8ssandra.io`
* `k8ssandraclusters.k8ssandra.io`
* `replicatedsecrets.replication.k8ssandra.io`
* `stargates.stargate.k8ssandra.io`


Check that there are two Deployments. The output should look similar to this:

```console
kubectl -n k8ssandra-operator get deployment
```

```console
NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE
k8ssandra-operator-cass-operator        1/1     1            1           85s
k8ssandra-operator-k8ssandra-operator   1/1     1            1           85s
```

#### Deploy a K8ssandraCluster
Now we will deploy a K8ssandraCluster that consists of a 3-node Cassandra cluster and a Stargate node.

```sh
cat <<EOF | kubectl -n k8ssandra-operator apply -f -
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    serverVersion: "4.0.1"
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

Monitor the status of the deployment, eventually resulting in all the resources being in the `Ready` state:

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

### Multi-Cluster

If you previously created a cluster with `setup-kind-multicluster.sh` we need to delete 
it in order to create the multi-cluster setup. The script currently does not support 
adding clusters to an existing setup (see[#128](https://github.com/k8ssandra/k8ssandra-operator/issues/128)).

We will create two kind clusters with 3 worker nodes per clusters. Remember that 
K8ssandra Operator requires clusters to have routable pod IPs. kind clusters by default 
will run on the same Docker network which means that they will have routable IPs.

#### Create kind clusters
Run `setup-kind-multicluster.sh` as follows:

```sh
./setup-kind-multicluster.sh --clusters 2 --kind-worker-nodes 4
```

When creating a cluster, kind generates a kubeconfig with the address of the API server 
set to localhost. We need a kubeconfig that has the API server address set to its 
internal ip address. `setup-kind-multi-cluster.sh` takes care of this for us. Generated 
files are written into a `build` directory.

Run `kubectx` without any arguments and verify that you see the following contexts 
listed in the output:

* kind-k8ssandra-0
* kind-k8ssandra-1

#### Install the control plane
We will install the control plane in `kind-k8ssandra-0`. Make sure your active context 
is configured correctly:

```console
kubectx kind-k8ssandra-0
```

Install the operator:

```console
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
```

This `helm install` command does the following:

* Create the `k8ssandra-operator` namespace if necessary
* Install Cass Operator in the `k8ssandra-operator` namespace
* Install K8ssandra Operator in the `k8ssandra-operator` namespace

This does not currently install Cert Manager. Cass Operator requires Cert Manager when
its webhook is enabled. This installs with Cass Operator's webhook disabled.

Verify that the Helm release is installed:

```console
helm ls -n k8ssandra-operator
```

```console
NAME                NAMESPACE           REVISION    UPDATED                                 STATUS      CHART                       APP VERSION
k8ssandra-operator  k8ssandra-operator  1           2021-09-30 16:28:08.722822 -0400 EDT    deployed    k8ssandra-operator-0.32.0   1.0.0
```

Verify that the following CRDs are installed:

* `cassandradatacenters.cassandra.datastax.com`
* `clientconfigs.k8ssandra.io`
* `k8ssandraclusters.k8ssandra.io`
* `replicatedsecrets.k8ssandra.io`
* `stargates.k8ssandra.io`

Check that there are two Deployments. The output should look similar to this:

```console
kubectl -n k8ssandra-operator get deployment
```

```console
NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE
k8ssandra-operator-cass-operator        1/1     1            1           85s
k8ssandra-operator-k8ssandra-operator   1/1     1            1           85s
```

The operator looks for an environment variable named `K8SSANDRA_CONTROL_PLANE`. When set 
to `true` the control plane is enabled. It is enabled by default.

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `true`:

```console
kubectl -n k8ssandra-operator get deployment k8ssandra-operator-k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

#### Install the data plane
Now we will install the data plane in `kind-k8ssandra-1`. Switch the active context:

```console
kubectx kind-k8ssandra-1
```

Install the operator:

```console
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace --set controlPlane=false
```

This `helm install` command does the following:

* Create the `k8ssandra-operator` namespace if necessary
* Install Cass Operator in the `k8ssandra-operator` namespace
* Install K8ssandra Operator in the `k8ssandra-operator` namespace
* Configures K8ssandra Operator to run in the data plane 

This does not currently install Cert Manager. Cass Operator requires Cert Manager when
its webhook is enabled. This installs with Cass Operator's webhook disabled.

Verify that the following CRDs are installed:

* `cassandradatacenters.cassandra.datastax.com`
* `clientconfigs.k8ssandra.io`
* `k8ssandraclusters.k8ssandra.io`
* `replicatedsecrets.k8ssandra.io`
* `stargates.k8ssandra.io`

Check that there are two Deployments. The output should look similar to this:

```console
kubectl -n k8ssandra-operator get deployment
```

```console
NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE
k8ssandra-operator-cass-operator        1/1     1            1           85s
k8ssandra-operator-k8ssandra-operator   1/1     1            1           85s
```

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `false`:

```console
kubectl -n k8ssandra-operator get deployment k8ssandra-operator-k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

#### Create a ClientConfig
Now we need to create a `ClientConfig` for the `kind-k8ssandra-1` cluster. We will use 
the `create-clientconfig.sh` script which can be found [here](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/create-clientconfig.sh).

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
./create-clientconfig.sh --namespace k8ssandra-operator --src-kubeconfig build/kubeconfigs/k8ssandra-1.yaml --dest-kubeconfig build/kubeconfigs/k8ssandra-0.yaml --in-cluster-kubeconfig build/kubeconfigs/updated/k8ssandra-1.yaml --output-dir clientconfig
```
The script stores all the artifacts that it generates in a directory which is specified with the `--output-dir` option. If not specified, a temp directory is created.

You can specify the namespace where the secret and ClientConfig are created with the `--namespace` option.

The `--in-cluster-kubeconfig` option is required for clusters that run locally like kind.

#### Restart the control plane

Make the active context `kind-k8ssandra-0`:

```console
kubectx kind-k8ssandra-0
```

Restart the operator:

```console
kubectl -n k8ssandra-operator rollout restart deployment k8ssandra-operator-k8ssandra-operator
```

**Note:** See https://github.com/k8ssandra/k8ssandra-operator/issues/178 for details on
why it is necessary to restart the control plane operator.

#### Deploy a K8ssandraCluster
Now we will create a K8ssandraCluster that consists of a Cassandra cluster with 2 DCs and 3 
nodes per DC, and a Stargate node per DC.

```sh
cat <<EOF | kubectl -n k8ssandra-operator apply -f -
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    serverVersion: "3.11.11"
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

## Kustomize
K8ssandra Operator can be installed with [Kustomize](https://kustomize.io/) which takes 
a declarative approach to configuring and deploying resources whereas Helm takes more of 
an imperative approach.

The following examples use `kubectl apply -k` to deploy resources. The `-k` option
essentially runs `kustomize build` over the specified directory followed by `kubectl
apply`. See this [doc](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/)
for details on the integration of Kustomize into `kubectl`.

{{% alert title="Tip" color="success" %}}
If `kubectl -k <dir>` does not work for, you can instead use 
`kustomize build <dir> | kubectl apply -f -`.
{{% /alert %}}

### Single Cluster
We will first look at a single cluster install to demonstrate that while K8ssandra 
Operator is designed for multi-cluster use, it can be used in a single cluster without 
any extra configuration.

#### Create kind cluster
Run `setup-kind-multicluster.sh` as follows:

```sh
./setup-kind-multicluster.sh --kind-worker-nodes 4
```

#### Install Cert Manager
We need to first install Cert Manager as it is a dependency of cass-operator:

```console
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.3/cert-manager.yaml
```

#### Install K8ssandra Operator
The GitHub Actions for the project are configured to build and push a new operator image 
to Docker Hub whenever commits are pushed to `main`. 

See [here](https://hub.docker.com/repository/docker/k8ssandra/k8ssandra-operator/tags?page=1&ordering=last_updated) 
on Docker Hub for a list of available images.

Install with kubectl:

```console
kubectl apply -k github.com/k8ssandra/k8ssandra-operator/config/deployments/control-plane
```

This installs the operator in the `k8ssandra-operator` namespace.

**Note:** This will deploy the `latest` operator image, i.e., 
`k8ssandra/k8ssandra-operator:latest`. In general it is best to avoid using `latest`. 

In case you want to customize the installation, create a kustomization directory that 
builds from the `main` branch and in this case we'll add namespace creation and define 
new namespace. Note the `namespace` property which we added. This property tells 
Kustomize to apply a transformation on all resources that specify a namespace.

```sh
K8SSANDRA_OPERATOR_HOME=$(mktemp -d)
cat <<EOF >$K8SSANDRA_OPERATOR_HOME/kustomization.yaml

namespace: k8ssandra-operator

resources:
- github.com/k8ssandra/k8ssandra-operator/config/deployments/default?ref=main

components:
- github.com/k8ssandra/k8ssandra-operator/config/components/namespace

images:
- name: k8ssandra/k8ssandra-operator
  newTag: v1.0.0-alpha.1
EOF
```

Now install the operator:

```console
kubectl apply -k $K8SSANDRA_OPERATOR_HOME
```

This installs the operator in the `k8ssandra-operator` namespace.

If you just want to generate the manifests then run:

```console
kustomize build $K8SSANDRA_OPERATOR_HOME
```

Verify that the following CRDs are installed:

* `cassandradatacenters.cassandra.datastax.com`
* `certificaterequests.cert-manager.io`
* `certificates.cert-manager.io`
* `challenges.acme.cert-manager.io`
* `clientconfigs.config.k8ssandra.io`
* `clusterissuers.cert-manager.io`
* `issuers.cert-manager.io`
* `k8ssandraclusters.k8ssandra.io`
* `orders.acme.cert-manager.io`
* `replicatedsecrets.replication.k8ssandra.io`
* `stargates.stargate.k8ssandra.io`


Check that there are two Deployments. The output should look similar to this:

```console
kubectl -n k8ssandra-operator get deployment
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
cass-operator        1/1     1            1           2m
k8ssandra-operator   1/1     1            1           2m
```

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `false`:

```console
kubectl -n k8ssandra-operator get deployment k8ssandra-operator-k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

#### Deploy a K8ssandraCluster
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
    serverVersion: "4.0.1"
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

### Multi-cluster

If you previously created a cluster with `setup-kind-multicluster.sh` we need to delete 
it in order to create the multi-cluster setup. The script currently does not support 
adding clusters to an existing setup (see [#128](https://github.com/k8ssandra/k8ssandra-operator/issues/128)).

We will create two kind clusters with 3 worker nodes per clusters. Remember that 
K8ssandra Operator requires clusters to have routable pod IPs. kind clusters by default 
will run on the same Docker network which means that they will have routable IPs.

#### Create kind clusters
Run `setup-kind-multicluster.sh` as follows:

```sh
./setup-kind-multicluster.sh --clusters 2 --kind-worker-nodes 4
```

When creating a cluster, kind generates a kubeconfig with the address of the API server 
set to localhost. We need a kubeconfig that has the API server address set to its 
internal ip address. `setup-kind-multi-cluster.sh` takes care of this for us. Generated 
files are written into a `build` directory.

Run `kubectx` without any arguments and verify that you see the following contexts 
listed in the output:

* kind-k8ssandra-0
* kind-k8ssandra-1

#### Install Cert Manager
Set the active context to `kind-k8ssandra-0`:

```console
kubectx kind-k8ssandra-0
```

Install Cert Manager:

```console
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.3/cert-manager.yaml
```

Set the active context to `kind-k8ssandra-1`:

```console
kubectx kind-k8ssandra-1
```

Install Cert Manager:

```console
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.3/cert-manager.yaml
```

#### Install the control plane
We will install the control plane in `kind-k8ssandra-0`. Make sure your active context is 
configured correctly:

```console
kubectx kind-k8ssandra-0
```
Now install the operator:

```console
kubectl apply -k github.com/k8ssandra/config/deployments/control-plane
```

This installs the operator in the `k8ssandra-operator` namespace.

Verify that the following CRDs are installed:

* `cassandradatacenters.cassandra.datastax.com`
* `certificaterequests.cert-manager.io`
* `certificates.cert-manager.io`
* `challenges.acme.cert-manager.io`
* `clientconfigs.k8ssandra.io`
* `clusterissuers.cert-manager.io`
* `issuers.cert-manager.io`
* `k8ssandraclusters.k8ssandra.io`
* `orders.acme.cert-manager.io`
* `replicatedsecrets.k8ssandra.io`
* `stargates.k8ssandra.io`


Check that there are two Deployments. The output should look similar to this:

```console
kubectl get deployment
```

```console
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
cass-operator        1/1     1            1           2m
k8ssandra-operator   1/1     1            1           2m
```

The operator looks for an environment variable named `K8SSANDRA_CONTROL_PLANE`. When set 
to `true` the control plane is enabled. It is enabled by default.

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `true`:

```sh
kubectl -n k8ssandra-operator get deployment k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

#### Install the data plane
Now we will install the data plane in `kind-k8ssandra-1`. Switch the active context:

```
kubectx kind-k8ssandra-1
```

Now install the operator:

```console
kubectl apply -k github.com/k8ssandra/config/deployments/data-plane
```

This installs the operator in the `k8ssandra-operator` namespace.

Verify that the following CRDs are installed:

* `cassandradatacenters.cassandra.datastax.com`
* `certificaterequests.cert-manager.io`
* `certificates.cert-manager.io`
* `challenges.acme.cert-manager.io`
* `clientconfigs.k8ssandra.io`
* `clusterissuers.cert-manager.io`
* `issuers.cert-manager.io`
* `k8ssandraclusters.k8ssandra.io`
* `orders.acme.cert-manager.io`
* `replicatedsecrets.k8ssandra.io`
* `stargates.k8ssandra.io`


Check that there are two Deployments. The output should look similar to this:

```console
kubectl -n k8ssandra-operator get deployment
```
```console
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
cass-operator        1/1     1            1           2m
k8ssandra-operator   1/1     1            1           2m
```

Verify that the `K8SSANDRA_CONTROL_PLANE` environment variable is set to `false`:

```console
kubectl -n k8ssandra-operator get deployment k8ssandra-operator -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

#### Create a ClientConfig
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
./create-clientconfig.sh --namespace k8ssandra-operator --src-kubeconfig build/kubeconfigs/k8ssandra-1.yaml --dest-kubeconfig build/kubeconfigs/k8ssandra-0.yaml --in-cluster-kubeconfig build/kubeconfigs/updated/k8ssandra-1.yaml --output-dir clientconfig
```
The script stores all of the artifacts that it generates in a directory which is specified with the `--output-dir` option. If not specified, a temp directory is created.

The `--in-cluster-kubeconfig` option is required for clusters that run locally like kind.

#### Restart the control plane

Make the active context `kind-k8ssandra-0`:

```console
kubectx kind-k8ssandra-0
```

Restart the operator:

```console
kubectl -n k8ssandra-operator rollout restart deployment k8ssandra-operator
```

**Note:** See https://github.com/k8ssandra/k8ssandra-operator/issues/178 for details on
why it is necessary to restart the control plane operator.

### Deploy a K8ssandraCluster
Now we will create a K8ssandraCluster that consists of a Cassandra cluster with 2 DCs and 3 
nodes per DC, and a Stargate node per DC.

```sh
cat <<EOF | kubectl -n k8ssandra-operator apply -f -
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: demo
spec:
  cassandra:
    serverVersion: "4.0.1"
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

* If you're a developer, and you'd like to get started coding using CQL or Stargate, see the [Quickstart for developers]({{< relref "/quickstarts/developer" >}}).
* If you're a Site Reliability Engineer, and you'd like to explore the K8ssandra administration environment including monitoring and maintenance utilities, see the [Quickstart for Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer" >}}).
