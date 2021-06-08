# K8ssandra Operator Design

Author: @jsanda

## Overview

GH issue: https://github.com/k8ssandra/k8ssandra/issues/683

This is a design document for an operator for k8ssandra. For background see 
https://github.com/k8ssandra/k8ssandra/issues/485.

Much of the behavior that is implemented in the Helm charts will be reimplemented in various controllers. The role of the Helm charts changes to focus primarily on installation and configuration of operators.

While we will continue to use and support Helm, it will not be required for installing k8ssandra from a purely technical perspective. Users can install manifests with kubectl apply for example. We can also provide support for kustomize.

## CustomResourceDefinitions

We will utilize a number of CRDs, some of which are developed and maintained by other projects.

There will be a separate controller for each CRD.

The goal of this document is not to be entirely prescriptive for each CRD. Rather, the intent is to give a general idea of how things will look and fit together. 

### K8ssandraCluster

This is the primary object which users will manage. The controller will be responsible for creating other constituent objects.

It will look something like this:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: k8ssandra
spec:
  cassandra:
    clusterName: k8ssandra
    auth:
      enabled: true
    datacenters:
    - name: dc1
      size: 3
  reaper:
    enabled: true
    # optionally set and override defaults for Reaper    
  stargate: 
    enabled: true  
    # optionally set and override defaults for Stargate 
  monitoring:
    prometheus:
      enabled: true
      # optionally set and override defaults for Prometheus
    grafana: 
      enabled: true
      # optionally set and override defaults for Grafana
```

Similar to the k8ssandra Helm chart, there is an enabled flag for each component. The Go type for this will be a *bool. In general as it is with the Helm chart, setting enabled: true should be all that is necessary. The controller will take care of configuring the respective objects. 

The user will be able to override defaults as necessary as illustrated in the next example:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: k8ssandra
spec:
  cassandra:
    clusterName: k8ssandra
    auth:
      enabled: true
    datacenters:
    - name: dc1
      size: 3
  reaper:
    enabled: true
    serverConfig:
      jmxUserSecretName: my-custom-secret
      cassandraBackend:
        replication:
          networkTopology:
            dc1: 5         
  stargate: 
    enabled: true  
    replicas: 3
```

### CassandraCluster

```
apiVersion: k8ssandra.io/v1alpha1
kind: CassandraCluster
metadata:
  name: k8ssandra
spec:
  clusterName: k8ssandra
  auth:
    enabled: true
  datacenters:
  - name: dc1
    size: 3
```

The controller will create a CassandraDatacenter object for each entry in the datacenters array. Similar to what we have done in the Helm charts, you will be able to specify things at cluster level and override them at the DC level. You will also be able to specify properties at the DC level without specifying them at the cluster level.

Multi-cluster is out of scope, but I want to point out that I think this is where we introduce it. Consider the following example:

```
apiVersion: k8ssandra.io/v1alpha1
kind: CassandraCluster
metadata:
  name: k8ssandra
spec:
  clusterName: k8ssandra
  auth:
    enabled: true
  datacenters:
  - name: dc1
    size: 3
  - name: dc2
    k8s:
      kubeconfigSecret: my-other-cluster
      namespace: dev
    size: 3
```

`kubeconfigSecret` specifies the name of a secret that contains the kubeconfig file needed to access a remote cluster. 

`namespace` indicates in which namespace the DC should be created.

### CassandraDatacenter

This is defined in the [cass-operator](https://github.com/k8ssandra/cass-operator) repo. The controller runs in its own Deployment.

### Reaper

```
apiVersion k8ssandra.io/v1alpha1
kind: Reaper
metadata:
  name: k8ssandra
spec:
  serverConfig:
    storageType: cassandra
    jmxUserSecretName: k8ssandra
    cassandraBackend:
      cassandraDatacenter
        name: dc1
      replication:
        networkTopology:
          dc1: 3
```

This is defined in the [reaper-operator](https://github.com/k8ssandra/reaper-operator) repo. It is possible today to use reaper-operator on its own without k8ssandra. In that regard, living in its own repo makes sense. There is a good bit of overhead for each repo that we need to maintain. In practice, I do not think reaper-operator will be used a whole lot outside of k8ssandra. For these reasons I am in favor of consolidating with the k8ssandra-operator repo.

From a technical standpoint, consolidation means that the reaper-operator controllers would run in the k8ssandra-operator binary. It would reduce overhead in that it would be able to use the same cache as other controllers.

### Stargate

This does not yet exist. Similar to reaper-operator, one could certainly make the case that a separate stargate-operator project makes sense. I am more on the fence about where this should live than I am with reaper-operator.

```
apiVersion: k8ssandra.io/v1alpha1
kind: Stargate
metadata:
  name: k8ssandra
spec:
  replicas: 3
  cassandraDatacenter:
    name: dc1
```

### CassandraBackup / CassandraRestore

These are defined in the [medusa-operator](https://github.com/k8ssandra/medusa-operator) repo. I propose consolidating this into the k8ssandra-operator repo.

See this [GitHub discussion](https://github.com/k8ssandra/k8ssandra/discussions/273) more for ideas around backup/restore and CRDs.

### MonitoringConfiguration

```
apiVersion: k8ssandra.io/v1alpha1
kind: MonitoringConfiguration
metadata:
  name: k8ssandra
spec:
  prometheus:
    # prometheus properties...
  grafana:
  # grafana properties...
```

With the Grafana license changes and discussion around Victoria Metrics, I think it is a good idea to encapsulate Prometheus and Grafana.

### Grafana, GrafanaDashboard, GrafanaDatasource

Defined in [grafana-operator](https://github.com/integr8ly/grafana-operator).

### Prometheus

Defined in [prometheus-operator](https://github.com/integr8ly/grafana-operator).

## Scenarios

Each scenario consists of a series of events. Each event lists the steps performed by controllers. Controllers run concurrently. When you see multiple controllers listed for an event assume that the steps are performed concurrently and potentially in parallel.

### Create K8ssandraCluster

Let's first consider creating a single DC cluster. User creates a K8ssandraCluster object via kubectl apply:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  cassandra:
    clusterName: "test"
    version: "4.0.0" 
    auth:
      enabled: true
    datacenters:
    - name: dc1
      size: 3
  reaper:
    # reaper properties...
  stargate:
    enabled: true
  medusa:
    enabled: true
  monitoring:
    prometheus:
      enabled: true
    grafana:
      enabled: true
```

#### Event: K8ssandraCluster created

##### k8ssandra controller

* Start reconciliation loop when K8ssandraCluster object is created
* Creates a CassandraCluster
  * Update status
* Creates a Reaper
  * Update status
* Creates a Stargate
  * Update status
* Creates a MonitoringConfiguration
  * Update status
* End reconciliation loop

#### Event: CassandraCluster created

##### cassandra cluster controller

* Start reconciliation loop when CassandraCluster object is created
* Create CassandraDatacenter
  * Update status
* End reconciliation loop

#### Event: Reaper created

##### reaper controller

* Start reconciliation loop when Reaper object is created
* Creates reaper Service
* Requeue request until CassandraDatacenter is ready

#### Event: Stargate created

##### stargate controller

* Start reconciliation loop when Stargate object is created
* Creates stargate Service
* Creates ServiceMonitor
* Creates a GrafanaDashboard
* Requeue request until CassandraDatacenter is ready

#### Event: MonitoringConfiguration created

##### monitoring controller

* Start reconciliation loop when MonitoringConfiguration created
* Create Grafana object
* Create GrafanaDataSource object
* Create Prometheus object
* End reconciliation loop

#### Event: CassandraDatacenter ready

##### cassandra cluster controller

* Update status to reflect state of CassandraCluster
* End reconciliation loop

##### reaper controller

* Start reconciliation loop after requeue
* Create Job to initialize reaper keyspace
* Requeue request until job finishes

##### stargate controller

* Start reconciliation loop after requeue
* Create stargate Deployment
* Requeue request until Deployment ready

#### Event: CassandraCluster status updated

##### k8ssandra controller

* Start reconciliation loop
  * Update status to reflect state of CassandraCluster
* End reconciliation loop

#### Event: Reaper schema job completed

##### reaper controller

* Create Deployment
* Requeue request until Deployment ready
* End reconciliation loop

#### Event: Stargate ready

##### k8ssandra controller

* Update status to reflect state of Stargate

#### Event: Reaper ready

##### k8ssandra controller

* Update status to reflect state of Reaper

#### Event: Grafana and Prometheus ready

##### monitoring controller

* Update status of MonitoringConfiguration to reflect state of Grafana and Prometheus

#### Event: MonitoringConfiguration status updated

##### k8ssandra controller

* Update status to reflect state of MonitoringConfiguration

#### Notes

The k8ssandra controller creates several objects:

* CassandraCluster
* Reaper
* Stargate
* MonitoringConfiguration

Whenever any of these objects is updated, it should trigger a reconciliation of the k8ssandra controller. In the reconciliation loop, the status of each object will essentially be pushed up into the K8ssandraCluster object.

### Change Cassandra heap settings

The user updates the heap settings in the K8ssandraCluster manifest and applies the changes with kubectl apply.

#### Event: K8ssandraCluster updated

##### k8ssandra controller

* Update CassandraCluster 
* End reconciliation loop

#### Event: CassandraCluster updated

##### cassandra cluster controller

* Update CassandraDatacenter(s)
  * Update CassandraCluster status
* Requeue reconciliation request

#### Event: CassandraCluster status updated

##### k8ssandra controller

* Update status to reflect state of CassandraCluster

#### Event: CassandraDatacenter update complete and ready

##### cassandra cluster controller

* Update status of CassandraCluster
* End reconciliation loop

#### Event: CassandraCluster status updated

##### k8ssandra controller

* Update status to reflect state of CassandraCluster

#### Notes

Notice that the user does not directly update the CassandraDatacenter. The CassandraDatacenter is owned and managed by the CassandraCluster controller. Updates made directly to the CassandraDatacenter would potentially be lost. Changes need to be propagated through the K8ssandraCluster object.

## Being Prescriptive

With the Helm charts we are prescriptive with what properties and settings of the CassandraDatacenter we expose. I have been thinking that each parent type would have a TemplateSpec property for each of its child types. For example, K8ssandraCluster would have a CassandraClusterTemplate property that is of type CassandraClusterSpec. CassandraCluster would have something like a []CassandraDatacenterTemplate where CassandraDatacenterTemplate is of type CassandraDatacenterSpec. This provides flexibility and is a pretty common design pattern. 

Aside from deviating away from the current, prescriptive approach there is something else that we have to consider. Suppose the user decides to add a volume or container that clases with one that K8ssandra adds. cass-operator already deals with this with the PodTemplateSpec property. I am not suggesting that we do one thing or another at this time. I just want to raise awareness.

## Upgrades

Let's assume for the moment that K8ssandra 2.0 is based on the k8ssandra-operator. We need to figure out what the upgrade and migration path looks like. Would it be possible to start incrementally introducing changes in K8ssandra 1.x?