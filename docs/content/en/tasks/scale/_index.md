---
title: "Scale your Cassandra cluster"
linkTitle: "Scale"
description: "Steps to provision and scale up/down an Apache Cassandra® cluster in Kubernetes."
---

This topic explains how to add and remove Cassandra nodes in a Kubernetes cluster, as well as insights into the underlying operations that occur with scaling. 

{{% alert title="Note" color="success" %}}
K8ssandra currently only supports a single-datacenter cluster.
{{% /alert %}}

## Prerequisites

* A Kubernetes environment.
* [Helm](https://helm.sh/docs/intro/install/) is installed.
* K8ssandra is installed and running in Kubernetes - see the [Quick starts]({{< relref "quickstarts" >}}).

## Create a cluster

Suppose we install K8ssandra as follows:

```bash
helm install my-k8ssandra k8ssandra/k8ssandra -f k8ssandra-values.yaml
```

Assume that `k8ssandra-values.yaml` has these properties:

```yaml
cassandra:
  clusterName: my-k8ssandra
  datacenters:
  - name: dc1
    size: 3
```

The `helm install` command will result in the creation of a `CassandraDatacenter` object with the size set to 3. The cass-operator deployment that's installed by K8ssandra will in turn create the underlying StatefulSet that has 3 Cassandra pods.

## Add nodes

Add nodes by updating the size property of the datacenter. Example values file:

```yaml
cassandra:
  clusterName: my-k8ssandra
  datacenters:
  - name: dc1
    size: 4
```

Apply the changes with `helm upgrade`:

```bash
helm upgrade my-k8ssandra k8ssandra/k8ssandra -f k8ssandra-values.yaml
```

{{% alert title="Tip" color="success" %}}
Another way to upgrade your K8ssandra cluster is by passing in a `--set` parameter. Also include a `--reuse-values` parameter so that Helm will reuse previous values (other than the one you're overriding with each `--set` parameter). Without `--reuse-values` it's easy to make a mistake if you have other, additional properties that you previously set.  

```bash
helm upgrade my-k8ssandra k8ssandra/k8ssandra --reuse-values --set cassandra.datacenters\[0\].size=4,cassandra.datacenters\[0\].name=dc1
```

{{% /alert %}}

## Underlying considerations when increasing size values

By default, cass-operator configures the Cassandra pods so that Kubernetes will not schedule multiple Cassandra pods on the same worker node. If you try to increase the cluster size beyond the number of available worker nodes, you may find that the additional pods do not deploy. 

Look at this example output from `kubectl get pods` with a test cluster whose size was increased to 6. Assume that this value is beyond the number of available worker nodes:

```bash
kubectl get pods
```

**Output:**

```bash
NAME                                   READY   STATUS      RESTARTS   AGE
test-dc1-default-sts-0                 2/2     Running     0          87m
test-dc1-default-sts-1                 2/2     Running     0          87m
test-dc1-default-sts-2                 2/2     Running     0          87m
test-dc1-default-sts-3                 2/2     Running     0          87m
test-dc1-default-sts-4                 2/2     Running     0          87m
test-dc1-default-sts-5                 2/2     Running     0          87m
test-dc1-default-sts-6                 0/2     Pending     0          3m6s
```

Notice that the `test-dc1-default-sts-6` pod has a status of `Pending`. We can use `kubectl describe pod` to get more details about the pod:

```bash
kubectl describe pod test-dc1-default-sts-6
```

**Output:**

```bash
...
Events:
  Type     Reason            Age                   From               Message
  ----     ------            ----                  ----               -------
  Warning  FailedScheduling  3m22s (x51 over 73m)  default-scheduler  0/6 nodes are available: 6 node(s) didn't match pod affinity/anti-affinity, 6 node(s) didn't satisfy existing pods anti-affinity rules.
```

The output reveals a `FailedScheduling` event.

To resolve the mismatch between the configured size and the available nodes, consider the following option to set the `allowMultipleNodesPerWorker` property to relax the constraint of only allowing one Cassandra pod per Kubernetes worker node.

Here is an updated k8ssandra-values.yaml with `allowMultipleNodesPerWorker`:

```yaml
cassandra:
  clusterName: my-k8ssandra
  allowMultipleNodesPerWorker: true
  # resources must be set when allowMultipleNodesPerWorker is true.   
  resources: 
    requests:
      cpu: 2
      memory: 2Gi
    limits:
      cpu: 2
      memory: 2Gi
  # It is not required to set the heap but is recommended.
  heap:
    size: 1024M
    newGenSize: 512M
  datacenters:
  - name: dc1
    size: 3
```

When applied to the test cluster, this configuration updates the size property of the `CassandraDatacenter`. Then cass-operator will in turn update the underlying `StatefulSet`.

If you check the status of the `CassandraDatacenter` object, there should be a `ScalingUp` condition with its status set to `true`. It should look like this:

```bash
 kubectl get cassandradatacenter dc1 -o yaml
```

**Output:**

```bash
...
status:
  cassandraOperatorProgress: Updating
  conditions:
  - lastTransitionTime: "2021-03-30T22:01:48Z"
    message: ""
    reason: ""
    status: "True"
    type: ScalingUp
...
```

After the new nodes are up and running, `nodetool cleanup` should run on all of the nodes except the new ones to remove keys and data that no longer belong to those nodes. There is no need to do this manually. The cass-operator deployment, which again is installed with K8ssandra, automatically runs `nodetool cleanup` for you.

## Remove nodes

Just like with adding nodes, removing nodes is simply a matter of changing the configured `size` property. The cass-operator does a few things when you decrease the datacenter size.

## Underlying considerations when lowering size values

First, cass-operator checks that the remaining nodes have enough capacity to handle the increased storage capacity. If cass-operator determines that there is insufficient capacity, it will log a message. Example:

```text
Not enough free space available to decommission. my-k8ssandra-dc1-default-sts-3 has 12345 free space, but 67891 is needed.
```

The reported units are in bytes.

The cass-operator deployment will also add a condition to the `CassandraDatacenter` status. Example:

```yaml
status:
  conditions:
  - lastTransitionTime: "2021-03-30T22:01:48Z"
    message: "Not enough free space available to decommission. my-k8ssandra-dc1-default-sts-3 has 12345 free space, but 67891 is needed."
    reason: "NotEnoughSpaceToScaleDown"
    status: "False"
    type: Valid
...
```

Next, cass-operator runs `nodetool decommission` on the node to be removed. This step is done automatically on your behalf.

Lastly, the pod is terminated.

{{% alert title="Note" color="success" %}}
The StatefulSet controller manages the deletion of Cassandra pods. It deletes one pod at a time, in reverse order with respect to its ordinal index. 
This means for example that `my-k8ssandra-dc1-default-sts-3` will be deleted before `my-k8ssandra-dc1-default-sts-2`.
{{% /alert %}}

## Next steps

Explore other K8ssandra [tasks]({{< relref "/tasks" >}}).

See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Helm charts, and a glossary. 