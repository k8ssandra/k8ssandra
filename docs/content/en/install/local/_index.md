---
title: "Install K8ssandra Operator on local K8s"
linkTitle: "Local"
no_list: true
weight: 1
description: "Details to install K8ssandra Operator on a local Kubernetes **kind** development environment."
---

This topic explains how to install K8ssandra Operator to manage Apache Cassandra&reg; in single- and multi-cluster Kubernetes. The examples use **kind** as a local (dev) Kubernetes environment.

## Introduction

In this topic for installs of K8ssandra Operator on your local development environment, using **kind**, we'll cover:

* [Prerequisites]({{< relref "#prerequisites" >}}): Required supporting software
* [Quick Start for single-cluster]({{< relref "#quick-start-for-a-single-cluster" >}}): Quick start K8ssandra Operator install in a single-cluster Kubernetes
* [Quick Start for multi-cluster]({{< relref "#quick-start-for-multi-cluster" >}}): Quick start K8ssandra Operator install in a multi-cluster Kubernetes, with examples for a control plane and three data planes 
* [Helm]({{< relref "#helm" >}}): Install via a single K8ssandra Operator Helm chart 
* [Kustomize]({{< relref "#kustomize" >}}): Install via Kustomize - a declarative approach 
* [Next steps]({{< relref "#next-steps" >}}): Quick start info including helper scripts

## Prerequisites

Make sure you have the following installed before going through this topic. 

* [kind](#kind)
* [kubectx](#kubectx)
* [yq (YAML processor)](#yq)
* [setup-kind-multicluster.sh](#setup-kind-multiclustersh)
* [create-clientconfig.sh](#create-clientconfigsh)

You'll also need [kubectl](https://kubernetes.io/docs/tasks/tools/) and [helm v3+](https://helm.sh/docs/intro/install/) on your preferred OS. 

### **kind**

The examples in this topic use [kind](https://kind.sigs.k8s.io/) clusters. Install it now if you have not already done so.

By default, kind clusters run on the same Docker network, which means we will have routable pod IPs across clusters.

{{% alert title="Note" color="success" %}}
Issues creating multiple kind clusters have been observed on various versions of Docker Desktop for macOS.  These issues seem to be resolved with the 4.5.0 release of Docker Desktop.  Please be sure to upgrade Docker Desktop if you plan to deploy using kind. Other options for local dev K8s environments include [minikube](https://minikube.sigs.k8s.io/docs/start/) or [K3D](https://k3d.io/v5.3.0/). 
{{% /alert %}}


### **kubectx**

[kubectx](https://github.com/ahmetb/kubectx) is a really handy tool when you are dealing with multiple clusters. The examples will use it so go ahead and install it now.

### **yq**

[yq](https://github.com/mikefarah/yq#install) is lightweight and portable command-line YAML processor.

### setup-kind-multicluster.sh

[setup-kind-multicluster.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/setup-kind-multicluster.sh) lives in the k8ssandra-operator repo. It is used extensively during development and testing. Not only does it configure and create kind clusters, it also generates kubeconfig files for each cluster.

{{% alert title="Tip" color="success" %}}
kind generates a kubeconfig with the IP address of the API server set to 
localhost because the cluster is intended for local development. We need a kubeconfig with the IP address set to the internal address of the api server. `setup-kind-mulitcluster.sh` takes care of this for us.
{{% /alert %}}

### create-clientconfig.sh

[create-clientconfig.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/create-clientconfig.sh) lives in the k8ssandra-operator repo. It is used to configure access to remote clusters. 

## Quick start for a single-cluster

Deploy K8ssandra with one Cassandra datacenter in a **single-cluster** kind environment.

### Clone the repo and use the setup script

If you haven't already, clone the https://github.com/k8ssandra/k8ssandra-operator repo to your local machine where you're already running a kind cluster. Example:

```bash
cd ~/github
git clone https://github.com/k8ssandra/k8ssandra-operator.git
cd k8ssandra-operator
```

Invoke `make` with the following parameters: 

```bash
make NUM_CLUSTERS=1 create-kind-multicluster
```

**Output:**

```bash
scripts/setup-kind-multicluster.sh --clusters 1 --kind-worker-nodes 4
Creating 1 clusters...
Creating cluster 1 out of 1
Creating cluster "k8ssandra-0" ...
 ✓ Ensuring node image (kindest/node:v1.22.4) 🖼
 ✓ Preparing nodes 📦 📦 📦 📦 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
Set kubectl context to "kind-k8ssandra-0"
You can now use your cluster with:

kubectl cluster-info --context kind-k8ssandra-0

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
Error response from daemon: endpoint with name kind-registry already exists in network kind

Generating kubeconfig
Generating in-cluster kubeconfig
```

Verify the depoyment:

```bash
 kubectl get nodes 
```

**Output:**

```bash
NAME                        STATUS   ROLES                  AGE   VERSION
k8ssandra-0-control-plane   Ready    control-plane,master   80s   v1.22.4
k8ssandra-0-worker          Ready    <none>                 42s   v1.22.4
k8ssandra-0-worker2         Ready    <none>                 42s   v1.22.4
k8ssandra-0-worker3         Ready    <none>                 42s   v1.22.4
k8ssandra-0-worker4         Ready    <none>                 42s   v1.22.4
```

### Deploy cert-manager

Update your helm repo and set the context:

```bash
helm repo add jetstack https://charts.jetstack.io

helm repo update

kubectx kind-k8ssandra-0
```

**The output includes:**

```bash
Switched to context "kind-k8ssandra-0".
```

Now install the `jetstack/cert-manager`:

```bash
helm install cert-manager jetstack/cert-manager \
     --namespace cert-manager --create-namespace --set installCRDs=true
```

**Output:**

```bash
NAME: cert-manager
LAST DEPLOYED: Mon Jan 31 12:29:43 2022
NAMESPACE: cert-manager
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
cert-manager v1.7.0 has been deployed successfully!

In order to begin issuing certificates, you will need to set up a ClusterIssuer
or Issuer resource (for example, by creating a 'letsencrypt-staging' issuer).

More information on the different types of issuers and how to configure them
can be found in our documentation:

https://cert-manager.io/docs/configuration/

For information on how to configure cert-manager to automatically provision
Certificates for Ingress resources, take a look at the `ingress-shim`
documentation:

https://cert-manager.io/docs/usage/ingress/
```

### Deploy K8ssandra Operator

You can install K8ssandra Operator for namespace-scoped operation (default), or cluster-scoped operation. The example in this section shows K8ssandra Operator installed as namespace-scoped:

```bash
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace
```

**Output:**

```bash
NAME: k8ssandra-operator
LAST DEPLOYED: Mon Jan 31 12:30:40 2022
NAMESPACE: k8ssandra-operator
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

{{% alert title="Tip" color="success" %}}
Optionally, you can use `--set global.clusterScoped=true` to install K8ssandra Operator cluster-scoped:

```bash
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \ 
     --set global.clusterScoped=true --create-namespace
```
{{% /alert %}}

### Verify the deployment

```bash
kubectl get pods -n k8ssandra-operator
```

**Output:**

```bash
NAME                                                READY   STATUS    RESTARTS   AGE
k8ssandra-operator-7f76579f94-7s2tw                 1/1     Running   0          60s
k8ssandra-operator-cass-operator-794f65d9f4-j9lm5   1/1     Running   0          60s
```

### Deploy the K8ssandraCluster

To deploy a `K8ssandraCluster`, we use a custom YAML file. In this example, k8c1.yml. Notice, there is just one datacenter, `dc1`.

```yaml
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
```

Apply the YAML to the already deployed K8ssandra Operator:

```bash
kubectl apply -n k8ssandra-operator -f k8c1.yml
```

**Output:**

```bash
k8ssandracluster.k8ssandra.io/demo created
```

### Verify pod deployment

```bash
$ kubectl get pods -n k8ssandra-operator
```

**Output:**

```
NAME                                                    READY   STATUS    RESTARTS   AGE
demo-dc1-default-stargate-deployment-7b6c9d8dcd-k65jx   1/1     Running   0          5m33s
demo-dc1-default-sts-0                                  2/2     Running   0          10m
demo-dc1-default-sts-1                                  2/2     Running   0          10m
demo-dc1-default-sts-2                                  2/2     Running   0          10m
k8ssandra-operator-7f76579f94-7s2tw                     1/1     Running   0          11m
k8ssandra-operator-cass-operator-794f65d9f4-j9lm5       1/1     Running   0          11m
```

### Verify `K8ssandraCluster` deployment

```bash
kubectl get k8cs -n k8ssandra-operator
```

**Output:**

```bash
NAME   AGE
demo   8m22s
```

```bash
kubectl describe k8cs demo -n k8ssandra-operator
```

**Output:**

```bash
Name:         demo
Namespace:    k8ssandra-operator
Labels:       <none>
Annotations:  k8ssandra.io/system-replication: {"datacenters":["dc1"],"replicationFactor":3}
API Version:  k8ssandra.io/v1alpha1
Kind:         K8ssandraCluster
Metadata:
  Creation Timestamp:  2022-01-31T17:32:18Z
  Finalizers:
    k8ssandracluster.k8ssandra.io/finalizer
  Generation:  2
  Managed Fields:
    API Version:  k8ssandra.io/v1alpha1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
      f:spec:
        .:
        f:auth:
        f:cassandra:
          .:
          f:datacenters:
          f:jmxInitContainerImage:
            .:
            f:name:
            f:registry:
            f:tag:
          f:serverVersion:
    Manager:      kubectl-client-side-apply
    Operation:    Update
    Time:         2022-01-31T17:32:18Z
    API Version:  k8ssandra.io/v1alpha1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          f:k8ssandra.io/system-replication:
        f:finalizers:
          .:
          v:"k8ssandracluster.k8ssandra.io/finalizer":
      f:spec:
        f:cassandra:
          f:superuserSecretRef:
            .:
            f:name:
    Manager:      manager
    Operation:    Update
    Time:         2022-01-31T17:32:18Z
    API Version:  k8ssandra.io/v1alpha1
    Fields Type:  FieldsV1
    fieldsV1:
      f:status:
        .:
        f:conditions:
        f:datacenters:
          .:
          f:dc1:
            .:
            f:cassandra:
              .:
              f:cassandraOperatorProgress:
              f:conditions:
              f:lastServerNodeStarted:
              f:nodeStatuses:
                .:
                f:demo-dc1-default-sts-0:
                  .:
                  f:hostID:
                f:demo-dc1-default-sts-1:
                  .:
                  f:hostID:
                f:demo-dc1-default-sts-2:
                  .:
                  f:hostID:
              f:observedGeneration:
              f:quietPeriod:
              f:superUserUpserted:
              f:usersUpserted:
            f:stargate:
              .:
              f:availableReplicas:
              f:conditions:
              f:deploymentRefs:
              f:progress:
              f:readyReplicas:
              f:readyReplicasRatio:
              f:replicas:
              f:serviceRef:
              f:updatedReplicas:
    Manager:         manager
    Operation:       Update
    Subresource:     status
    Time:            2022-01-31T17:37:52Z
  Resource Version:  3385
  UID:               bee3e4c9-59df-486c-b5ac-c83b65162b2c
Spec:
  Auth:  true
  Cassandra:
    Datacenters:
      Config:
        Jvm Options:
          Heap Size:  512M
      Jmx Init Container Image:
        Name:      busybox
        Registry:  docker.io
        Tag:       1.34.1
      Metadata:
        Name:  dc1
      Size:    3
      Stargate:
        Allow Stargate On Data Nodes:  false
        Container Image:
          Registry:       docker.io
          Repository:     stargateio
          Tag:            v1.0.45
        Heap Size:        256M
        Service Account:  default
        Size:             1
      Storage Config:
        Cassandra Data Volume Claim Spec:
          Access Modes:
            ReadWriteOnce
          Resources:
            Requests:
              Storage:         5Gi
          Storage Class Name:  standard
    Jmx Init Container Image:
      Name:          busybox
      Registry:      docker.io
      Tag:           1.34.1
    Server Version:  4.0.1
    Superuser Secret Ref:
      Name:  demo-superuser
Status:
  Conditions:
    Last Transition Time:  2022-01-31T17:37:04Z
    Status:                True
    Type:                  CassandraInitialized
  Datacenters:
    dc1:
      Cassandra:
        Cassandra Operator Progress:  Ready
        Conditions:
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  False
          Type:                    ScalingUp
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  False
          Type:                    Stopped
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  False
          Type:                    ReplacingNodes
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  False
          Type:                    Updating
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  False
          Type:                    RollingRestart
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  False
          Type:                    Resuming
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  False
          Type:                    ScalingDown
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  True
          Type:                    Valid
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  True
          Type:                    Initialized
          Last Transition Time:    2022-01-31T17:37:00Z
          Message:
          Reason:
          Status:                  True
          Type:                    Ready
        Last Server Node Started:  2022-01-31T17:35:39Z
        Node Statuses:
          demo-dc1-default-sts-0:
            Host ID:  61dfa8cc-2a8b-4e8f-ae82-01c51833e0ba
          demo-dc1-default-sts-1:
            Host ID:  369aa179-d96e-4f21-a893-f6e6dc84b396
          demo-dc1-default-sts-2:
            Host ID:          bbdb6a9a-063b-4565-9704-f4caa6fd80f1
        Observed Generation:  1
        Quiet Period:         2022-01-31T17:37:06Z
        Super User Upserted:  2022-01-31T17:37:00Z
        Users Upserted:       2022-01-31T17:37:00Z
      Stargate:
        Available Replicas:  1
        Conditions:
          Last Transition Time:  2022-01-31T17:37:48Z
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

### Extract credentials

Use the following commands to extract the username and password:

```bash
CASS_USERNAME=$(kubectl get secret demo-superuser -n k8ssandra-operator -o=jsonpath='{.data.username}' | base64 --decode)

echo $CASS_USERNAME
```

**Output:**
```bash
demo-superuser
```

Now obtain the password secret:

```bash
CASS_PASSWORD=$(kubectl get secret demo-superuser -n k8ssandra-operator -o=jsonpath='{.data.password}' | base64 --decode)

echo $CASS_PASSWORD
```

**Output example - your value will be different:**
```bash
ACK7dO9qpsghIme-wvfI
```
{{% alert title="Tip" color="success" %}}
You'll use the extract credentials for subsequent authentication in deployed containers.
{{% /alert %}}

### Verify cluster status

```bash
kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
```

**Output plus nodetool example:**

```bash
Defaulted container "cassandra" out of: cassandra, server-system-logger, jmx-credentials (init), server-config-init (init)
cassandra@demo-dc1-default-sts-0:/$ nodetool -u demo-superuser -pw ACK7dO9qpsghIme-wvfI status
Datacenter: dc1
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load       Tokens  Owns (effective)  Host ID                               Rack
UN  10.244.4.4  96.71 KiB  16      100.0%            bbdb6a9a-063b-4565-9704-f4caa6fd80f1  default
UN  10.244.2.5  96.72 KiB  16      100.0%            61dfa8cc-2a8b-4e8f-ae82-01c51833e0ba  default
UN  10.244.3.4  96.63 KiB  16      100.0%            369aa179-d96e-4f21-a893-f6e6dc84b396  default
```

{{% alert title="Tip" color="success" %}}
All nodes should have the status UN, which stands for "Up Normal".
{{% /alert %}}

### Test a few operations

```bash
kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
```

**Output and cqlsh example:**

In the launched container's `cqlsh` session, notice we provide the extracted password for `demo-superuser`.

```bash
Defaulted container "cassandra" out of: cassandra, server-system-logger, jmx-credentials (init), server-config-init (init)
cassandra@demo-dc1-default-sts-0:/$ cqlsh -u demo-superuser -p ACK7dO9qpsghIme-wvfI
Connected to demo at 127.0.0.1:9042
[cqlsh 6.0.0 | Cassandra 4.0.1 | CQL spec 3.4.5 | Native protocol v5]
Use HELP for help.
demo-superuser@cqlsh> CREATE KEYSPACE test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};
demo-superuser@cqlsh> USE test;
demo-superuser@cqlsh:test> CREATE TABLE users (email text primary key, name text, state text);
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('john@gamil.com', 'John Smith', 'NC');
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('joe@gamil.com', 'Joe Jones', 'VA');
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('sue@help.com', 'Sue Sas', 'CA');
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('tom@yes.com', 'Tom and Jerry', 'NV');
demo-superuser@cqlsh:test> select * from users;

 email          | name          | state
----------------+---------------+-------
 john@gamil.com |    John Smith |    NC
  joe@gamil.com |     Joe Jones |    VA
   sue@help.com |       Sue Sas |    CA
    tom@yes.com | Tom and Jerry |    NV

(4 rows)
```

Now test an operation via the open-source Stargate API.

```bash
kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
```

**Output plus cqlsh &amp; stargate-service example:**

```bash
Defaulted container "cassandra" out of: cassandra, server-system-logger, server-config-init (init)
cassandra@k8ssandra-3-worker:/$ ping demo-dc3-stargate-service
cassandra@demo-dc1-default-sts-0:/$ cqlsh -u demo-superuser -p ACK7dO9qpsghIme-wvfI demo-dc1-stargate-service
Connected to demo at demo-dc1-stargate-service:9042
[cqlsh 6.0.0 | Cassandra 4.0.1 | CQL spec 3.4.5 | Native protocol v4]
Use HELP for help.
demo-superuser@cqlsh> use test;
demo-superuser@cqlsh:test> select * from users;

 email          | name          | state
----------------+---------------+-------
 john@gamil.com |    John Smith |    NC
  joe@gamil.com |     Joe Jones |    VA
   sue@help.com |       Sue Sas |    CA
    tom@yes.com | Tom and Jerry |    NV

(4 rows)
```

## Quick start for multi-cluster

Follow these steps to deploy K8ssandra Operator with multiple Cassandra datacenters in a **multi-cluster** kind environment.

### Clone the repo and use the setup script

If you haven't already, clone the https://github.com/k8ssandra/k8ssandra-operator repo to your local machine where you're already running a kind cluster. Example:

```bash
cd ~/github
git clone https://github.com/k8ssandra/k8ssandra-operator.git
cd k8ssandra-operator
```

Invoke `make` with the following parameters: 

```bash
make NUM_CLUSTERS=4 create-kind-multicluster
```

### Verify the deployments 

Set the context to each of the four created clusters, and get node information for each cluster. Examples:

```bash
kubectx kind-k8ssandra-0
```

**Output:**

```bash
Switched to context "kind-k8ssandra-0".
```

Then enter:

```bash
kubectl get nodes
```

**Output:**

```bash
NAME                        STATUS   ROLES                  AGE     VERSION
k8ssandra-0-control-plane   Ready    control-plane,master   5h14m   v1.22.1
k8ssandra-0-worker          Ready    <none>                 5h14m   v1.22.1
k8ssandra-0-worker2         Ready    <none>                 5h14m   v1.22.1
k8ssandra-0-worker3         Ready    <none>                 5h14m   v1.22.1
k8ssandra-0-worker4         Ready    <none>                 5h14m   v1.22.1
```

Then enter:

```bash
kubectx kind-k8ssandra-1

kubectl get nodes
```

**Output:**

```bash
NAME                        STATUS   ROLES                  AGE     VERSION
k8ssandra-1-control-plane   Ready    control-plane,master   5h13m   v1.22.1
k8ssandra-1-worker          Ready    <none>                 5h13m   v1.22.1
k8ssandra-1-worker2         Ready    <none>                 5h13m   v1.22.1
k8ssandra-1-worker3         Ready    <none>                 5h13m   v1.22.1
k8ssandra-1-worker4         Ready    <none>                 5h13m   v1.22.1
```

Then enter:

```bash
kubectx kind-k8ssandra-2

kubectl get nodes
```

**Output:**

```bash
NAME                        STATUS   ROLES                  AGE     VERSION
k8ssandra-2-control-plane   Ready    control-plane,master   5h13m   v1.22.1
k8ssandra-2-worker          Ready    <none>                 5h12m   v1.22.1
k8ssandra-2-worker2         Ready    <none>                 5h12m   v1.22.1
k8ssandra-2-worker3         Ready    <none>                 5h12m   v1.22.1
k8ssandra-2-worker4         Ready    <none>                 5h12m   v1.22.1
```

Then enter:

```bash
kubectx kind-k8ssandra-3

kubectl get nodes
```

**Output:**

```bash
NAME                        STATUS   ROLES                  AGE     VERSION
k8ssandra-3-control-plane   Ready    control-plane,master   5h12m   v1.22.1
k8ssandra-3-worker          Ready    <none>                 5h12m   v1.22.1
k8ssandra-3-worker2         Ready    <none>                 5h12m   v1.22.1
k8ssandra-3-worker3         Ready    <none>                 5h12m   v1.22.1
k8ssandra-3-worker4         Ready    <none>                 5h12m   v1.22.1
```

### Install cert-manager in each cluster

If you haven't already, update your helm repo with the jetstack cert-manager. 

```bash
helm repo add jetstack https://charts.jetstack.io

helm repo update
```

Set the per-cluster context and install `jetstack/cert-manager`. Examples:

```bash
kubectx kind-k8ssandra-0

helm install cert-manager jetstack/cert-manager --namespace cert-manager \
     --create-namespace --set installCRDs=true
```

**Output:**

```bash
NAME: cert-manager
LAST DEPLOYED: Thu Jan 27 15:28:59 2022
NAMESPACE: cert-manager
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
cert-manager v1.7.0 has been deployed successfully!
```

Then enter:

```bash
kubectx kind-k8ssandra-1

helm install cert-manager jetstack/cert-manager --namespace cert-manager \
     --create-namespace --set installCRDs=true
```

**Output:**

```bash
NAME: cert-manager
LAST DEPLOYED: Thu Jan 27 15:28:59 2022
NAMESPACE: cert-manager
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
cert-manager v1.7.0 has been deployed successfully!
```

Then enter:

```bash
kubectx kind-k8ssandra-2

helm install cert-manager jetstack/cert-manager --namespace cert-manager \
     --create-namespace --set installCRDs=true
```

**Output:**

```bash
NAME: cert-manager
LAST DEPLOYED: Thu Jan 27 15:28:59 2022
NAMESPACE: cert-manager
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
cert-manager v1.7.0 has been deployed successfully!
```

Then enter:

```bash
kubectx kind-k8ssandra-3

helm install cert-manager jetstack/cert-manager --namespace cert-manager \
     --create-namespace --set installCRDs=true
```

**Output:**

```bash
NAME: cert-manager
LAST DEPLOYED: Thu Jan 27 15:28:59 2022
NAMESPACE: cert-manager
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
cert-manager v1.7.0 has been deployed successfully!
```

### Install K8ssandra Operator in the control-plane

In this example, of the four clusters we've created, we'll use `kind-k8ssandra-0` as our control-plane.

```bash
kubectx kind-k8ssandra-0

helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \
     --create-namespace
```

### Install K8ssandra Operator in the data-planes

In this example, we'll use the three other clusters as data-planes.

```bash
kubectx kind-k8ssandra-1
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \
 --create-namespace

kubectx kind-k8ssandra-2
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \
 --create-namespace

kubectx kind-k8ssandra-3
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \
 --create-namespace
```

### Verify control-plane configuration

```
kubectx kind-k8ssandra-0

kubectl -n k8ssandra-operator get deployment k8ssandra-operator \
 -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

**Output:**

```bash
true
```

### Verify data-plane configuration

We could test for `K8SSANDRA_CONTROL_PLANE`, which for each of the three clusters in our example serving as data-planes, should return `false`. Just one example:

```
kubectx kind-k8ssandra-1

kubectl -n k8ssandra-operator get deployment k8ssandra-operator \
 -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

**Output:**

```bash
false
```

### Generate and install ClientConfigs

[create-clientconfig.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/create-clientconfig.sh) lives in the k8ssandra-operator repo. It is used to configure access to remote clusters. 

First, set the context to `kind-k8ssandra-0`, the control plane cluster. 

```bash
kubectx kind-k8ssandra-0
```

Run the create-clientconfig.sh script, once per data plane cluster.  

```bash
./scripts/create-clientconfig.sh --namespace k8ssandra-operator \
  --src-kubeconfig build/kubeconfigs/k8ssandra-1.yaml \
  --dest-kubeconfig build/kubeconfigs/k8ssandra-0.yaml \
  --in-cluster-kubeconfig build/kubeconfigs/updated/k8ssandra-1.yaml \
  --output-dir clientconfig
```

**Output:**

```bash
Creating clientconfig/kubeconfig
Creating secret kind-k8ssandra-1-config
Error from server (NotFound): secrets "kind-k8ssandra-1-config" not found
secret/kind-k8ssandra-1-config created
Creating ClientConfig clientconfig/kind-k8ssandra-1.yaml
clientconfig.config.k8ssandra.io/kind-k8ssandra-1 created
```

Then enter:

```bash
./scripts/create-clientconfig.sh --namespace k8ssandra-operator \
 --src-kubeconfig build/kubeconfigs/k8ssandra-2.yaml \
 --dest-kubeconfig build/kubeconfigs/k8ssandra-0.yaml \
 --in-cluster-kubeconfig build/kubeconfigs/updated/k8ssandra-2.yaml 
 --output-dir clientconfig
```

**Output:**

```bash
Creating clientconfig/kubeconfig
Creating secret kind-k8ssandra-2-config
Error from server (NotFound): secrets "kind-k8ssandra-2-config" not found
secret/kind-k8ssandra-2-config created
Creating ClientConfig clientconfig/kind-k8ssandra-2.yaml
clientconfig.config.k8ssandra.io/kind-k8ssandra-2 created
```

Then enter:

```bash
./scripts/create-clientconfig.sh --namespace k8ssandra-operator \
 --src-kubeconfig build/kubeconfigs/k8ssandra-3.yaml \
 --dest-kubeconfig build/kubeconfigs/k8ssandra-0.yaml \
 --in-cluster-kubeconfig build/kubeconfigs/updated/k8ssandra-3.yaml \
 --output-dir clientconfig
```

**Output:**

```bash
Creating clientconfig/kubeconfig
Creating secret kind-k8ssandra-3-config
Error from server (NotFound): secrets "kind-k8ssandra-3-config" not found
secret/kind-k8ssandra-3-config created
Creating ClientConfig clientconfig/kind-k8ssandra-3.yaml
clientconfig.config.k8ssandra.io/kind-k8ssandra-3 created
```

### Deploy the K8ssandraCluster

To deploy the `K8ssandraCluster`, we use a custom YAML file. In this example, k8cm1.yml. Notice, there are three Cassandra 4.0.1 datacenters, `dc1`, `dc2`, and `dc3` that are associated with the three data plane clusters.

```yaml
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
        k8sContext: kind-k8ssandra-1
        size: 3
        stargate:
          size: 1
          heapSize: 256M
      - metadata:
          name: dc2
        k8sContext: kind-k8ssandra-2
        size: 3
        stargate:
          size: 1
          heapSize: 256M
      - metadata:
          name: dc3
        k8sContext: kind-k8ssandra-3
        size: 3
        stargate:
          size: 1
          heapSize: 256M
```

Verify again that your context is set to the control plane cluster, which is in this example:

```bash
kubectx kind-k8ssandra-0
```

Apply the YAML to the already deployed K8ssandra Operator. 

```bash
kubectl apply -n k8ssandra-operator -f k8cm1.yml
```

### Verify pod deployment

Initially the rollout will begin in dc1 and work across the full cluster:

```bash
kubectx kind-k8ssandra-1

kubectl get pods -n k8ssandra-operator
```

Do the same on each of the other two clusters by setting the kubectx context to kind-k8ssandra-2, check the pods status; then kind-k8ssandra-3, and check the pods status.

Eventually the datacenters will be fully deployed:

```bash
kubectx kind-k8ssandra-0

kubectl get pods -n k8ssandra-operator
```

**Output:**

```bash
NAME                                                READY   STATUS    RESTARTS   AGE
k8ssandra-operator-68568ffbd5-l6t2f                 1/1     Running   0          93m
k8ssandra-operator-cass-operator-794f65d9f4-kqrpf   1/1     Running   0          97m
```

```bash
kubectx kind-k8ssandra-1

kubectl get pods -n k8ssandra-operator
```

**Output:**

```bash
NAME                                                    READY   STATUS    RESTARTS   AGE
demo-dc1-default-stargate-deployment-547df5877d-bvnz2   1/1     Running   0          66m
demo-dc1-default-sts-0                                  2/2     Running   0          80m
demo-dc1-default-sts-1                                  2/2     Running   0          80m
demo-dc1-default-sts-2                                  2/2     Running   0          80m
k8ssandra-operator-7cfd7977cb-wxww5                     1/1     Running   0          97m
k8ssandra-operator-cass-operator-794f65d9f4-s697p       1/1     Running   0          97m
```

```bash
kubectx kind-k8ssandra-2

kubectl get pods -n k8ssandra-operator
```

**Output:**

```bash
NAME                                                    READY   STATUS    RESTARTS   AGE
demo-dc2-default-stargate-deployment-86c5fc44ff-lt9ts   1/1     Running   0          65m
demo-dc2-default-sts-0                                  2/2     Running   0          76m
demo-dc2-default-sts-1                                  2/2     Running   0          76m
demo-dc2-default-sts-2                                  2/2     Running   0          76m
k8ssandra-operator-7cfd7977cb-59nld                     1/1     Running   0          96m
k8ssandra-operator-cass-operator-794f65d9f4-79z6z       1/1     Running   0          96m
```

```bash
kubectx kind-k8ssandra-3

kubectl get pods -n k8ssandra-operator
```

**Output:**

```bash
NAME                                                   READY   STATUS    RESTARTS   AGE
demo-dc3-default-stargate-deployment-6bd8f87b4-ztxb8   1/1     Running   0          65m
demo-dc3-default-sts-0                                 2/2     Running   0          71m
demo-dc3-default-sts-1                                 2/2     Running   0          71m
demo-dc3-default-sts-2                                 2/2     Running   0          71m
k8ssandra-operator-7cfd7977cb-g6hcz                    1/1     Running   0          96m
k8ssandra-operator-cass-operator-794f65d9f4-prfd8      1/1     Running   0          96m
```

### Verify K8ssandraCluster status

While deployment is still in progress, you can check the status:

```bash
kubectx kind-k8ssandra-0

kubectl describe k8cs demo -n k8ssandra-operator
```

In the **earlier** deployment phases, you may notice statuses such as:

**Output:**

```bash
   .
   .
   .
dc1:
      Cassandra:
        Cassandra Operator Progress:  Ready
        Conditions:
          Last Transition Time:    2022-01-31T19:02:40Z
          Message:
          Reason:
          Status:                  False
          Type:                    ScalingUp
          Last Transition Time:    2022-01-31T19:02:40Z
          Message:
          Reason:
          Status:                  False
          Type:                    Stopped
          Last Transition Time:    2022-01-31T19:02:40Z
          Message:
          Reason:
          Status:                  False
          Type:                    ReplacingNodes
          Last Transition Time:    2022-01-31T19:02:40Z
          Message:
          Reason:
   .
   .
   .
```

This behavior is expected for the deployments-in-progress. If you're curious, you can continue to check status bu submitting the command again. When the deployments have been completed, for example, here's the command again and a portion of its completed output:

```bash
kubectl describe k8cs demo -n k8ssandra-operator
```

**Output subset:**

```bash
   .
   .
   .
dc3:
      Cassandra:
        Cassandra Operator Progress:  Ready
        Conditions:
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  False
          Type:                    ScalingUp
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  False
          Type:                    Stopped
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  False
          Type:                    ReplacingNodes
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  False
          Type:                    Updating
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  False
          Type:                    RollingRestart
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  False
          Type:                    Resuming
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  False
          Type:                    ScalingDown
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  True
          Type:                    Valid
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  True
          Type:                    Initialized
          Last Transition Time:    2022-01-31T19:12:33Z
          Message:
          Reason:
          Status:                  True
          Type:                    Ready
        Last Server Node Started:  2022-01-31T19:11:12Z
        Node Statuses:
          demo-dc3-default-sts-0:
            Host ID:  2cceff49-6df2-4045-8e04-6ce262bd6fc4
          demo-dc3-default-sts-1:
            Host ID:  018bfd17-8a77-43c6-859b-ab69c1fc8a66
          demo-dc3-default-sts-2:
            Host ID:          38438f65-b10a-4b3f-a56f-488536bf4cd3
        Observed Generation:  1
        Quiet Period:         2022-01-31T19:12:39Z
        Super User Upserted:  2022-01-31T19:12:34Z
        Users Upserted:       2022-01-31T19:12:34Z
Events:                       <none>
```

### Extract credentials

On the control plane, use the following commands to extract the username and password.

```bash
kubectx kind-k8ssandra-0

CASS_USERNAME=$(kubectl get secret demo-superuser -n k8ssandra-operator -o=jsonpath='{.data.username}' | base64 --decode)

echo $CASS_USERNAME
```

**Output:**
```bash
demo-superuser
```

Now obtain the password secret:

```bash
CASS_PASSWORD=$(kubectl get secret demo-superuser -n k8ssandra-operator -o=jsonpath='{.data.password}' | base64 --decode)

echo $CASS_PASSWORD
```

**Output example - your value will be different:**
```bash
KT-ROFfbD-O9BzWS3Lxq
```

{{% alert title="Tip" color="success" %}}
You'll use the extract credentials for subsequent authentication in deployed containers.
{{% /alert %}}

### Verify cluster status

On one of the data plane clusters, verify the cluster status. Example:

```bash
kubectx kind-k8ssandra-1

$ kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
```

**Output plus nodetool example:**

```bash
Defaulted container "cassandra" out of: cassandra, server-system-logger, jmx-credentials (init), server-config-init (init)
cassandra@k8ssandra-1-worker2:/$ nodetool -u demo-superuser -pw KT-ROFfbD-O9BzWS3Lxq status
Datacenter: dc1
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address      Load        Tokens  Owns (effective)  Host ID                               Rack
UN  172.18.0.8   106.28 KiB  16      100.0%            676592e3-72f2-40c7-893e-932ab7f99fdd  default
UN  172.18.0.9   101.93 KiB  16      100.0%            a6b62d88-0223-432d-a040-e10ce93a3342  default
UN  172.18.0.7   106.26 KiB  16      100.0%            1228c0ec-b5f3-4bbc-9b32-e32a8600e58b  default

Datacenter: dc2
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address      Load        Tokens  Owns (effective)  Host ID                               Rack
UN  172.18.0.14  106.2 KiB   16      100.0%            0836237c-66d7-453a-b03a-91d64aaa0f77  default
UN  172.18.0.12  106.18 KiB  16      100.0%            592a453f-7abb-483c-ab4e-2f37c8cd4713  default
UN  172.18.0.16  106.2 KiB   16      100.0%            d19b5028-b28e-450b-a417-1cd2e42aa179  default

Datacenter: dc3
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address      Load        Tokens  Owns (effective)  Host ID                               Rack
UN  172.18.0.21  101.74 KiB  16      100.0%            2cceff49-6df2-4045-8e04-6ce262bd6fc4  default
UN  172.18.0.20  106.11 KiB  16      100.0%            38438f65-b10a-4b3f-a56f-488536bf4cd3  default
UN  172.18.0.19  106.1 KiB   16      100.0%            018bfd17-8a77-43c6-859b-ab69c1fc8a66  default
```

{{% alert title="Tip" color="success" %}}
All nodes should have the status UN, which stands for "Up Normal".
{{% /alert %}}

### Test a few operations

```bash
kubectx kind-k8ssandra-1

kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
```

**Output and cqlsh example:**

In the launched container's `cqlsh` session, notice we provide the extracted password for `demo-superuser`.

```bash
Defaulted container "cassandra" out of: cassandra, server-system-logger, server-config-init (init)
cassandra@k8ssandra-1-worker2:/$ cqlsh -u demo-superuser -p KT-ROFfbD-O9BzWS3Lxq
Connected to demo at 127.0.0.1:9042
[cqlsh 6.0.0 | Cassandra 4.0.1 | CQL spec 3.4.5 | Native protocol v5]
Use HELP for help.
demo-superuser@cqlsh> describe keyspaces;

data_endpoint_auth  system_auth         system_schema  system_views
system              system_distributed  system_traces  system_virtual_schema

demo-superuser@cqlsh> CREATE KEYSPACE test WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1' : 3, 'dc2' : 3, 'dc3': 3};
demo-superuser@cqlsh> USE test;
demo-superuser@cqlsh:test> CREATE TABLE users (email text primary key, name text, state text);
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('john@gamil.com', 'John Smith', 'NC');
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('joe@gamil.com', 'Joe Jones', 'VA');
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('sue@help.com', 'Sue Sas', 'CA');
demo-superuser@cqlsh:test> insert into users (email, name, state) values ('tom@yes.com', 'Tom and Jerry', 'NV');
demo-superuser@cqlsh:test> select * from users;

 email          | name          | state
----------------+---------------+-------
 john@gamil.com |    John Smith |    NC
  joe@gamil.com |     Joe Jones |    VA
   sue@help.com |       Sue Sas |    CA
    tom@yes.com | Tom and Jerry |    NV

(4 rows)
demo-superuser@cqlsh:test> exit
cassandra@k8ssandra-1-worker2:/$ exit
exit
```

Try another cqlsh operation on a different cluster.


```bash
kubectx kind-k8ssandra-3

kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
```

**Output:**

```bash
Defaulted container "cassandra" out of: cassandra, server-system-logger, jmx-credentials (init), server-config-init (init)
cassandra@k8ssandra-3-worker3:/$ cqlsh -u demo-superuser -p KT-ROFfbD-O9BzWS3Lxq
Connected to demo at 127.0.0.1:9042
[cqlsh 6.0.0 | Cassandra 4.0.1 | CQL spec 3.4.5 | Native protocol v5]
Use HELP for help.
demo-superuser@cqlsh> USE test;
demo-superuser@cqlsh:test> SELECT * FROM users;

 email          | name          | state
----------------+---------------+-------
 john@gamil.com |    John Smith |    NC
  joe@gamil.com |     Joe Jones |    VA
   sue@help.com |       Sue Sas |    CA
    tom@yes.com | Tom and Jerry |    NV

(4 rows)
demo-superuser@cqlsh:test> exit
cassandra@k8ssandra-3-worker3:/$ exit
exit
```

Now try using the Stargate API. 

**Output plus cqlsh &amp; stargate-service example:**

```bash
kubectx kind-k8ssandra-3

kubectl exec --stdin --tty demo-dc3-default-sts-0 -n k8ssandra-operator -- /bin/bash

Defaulted container "cassandra" out of: cassandra, server-system-logger, jmx-credentials (init), server-config-init (init)
cassandra@k8ssandra-3-worker3:/$ cqlsh -u demo-superuser -p KT-ROFfbD-O9BzWS3Lxq demo-dc3-stargate-service
Connected to demo at demo-dc3-stargate-service:9042
[cqlsh 6.0.0 | Cassandra 4.0.1 | CQL spec 3.4.5 | Native protocol v4]
Use HELP for help.
demo-superuser@cqlsh> use test;
demo-superuser@cqlsh:test> select * from users;

 email          | name          | state
----------------+---------------+-------
 john@gamil.com |    John Smith |    NC
  joe@gamil.com |     Joe Jones |    VA
   sue@help.com |       Sue Sas |    CA
    tom@yes.com | Tom and Jerry |    NV

(4 rows)
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
