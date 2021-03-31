---
title: "Scale your Apache Cassandra® cluster in K8ssandra"
linkTitle: "Scale Cassandra"
description: "Steps to provision and scale up or scale down an Apache Cassandra® cluster in Kubernetes"
---

## Tools

[helm](https://helm.sh/docs/intro/install/)

## Prerequisites

* A Kubernetes environment
* K8ssandra installed and running in Kubernetes - see [Quick start]({{< ref "getting-started" >}})

## Steps

### Use helm to get the running configuration

For many basic configuration options, you may change values in the deployed YAML files. For example, you can scale up or scale down, as needed, by updating the YAML via `helm` command `--set` parameters.

Let's check the currently running values. First get the list of the installed K8ssandra chart. In this example, assume the `releaseName` was defined as `k8ssandra` on the `helm install` command.

```bash
helm list
```

**Output**:

```bash
NAME     	  NAMESPACE	 REVISION   UPDATED                               STATUS  	CHART      APP VERSION
k8ssandra	  default  	 1          2021-03-04 20:49:32.975090399 +0000   UTC	      deployed	 k8ssandra-1.0.0	              
```

You can specify the name of the installed cluster's `releaseName` to get the full manifest. 

`helm get manifest k8ssandra`

Helm displays full details of the properties defined in each deployed YAML file. 

### Scale up the cluster

To scale up, focus on the `size` property. Let's find the current value:

```bash
helm get manifest k8ssandra | grep size
```

**Output**:

```yaml
.
.
.
    size: 1
      initial_heap_size: 1G
      max_heap_size: 1G
      heap_size_young_generation: 1G
```

The value of `size: 1` is from cassdc.yaml, which is the CassandraDatacenter definition. 

To scale up, you could change the `size` to 3. In the following example, we'll also set the name `dc1`:

```bash
helm upgrade k8ssandra k8ssandra/k8ssandra --set cassandra.datacenters\[0\].size=3,cassandra.datacenters\[0\].name=dc1
```

**Output:**

```bash
Release "k8ssandra" has been upgraded. Happy Helming!
NAME: k8ssandra
LAST DEPLOYED: Thu Mar  4 21:12:01 2021
NAMESPACE: default
STATUS: deployed
REVISION: 2
```

Verify the upgrade:

```bash
helm get manifest k8ssandra | grep size
```

**Output**:

```yaml
.
.
.
                   "description": "Total sizes of the data on distinct nodes",
                   "description": "Maximum JVM Heap Memory size (worst node) and minimum available heap size",
  size: 3
```

### Scale down the cluster

Similarly, to scale down, lower the current `size` to conserve cloud resource costs, if the new value can support your computing requirements in Kubernetes. For example, this time we'll lower the size to 1, and again set the CassandraDatacenter name `dc1` (currently required each time) with the command:

```bash
helm upgrade k8ssandra k8ssandra/k8ssandra --set cassandra.datacenters\[0\].size=1,cassandra.datacenters\[0\].name=dc1
```

**Output**:

```bash
Release "k8ssandra" has been upgraded. Happy Helming!
NAME: k8ssandra
LAST DEPLOYED: Thu Mar  4 21:14:05 2021
NAMESPACE: default
STATUS: deployed
REVISION: 3
```

Again, verify the upgrade:

```bash
helm get manifest k8ssandra | grep size
```

**Output**:

```yaml
.
.
.
                   "description": "Total sizes of the data on distinct nodes",
                   "description": "Maximum JVM Heap Memory size (worst node) and minimum available heap size",
  size: 1
```

## Next

Use Medusa for Apache Cassandra to [backup and restore]({{< ref "/docs/topics/restore-a-backup/" >}}) data from/to a Cassandra database.
