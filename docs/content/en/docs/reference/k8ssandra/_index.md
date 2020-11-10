---
title: "k8ssandra Helm Chart"
linkTitle: "k8ssandra"
weight: 1
description: >
  Handles installation of all required operators for a K8ssandra stack.
---

## `cass-operator`

* `clusterWideInstall`
  _boolean_
  default: `true`

* `image`
  _string_
  default: `datastax/cass-operator:1.5.0`

## `reaper-operator`

* `enabled`
  _boolean_
  default: `true`

## `kube-prometheus-stack`
