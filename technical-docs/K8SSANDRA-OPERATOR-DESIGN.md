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

This is the primary object which users will manage. The controllers will be responsible for creating other constituent objects.

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

Similar to what we have done in the Helm charts, you will be able to specify things at cluster level and override them at the DC level. You will also be able to specify properties at the DC level without specifying them at the cluster level.


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

This is defined in the [reaper-operator](https://github.com/k8ssandra/reaper-operator) repo. It is possible today to use reaper-operator on its own without k8ssandra. In that regard, living in its own repo makes sense. There is a good bit of overhead for each repo that we need to maintain. In practice, I do not think reaper-operator will be used a whole lot outside of k8ssandra. For these reasons we will consolidate the reaper-operator code into the k8ssandra-operator repo.

From a technical standpoint, consolidation means that the reaper-operator controllers would run in the k8ssandra-operator binary. This will reduce overhead in that it would be able to reuse the same cache as other controllers.

### Stargate

This does not yet exist. Similar to reaper-operator, one could certainly make the case that a separate stargate-operator project makes sense; however, we will create this in the k8ssandra-operator repo.

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

These are defined in the [medusa-operator](https://github.com/k8ssandra/medusa-operator) repo. These will be consolidated into the k8ssandra-operator repo.

See this [GitHub discussion](https://github.com/k8ssandra/k8ssandra/discussions/273) more for ideas around backup/restore and CRDs.

### Monitoring

With the Grafana license changes and discussion around Victoria Metrics, we only do the following for the first iteration:

* Expose metrics over Prometheus endpoints (i.e., can be scraped by Prometheus or anything else that handles Prometheus metrics)
* Provide ServiceMonitors
* Provide Grafana dashboards


## Go Types

### Stargate
```go
type Stargate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec StargateSpec     `json:"spec,omitempty"`
	Status StargateStatus `json:"status,omitempty"`
}

type StargateSpec struct {
	// properties to be defined
}

type StargateStatus struct {
	// properties to be defined
}
```

### Reaper

```go
type Reaper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec ReaperSpec     `json:"spec,omitempty"`	Status ReaperStatus `json:"status,omitempty"`
}

type ReaperSpec struct {
	// properties to be defined
}

type ReaperStatus struct {
	// properties to be defined
}
```

### CassandraCluster

```go
type CassandraCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec CassandraClusterSpec     `json:"spec,omitempty"`
	Status CassandraClusterStatus `json:"status,omitempty"`
}

type CassandraClusterSpec struct {
	DatacenterTemplates []CassandraDatacenterTemplate `json:"datacenters,omitempty"`
	
	// The following properties are also defined in CassandraDatacenterTemplatespec.
	// The idea is to be able to define things at the cluster level and let them be
	// inherited at the datacenter level with the option of overriding. 
	
	StargateTemplate *StargateTemplateSpec `json:"stargate,omitempty"`
	
	ReaperTemplate *ReaperTemplateSpec `json:"reaper,omitempty"`
	
	MonitoringTemplate *MonitoringTemplateSpec `json:"monitoring,omitempty"`
}

type CassandraClusterStatus struct {
}


type MonitoringTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec MonitoringSpec `json:"spec,omitempty"`

}
```

### K8ssandra

```go
type K8ssandraCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec K8ssandraClusterSpec     `json:"spec,omitempty"`
	Status K8ssandraClusterStatus `json:"status,omitempty"`
}

type K8ssandraClusterSpec struct {
	DatacenterTemplates []CassandraDatacenterTemplateSpec *CassandraClusterTemplateSpec `json:"cassandra,omitempty"`

    StargateTemplate *StargateTemplateSpec `json:"stargate,omitempty"`

    ReaperTemplate *ReaperTemplateSpec `json:"reaper,omitempty"`
}

type CassandraDatacenterTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CassandraDatacenterSpec `json:"spec,omitempty"`
	
	StargateTemplate *StargateTemplateSpec `json:"stargate,omitempty"`
	
	ReaperTemplate *ReaperTemplateSpec `json:"reaper,omitempty"	
}

type StargateTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec StargateSpec `json:"spec,omitempty"`
}

type ReaperTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec ReaperSpec `json:"spec,omitempty"`
}
```

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
```

#### Event: K8ssandraCluster created

##### k8ssandra controller

* Start reconciliation loop when K8ssandraCluster object is created
* Creates a CassandraDatacenter
  * Update status
* Creates a Reaper
  * Update status
* Creates a Stargate
  * Update status
* End reconciliation loop

#### Event: CassandraCluster created

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
* Requeue request until CassandraDatacenter is ready

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

#### Event: CassandraDatacenter status updated

##### k8ssandra controller

* Start reconciliation loop
  * Update status to reflect state of CassandraDatacenters
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

#### Notes

The k8ssandra controller creates several objects:

* CassandraDatacenter
* Reaper
* Stargate

Whenever any of these objects is updated, it should trigger a reconciliation of the k8ssandra controller. In the reconciliation loop, the status of each object will essentially be pushed up into the K8ssandraCluster object.

### Change Cassandra heap settings

The user updates the heap settings in the K8ssandraCluster manifest and applies the changes with kubectl apply.

#### Event: K8ssandraCluster updated

##### k8ssandra controller

* Update CassandraDatacenter(s)
  * Update K8ssandraCluster status
* Requeue reconciliation request

#### Event: CassandraCluster status updated

#### Event: CassandraDatacenter update complete and ready

##### k8sandra cluster controller

* Update status of K8ssandraCluster
* End reconciliation loop

#### Event: CassandraDatacenter status updated

##### k8ssandra controller

* Update status to reflect state of CassandraDatacenter

#### Notes

Notice that the user does not directly update the CassandraDatacenter. The CassandraDatacenter is owned and managed by the K8ssandraCluster controller. Updates made directly to the CassandraDatacenter would potentially be lost. Changes need to be propagated through the K8ssandraCluster object.

## Being Prescriptive

With the Helm charts we are prescriptive with what properties and settings of the CassandraDatacenter we expose. Each parent type should have a TemplateSpec property for each of its child types that need to be made configurable. This provides flexibility and is a pretty common design pattern. 

## Installation

There are a lot of CRDs and multiple operators covered in this doc. The k8ssandra operator itself should not be responsible for installing and managing all of the operators and CRDs. We will provide Helm charts that manage the installation of the k8ssandra operator and its various dependencies, e.g., prometheus-operator.

We should also provide first class support for kustomize. kustomize is integrated into kubectl as well as operator-sdk and kubebuilder. It therefore makes sense to provide support for kustomize.

Lastly, we should provide integration with Operator Lifecycle Manager (OLM).

## Migrations

We need to provide a migration path for k8ssandra 1.x. This should be external tooling and documentation, not something done within the operator itself.

We also need to provide a migration for people using only cass-operator and not k8ssandra.