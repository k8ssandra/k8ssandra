---
title: "Scale your Apache Cassandra® Cluster"
linkTitle: "Scale Cassandra"
weight: 4
description: Steps to provision a Cassandra cluster in Kubernetes
---

## Tools

[helm](https://helm.sh/docs/intro/install/)

## Prerequisites

* A Kubernetes environment
* K8ssandra installed and running in Kubernetes - see [Quick start]({{< ref "getting-started" >}})

## Steps

### Use helm to get the running configuration

For many basic configuration options, you may change values in the deployed YAML files. For example, you can scale up or scale down, as needed, by updating the YAML via `helm` command `--set` parameters.

Let's check the currently running values. First get the list of the installed K8ssandra chart. In this example, assume the `releaseName` was defined as `demo` on the `helm install demo k8ssandra/k8ssandra` command.

```bash
helm list
```

**Output**:

```bash
NAME	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART           	APP VERSION
demo	default  	1       	2021-03-04 15:21:41.573192403 +0000 UTC	deployed	k8ssandra-0.55.0	3.11.10     
```

You can specify the name of the installed cluster's `releaseName` to get the full manifest. 

`helm get manifest demo`

Helm displays full details of the properties defined in each deployed YAML file. 

### Scale up the cluster

To scale up, focus on the `size` property. Again, in this example `releaseName` was defind at kubectl install time as `demo`:

```bash
helm get manifest demo | grep size
```

**Output**:

```yaml
.
.
.
  size: 1
```

The value of `size: 1` is from cassdc.yaml, which is the Cassandra DataCenter definition. 

To scale up, you could change the `size` to 3. Example:

```bash
helm upgrade demo k8ssandra/k8ssandra --set cassandra.datacenters\[0\].size=3,cassandra.datacenters\[0\].name=dc1
```

**Output:**

```bash
Release "demo" has been upgraded. Happy Helming!
NAME: demo
LAST DEPLOYED: Thu Mar  4 15:39:58 2021
NAMESPACE: default
STATUS: deployed
REVISION: 2
```

Verify the upgrade:

```bash
helm get manifest demo | grep size
```

**Output**:

```yaml
.
.
.
  size: 3
```

### Scale down the cluster

Similarly, to scale down, lower the current `size` to conserve cloud resource costs, if the new value can support your computing requirements in Kubernetes. Example:

```bash
helm upgrade demo k8ssandra/k8ssandra --set cassandra.datacenters\[0\].size=1,cassandra.datacenters\[0\].name=dc1
```

**Output**:

```bash
Release "demo" has been upgraded. Happy Helming!
NAME: demo
LAST DEPLOYED: Thu Mar  4 15:42:39 2021
NAMESPACE: default
STATUS: deployed
REVISION: 3
```

Again, verify the upgrade:

```bash
helm get manifest demo | grep size
```

**Output**:

```yaml
.
.
.
  size: 1
```

## Next

Use Medusa for Apache Cassandra to [backup and restore]({{< ref "/docs/topics/restore-a-backup/" >}}) data from/to a Cassandra database.
