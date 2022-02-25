---
title: "Quickstart for Site Reliability Engineers"
linkTitle: "SREs"
weight: 2
description: "Familiarize yourself with K8ssandra utilities and procedures for repair, upgrade, and backup/restore operations for your Apache Cassandra® database."
---

{{% alert title="Tip" color="success" %}}
Before performing these post-install steps, complete at least one K8ssandra Operator [cluster deployment]({{< relref "/install/local" >}}) in Kubernetes. 
{{% /alert %}}

In this quickstart for Site Reliability Engineers (SREs), we'll cover:

* [Accessing nodetool commands]({{< relref "#nodetool" >}}) like status, ring, and info.
* [Configure port forwarding]({{< relref "#port-forwarding" >}}) for the Prometheus and Grafana monitoring utilities as well as Reaper for Apache Cassandra® (Reaper).
* [Accessing the K8ssandra Operator monitoring utilities]({{< relref "#monitoring" >}}), Prometheus and Grafana.
* [Accessing Reaper]({{< relref "#reaper" >}}), an easy to use repair interface.
* [Upgrading a K8ssandra Operator cluster]({{< relref "#upgrade" >}}): to ensure you're using the latest K8ssandra Operator, or to apply new settings.

## Access the Apache Cassandra® nodetool utility {#nodetool}

Cassandra's nodetool utility is commonly used for a variety of monitoring and management tasks. You'll need to run nodetool on your K8ssandra cluster using the `kubectl exec` command, because there's no external standalone option available.

To run `nodetool` commands:

1. Get a list of the running K8ssandra pods using `kubectl get`:

    ```bash
    kubectl get pods -n k8ssandra-operator
    ```

    **Output:**

    ```
    NAME                                                    READY   STATUS    RESTARTS   AGE
    demo-dc1-default-stargate-deployment-7b6c9d8dcd-k65jx   1/1     Running   0          5m33s
    demo-dc1-default-sts-0                                  2/2     Running   0          10m
    demo-dc1-default-sts-1                                  2/2     Running   0          10m
    demo-dc1-default-sts-2                                  2/2     Running   0          10m
    k8ssandra-operator-7f76579f94-7s2tw                     1/1     Running   0          11m
    k8ssandra-operator-cass-operator-794f65d9f4-j9lm5       1/1     Running   0          11m
```

    The K8ssandra pod running Cassandra takes the form `<k8ssandra-cluster-name>-<datacenter-name>-default-sts-<n>` and, in the example above is `demo-dc1-default-sts-0` which we'll use throughout the following sections.

    {{% alert title="Tip" color="success" %}}
Although not applicable to this quick start, additional K8ssandra Operator Cassandra nodes will increment the final `<n>` but the rest of the name will remain the same.
    {{% /alert %}}

1. Run [`nodetool status`](https://docs.datastax.com/en/cassandra-oss/3.x/cassandra/tools/toolsStatus.html), using the Cassandra node name `demo-dc1-default-sts-0`, and replacing `<k8ssandra-username>` and `<k8ssandra-password>` with the values you retrieved in the local install topic's [Extract credentials]({{< relref "/install/local#extract-credentials" >}}) section:

Hint: in that topic, the first single-cluster example returned this password: `ACK7dO9qpsghIme-wvfI`.

With known `-u` and `-p` credentials, you can enter a command that invokes nodetool. Here's an example. Your credentials will be different:

    ```bash
    kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- nodetool -u demo-superuser -pw ACK7dO9qpsghIme-wvfI status

    ```

    ```bash
    Datacenter: dc1
    ===============
    Status=Up/Down
    |/ State=Normal/Leaving/Joining/Moving
    --  Address     Load       Tokens  Owns (effective)  Host ID                               Rack
    UN  10.244.1.5  96.71 KiB  16      100.0%            4b95036b-1603-464f-bdee-b519fa28a079  default
    UN  10.244.2.4  96.62 KiB  16      100.0%            ade61d9f-90f4-464c-8e18-dd3522c2bf3c  default
    UN  10.244.3.4  96.7 KiB   16      100.0%            0f75a6fe-c91d-4c0e-9253-2235b6c9a206  default
    ```

{{% alert title="Tip" color="success" %}}
All nodes should have the status UN, which stands for "Up Normal".
{{% /alert %}}

Other useful nodetool commands include:

* [`nodetool ring`](https://docs.datastax.com/en/cassandra-oss/3.x/cassandra/tools/toolsRing.html) which outputs all the tokens in the node. Example:

    ```bash
    kubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- nodetool -u demo-superuser -pw ACK7dO9qpsghIme-wvfI ring
    ```

    **Output**:

    ```bash
    Datacenter: dc1
    ==========
    Address      Rack        Status State   Load            Owns                Token
                                                                                9126546575375666475
    172.17.0.13  default     Up     Normal  597.42 KiB      ?                   -9138166261715795932
    172.17.0.13  default     Up     Normal  597.42 KiB      ?                   -9120920057340937901
    172.17.0.13  default     Up     Normal  597.42 KiB      ?                   -9117737800555727340
    172.17.0.13  default     Up     Normal  597.42 KiB      ?                   -9058127181143818684
    172.17.0.13  default     Up     Normal  597.42 KiB      ?                   -8998548020695455271
    ...
    ```

* [`nodetool info`](https://docs.datastax.com/en/cassandra-oss/3.x/cassandra/tools/toolsInfo.html) which provides load and uptime information:

    ```bash
    kkubectl exec --stdin --tty demo-dc1-default-sts-0 -n k8ssandra-operator -c cassandra -- nodetool -u demo-superuser -pw ACK7dO9qpsghIme-wvfI info
    ```

    **Output**:

    ```bash
    ID                     : dec6a537-f00c-458a-bbc0-26b173675cc7
    Gossip active          : true
    Thrift active          : true
    Native Transport active: true
    Load                   : 597.42 KiB
    Generation No          : 1614265335
    Uptime (seconds)       : 9232
    Heap Memory (MB)       : 567.72 / 1024.00
    Off Heap Memory (MB)   : 0.00
    Data Center            : dc1
    Rack                   : default
    Exceptions             : 0
    Key Cache              : entries 39, size 3.46 KiB, capacity 51 MiB, 199 hits, 240 requests, 0.829 recent hit rate, 14400 save period in seconds
    Row Cache              : entries 0, size 0 bytes, capacity 0 bytes, 0 hits, 0 requests, NaN recent hit rate, 0 save period in seconds
    Counter Cache          : entries 0, size 0 bytes, capacity 25 MiB, 0 hits, 0 requests, NaN recent hit rate, 7200 save period in seconds
    Chunk Cache            : entries 6, size 384 KiB, capacity 224 MiB, 111 misses, 3472 requests, 0.968 recent hit rate, NaN microseconds miss latency
    Percent Repaired       : 100.0%
    Token                  : (invoke with -T/--tokens to see all 256 tokens)
    ```

For details on all nodetool commands, see [The nodetool utility](https://docs.datastax.com/en/cassandra-oss/3.x/cassandra/tools/toolsNodetool.html).

## Configure port forwarding {#port-forwarding}

In order to access Cassandra utilities outside of the K8s cluster, if you don't have an Ingress setup as described in [Configure Ingress]({{< relref "tasks/connect/ingress" >}}), you'll need to configure port forwarding.

Begin by getting a list of your K8ssandra K8s services and ports:

```bash
kubectl get services
```

**Output**:

```bash
(TODO: NEED SAMPLE OUTPUT HERE)
```

To configure port forwarding:

1. Open a new terminal.

2. Run the following 3 `kubectl port-forward` commands in the background. Example:

    ```bash
    kubectl port-forward svc/k8ssandra-grafana 9191:80 &
    kubectl port-forward svc/prometheus-operated 9292:9090 &
    kubectl port-forward svc/k8ssandra-reaper-reaper-service 9393:8080 &
    ```

    **Output**:

    ```bash
    [1] 29211
    [2] 29212
    [3] 29213

    ~/
    Forwarding from 127.0.0.1:9292 -> 9090
    Forwarding from [::1]:9292 -> 9090
    Forwarding from 127.0.0.1:9393 -> 8080
    Forwarding from [::1]:9393 -> 8080
    Forwarding from 127.0.0.1:9191 -> 3000
    Forwarding from [::1]:9191 -> 3000
    ```

The K8ssandra Operator services are now available at:

* Prometheus: <http://127.0.0.1:9292>
* Grafana: <http://127.0.0.1:9191>
* Reaper: <http://127.0.0.1:9393/webui>

### Terminate port forwarding

To terminate a particular forwarded port:

1. Get the process ID:

    ```bash
    jobs -l
    ```

    **Output**:

    ```bash
    [3]  + 29213 running    kubectl port-forward svc/k8ssandra-reaper-k8ssandra-reaper-service 9393:8080
    ```

1. Kill the process

    ```bash
    kill 80940
    ```

    **Output**:

    ```bash
    [3]  + terminated  kubectl port-forward svc/k8ssandra-reaper-k8ssandra-reaper-service 9393:8080
    ```

{{% alert title="Tip" color="success" %}}
Exiting the terminal instance will terminate all port forwarding services.
{{% /alert %}}

## Access K8ssandra Operator monitoring utilities {#monitoring}

K8ssandra Operator allows you to integrate with monitoring tools such as the following utilities:

* [Prometheus](https://prometheus.io/) a standard metrics collection and alerting tool.
* [Grafana](https://grafana.com/) a set of preconfigured dashboards displaying important K8ssandra metrics.

### Prometheus

To check on the health of your K8ssandraCluster using the K8ssandra Operator Prometheus interface:

1. Access the Prometheus home page at <http://127.0.0.1:9292>:

    ![Prometheus home page](prom-home.png)

1. From the **Status** menu, choose **Targets**.

1. Verify that the `stargate/0` and `k8ssandra/0` are in the state `UP`:

    ![Prometheus targets](prom-targets.png)

For more details on Prometheus, see the [Prometheus](https://prometheus.io/) web site.

### Grafana

To monitor the health and performance of your K8ssandraCluster using pre-configured Grafana dashboards:

1. Retrieve the Grafana login username using the `helm show` command:

    ```bash
    helm show values k8ssandra/k8ssandra | grep "adminUser"
    ```

    **Output**:

    ```bash
    admin
    ```

1. Retrieve the Grafana login password using the `helm show` command:

    ```bash
    helm show values k8ssandra/k8ssandra | grep "adminPassword"
    ```

    **Output**:

    ```bash
    secret
    ```

1. Access the Grafana login screen at <http://127.0.0.1:9191> and login using the username and password:

    ![Grafana login page](grafana-login.png)

1. Click the home button indicated by the arrow:

    ![Grafana home page](grafana-home.png)

1. Click the `K8ssandra Overview` dashboard:

    ![Grafana dashboards](grafana-dashboards.png)

1. The `K8ssandra Overview` dashboard is displayed:

    ![Grafana K8ssandra overview](grafana-k8overview.png)

1. Explore the other K8ssandra dashboards.

For more information see the [Grafana](https://grafana.com/) web site.

## Access Reaper {#reaper}

[Reaper](http://cassandra-reaper.io/) is an easy interface for managing K8ssandra cluster repairs.  Reaper is deployed as part of the K8ssandra Operator [install]({{< relref "/install/local" >}}). 

![Reaper](cass-reaper.png)

For details, start in the [Reaper]({{< relref "/components/reaper" >}}) topic. Then read about the [repair]({{< relref "/tasks/repair" >}}) tasks you can perform with Reaper.

## Upgrade K8ssandra Operator {#upgrade}

You can easily upgrade your K8ssandra software with the `helm repo update` command, or apply new settings with the `helm upgrade` command. For details, see [Upgrade K8ssandra]({{< relref "upgrade" >}}).

## Next steps

* [Components]({{< relref "components" >}}): Dig in to each deployed component of the K8ssandra stack and see how it communicates with the others.
* [Tasks]({{< relref "tasks" >}}): Need to get something done? Check out the Tasks topics for a helpful collection of outcome-based solutions.
* [Reference]({{< relref "reference" >}}): Explore the Custom Resource Definitions (CRDs) used by K8ssandra Operator.

We encourage developers and SREs to actively participate in the [K8ssandra community](https://k8ssandra.io/community/).
