---
title: "Repair UI"
linkTitle: "Repair UI"
weight: 3
description: |
  Follow these steps to access the Repair Web Interface (Reaper).
---

Repairs are a critical anti-entropy operation in Apache Cassandra&reg;. In the past, there have been many custom solutions to manage them outside of your main Cassandra Installation. K8ssandra provides the Repair Web Interface (also known as Reaper) that eliminates the need for a custom solution. Just like K8ssandra makes Cassandra setup easy, Reaper makes configuration of repairs even easier.

**Note:** The requirement for your environment may vary considerably, however the general recommendation is to run a repair operation on your Cassandra clusters about once a week. 

## Tools

* Web Browser
* values.yaml configuration, or use `--set` flags on the command line

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

For example, to upgrade a previously installed `k8ssandra` that's running locally:

`helm upgrade k8ssandra k8ssandra/k8ssandra --set ingress.traefik.enabled=true --set ingress.traefik.repair.host=repair.localhost`

Notice how in this example, the DNS host name is specified on the command line as `repair.localhost`.

After a few minutes, check that the pods are running. Example:

```
kubectl get pods
NAME                                                            READY   STATUS      RESTARTS   AGE
cass-operator-86d4dc45cd-pgcs8                                  1/1     Running     0          12m
grafana-deployment-6bb9bc6d89-ghc4s                             1/1     Running     0          4m8s
k8ssandra-dc1-default-sts-0                                     2/2     Running     0          4m48s
k8ssandra-tools-grafana-operator-k8ssandra-54fbbc799c-68htn     1/1     Running     0          12m
k8ssandra-tools-kube-prome-operator-f87955c85-t2s9k             2/2     Running     0          12m
k8ssandra-reaper-k8ssandra-64b6b4c58-mkfxw                      1/1     Running     0          2m52s
k8ssandra-reaper-operator-k8ssandra-799bd4568f-lk4hv            1/1     Running     0          4m49s
prometheus-mycluster-prometheus-k8ssandra-0                     3/3     Running     1          4m48s
```

## What can I do in Reaper?

To access Reaper, if you are running locally, navigate to [http://repair.localhost:8080/webui/](http://repair.localhost:8080/webui/).

### Check the cluster’s health

In the Reaper UI, notice how the nodes are displayed inside the datacenter for the cluster.

![OK](https://github.com/DataStax-Academy/kubecon2020/blob/main/Images/reaper1.png?raw=true)

The color of the nodes indicates the overall load the nodes are experiencing at the current moment. 

See [Check a cluster's health](http://cassandra-reaper.io/docs/usage/health/).

### Schedule a cluster repair

On the UI's left sidebar, notice the **Schedule** option.

![OK](https://github.com/DataStax-Academy/kubecon2020/blob/main/Images/reaper2.png?raw=true)

Click **Schedules**

![OK](https://github.com/DataStax-Academy/kubecon2020/blob/main/Images/reaper3.png?raw=true)

Click **Add schedule** and fill out the details when you are done click the final _add schedule_ to apply the new repair job.  A Cassandra best practice is to have one repair complete per week to prevent zombie data from coming back after a deletion. 

![OK](https://github.com/DataStax-Academy/kubecon2020/blob/main/Images/reaper4.png?raw=true)

Notice the new repair added to the list.

See [Schedule a cluster repair](http://cassandra-reaper.io/docs/usage/schedule/).

### Run a cluster repair

On the repair job you just configured, click **Run now**.  

![OK](https://github.com/DataStax-Academy/kubecon2020/blob/main/Images/reaper5.png?raw=true)

Notice the repair job kicking off.

## Recommended reading

* [Running a cluster repair](http://cassandra-reaper.io/docs/usage/single/)
* [Reaper details](http://cassandra-reaper.io/)
* [Blog articles](https://thelastpickle.com/blog/)
