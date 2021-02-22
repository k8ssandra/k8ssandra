---
title: "K8ssandra quick start"
linkTitle: "Quick start"
weight: 1
description: |
  Kick the tires and take K8ssandra for a spin!
---

Welcome to K8ssandra! This guide gets you up and running with a single-node Apache Cassandra&reg; cluster on Kubernetes. If you're interested in a more detailed component walkthroughs check out the [tasks]({{< ref "topics">}}) section.

**Completion time**: 15 to 20 minutes.

In this quick start, we'll cover the following topics:

* [Prerequisites]({{< relref "#prerequisites" >}}): Required supporting software including resource recommendations.
* [K8ssandra Helm repository configuration]({{< relref "#configure-the-k8ssandra-helm-repository" >}}): Accessing the Helm charts that install K8ssandra.
* [K8ssandra installation]({{< relref "#install-k8ssandra" >}}): Getting K8ssandra up and running using the Helm chart repo.
* [Verifying K8ssandra functionality]({{< relref "#verify-your-k8ssandra-installation" >}}): Making sure K8ssandra is working as expected.
* [Starting and stopping K8ssandra]({{< relref "#cassandra-operations" >}}): Cleanly stopping and restarting the K8ssandra pod.

## Prerequisites

In your local environment the following tools are required for provisioning a K8ssandra cluster.

* [Helm v3+](https://helm.sh/docs/intro/install/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

As K8ssandra deploys on a Kubernetes cluster, one must be available to target for installation. The K8 environment may be a local version running on your development machine, an on-premises self-hosted environment, or a managed cloud offering. To that end the cluster must be up and available to your `kubectl` command:

```bash
kubectl cluster-info
Kubernetes control plane is running at https://127.0.0.1:55017
KubeDNS is running at https://127.0.0.1:55017/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

If you do not have a Kubernetes cluster available, you can use [OpenShift CodeReady Containers](https://developers.redhat.com/products/codeready-containers/overview) that run within a VM, or one of the following local versions that run within Docker:

* [K3D](https://k3d.io/)
* [Minikube](https://minikube.sigs.k8s.io/docs/start/)
* [Kind](https://kind.sigs.k8s.io/)

The instructions in this section focus on the Docker container solutions above, but the general instructions should work for other environments as well.

### Kubernetes version support

K8ssandra works with the following versions of Kubernetes either standalone or via a cloud provider:

* 1.16
* 1.17
* 1.18
* 1.19
* 1.20

### Cassandra version support

K8ssandra can install the following versions of Apache Cassandra:

* 3.11.7
* 3.11.8
* 3.11.9
* 3.11.10
* 4.0-beta4

### Resource recommendations

The minimum recommended development configuration for a single Cassandra node is 8 gigs of RAM and 2 virtual processor cores (2 physical cores). Given that, we recommend a machine specification of **no less** than 16 gigs of RAM and 8 virtual processor cores (4 physical cores). You'll want to adjust your Docker resource preferences accordingly. For this quick start we're allocating 4 virtual processors and 8 gigs of RAM to the Docker environment.

An ideal environment, enabling you to run a 3 node Cassandra cluster, enforcing QUORUM consistency, would consist of 32 gigs of RAM and 12 virtual cores (6 physical cores).

See the documentation for your particular flavor of Docker for instructions on configuring resource limits.

### Example K8s container configuration

The following Minikube example creates a K8s cluster named `k8ssandra` running K8s version 1.18.16 with 4 virtual processor cores and 8 gigs of RAM:

```bash
minikube start --cpus=4 --memory='8128m' --kubernetes-version=1.18.16 -p k8ssandra
😄  [k8ssandra] minikube v1.17.1 on Darwin 11.2.1
✨  Automatically selected the docker driver. Other choices: hyperkit, ssh
👍  Starting control plane node k8ssandra in cluster k8ssandra
🔥  Creating docker container (CPUs=4, Memory=8128MB) ...
🐳  Preparing Kubernetes v1.18.16 on Docker 20.10.2 ...
    ▪ Generating certificates and keys ...
    ▪ Booting up control plane ...
    ▪ Configuring RBAC rules ...
🔎  Verifying Kubernetes components...
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "k8ssandra" cluster and "default" namespace by default
```

## Configure the K8ssandra Helm repository

K8ssandra is delivered via a collection of Helm charts for easy installation, so once you've got a suitable K8s environment configured, you'll need to add the K8ssandra Helm chart repositories.

To add the K8ssandra helm chart repos:

1. Add the main K8ssandra Helm chart repo:

    ```bash
    helm repo add k8ssandra https://helm.k8ssandra.io/
    ```

2. If you want to access K8ssandra services from outside of the Kubernetes cluster, also add the Traefik Ingress repo (highly recommended):

    ```bash
    helm repo add traefik https://helm.traefik.io/traefik
    ```

3. Finally, update your helm repository listing:

    ```bash
    helm repo update
    ```

{{% alert title="Tip" color="success" %}}
Alternatively, you can download the individual charts directly from the project's [releases](https://github.com/k8ssandra/k8ssandra/releases) page.
{{% /alert %}}

## Install K8ssandra

The K8ssandra helm charts make installation a snap. You can override chart configurations during installation as necessary if you're an advanced user, or make changes after a default installation using `helm upgrade` at a later time.

{{% alert title="Important" color="warning" %}}
K8ssandra comes out of the box with a set of default values tailored to getting up and running quickly.  Those defaults are intended to be a great starting point for smaller-scale local development but are **not** intended for production deployments.
{{% /alert %}}

To install a single node K8ssandra cluster:

1. Copy the following YAML to a file named `k8ssandra.yaml`:

    ```yaml
    cassandra:
      version: "3.11.7"
      clusterName: k8ssandra
      datacenters:
      - name: dc1
        size: 1
    kube-prometheus-stack:
      grafana:
        adminUser: admin
        adminPassword: admin123
    ```

    The configuration file creates a K8ssandra cluster named `k8ssandra` with a datacenter, `dc1` containing a single Cassandra node.

2. Use `helm install` to install K8ssandra, referring to the example configuration file:

    ```bash
    helm install -f k8ssandra.yaml k8ssandra k8ssandra/k8ssandra
    NAME: k8ssandra
    LAST DEPLOYED: Thu Feb 18 10:05:44 2021
    NAMESPACE: default
    STATUS: deployed
    REVISION: 1
    ```

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

Depending upon your K8 configuration, initialization of your K8ssandra installation can take a few minutes. To check the status of your K8ssandra deployment, use the `kubectl` command:

```bash
kubectl get pods
NAMESPACE     NAME                                                  READY   STATUS      RESTARTS   AGE
default       k8ssandra-cass-operator-6666588dc5-rrbwx              1/1     Running     0          15m
default       k8ssandra-dc1-default-sts-0                           2/2     Running     0          15m
default       k8ssandra-dc1-stargate-7db4dbfdd5-dxk5v               1/1     Running     0          15m
default       k8ssandra-grafana-b6f7978c4-cwdwg                     2/2     Running     0          15m
default       k8ssandra-kube-prometheus-operator-5556885bd6-2tn4t   1/1     Running     0          15m
default       k8ssandra-reaper-k8ssandra-64cc9d57d9-rc2z4           1/1     Running     0          12m
default       k8ssandra-reaper-k8ssandra-schema-pxdcz               0/1     Completed   0          12m
default       k8ssandra-reaper-operator-cc46fd5f4-hlmck             1/1     Running     0          15m
default       prometheus-k8ssandra-kube-prometheus-prometheus-0     2/2     Running     1          15m
```

NOTE: k8ssandra prefixes the pod names. Stargate takes 4 minutes to register as ready.




Once all the pods are in the `Running` or `Completed` state, you can check the health of your K8ssandra cluster.

{{% alert title="Tip" color="success" %}}
The pod `k8ssandra-reaper-k8ssandra-schema-xxxxx` runs once and does not persist.
{{% /alert %}}

The actual Cassandra node name from the listing above is `k8ssandra-dc1-default-sts-0`. If you've configured multiple nodes, `sts-0` will increment.

We'll use node name `k8ssandra-dc1-default-sts-0` throughout the following sections.

To check the health of your K8ssandra cluster:

1. Confirm that the Cassandra operator is `Ready`:

    ```bash
    kubectl describe CassandraDataCenter dc1 | grep "Cassandra Operator Progress:"
       Cassandra Operator Progress:  Ready
    ```

1. Get K8ssandra superuser credentials:

    * K8ssandra superuser name:

        ```bash
        kubectl get secret k8ssandra-superuser -o jsonpath="{.data.username}" | base64 --decode | more
        k8ssandra-superuser
        (END)
        ```

    * K8ssandra superuser password:

        ```bash
        kubectl get secret k8ssandra-superuser -o jsonpath="{.data.password}" | base64 --decode | more
        PGo8kROUgAJOa8vhjQrE49Lgruw7s32HCPyVvcfVmmACW8oUhfoO9A
        (END)
        ```

1. Run `nodetool status`, using the Cassandra node name `k8ssandra-dc1-default-sts-0`, and passing the superuser name and password. Verify that the node is in the state `UN` or Up Normal:

    ```bash
    kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- nodetool -u k8ssandra-superuser -pw PGo8kROUgAJOa8vhjQrE49Lgruw7s32HCPyVvcfVmmACW8oUhfoO9A status
    Datacenter: dc1
    ===============
    Status=Up/Down
    |/ State=Normal/Leaving/Joining/Moving
    --  Address      Load       Owns    Host ID                               Token                                    Rack
    UN  10.244.1.12  215.3 KiB  ?       75e52e51-edc9-49f8-84f6-f044999ac130  -1080085985719557225                     default

    Note: Non-system keyspaces don't have the same replication settings, effective ownership information is meaningless
    ```

{{% alert title="Tip" color="success" %}}
Save the superuser name and password for use in future steps.
{{% /alert %}}

## Stopping and starting Cassandra {#cassandra-operations}

Before shutting down your Kubernetes cluster, you'll want to make sure you cleanly shut down your Cassandra datacenters. You can do that using the `kubectl patch` command and setting the `spec:stopped` property to either `true` (stopped) or `false` (running).

### Shut down Cassandra

To shut down a Cassandra datacenter:

```bash
kubectl patch cassdc <datacenter-name> --type merge -p '{"spec":{"stopped":true}}'
```

Example:

```bash
kubectl patch cassdc dc1 --type merge -p '{"spec":{"stopped":true}}'
cassandradatacenter.cassandra.datastax.com/dc1 patched
```

### Start up Cassandra

To start up a Cassandra datacenter

```bash
kubectl patch cassdc <datacenter-name> --type merge -p '{"spec":{"stopped":false}}'
```

Example:

```bash
kubectl patch cassdc dc1 --type merge -p '{"spec":{"stopped":false}}'
cassandradatacenter.cassandra.datastax.com/dc1 patched
```

## Next

* For detailed information on additional K8ssandra tasks, see [Tasks]({{< relref "docs/topics" >}}).
* For a list of frequently asked questions, see the [FAQs]({{< relref "docs/faqs" >}}).
* For detailed information on K8ssandra, see [Architecture]({{< relref "docs/architecture" >}}).
* For information on the various K8ssandra Helm charts, see [Reference]({{< relref "docs/reference" >}}).
* If you'd like to contribute to K8ssandra, see [Contribution guidelines]({{< relref "docs/contribution-guidelines" >}}).
