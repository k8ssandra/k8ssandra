---
title: "Repair Cassandra with Reaper"
linkTitle: "Repair"
description: "Use the Reaper for Apache Cassandra&reg; web interface to perform repairs."
---

Repairs are a critical anti-entropy operation in Cassandra. In the past, there have been many custom solutions to manage them outside of your main Cassandra installation. K8ssandra Operator provides the Reaper web interface that eliminates the need for a custom solution. Just as K8ssandra Operator makes Cassandra setup easy, Reaper makes configuration of repairs even easier.

{{% alert title="Tip" color="success" %}}
The requirement for your environment may vary considerably, however the general recommendation is to run a repair operation on your Cassandra clusters once a week.
{{% /alert %}}

## Tools

* Web Browser
* kubectl

## Prerequisites

1. Kubernetes cluster with the following elements deployed:
   * If you haven't already installed a K8ssandraCluster using K8ssandra Operator, see the [local install]({{< relref "/install/local" >}}) topic.

Access to the Reaper web interface requires either:  

* setting up a custom Ingress resource
* or modifying Reaper's Kubernetes service as a LoadBalancer (in cloud environments), which will expose it over a public IP
* or using port forwarding, which is another way to provide external access to resources that have been deployed by K8ssandra Operator in your Kubernetes environment:  
  * Developers, see [Set up port forwarding]({{< relref "/quickstarts/developer/#set-up-port-forwarding" >}}).  
  * Site reliability engineers, see [Configure port forwarding]({{< relref "/quickstarts/site-reliability-engineer/#port-forwarding" >}}).

## Access the Reaper web interface

![Reaper UI](reaper-main-ui.png)

With the prerequisites satisfied the Reaper web interface should be available at the following address:

http://REAPER_DOMAIN/webui

Check that the pods are running. Example:

```bash
kubectl get pods
```

**Output:**

```bash
NAME                                                         READY   STATUS    RESTARTS   AGE
cass-operator-controller-manager-55f6b84454-zpcfd            1/1     Running   1          10d
k8ssandra-dc1-default-stargate-deployment-7847d945b4-vxth8   1/1     Running   0          10d
k8ssandra-dc1-default-sts-0                                  3/3     Running   0          10d
k8ssandra-dc1-default-sts-1                                  3/3     Running   0          10d
k8ssandra-dc1-default-sts-2                                  3/3     Running   0          10d
k8ssandra-dc1-reaper-58cd6b795b-dw9dw                        1/1     Running   0          10d
k8ssandra-operator-6d4dd9fb8f-5kzl6                          1/1     Running   0          10d
```

## What can I do in Reaper?

To access Reaper, navigate to [http://localhost:8080/webui/](http://localhost:8080/webui/). 

{{% alert title="Tip" color="success" %}}
If you are not running locally, use the IP address provided by your cluster for the ingress resource or the load balancer service.
{{% /alert %}}

### Check the cluster’s health

In the Reaper UI, notice how the nodes are displayed inside the datacenter for the cluster.

![Reaper cluster](reaper-cluster.png)

The color of the nodes indicates the overall load the nodes are experiencing at the current moment.

See [Check a cluster's health](http://cassandra-reaper.io/docs/usage/health/).

### Schedule a cluster repair

On the UI's left sidebar, notice the **Schedule** option.

![Reaper schedule](reaper-schedule.png)

Click **Schedules**

![Reaper add schedule](reaper-add-schedule1.png)

Click **Add schedule** and fill out the details when you are done click the final _add schedule_ to apply the new repair job.  A Cassandra best practice is to have one repair complete per week to prevent zombie data from coming back after a deletion.

![Reaper add schedule part 2](reaper-add-schedule2.png)

Enter values for the keyspace, tables, owner, and other fields. Then click **Add Schedule**. The details for adding a schedule are similar to the details for the Repair form, except the “Clause” field is replaced with two fields:

* “Start time”
* “Interval in days”

After creating a scheduled repair, the page is updated with a list of Active and Paused repair schedules.

{{% alert title="Important" color="info" %}}
When choosing to add a new repair schedule, we recommended that you limit the repair schedules to specific tables, instead of scheduling repairs for an entire keyspace. Creating different repair schedules will allow for simpler scheduling, fine-grain tuning for more valuable data, and easily grouping tables with smaller data load into different repair cycles. For example, if there are certain tables that contain valuable data or a business requirement for high consistency and high availability, they could be scheduled for repairs during low-traffic periods.
{{% /alert %}}

For additional information, see [Schedule a cluster repair](http://cassandra-reaper.io/docs/usage/schedule/) on the Reaper site.

{{% alert title="Warning" color="warning" %}}
Users with access to the Reaper web interface can pause or delete scheduled repairs. Authentication security in the UI is automatically added if authentication is enabled in the K8ssandraCluster object. A secret can be referenced that contains the credentials under `reaper.uiUserSecretRef`, otherwise it will be generated by the k8ssandra-operator.
{{% /alert %}}

### Autoschedule repairs

When you enable the autoschedule feature, Reaper dynamically schedules repairs for all non-system keyspaces in a cluster. A cluster's keyspaces are monitored and any modifications (additions or removals) are detected. When a new keyspace is created, a new repair schedule is created automatically for that keyspace. Conversely, when a keyspace is removed, the corresponding repair schedule is deleted.

To enable autoschedule in Reaper, set the property `reaper.autoScheduling.enabled` to `true`. 

### Run a cluster repair

On the repair schedule you just configured, click **Run now**.

![Reaper run now](reaper-schedule-run-now.png)

Notice the repair job kicking off.

## Recommended reading

* [Running a cluster repair](http://cassandra-reaper.io/docs/usage/single/)
* [Reaper details](http://cassandra-reaper.io/)
* [Blog articles](https://thelastpickle.com/blog/)

## Next steps

* Explore other K8ssandra Operator [tasks]({{< relref "/tasks" >}}).
* See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Operator Custom Resource Definitions (CRDs) and the single K8ssandra Operator Helm chart. 
* See the Reaper [Custom Resource Definition (CRD)]({{< relref "/reference/crd" >}}) reference.
