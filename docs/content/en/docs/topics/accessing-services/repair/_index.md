---
title: "Repair UI"
linkTitle: "Repair UI"
weight: 1
date: 2020-11-13
description: |
  Follow these steps to access the Repair Web Interface (Reaper).
---

## Tools

* Web Browser
* values.yaml configuration, or use --set flags on the command line

## Prerequisites

1. Kubernetes cluster with the following elements deployed:
   * [K8ssandra Operators]({{< ref "getting-started#install-k8ssandra" >}}) Helm
     Chart
   * [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}}) Helm
     Chart
   * [Ingress Controller]({{< ref "ingress" >}})
1. DNS name configured for the repair interface, referred to as _REPAIR DOMAIN_
   below.

## Access Repair Interface

![Reaper UI](reaper-ui.png)

With the prerequisites satisfied the repair GUI should be available at the
following address:

http://REPAIR_DOMAIN/webui

For example, if you installed a `k8ssandra-cluster` instance named `mycluster` with the following command:

```
helm install mycluster k8ssandra/k8ssandra-cluster \
  --set ingress.traefik.enabled=true \
  --set ingress.traefik.repair.host=repair.127.0.0.1.xip.io
NAME: mycluster
LAST DEPLOYED: Mon Nov 16 15:53:01 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

Notice how in this example, the DNS host name is specified on the command line as `repair.127.0.0.1.xip.io`.

After a few minutes, check that the pods are running. Example:

```
kubectl get pods
NAME                                                          READY   STATUS      RESTARTS   AGE
cass-operator-86d4dc45cd-pgcs8                                1/1     Running     0          12m
grafana-deployment-6bb9bc6d89-ghc4s                           1/1     Running     0          4m8s
k8ssandra-dc1-default-sts-0                                   2/2     Running     0          4m48s
k8ssandra-tools-grafana-operator-k8ssandra-54fbbc799c-68htn   1/1     Running     0          12m
k8ssandra-tools-kube-prome-operator-f87955c85-t2s9k           2/2     Running     0          12m
mycluster-reaper-k8ssandra-64b6b4c58-mkfxw                    1/1     Running     0          2m52s
mycluster-reaper-k8ssandra-schema-8lm4m                       0/1     Completed   4          4m27s
mycluster-reaper-operator-k8ssandra-799bd4568f-lk4hv          1/1     Running     0          4m49s
prometheus-mycluster-prometheus-k8ssandra-0                   3/3     Running     1          4m48s
```

With the pods running, access http://repair.127.0.0.1.xip.io/webui in a Web browser.

## What can I do in Reaper?

For details about the tasks you can perform in Reaper, see these topics in the
Cassandra Reaper documentation:

* [Check a cluster's health](http://cassandra-reaper.io/docs/usage/health/)
* [Run a cluster repair](http://cassandra-reaper.io/docs/usage/single/)
* [Schedule a cluster repair](http://cassandra-reaper.io/docs/usage/schedule/)
* [Monitor Cassandra diagnostic
  events](http://cassandra-reaper.io/docs/usage/cassandra-diagnostics/)
