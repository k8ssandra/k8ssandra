---
title: "Scale your Cassandra Cluster"
linkTitle: "Scale Cassandra"
weight: 4
description: Steps to provision a cluster in Kubernetes
---

## Tools

[helm](https://helm.sh/docs/intro/install/)

## Prerequisites

* A Kubernetes environment
* k8ssandra installed and running in Kubernetes - see [Quick start]({{< ref "getting-started" >}})

## Steps

### Use helm to get the running configuration

For many basic configuration options, you may change values in the deployed YAML files. For example, you can scale up or scale down, as needed, by updating the YAML via `helm` command `--set` parameters.

Let's check the currently running values. First get the list of the installed K8ssandra chart. In this example, assume the `releaseName` was defined as `demo` on the `helm install demo k8ssandra/k8ssandra` command.

`helm list`
 
```
NAME	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART           	APP VERSION
demo	default  	1       	2021-02-18 23:34:14.547364974 +0000 UTC	deployed	k8ssandra-0.51.0	3.11.10     
```

You can specify the name of the installed cluster's `releaseName` to get the full manifest. 

`helm get manifest demo`

Helm displays full details of the properties defined in each deployed YAML file. 

### Scale up the cluster

To scale up, focus on the `size` property. Again, in this example `releaseName` was defind at kubectl install time as `demo`:

`helm get manifest demo | grep size`

The output includes:
   
```
.
.
.
  size: 1
      initial_heap_size: "800M"
      max_heap_size: "800M"
    :[{\"expr\":\"sum(mcac_table_memtable_off_heap_size{cluster=~\\\"$cluster\\\"\
    :\"A\"},{\"expr\":\"sum(mcac_table_memtable_on_heap_size{cluster=~\\\"$cluster\\\
    ,\"description\":\"Total sizes of the data on distinct nodes\",\"fill\":0,\"gridPos\"\
    datasource\":\"$PROMETHEUS_DS\",\"description\":\"Maximum JVM Heap Memory size\
    \ (worst node) and minimum available heap size\",\"fill\":1,\"gridPos\":{},\"\
```

Notice the value of `size: 1` from cassdc.yaml. This is the Cassandra DataCenter definition. 

To scale up, you could change the `size` to 3. Example with helm:

`helm upgrade demo k8ssandra/k8ssandra --set k8ssandra.size=3 --reuse-values`

{{% alert title="Tip" color="success" %}}
Use `--reuse-values` to ensure keeping settings from a previous `helm upgrade`.
{{% /alert %}}

```
Release "demo" has been upgraded. Happy Helming!
NAME: demo
LAST DEPLOYED: Thu Feb 18 23:35:12 2021
NAMESPACE: default
STATUS: deployed
REVISION: 2
TEST SUITE: None
```

Verify the upgrade:

`helm get manifest demo | grep size`
 
```
.
.
.
  size: 3
      initial_heap_size: "800M"
      max_heap_size: "800M"
    :[{\"expr\":\"sum(mcac_table_memtable_off_heap_size{cluster=~\\\"$cluster\\\"\
    :\"A\"},{\"expr\":\"sum(mcac_table_memtable_on_heap_size{cluster=~\\\"$cluster\\\
    ,\"description\":\"Total sizes of the data on distinct nodes\",\"fill\":0,\"gridPos\"\
    datasource\":\"$PROMETHEUS_DS\",\"description\":\"Maximum JVM Heap Memory size\
    \ (worst node) and minimum available heap size\",\"fill\":1,\"gridPos\":{},\"\
```

### Scale down the cluster

Similarly, to scale down, lower the current `size` to conserve cloud resource costs, if the new value can support your computing requirements in Kubernetes. Example:

`helm upgrade demo k8ssandra/k8ssandra --set k8ssandra.size=1 --reuse-values`
 
```
Release "demo" has been upgraded. Happy Helming!
NAME: demo
LAST DEPLOYED: Thu Feb 18 23:37:25 2021
NAMESPACE: default
STATUS: deployed
REVISION: 3
TEST SUITE: None
```

Again, verify the upgrade:

`helm get manifest demo | grep size`
```
.
.
.
  size: 1
      initial_heap_size: "800M"
      max_heap_size: "800M"
    :[{\"expr\":\"sum(mcac_table_memtable_off_heap_size{cluster=~\\\"$cluster\\\"\
    :\"A\"},{\"expr\":\"sum(mcac_table_memtable_on_heap_size{cluster=~\\\"$cluster\\\
    ,\"description\":\"Total sizes of the data on distinct nodes\",\"fill\":0,\"gridPos\"\
    datasource\":\"$PROMETHEUS_DS\",\"description\":\"Maximum JVM Heap Memory size\
    \ (worst node) and minimum available heap size\",\"fill\":1,\"gridPos\":{},\"\
```

## Next

Use Medusa to [backup and restore]({{< ref "/docs/topics/restore-a-backup/" >}}) data from/to a Cassandra database. 
