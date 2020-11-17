---
title: "K8ssandra Helm Chart"
linkTitle: "k8ssandra"
weight: 1
description: >
  Handles installation of all required operators for K8ssandra stacks.
---

```yaml
# Parameters related to the cass-operator deployment
cass-operator:
  # Namespace-scoped installs are currently not (well) supported.
  # See K8C-19 for details.
  clusterWideInstall: true

  # We need to use a patched version of cass-operator for now that has changes needed in
  # for Reaper and Medusa integration. Images will be built from
  # https://github.com/jsanda/cass-operator/tree/k8ssandra.
  image: jsanda/cass-operator:91205f4d8f1e

# Configuration for the dependent kube-prometheus-stack Chart
kube-prometheus-stack:
  alertmanager:
    enabled: false

  kubeStateMetrics:
    enabled: false

  grafana:
    enabled: false

  nodeExporter:
    enabled: false

  kubeletService:
    enabled: false

  kubeControllerManager:
    enabled: false

  kubelet:
    enabled: false

  kubeApiServer:
    enabled: false

  prometheus:
    enabled: false
```
