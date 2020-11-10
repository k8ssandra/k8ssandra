---
title: "cass-operator Helm Chart"
linkTitle: "cass-operator"
weight: 1
description: >
  Installs cass-operator
---

* `clusterWideInstall`
  _boolean_
  default: `false`
* `serviceAccountName`
  _string_
  default: `cass-operator`
* `clusterRoleName`
  _string_
  default: `cass-operator-cr`
* `clusterRoleBindingName`
  _string_
  default: `cass-operator-crb`
* `roleName`
  _string_
  default: `cass-operator`
* `roleBindingName`
  _string_
  default: `cass-operator`
* `webhookClusterRoleName`
  _string_
  default: `cass-operator-webhook`
* `webhookClusterRoleBindingName`
  _string_
  default: `cass-operator-webhook`
* `deploymentName`
  _string_
  default: `cass-operator`
* `deploymentReplicas`
  _integer_
  default: `1`
* `defaultImage`
  _string_
  default: `datastax/cass-operator:1.4.1`
* `imagePullPolicy`
  _string_
  default: `IfNotPresent`
* `imagePullSecret`
  _string_
  default: ``
