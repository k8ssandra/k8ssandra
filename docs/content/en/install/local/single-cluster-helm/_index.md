---
title: "Single-cluster install with helm"
linkTitle: "Single-cluster/helm"
no_list: false
weight: 1
description: "Quickstart with a unified helm chart to install K8ssandraCluster in single-cluster Kubernetes."
---

This topic shows how you can use `helm` and other tools to install and configure the `K8ssandraCluster` custom resource in **single-cluster** local Kubernetes, using K8ssandra Operator. 

## Prerequisites

If you haven't already, see the install [prerequisites]({{< relref "install/local/" >}}).

## Quick start for a single-cluster

Deploy K8ssandra with one Cassandra datacenter in a **single-cluster** kind environment.

### Add the K8ssandra Helm chart repo

If you haven't already, add the main K8ssandra stable Helm chart repo:

```bash
helm repo add k8ssandra https://helm.k8ssandra.io/stable
helm repo update
```

### Clone the K8ssandra Operator's GitHub repo and use the setup script

Also clone the https://github.com/k8ssandra/k8ssandra-operator GitHub repo to your local machine where you're already running a kind cluster. Example:

```bash
cd ~/github
git clone https://github.com/k8ssandra/k8ssandra-operator.git
cd k8ssandra-operator
```

Invoke `make` with the following parameters for a single cluster: 

```bash
scripts/setup-kind-multicluster.sh --clusters 1 --kind-worker-nodes 4
```

**Output:**

```bash
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

Verify the deployment:

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

K8ssandra Operator has a dependency on `cert-manager`, which must be installed in each cluster, if not already available.

Update your helm repo and set the context:

```bash
helm repo add jetstack https://charts.jetstack.io

helm repo update

kubectl config use-context kind-k8ssandra-0
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

You can deploy K8ssandra Operator for namespace-scoped operations (the default), or cluster-scoped operations. 

* Deploying a namespace-scoped K8ssandra Operator means its operations -- watching for resources to deploy in Kubernetes -- are specific only to the **identified namespace** within a cluster. 
* Deploying a cluster-scoped operator means its operations -- again, watching for resources to deploy in Kubernetes -- are **global to all namespace(s)** in the cluster. The example in this section shows K8ssandra Operator deployed as namespace scoped:

Namespace-scoped example:

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
Optionally, you can use `--set global.clusterScoped=true` to install K8ssandra Operator cluster-scoped. Example:

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

Create the K8ssandraCluster with `kubectl apply`:

```bash
kubectl apply -n k8ssandra-operator -f k8c1.yml
```

**Output:**

```bash
k8ssandracluster.k8ssandra.io/demo created
```

### Verify pod deployment

```bash
kubectl get pods -n k8ssandra-operator
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
You'll use the extracted credentials for subsequent authentication in deployed containers.
{{% /alert %}}

### Verify cluster status

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- nodetool -u $CASS_USERNAME -pw $CASS_PASSWORD status
```

**Output plus nodetool example:**

```bash
Datacenter: dc1
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load       Tokens  Owns (effective)  Host ID                               Rack
UN  10.244.1.5  96.71 KiB  16      100.0%            4b95036b-1603-464f-bdee-b519fa28a079  default
UN  10.244.2.4  96.62 KiB  16      100.0%            ade61d9f-90f4-464c-8e18-dd3522c2bf3c  default
UN  10.244.3.4  96.7 KiB   16      100.0%            0f75a6fe-c91d-4c0e-9253-2235b6c9a206  default
```

{{% alert title="Tip" color="success" %}}
All nodes should have the status UN, which stands for "Up Normal".
{{% /alert %}}

### Test a few operations

To make it easier for you to copy the commands, we've listed them individually below:

Create a keyspace in the deployed Cassandra database, which is managed by K8ssandra Operator in the Kubernetes environment:

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- cqlsh -u $CASS_USERNAME -p $CASS_PASSWORD -e "CREATE KEYSPACE test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};"
```

Create a `test.users` table in the deployed Cassandra database:

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- cqlsh -u $CASS_USERNAME -p $CASS_PASSWORD  -e "CREATE TABLE test.users (email text primary key, name text, state text);"
```

Insert some data in the table:

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- cqlsh -u $CASS_USERNAME -p $CASS_PASSWORD -e "insert into test.users (email, name, state) values ('john@gamil.com', 'John Smith', 'NC');"
```

Insert another row of data in the table:

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- cqlsh -u $CASS_USERNAME -p $CASS_PASSWORD -e "insert into test.users (email, name, state) values ('joe@gamil.com', 'Joe Jones', 'VA');"
```

Insert another row of data in the table:

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- cqlsh -u $CASS_USERNAME -p $CASS_PASSWORD -e "insert into test.users (email, name, state) values ('sue@help.com', 'Sue Sas', 'CA');"
```

Insert another row of data in the table:

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- cqlsh -u $CASS_USERNAME -p $CASS_PASSWORD -e "insert into test.users (email, name, state) values ('tom@yes.com', 'Tom and Jerry', 'NV');"
```

Select data from the table:

```bash
% kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- cqlsh -u $CASS_USERNAME -p $CASS_PASSWORD -e "select * from test.users;"
```

**Output of the SELECT:**

In the launched container's `cqlsh` session, notice we provide the extracted password for `demo-superuser`.

```cqlsh
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
kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
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

## Next steps

* See other [local install]({{< relref "install/local/" >}}) options, including K8ssandra Operator in multi-cluster Kubernetes.
* Also, dig into the K8ssandra Operator [components]({{< relref "components" >}}).
