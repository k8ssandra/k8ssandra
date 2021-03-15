# Overview
This document describes the steps involved to upgrade the steps involves to apply a cass-operator upgrade in k8ssandra. There are a number of steps involved in the process. Here is a quick rundown:

* [medusa-operator](#medusa-operator)
  * [cass-operator manifests](#cass-operator-manifests)
  * [cass-operator image](#cass-operator-image)
  * [CassandraDatacenter CRD](#cassandradatacenter-crd)
  * [go.mod](#go-mod)
* Update reaper-operator
  * Update cass-operator manifests include CassandraDatacenter CRD
  * Update cass-operator image
  * Update CassandraDatacenter manifests if necessary
  * Update version in go.mod
* Update k8ssandra
  * Update templates in the cass-operator chart
  * Update default image cass-operator chart
  * Update version in go.mod 


# medusa-operator
medusa-operator uses kustomize to generate manifests for tests and to be used in the k8ssandra Helm charts.

## cass-operator manifests
Let's say we are upgrading to cass-operator 1.6.0. The cass-operator maniests should be taken from the v1.6.0 tag. The manifest (for this example upgrade) can be found [here](https://github.com/datastax/cass-operator/tree/v1.6.0/operator/deploy).

In medusa-operator, the cass-operator manifests, minus the CassandraDatacenter CRD, are bundled together in [cass-operator.yaml](https://github.com/k8ssandra/medusa-operator/blob/master/test/config/cass-operator/cass-operator.yaml).

**Note:** When we add kustomize support in cass-operator, we won't have to copy the manifests any more.

The CassandraDatacenter CRD to be updates lives at [medusa-operator/test/config/cass-operator/crd/bases/cassandradatacenter.yaml](https://github.com/k8ssandra/medusa-operator/blob/master/test/config/cass-operator/crd/bases/cassandradatacenter.yaml).

## cass-operator image
The cass-operator image needs to be udpated in [medusa-operator/test/config/cass-operator/kustomization.yaml](https://github.com/k8ssandra/medusa-operator/blob/master/test/config/cass-operator/kustomization.yaml)

## CassandraDatacenter CRD
Depending on the changes in cass-operator we might need to update the CassandraDatacenter manifest in [medusa-operator/test/config/cassdc](https://github.com/k8ssandra/medusa-operator/tree/master/test/config/cassdc).

## go.mod
Update the cass-operator version in [medusa-operator/go.mod](https://github.com/k8ssandra/medusa-operator/blob/master/go.mod)

## Release new version
Create a new release of medusa-operator. This will be needed for the medusa-operator chart in the k8ssandra repo.

**TODO:** Create doc on how to release new version and link to it.

# Update reaper-operator
reaper-operator uses kustomize to generarte manifests for tests and to be used in the k8ssandra Helm charts.

## Update cass-operator manifests
Let's say we are upgrading to cass-operator 1.6.0. The cass-operator maniests should be taken from the v1.6.0 tag. The manifest (for this example upgrade) can be found [here](https://github.com/datastax/cass-operator/tree/v1.6.0/operator/deploy).

In reaper-operator, the cass-operator manifests, minus the CassandraDatacenter CRD, are bundled together in [cass-operator.yaml](https://github.com/k8ssandra/reaper-operator/blob/master/test/config/cass-operator/cass-operator.yaml).

**Note:** When we add kustomize support in cass-operator, we won't have to copy the manifests any more.

The CassandraDatacenter CRD to be updates lives at [reaper-operator/test/config/cass-operator/crd/bases/cassandradatacenter.yaml](https://github.com/k8ssandra/reaper-operator/blob/master/test/config/cass-operator/crd/bases/cassandradatacenter.yaml).

## Update cass-operator image
The cass-operator image needs to be udpated in [reaper-operator/test/config/cass-operator/kustomization.yaml](https://github.com/k8ssandra/reaper-operator/blob/master/test/config/cass-operator/kustomization.yaml)

## Update CassandraDatacenter manifests if necessary
Depending on the changes in cass-operator we might need to update the CassandraDatacenter manifest in [reaper-operator/test/config/cassdc](https://github.com/k8ssandra/reaper-operator/tree/master/test/config/cassdc).

## Update go.mod
Update the cass-operator version in [reaper-operator/go.mod](https://github.com/k8ssandra/reaper-operator/blob/master/go.mod)

## Release new version
Create a new release of reaper-operator. This will be needed for the reaper-operator chart in the k8ssandra repo.

**TODO:** Create doc on how to release new version and link to it.

# Update k8ssandra
Several changes need to be made in the k8ssandra project including:

* Update cass-operator chart
* Update medusa-operator chart
* Update reaper-operator chart
* Update go.mod

## Update cass-operator chart
Update the templates in the cass-operator chart [here](https://github.com/k8ssandra/k8ssandra/tree/main/charts/cass-operator).

Update the default image in [k8ssandra/charts/cass-operator/values.yaml](https://github.com/k8ssandra/k8ssandra/blob/main/charts/cass-operator/values.yaml).

## Update medusa-operator chart
Update the default image in [k8ssandra/charts/medusa-operator/values.yaml](https://github.com/k8ssandra/k8ssandra/blob/main/charts/medusa-operator/values.yaml).

## Update reaper-operator chart
Update the default image in [k8ssandra/charts/reaper-operator/values.yaml](https://github.com/k8ssandra/k8ssandra/blob/main/charts/reaper-operator/values.yaml).