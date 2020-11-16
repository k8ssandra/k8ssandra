---
title: "k8ssandra-cluster Helm Chart"
linkTitle: "k8ssandra-cluster"
weight: 2
description: >
  Provisions an instance of the k8ssandra stack
---

* `name`
  _string_
  default: `k8ssandra`
  Name of the cluster instance

* `clusterName`
  _string_
  default: `k8ssandra`
  validation: lowercase and consist of characters [a-z0-9\-]
  Name of the C\* cluster

* `datacenterName`
  _string_
  default: `dc1`
  validation: lowercase and consist of characters [a-z0-9\-]
  Name of the datacenter

* `size`
  _integer_
  default: 1
  Number of nodes in the datacenter

## `reaper`

* `enabled`
  _boolean_
  default: true
  Enables support for the Reaper repair service

* `jmx`
  _object_

  * `username`
    _string_
  * `password`
    _string_
