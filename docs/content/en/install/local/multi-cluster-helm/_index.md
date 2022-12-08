---
title: "Multi-cluster install with helm"
linkTitle: "Multi-cluster/helm"
no_list: false
weight: 2
description: "Quickstart with a unified helm chart to install K8ssandraCluster in multi-cluster Kubernetes."
---

This topic shows how you can use `helm` and other tools to install and configure the `K8ssandraCluster` custom resource in **multi-cluster** local Kubernetes, using K8ssandra Operator. 

## Prerequisites

If you haven't already, see the install [prerequisites]({{< relref "install/local/" >}}).

## Quick start for multi-cluster

Deploy K8ssandra with multiple Cassandra datacenters in a **multi-cluster** kind environment.

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

Invoke `make` with the following parameters: 

```bash
scripts/setup-kind-multicluster.sh --clusters 3 --kind-worker-nodes 4
```

### Verify the deployments 

Set the context to each of the four created clusters, and get node information for each cluster. Examples:

```bash
kubectl config use-context kind-k8ssandra-0
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
kubectl config use-context kind-k8ssandra-1

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
kubectl config use-context kind-k8ssandra-2

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

### Install cert-manager in each cluster

If you haven't already, update your helm repo with the jetstack cert-manager. 

```bash
helm repo add jetstack https://charts.jetstack.io

helm repo update
```

Set the per-cluster context and install `jetstack/cert-manager`. Examples:

```bash
kubectl config use-context kind-k8ssandra-0

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
kubectl config use-context kind-k8ssandra-1

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
kubectl config use-context kind-k8ssandra-2

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

First, you'll need to have [Helm v3+](https://helm.sh/docs/intro/install/) installed.

Then configure the K8ssandra Helm repository:

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

In the following example, of the three clusters we've created in the section above, we'll use `kind-k8ssandra-0` as our control-plane.

```bash
kubectl config use-context kind-k8ssandra-0

helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \
     --create-namespace
```

### Install K8ssandra Operator in the data-planes

In this example, we'll use the three other clusters as data-planes.

```bash
kubectl config use-context kind-k8ssandra-1
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \
 --create-namespace --set controlPlane=false

kubectl config use-context kind-k8ssandra-2
helm install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator \
 --create-namespace --set controlPlane=false
 ```

### Verify control-plane configuration

```bash
kubectl config use-context kind-k8ssandra-0

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
kubectl config use-context kind-k8ssandra-1

kubectl -n k8ssandra-operator get deployment k8ssandra-operator \
 -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="K8SSANDRA_CONTROL_PLANE")].value}'
```

**Output:**

```bash
false
```

### Generate and install ClientConfigs

[create-clientconfig.sh](https://github.com/k8ssandra/k8ssandra-operator/blob/main/scripts/create-clientconfig.sh) lives in the k8ssandra-operator repo. It is used to configure access to remote clusters. 

**Note:** K8ssandra Operator restarts automatically whenever there is a change to a `ClientConfig` (a create, update, or delete operation). This restart is done in order to update connections to remote clusters.

First, set the context to `kind-k8ssandra-0`, the control plane cluster. 

```bash
kubectl config use-context kind-k8ssandra-0
```

Run the create-clientconfig.sh script, once per data plane cluster.  

```bash
scripts/create-clientconfig.sh --namespace k8ssandra-operator \
    --src-kubeconfig ./build/kind-kubeconfig \
    --dest-kubeconfig ./build/kind-kubeconfig \
    --src-context kind-k8ssandra-1 \
    --dest-context kind-k8ssandra-0 \
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
scripts/create-clientconfig.sh --namespace k8ssandra-operator \
    --src-kubeconfig ./build/kind-kubeconfig \
    --dest-kubeconfig ./build/kind-kubeconfig \
    --src-context kind-k8ssandra-2 \
    --dest-context kind-k8ssandra-0 \
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

### Deploy the K8ssandraCluster

To deploy the `K8ssandraCluster`, we use a custom YAML file. In this example, k8cm1.yml. Notice, there are two Cassandra 4.0.1 datacenters, `dc1` and `dc2` that are associated with the two data plane clusters.

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
      - metadata:
          name: dc2
        k8sContext: kind-k8ssandra-2
        size: 3
  stargate:
     size: 1
     heapSize: 512M
```

Verify again that your context is set to the control plane cluster, which is in this example:

```bash
kubectl config use-context kind-k8ssandra-0
```

Apply the YAML to the already deployed K8ssandra Operator. 

```bash
kubectl apply -n k8ssandra-operator -f k8cm1.yml
```

### Verify pod deployment

Initially the rollout will begin in dc1 and work across the full cluster:

```bash
kubectl config use-context kind-k8ssandra-1

kubectl get pods -n k8ssandra-operator
```

Do the same the other cluster by setting the kubectl config use-context context to kind-k8ssandra-2, check the pods status.

Eventually the datacenters will be fully deployed:

```bash
kubectl config use-context kind-k8ssandra-0

kubectl get pods -n k8ssandra-operator
```

**Output:**

```bash
NAME                                                READY   STATUS    RESTARTS   AGE
k8ssandra-operator-68568ffbd5-l6t2f                 1/1     Running   0          93m
k8ssandra-operator-cass-operator-794f65d9f4-kqrpf   1/1     Running   0          97m
```

```bash
kubectl config use-context kind-k8ssandra-1

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
kubectl config use-context kind-k8ssandra-2

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

### Verify K8ssandraCluster status

While deployment is still in progress, you can check the status:

```bash
kubectl config use-context kind-k8ssandra-0

kubectl describe k8cs demo -n k8ssandra-operator
```

In the **earlier** deployment phases, you may notice statuses such as:

**Output:**

```bash
   ...
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
   ...
```

This behavior is expected for the deployments-in-progress. If you're curious, you can continue to check status by submitting the command again. When the deployments have been completed, for example, here's the command again and a portion of its completed output:

```bash
kubectl describe k8cs demo -n k8ssandra-operator
```

**Output subset:**

```bash
   ...
dc2:
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
          demo-dc2-default-sts-0:
            Host ID:  f2bd31ef-5ca5-4c28-a9b2-bac28f76af4f
          demo-dc2-default-sts-1:
            Host ID:  878e519b-a6f4-4aff-b8ab-c1fb30679847
          demo-dc2-default-sts-2:
            Host ID: 3efa6a2f-47d1-49e3-ba93-0a58870e7c0f
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
kubectl config use-context kind-k8ssandra-0

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
You'll use the extracted credentials for subsequent authentication in deployed containers.
{{% /alert %}}

### Verify cluster status

On one of the data plane clusters, verify the cluster status. Example:

```bash
kubectl config use-context kind-k8ssandra-1

kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- nodetool -u $CASS_USERNAME -pw $CASS_PASSWORD status
```

**Output plus nodetool example:**

```bash
Datacenter: dc1
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address      Load       Tokens       Owns (effective)  Host ID                               Rack
UN  172.18.0.10  327.02 KiB  256          100.0%            9e3e48ee-529e-4b2a-9bf2-39575d32f578  default
UN  172.18.0.11  338.79 KiB  256          100.0%            305616ce-3440-4d37-b9be-32bca624abb7  default
UN  172.18.0.8   304.01 KiB  256          100.0%            0a0864b7-968a-4e07-839f-a5abb2e3d0a4  default
Datacenter: dc2
===============
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address      Load       Tokens       Owns (effective)  Host ID                               Rack
UN  172.18.0.16  355.6 KiB  256          100.0%            f2bd31ef-5ca5-4c28-a9b2-bac28f76af4f  default
UN  172.18.0.15  368.66 KiB  256          100.0%            878e519b-a6f4-4aff-b8ab-c1fb30679847  default
UN  172.18.0.13  344.76 KiB  256          100.0%            3efa6a2f-47d1-49e3-ba93-0a58870e7c0f  default
```

{{% alert title="Tip" color="success" %}}
All nodes should have the status UN, which stands for "Up Normal".
{{% /alert %}}

### Test a few operations

```bash
kubectl config use-context kind-k8ssandra-1

kubectl exec -it demo-dc1-default-sts-0 -n k8ssandra-operator -- /bin/bash
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

demo-superuser@cqlsh> CREATE KEYSPACE test WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1' : 3, 'dc2' : 3};
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
kubectl config use-context kind-k8ssandra-2

kubectl exec -it demo-dc2-default-sts-0 -n k8ssandra-operator -- /bin/bash
```

**Output:**

```bash
Defaulted container "cassandra" out of: cassandra, server-system-logger, jmx-credentials (init), server-config-init (init)
cassandra@k8ssandra-2-worker3:/$ cqlsh -u demo-superuser -p KT-ROFfbD-O9BzWS3Lxq
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

Now use the Stargate API. 

**Output plus cqlsh &amp; stargate-service example:**

```bash
kubectl config use-context kind-k8ssandra-2

kubectl exec -it demo-dc2-default-sts-0 -n k8ssandra-operator -- /bin/bash

Defaulted container "cassandra" out of: cassandra, server-system-logger, jmx-credentials (init), server-config-init (init)
cassandra@k8ssandra-2-worker3:/$ cqlsh -u demo-superuser -p KT-ROFfbD-O9BzWS3Lxq demo-dc3-stargate-service
Connected to demo at demo-dc2-stargate-service:9042
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

* See other [local install]({{< relref "install/local/" >}}) options.
* Also, dig into the K8ssandra Operator [components]({{< relref "components" >}}).




