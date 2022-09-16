---
title: "Scale your Cassandra cluster"
linkTitle: "Scale"
description: "Steps to provision and scale up/down an Apache Cassandra® cluster in Kubernetes."
---

This topic explains how to add and remove Cassandra nodes in a Kubernetes cluster, as well as insights into the underlying operations that occur with scaling. 

## Prerequisites

1. Kubernetes cluster with the K8ssandra operators deployed:
    * If you haven't already installed a `K8ssandraCluster` using K8ssandra Operator, see the [local install]({{< relref "/install/local" >}}) topic.

## Create a cluster

Suppose we have the below `K8ssandraCluster` object, stored in a k8ssandra.yaml file:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
   name: my-k8ssandra
spec:
   cassandra:
      serverVersion: "4.0.5"
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

```

We can create this cluster with the following command:

```bash
kubectl apply -f k8ssandra.yaml
```

The above definition will result in the creation of a `CassandraDatacenter` object named `dc1` with
the size set to 3. 

```bash
kubectl get cassandradatacenter dc1
```

**Output:**

```text
NAME   AGE
dc1    3m35s
```

The cass-operator deployment that's installed by K8ssandra will in turn create the underlying
`StatefulSet` that has 3 Cassandra pods:

```bash
kubectl get statefulsets -l cassandra.datastax.com/datacenter=dc1
```

**Output:**

```text
NAME                           READY   AGE
my-k8ssandra-dc1-default-sts   3/3     7m35s
```

## Add nodes

Add nodes by updating the size property of the `K8ssandraCluster` spec:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
   name: my-k8ssandra
spec:
   cassandra:
      serverVersion: "4.0.5"
      datacenters:
         - metadata:
              name: dc1
           size: 4 # change this from 3 to 4  
      storageConfig:
         cassandraDataVolumeClaimSpec:
            storageClassName: standard
            accessModes:
               - ReadWriteOnce
            resources:
               requests:
                  storage: 5Gi
```

We can then update the cluster by running `kubectl apply` again:

```bash
kubectl apply -f k8ssandra.yaml
```

Alternatively, we can also patch the existing object with a `kubectl patch` command:

```bash
kubectl patch k8c my-k8ssandra --type='json' -p='[{"op": "replace", "path": "/spec/cassandra/datacenters/0/size", "value": 4}]'
```

### Underlying considerations when scaling up

After the new nodes are up and running, `nodetool cleanup` should run on all of the nodes except the
new ones to remove keys and data that no longer belong to those nodes. There is no need to do this
manually. The cass-operator deployment, which again is installed with K8ssandra, automatically runs
`nodetool cleanup` for you.

By default, cass-operator configures the Cassandra pods so that Kubernetes will not schedule
multiple Cassandra pods on the same worker node. If you try to increase the cluster size beyond the
number of available worker nodes, you may find that the additional pods do not deploy.

Look at this example output from `kubectl get pods` with a `K8ssandraCluster` whose size was
increased to 6. Assume that this value is beyond the number of available worker nodes:

```bash
kubectl get pods
```

**Output:**

```text
NAME                                   READY   STATUS      RESTARTS   AGE
my-k8ssandra-dc1-default-sts-0         2/2     Running     0          87m
my-k8ssandra-dc1-default-sts-1         2/2     Running     0          87m
my-k8ssandra-dc1-default-sts-2         2/2     Running     0          87m
my-k8ssandra-dc1-default-sts-3         2/2     Running     0          87m
my-k8ssandra-dc1-default-sts-4         2/2     Running     0          87m
my-k8ssandra-dc1-default-sts-5         2/2     Running     0          87m
my-k8ssandra-dc1-default-sts-6         0/2     Pending     0          3m6s
```

Notice that the `my-k8ssandra-dc1-default-sts-6` pod has a status of `Pending`. We can use `kubectl
describe pod` to get more details about the pod:

```bash
kubectl describe pod my-k8ssandra-dc1-default-sts-6
```

**Output:**

```text
...
Events:
  Type     Reason            Age                   From               Message
  ----     ------            ----                  ----               -------
  Warning  FailedScheduling  3m22s (x51 over 73m)  default-scheduler  0/6 nodes are available: 6 node(s) didn't match pod affinity/anti-affinity, 6 node(s) didn't satisfy existing pods anti-affinity rules.
```

The output reveals a `FailedScheduling` event.

To resolve the mismatch between the configured size and the available nodes, consider setting the
`softPodAntiAffinity` property to true in order to relax the constraint of only allowing one
Cassandra pod per Kubernetes worker node. This is useful in test/dev scenarios to minimise the
number of nodes required, but should not be done in production clusters.

Here is an updated `K8ssandraCluster` with `softPodAntiAffinity`:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
   name: my-k8ssandra
spec:
   cassandra:
      serverVersion: "4.0.5"
      datacenters:
         - metadata:
              name: dc1
           size: 6
      softPodAntiAffinity: true
      # Resources must be specified for each Cassandra node when using softPodAntiAffinity
      resources:
         requests:
            cpu: 1
            memory: 2Gi
         limits:
            cpu: 2
            memory: 2Gi
      # It is also recommended to set the JVM heap size
      config:
         jvmOptions:
            heap_initial_size: 1G
            heap_max_size: 1G
      storageConfig:
         cassandraDataVolumeClaimSpec:
            storageClassName: standard
            accessModes:
               - ReadWriteOnce
            resources:
               requests:
                  storage: 5Gi
```

When applied, this configuration updates the size property of the `CassandraDatacenter` to 6. Then
cass-operator will in turn update the underlying `StatefulSet`, allowing more than one Cassandra
node to sit on the same worker node.

### Conditions and events when scaling up

If you describe the `CassandraDatacenter` object, there should be a `ScalingUp` condition with its
status set to `true`. You should also see one or more `ScalingUpRack` events, one per rack; and
finally, you should also see `StartedCassandra` events, one per pod:

```bash

```bash
 kubectl describe cassandradatacenter dc1
```

**Output:**

```text
...
Status:
  Cassandra Operator Progress:  Updating
  Conditions:
    Last Transition Time:    2022-09-16T11:39:35Z
    Message:                 
    Reason:                  
    Status:                  True
    Type:                    ScalingUp

Events:
  Type     Reason                       Age       From           Message
  ----     ------                       ----      ----           -------
  Normal   ScalingUpRack                8m45s     cass-operator  Scaling up rack default
  Normal   StartedCassandra             2m11s     cass-operator  Started Cassandra for pod my-k8ssandra-dc1-default-sts-1
...
```

You can wait for the scale up operation to be finished by checking the value of the `ScalingUp`
condition. This can also be done with the following `kubectl wait` command:

```bash
kubectl wait --for=condition=ScalingUp=false cassandradatacenter dc1 --timeout=30m
```

## Remove nodes

Just like with adding nodes, removing nodes is simply a matter of changing the configured `size`
property. Then cass-operator does a few things when you decrease the datacenter size (see below).

### Underlying considerations when scaling down

First, cass-operator checks that the remaining nodes have enough capacity to handle the increased
storage capacity. If cass-operator determines that there is insufficient capacity, it will log a
message. Example:

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

Next, cass-operator runs `nodetool decommission` on the node to be removed. This step is done
automatically on your behalf.

Lastly, the pod is terminated.

{{% alert title="Note" color="success" %}} The StatefulSet controller manages the deletion of
Cassandra pods. It deletes one pod at a time, in reverse order with respect to its ordinal index.
This means for example that `my-k8ssandra-dc1-default-sts-3` will be deleted before
`my-k8ssandra-dc1-default-sts-2`. {{% /alert %}}

### Conditions and events when scaling down

If you describe the `CassandraDatacenter` object, there should be a `ScalingDown` condition with its
status set to `true`. You should also see one or more `ScalingDownRack` events, one per rack; and
finally, for each pod being decommissioned, you should also see one `LabeledPodAsDecommissioning`
and one `DeletedPvc` event:

```bash
 kubectl describe cassandradatacenter dc1
```

**Output:**

```text
...
Status:
  Cassandra Operator Progress:  Updating
  Conditions:
    Last Transition Time:    2022-09-16T11:51:20Z
    Message:                 
    Reason:                  
    Status:                  True
    Type:                    ScalingDown

Events:
  Type     Reason                       Age        From           Message
  ----     ------                       ----       ----           -------
  Normal   LabeledPodAsDecommissioning  82m        cass-operator  Labeled node as decommissioning my-k8ssandra-dc1-default-sts-1
  Normal   ScalingDownRack              82m        cass-operator  Scaling down rack rack1
  Normal   DeletedPvc                   80m        cass-operator  Claim Name: server-data-my-k8ssandra-dc1-default-sts-1
...
```

You can wait for the scale down operation to be finished by checking the value of the `ScalingDown`
condition. This can also be done with the following `kubectl wait` command:

```bash
kubectl wait --for=condition=ScalingDown=false cassandradatacenter dc1 --timeout=30m
```

## Scaling a multi-rack datacenter

The above `K8ssandraCluster` example has a single rack. When multiple racks are present, the
operator will always bootstrap new nodes and decommission existing nodes in a strictly deterministic
order that strives to achieve balanced racks.

### Scaling up with multiple racks

When scaling up, the operator will add new nodes to the racks with the fewest nodes. If two racks
have the same number of nodes, the operator will start with the rack that comes *first* in the
datacenter definition.

Let's take a look at an example. Suppose we create a 4-node `K8ssandraCluster` with 3 racks:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: my-k8ssandra
spec:
  cassandra:
    serverVersion: "4.0.5"
    datacenters:
      - metadata:
          name: dc1
        size: 4
        racks:
          - name: rack1
          - name: rack2
          - name: rack3
```

When this cluster is created, the target topology according to the rules above will be as follows:

* `rack1` should have 2 nodes
* `rack2` should have 1 node
* `rack3` should have 1 node

To achieve this, the operator will assign each rack its own `StatefulSet`:

```bash
kubectl get statefulsets
```

**Output:**

```text
NAME                         READY   AGE
my-k8ssandra-dc1-rack1-sts   2/2     7m21s
my-k8ssandra-dc1-rack2-sts   1/1     7m21s
my-k8ssandra-dc1-rack3-sts   1/1     7m21s
```

When all the pods in all StatefulSets are ready to be started, the operator will then bootstrap the
Cassandra nodes rack by rack, in this exact order:

* `rack1`: `my-k8ssandra-dc1-rack1-sts-0` gets bootstrapped first and becomes a seed;
* `rack2`: `my-k8ssandra-dc1-rack2-sts-0` gets bootstrapped next and becomes a seed;
* `rack3`: `my-k8ssandra-dc1-rack3-sts-0` gets bootstrapped next.
* `rack1`: `my-k8ssandra-dc1-rack1-sts-1` gets bootstrapped last.

This is indeed the most balanced topology we could have achieved with 4 nodes and 3 racks, and the
bootstrap order makes it so that no rack could have 2 nodes bootstrapped before the others.

Now let's scale up to, say, 8 nodes. The operator will bootstrap 4 new nodes, exactly in the
following order:

* `rack2` will scale up from 1 to 2 nodes and `my-k8ssandra-dc1-rack2-sts-1` gets bootstrapped;
* `rack3` will scale up from 1 to 2 nodes and `my-k8ssandra-dc1-rack3-sts-1` gets bootstrapped;
* `rack1` will scale up from 2 to 3 nodes and `my-k8ssandra-dc1-rack1-sts-2` gets bootstrapped;
* `rack2` will scale up from 2 to 3 nodes and `my-k8ssandra-dc1-rack2-sts-2` gets bootstrapped.

When the scale up operation is complete, the target topology will be as follows:

* `rack1` will have 3 nodes, among which 1 seed;
* `rack2` will have 3 nodes, among which 1 seed;
* `rack3` will have 2 nodes.

The StatefulSets will be updated accordingly:

```bash
kubectl get statefulsets                      
```

**Output:**

```text
NAME                         READY   AGE
my-k8ssandra-dc1-rack1-sts   3/3     19m
my-k8ssandra-dc1-rack2-sts   3/3     19m
my-k8ssandra-dc1-rack3-sts   2/2     19m
```

This is again the most balanced topology we could have achieved with 8 nodes and 3 racks.

### Scaling down with multiple racks

When scaling down, the operator will remove nodes from the racks with the most nodes. If two racks
have the same number of nodes, the operator will start with the rack that comes *last* in the
datacenter definition.

Now let's scale down the above cluster to 3 nodes. The operator will decommission 5 nodes, exactly
in the following order:

* `rack2` will scale down from 3 to 2 nodes and `my-k8ssandra-dc1-rack2-sts-2` gets decommissioned;
* `rack1` will scale down from 3 to 2 nodes and `my-k8ssandra-dc1-rack1-sts-2` gets decommissioned;
* `rack3` will scale down from 2 to 1 node and `my-k8ssandra-dc1-rack3-sts-1` gets decommissioned;
* `rack2` will scale down from 2 to 1 node and `my-k8ssandra-dc1-rack2-sts-1` gets decommissioned;
* `rack1` will scale down from 2 to 1 node and `my-k8ssandra-dc1-rack1-sts-1` gets decommissioned.

When the scale down operation is complete, the cluster will reach the following topology:

* `rack1` will have 1 node, among which 1 seed;
* `rack2` will have 1 node, among which 1 seed;
* `rack3` will have 1 node.

The StatefulSets will be updated accordingly:

```bash
kubectl get statefulsets                      
```

**Output:**

```text
NAME                         READY   AGE
my-k8ssandra-dc1-rack1-sts   1/1     30m
my-k8ssandra-dc1-rack2-sts   1/1     30m
my-k8ssandra-dc1-rack3-sts   1/1     30m
```

## Next steps

* Explore other K8ssandra Operator [tasks]({{< relref "/tasks" >}}).
* See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Operator
  Custom Resource Definitions (CRDs) and the single K8ssandra Operator Helm chart.
