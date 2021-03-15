# Overview
This document describes the steps involved to upgrade the steps involves to apply a cass-operator upgrade in k8ssandra. There are a number of steps involved in the process. Here is a quick rundown:

* [medusa-operator](#medusa-operator)
  * [cass-operator manifests](#cass-operator-manifests-1)
  * [cass-operator image](#cass-operator-image-1)
  * [CassandraDatacenter CRD](#cassandradatacenter-crd-1)
  * [go.mod](#gomod-1)
  * [Release new version](#release-new-version-1)
* [reaper-operator](#reaper-operator)
  * [cass-operator manifests](#cass-operator-manifests-2)
  * [cass-operator image](#cass-operator-image-2)
  * [CassandraDatacenter CRD](#cassandradatacenter-crd-2)
  * [go.mod](#gomod-2)
  * [Release new version](#release-new-version-2)
* [k8ssandra](#k8ssandra)
  * [cass-operator chart](#cass-operator-chart)
  * [medusa-operator chart](#medusa-operator-chart)
  * [reaper-operator chart](#reaper-operator-chart)
  * [go.mod](#gomod)


# medusa-operator
medusa-operator uses kustomize to generate manifests for tests and to be used in the k8ssandra Helm charts.

## cass-operator manifests
Let's say we are upgrading to cass-operator 1.6.0. The cass-operator manifests should be taken from the v1.6.0 tag in the cass-operator repo. The manifest (for this example upgrade) can be found [here](https://github.com/datastax/cass-operator/tree/v1.6.0/operator/deploy).

In medusa-operator, the cass-operator manifests, minus the CassandraDatacenter CRD, are bundled together in [cass-operator.yaml](https://github.com/k8ssandra/medusa-operator/blob/master/test/config/cass-operator/cass-operator.yaml).

**Note:** When we add kustomize support in cass-operator, we won't have to copy the manifests any more.

## cass-operator image
The cass-operator image needs to be udpated in [medusa-operator/test/config/cass-operator/kustomization.yaml](https://github.com/k8ssandra/medusa-operator/blob/master/test/config/cass-operator/kustomization.yaml)

## CassandraDatacenter CRD
The CassandraDatacenter CRD lives at [medusa-operator/test/config/cass-operator/crd/bases/cassandradatacenter.yaml](https://github.com/k8ssandra/medusa-operator/blob/master/test/config/cass-operator/crd/bases/cassandradatacenter.yaml).

## go.mod
Update the cass-operator version in [medusa-operator/go.mod](https://github.com/k8ssandra/medusa-operator/blob/master/go.mod)

## Release new version
Create a new release of medusa-operator. This is needed for the medusa-operator chart in the k8ssandra repo.

**TODO:** Create doc on how to release new version and link to it.

# reaper-operator
reaper-operator uses kustomize to generarte manifests for tests and to be used in the k8ssandra Helm charts.

## cass-operator manifests
Let's say we are upgrading to cass-operator 1.6.0. The cass-operator maniests should be taken from the v1.6.0 tag. The manifest (for this example upgrade) can be found [here](https://github.com/datastax/cass-operator/tree/v1.6.0/operator/deploy).

In reaper-operator, the cass-operator manifests, minus the CassandraDatacenter CRD, are bundled together in [cass-operator.yaml](https://github.com/k8ssandra/reaper-operator/blob/master/test/config/cass-operator/cass-operator.yaml).

**Note:** When we add kustomize support in cass-operator, we won't have to copy the manifests any more.

## cass-operator image
The cass-operator image needs to be udpated in [reaper-operator/test/config/cass-operator/kustomization.yaml](https://github.com/k8ssandra/reaper-operator/blob/master/test/config/cass-operator/kustomization.yaml)

## CassandraDatacenter CRD
The CassandraDatacenter CRD lives at [reaper-operator/test/config/cass-operator/crd/bases/cassandradatacenter.yaml](https://github.com/k8ssandra/reaper-operator/blob/master/test/config/cass-operator/crd/bases/cassandradatacenter.yaml).

## go.mod
Update the cass-operator version in [reaper-operator/go.mod](https://github.com/k8ssandra/reaper-operator/blob/master/go.mod)

## Release new version
Create a new release of reaper-operator. This is be needed for the reaper-operator chart in the k8ssandra repo.

**TODO:** Create doc on how to release new version and link to it.

# k8ssandra
Several chart updates are needed as well as `go.mod`.

## cass-operator chart
Update the templates in the cass-operator chart [here](https://github.com/k8ssandra/k8ssandra/tree/main/charts/cass-operator).

Update the default image in [k8ssandra/charts/cass-operator/values.yaml](https://github.com/k8ssandra/k8ssandra/blob/main/charts/cass-operator/values.yaml).

## medusa-operator chart
Update the default image in [k8ssandra/charts/medusa-operator/values.yaml](https://github.com/k8ssandra/k8ssandra/blob/main/charts/medusa-operator/values.yaml).

## reaper-operator chart
Update the default image in [k8ssandra/charts/reaper-operator/values.yaml](https://github.com/k8ssandra/k8ssandra/blob/main/charts/reaper-operator/values.yaml).

## go.mod
Update the cass-operator version in [k8ssandra/go.mod](https://github.com/k8ssandra/k8ssandra/blob/main/go.mod).