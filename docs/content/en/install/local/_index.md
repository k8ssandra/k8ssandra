---
title: "Local installs"
linkTitle: "Local"
no_list: true
weight: 1
description: "Details to install K8ssandra on a local Kubernetes **development** environment."
---

This topic gets you up and running with a single-node Apache Cassandra® cluster on Kubernetes (K8s). 

If you want to install K8ssandra on a cloud provider's Kubernetes environment, see:

* [K8ssandra installs on Azure Kubernetes Service (AKS)]({{< relref "/install/aks" >}})
* [K8ssandra installs on DigitalOcean Managed Kubernetes Service (DOKS)]({{< relref "/install/doks" >}})
* [K8ssandra installs on Amazon Elastic Kubernetes Service (EKS)]({{< relref "/install/eks" >}})
* [K8ssandra installs on Google Kubernetes Engine (GKE)]({{< relref "/install/gke" >}})

{{% alert title="Tip" color="success" %}}
Also available in followup topics are post-install steps and role-based considerations for [developers]({{< relref "/quickstarts/developer">}}) or [Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer">}}) (SREs).
{{% /alert %}}

## Introduction

In this quickstart for installs of K8ssansdra on your local **DEV** environment, we'll cover:

* [Prerequisites]({{< relref "#prerequisites" >}}): Required supporting software including resource recommendations.
* [K8ssandra Helm repository configuration]({{< relref "#configure-the-k8ssandra-helm-repository" >}}): Accessing the Helm charts that install K8ssandra.
* [K8ssandra installation]({{< relref "#install-k8ssandra" >}}): Getting K8ssandra up and running locally using the Helm chart repo.
* [Verifying K8ssandra functionality]({{< relref "#verify-your-k8ssandra-installation" >}}): Making sure K8ssandra is working as expected.
* [Retrieve K8ssandra superuser credentials]({{< relref "#superuser" >}}): Getting the K8ssandra superuser name and password so you can access common utilities as well as the Stargate API.

## Prerequisites

In your local environment the following tools are required for provisioning a K8ssandra cluster:

* [Helm v3+](https://helm.sh/docs/intro/install/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

As K8ssandra deploys on a K8s cluster, one must be available to target for installation. The K8s environment may be a local version running on your development machine, an on-premises self-hosted environment, or a managed cloud offering.

K8ssandra works with the following versions of Kubernetes either standalone or via a cloud provider:

* 1.16
* 1.17
* 1.18
* 1.19
* 1.20

To verify your K8s server version:

```bash
kubectl version
```

**Output**:

```json
Client Version: version.Info{Major:"1", Minor:"20", GitVersion:"v1.20.3", GitCommit:"01849e73f3c86211f05533c2e807736e776fcf29", GitTreeState:"clean", BuildDate:"2021-02-18T12:10:55Z", GoVersion:"go1.15.8", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.16", GitCommit:"7a98bb2b7c9112935387825f2fce1b7d40b76236", GitTreeState:"clean", BuildDate:"2021-02-17T11:52:32Z", GoVersion:"go1.13.15", Compiler:"gc", Platform:"linux/amd64"}
```

Your K8s server version is a combination of the `Major:` and `Minor:` key/value pairs following `Server Version:`, in the example above, `1.18`.

If you don't have a K8s cluster available, you can use [OpenShift CodeReady Containers](https://developers.redhat.com/products/codeready-containers/overview) that run within a VM, or one of the following local versions that run within Docker:

* [K3D](https://k3d.io/)
* [Minikube](https://minikube.sigs.k8s.io/docs/start/)
* [Kind](https://kind.sigs.k8s.io/)

The instructions in this section focus on the Docker container solutions above, but the general instructions should work for other environments as well.

### Resource recommendations for local Kubernetes installations

We recommend a machine specification of **no less** than 16 gigs of RAM and 8 virtual processor cores (4 physical cores). You'll want to adjust your Docker resource preferences accordingly. For this quick start we're allocating 4 virtual processors and 8 gigs of RAM to the Docker environment.

{{% alert title="Tip" color="success" %}}
See the documentation for your particular flavor of Docker for instructions on configuring resource limits.
{{% /alert %}}

The following Minikube example creates a K8s cluster running K8s version 1.18.16 with 4 virtual processor cores and 8 gigs of RAM:

```bash
minikube start --cpus=4 --memory='8128m' --kubernetes-version=1.18.16
```

**Output**:

```bash
😄  minikube v1.17.1 on Darwin 11.2.1
✨  Automatically selected the docker driver. Other choices: hyperkit, ssh
👍  Starting control plane node k8ssandra in cluster k8ssandra
🔥  Creating docker container (CPUs=4, Memory=8128MB) ...
🐳  Preparing Kubernetes v1.18.16 on Docker 20.10.2 ...
    ▪ Generating certificates and keys ...
    ▪ Booting up control plane ...
    ▪ Configuring RBAC rules ...
🔎  Verifying Kubernetes components...
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```

### Verify your Kubernetes environment

To verify your Kubernetes environment:

1. Verify that your K8s instance is up and running in the `READY` status:

    ```bash
    kubectl get nodes
    ```

    **Output**:

    ```bash
    NAME        STATUS   ROLES    AGE   VERSION
    k8ssandra   Ready    master   21m   v1.18.16
    ```

### Validate the available Kubernetes StorageClasses {#storage-classes}

Your K8s instance **must** support a storage class with a `VOLUMEBINDINGMODE` of `WaitForFirstConsumer`.

To list the available K8s storage classes for your K8s instance:

```bash
kubectl get storageclasses
```

**Output**:

```bash
NAME                 PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
standard (default)   k8s.io/minikube-hostpath   Delete          Immediate           false                  2m25s
```

If you don't have a storage class with a `VOLUMEBINDINGMODE` of `WaitForFirstConsumer` as in the Minikube example above, you can install the [Rancher Local Path Provisioner](https://github.com/rancher/local-path-provisioner):

```bash
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
```

**Output**:

```bash
namespace/local-path-storage created
serviceaccount/local-path-provisioner-service-account created
clusterrole.rbac.authorization.k8s.io/local-path-provisioner-role created
clusterrolebinding.rbac.authorization.k8s.io/local-path-provisioner-bind created
deployment.apps/local-path-provisioner created
storageclass.storage.k8s.io/local-path created
configmap/local-path-config created
```

Rechecking the available storage classes, you should see that a new `local-path` storage class is available with the required `VOLUMEBINDINGMODE` of `WaitForFirstConsumer`:

```bash
kubectl get storageclasses
```

**Output**:

```bash
NAME                 PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
local-path           rancher.io/local-path      Delete          WaitForFirstConsumer   false                  3s
standard (default)   k8s.io/minikube-hostpath   Delete          Immediate              false                  39s
```

## Configure the K8ssandra Helm repository

K8ssandra is delivered via a collection of Helm charts for easy installation, so once you've got a suitable K8s environment configured, you'll need to add the K8ssandra Helm chart repositories.

To add the K8ssandra helm chart repos:

1. Install [Helm v3+](https://helm.sh/docs/intro/install/) if you haven't already.

1. Add the main K8ssandra stable Helm chart repo:

    ```bash
    helm repo add k8ssandra https://helm.k8ssandra.io/stable
    ```

1. If you want to access K8ssandra services from outside of the Kubernetes cluster, also add the Traefik Ingress repo:

    ```bash
    helm repo add traefik https://helm.traefik.io/traefik
    ```

1. Finally, update your helm repository listing:

    ```bash
    helm repo update
    ```

{{% alert title="Tip" color="success" %}}
Alternatively, you can download the individual charts directly from the project's [releases](https://github.com/k8ssandra/k8ssandra/releases) page.
{{% /alert %}}

## Install K8ssandra

The K8ssandra helm charts make installation a snap. You can override chart configurations during installation as necessary if you're an advanced user, or make changes after a default installation using `helm upgrade` at a later time.

K8ssandra can install the following versions of Apache Cassandra:

* 3.11.7
* 3.11.8
* 3.11.9
* 3.11.10
* 4.0-beta4

{{% alert title="Important" color="warning" %}}
K8ssandra comes out of the box with a set of [default values](https://github.com/k8ssandra/k8ssandra/blob/main/charts/k8ssandra/values.yaml) tailored to getting up and running quickly.  Those defaults are intended to be a great starting point for smaller-scale local development but are **not** intended for production deployments.
{{% /alert %}}

To install a single node K8ssandra cluster:

1. Copy the following YAML to a file named `k8ssandra.yaml`:

    ```yaml
    cassandra:
      version: "3.11.10"
      cassandraLibDirVolume:
        storageClass: local-path
        size: 5Gi
      allowMultipleNodesPerWorker: true
      heap:
       size: 1G
       newGenSize: 1G
      resources:
        requests:
          cpu: 1000m
          memory: 2Gi
        limits:
          cpu: 1000m
          memory: 2Gi
      datacenters:
      - name: dc1
        size: 1
        racks:
        - name: default
    kube-prometheus-stack:
      grafana:
        adminUser: admin
        adminPassword: admin123
    stargate:
      enabled: true
      replicas: 1
      heapMB: 256
      cpuReqMillicores: 200
      cpuLimMillicores: 1000
    ```

    That configuration file creates a K8ssandra cluster with a datacenter, `dc1`, containing a single Cassandra node, `size: 1` version `3.11.10` with the following specifications:

    * 1 GB of heap
    * 2 GB of RAM for the container
    * 1 CPU core
    * 5 GB of storage
    * 1 Stargate node with
      * 1 CPU core
      * 256 MB of heap

    {{% alert title="Important" color="warning" %}}
The `storageClass:` parameter must be a storage class with a `VOLUMEBINDINGMODE` of `WaitForFirstConsumer` as described in [Validate the available Kubernetes StorageClasses]({{< relref "#storage-classes" >}}).
    {{% /alert %}}

1. Use `helm install` to install K8ssandra, pointing to the example configuration file using the `-f` flag:

    ```bash
    helm install -f k8ssandra.yaml k8ssandra k8ssandra/k8ssandra
    ```

    **Output**:

    ```bash
    NAME: k8ssandra
    LAST DEPLOYED: Thu Feb 18 10:05:44 2021
    NAMESPACE: default
    STATUS: deployed
    REVISION: 1
    ```

    {{% alert title="Tip" color="success" %}}
In the example above, the K8ssandra pods will have the cluster name `k8ssandra` prefixed or appended inline.
    {{% /alert %}}

    {{% alert title="Note" color="primary" %}}
When installing K8ssandra on newer versions of Kubernetes (v1.19+), some warnings may be visible on the command line related to deprecated API usage.  This is currently a known issue and will not impact the provisioning of the cluster.

```bash
W0128 11:24:54.792095  27657 warnings.go:70] 
apiextensions.k8s.io/v1beta1 CustomResourceDefinition is 
deprecated in v1.16+, unavailable in v1.22+; 
use apiextensions.k8s.io/v1 CustomResourceDefinition
```

For more information, check out issue [#267](https://github.com/k8ssandra/k8ssandra/issues/267).
    {{% /alert %}}

## Verify your K8ssandra installation

Depending upon your K8s configuration, initialization of your K8ssandra installation can take a few minutes. To check the status of your K8ssandra deployment, use the `kubectl get pods` command:

```bash
kubectl get pods
```

**Output**:

```bash
NAME                                                READY   STATUS      RESTARTS   AGE
k8ssandra-cass-operator-766849b497-klgwf            1/1     Running     0          7m33s
k8ssandra-dc1-default-sts-0                         2/2     Running     0          7m5s
k8ssandra-dc1-stargate-5c46975f66-pxl84             1/1     Running     0          7m32s
k8ssandra-grafana-679b4bbd74-wj769                  2/2     Running     0          7m32s
k8ssandra-kube-prometheus-operator-85695ffb-ft8f8   1/1     Running     0          7m32s
k8ssandra-reaper-655fc7dfc6-n9svw                   1/1     Running     0          4m52s
k8ssandra-reaper-operator-79fd5b4655-748rv          1/1     Running     0          7m33s
k8ssandra-reaper-schema-dxvmm                       0/1     Completed   0          5m3s
prometheus-k8ssandra-kube-prometheus-prometheus-0   2/2     Running     1          7m27s
```

The K8ssandra pods in the example above have the identifier `k8ssandra` either prefixed or inline, since that's the name that was specified when the cluster was created using Helm. If you choose a different cluster name during installation, your pod names will be different.

The actual Cassandra node name from the listing above is `k8ssandra-dc1-default-sts-0` which we'll use throughout the following sections.

Verify the following:

* The K8ssandra pod running Cassandra, `k8ssandra-dc1-default-sts-0` in the example above should show `2/2` as `Ready`.
* The Stargate pod, `k8ssandra-dc1-stargate-5c46975f66-pxl84` in the example above should show `1/1` as `Ready`.

{{% alert title="Important" color="warning" %}}

* The Stargate pod will not show `Ready` until at least 4 minutes have elapsed.
* The pod `k8ssandra-reaper-k8ssandra-schema-xxxxx` runs once as part of a job and does not persist.

{{% /alert %}}

Once all the pods are in the `Running` or `Completed` state, you can check the health of your K8ssandra cluster. There must be **no `PENDING` pods**.

To check the health of your K8ssandra cluster:

1. Verify the name of the Cassandra datacenter:

    ```bash
    kubectl get cassandradatacenters
    ```

    **Output**:

    ```bash
    NAME   AGE
    dc1    51m
    ```

1. Confirm that the Cassandra operator for the datacenter is `Ready`:

    ```bash
    kubectl describe CassandraDataCenter dc1 | grep "Cassandra Operator Progress:"
    ```

    **Output**:

    ```bash
       Cassandra Operator Progress:  Ready
    ```

1. Verify the list of available services:

    ```bash
    kubectl get services
    ```

    **Output**:

    ```bash
    NAME                                   TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)                                                 AGE
    cass-operator-metrics                  ClusterIP   10.80.3.92     <none>        8383/TCP,8686/TCP                                       47m
    k8ssandra-dc1-all-pods-service         ClusterIP   None           <none>        9042/TCP,8080/TCP,9103/TCP                              47m
    k8ssandra-dc1-service                  ClusterIP   None           <none>        9042/TCP,9142/TCP,8080/TCP,9103/TCP,9160/TCP            47m
    k8ssandra-dc1-stargate-service         ClusterIP   10.80.13.197   <none>        8080/TCP,8081/TCP,8082/TCP,8084/TCP,8085/TCP,9042/TCP   47m
    k8ssandra-grafana                      ClusterIP   10.80.7.168    <none>        80/TCP                                                  47m
    k8ssandra-kube-prometheus-operator     ClusterIP   10.80.8.109    <none>        443/TCP                                                 47m
    k8ssandra-kube-prometheus-prometheus   ClusterIP   10.80.2.44     <none>        9090/TCP                                                47m
    k8ssandra-reaper-reaper-service        ClusterIP   10.80.5.77     <none>        8080/TCP                                                47m
    k8ssandra-seed-service                 ClusterIP   None           <none>        <none>                                                  47m
    kubernetes                             ClusterIP   10.80.0.1      <none>        443/TCP                                                 47m
    prometheus-operated                    ClusterIP   None           <none>        9090/TCP                                                47m
    ```

    Verify that the following services are present:

    * <cluster-name>-<datacenter-name>-all-pods-service
    * <cluster-name>-<datacenter-name>-dc1-service
    * <cluster-name>-<datacenter-name>-stargate-service
    * <cluster-name>-<datacenter-name>-seed-service

## Retrieve K8ssandra superuser credentials {#superuser}

You'll need the K8ssandra superuser name and password in order to access Cassandra utilities and do things like generate a Stargate access token.

To retrieve K8ssandra superuser credentials:

1. Retrieve the K8ssandra superuser name:

    ```bash
    kubectl get secret k8ssandra-superuser -o jsonpath="{.data.username}" | base64 --decode ; echo
    ```

    **Output**:

    ```bash
    k8ssandra-superuser
    ```

1. Retrieve the K8ssandra superuser password:

    ```bash
    kubectl get secret k8ssandra-superuser -o jsonpath="{.data.password}" | base64 --decode ; echo
    ```

    **Output**:

    ```bash
    PGo8kROUgAJOa8vhjQrE49Lgruw7s32HCPyVvcfVmmACW8oUhfoO9A
    ```

{{% alert title="Tip" color="success" %}}
Save the superuser name and password for use in the [Quickstarts]({{< relref "/quickstarts" >}}), if you decide to follow those steps.
{{% /alert %}}

## Next steps

* If you're a developer, and you'd like to get started coding using CQL or Stargate, see the [Quickstart for developers]({{< relref "/quickstarts/developer" >}}).
* If you're a Site Reliability Engineer, and you'd like to explore the K8ssandra administration environment including monitoring and maintenance utilities, see the [Quickstart for Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer" >}}).

For details that are specific to cloud providers, see:

* K8ssandra installs on [Google Kubernetes Engine]({{< relref "/install/gke" >}}) (GKE)
* K8ssandra installs on [Amazon Elastic Kubernetes Service]({{< relref "/install/eks" >}}) (EKS)
