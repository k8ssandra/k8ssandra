---
title: "K8ssandra site engineer quick start"
linkTitle: "Site engineers"
weight: 2
description: |
  Take a tour of your K8ssandra cluster using all the utilities!
---

**Completion time**: **10 minutes**.

{{% alert title="Important" color="warning" %}}
You **must** complete the steps in [Quick start]({{< relref "docs/getting-started" >}}) before continuing.
{{% /alert %}}

In this quick start we'll cover common tasks of interest for site reliability engineers including:

* [Accessing nodetool commands]({{< relref "#access-nodetool" >}}) like status, ring, and info.
* [Accessing K8ssandra utilities]({{< relref "#access-k8ssandra-utilities" >}}) like the Cassandra Reaper repair tool and Grafana metrics reporting using port forwarding.
* [Upgrading a K8ssandra cluster]({{< relref "#upgrade-your-k8ssandra-cluster" >}}): to enable access to K8ssandra from outside the K8s cluster via Traefik.

## Access the Cassandra nodetool utility {#access-nodetool}

Cassandra's nodetool utility is commonly used for a variety of monitoring and management tasks. You'll need to run nodetool on your K8ssandra cluster using the `kubectl exec` command, since there's no external stand alone option available.

To run nodetool commands:

1. Get a list of the running K8ssandra pods using `kubectl get`:

    ```bash
    kubectl get pods
    ```

    **Output**:

    ```bash
    NAME                                                  READY   STATUS      RESTARTS   AGE
    k8ssandra-cass-operator-6666588dc5-qpvzg              1/1     Running     3          45h
    k8ssandra-dc1-default-sts-0                           2/2     Running     0          115m
    k8ssandra-dc1-stargate-6f7f5d6fd6-sblt8               1/1     Running     12         45h
    k8ssandra-grafana-6c4f6577d8-fxfsd                    2/2     Running     6          45h
    k8ssandra-kube-prometheus-operator-5556885bd6-st4fp   1/1     Running     3          45h
    k8ssandra-reaper-k8ssandra-5b6cc959b7-zzlzr           1/1     Running     15         45h
    k8ssandra-reaper-k8ssandra-schema-47qzk               0/1     Completed   0          45h
    k8ssandra-reaper-operator-cc46fd5f4-85mk5             1/1     Running     4          45h
    prometheus-k8ssandra-kube-prometheus-prometheus-0     2/2     Running     7          45h
    ```

    The K8ssandra pod running Cassandra takes the form `<k8ssandra-cluster-name>-<datacenter-name>-default-sts-<n>` and, in the example above is `k8ssandra-dc1-default-sts-0` which we'll use throughout the following sections.

    {{% alert title="Tip" color="success" %}}
Although not applicable to this quick start, additional K8ssandra Cassandra nodes will increment the final `<n>` but the rest of the name will remain the same.
    {{% /alert %}}

1. Run [`nodetool status`](https://docs.datastax.com/en/cassandra-oss/3.x/cassandra/tools/toolsStatus.html), using the Cassandra node name `k8ssandra-dc1-default-sts-0`, and replacing `<k8ssandra-username>` and `<k8ssandra-password>` with the values you retrieved in [Retrieve K8ssandra superuser credentials]({{< relref "/docs/getting-started#superuser" >}}):

    ```bash
    kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- nodetool -u <k8ssandra-username> -pw <k8ssandra-password> status
    ```

    **Output**:

    ```bash
    Datacenter: dc1
    ===============
    Status=Up/Down
    |/ State=Normal/Leaving/Joining/Moving
    --  Address      Load       Owns    Host ID                               Token                                    Rack
    UN  10.244.1.12  215.3 KiB  ?       75e52e51-edc9-49f8-84f6-f044999ac130  -1080085985719557225                     default

    Note: Non-system keyspaces don't have the same replication settings, effective ownership information is meaningless
    ```

    {{% alert title="Tip" color="success" %}}
All nodes should have the status `UN` or "Up Normal."
    {{% /alert %}}

Other useful nodetool commands include:

* [`nodetool ring`](https://docs.datastax.com/en/cassandra-oss/3.x/cassandra/tools/toolsRing.html) which outputs all the tokens in the node:

    ```bash
    kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- nodetool -u <k8ssandra-username> -pw <k8ssandra-password> ring
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
    kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- nodetool -u <k8ssandra-username> -pw <k8ssandra-password> ring
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

## Upgrade your K8ssandra cluster

While you can use port forwarding as described above, to enable _persistent_ external access for applications, and K8ssandra features, you'll need to configure Ingress. In this section we'll demonstrate upgrading your existing K8ssandra installation to support Ingress using a Traefik Helm chart.

To upgrade K8ssandra with Traefik Ingress support:

1. Copy the following YAML into a file named `traefik.values.yaml`:

    ```yaml
    ---
    providers:
      kubernetesCRD:
        namespaces:
          - default
          - traefik
      kubernetesIngress:
        namespaces:
          - default
          - traefik

    ports:
      traefik:
        expose: true
        nodePort: 32090
      web:
        nodePort: 32080
      websecure:
        nodePort: 32443
      cassandra:
        expose: true
        port: 9042
        nodePort: 32091
      cassandrasecure:
        expose: true
        port: 9142
        nodePort: 32092
      sg-graphql:
        expose: true
        port: 8080
        nodePort: 30080
      sg-auth:
        expose: true
        port: 8081
        nodePort: 30081
      sg-rest:
        expose: true
        port: 8082
        nodePort: 30082

    service:
      type: NodePort
    ```

2. Install Traefik:

    ```bash
    helm install -f traefik.values.yaml traefik traefik/traefik -n traefik --create-namespace
    NAME: traefik
    LAST DEPLOYED: Fri Feb 19 12:40:36 2021
    NAMESPACE: traefik
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None
    ```

3. Copy the following YAML into a new file named `k8ssandra-traefik.values.yaml`:

    ```yaml
    ingress:
      traefik:
        enabled: true
        repair:
          enabled: true
          host: repair.127.0.0.1.nip.io
        stargate:
          enabled: true
          host: stargate.127.0.0.1.nip.io
    kube-prometheus-stack:
      prometheus:
        enabled: true
        prometheusSpec:
          externalUrl: http://localhost:9090/prometheus
          routePrefix: /prometheus
        ingress:
          enabled: true
          paths:
            - /prometheus
      grafana:
        enabled: true
        ingress:
          enabled: true
          path: /grafana
        adminUser: admin
        adminPassword: admin123
        grafana.ini:
          server:
            root_url: http://localhost:3000/grafana
           serve_from_sub_path: true
    ```

4. Upgrade the K8ssandra installation with Traefik support:

    ```bash
    helm upgrade -f k8ssandra-traefik.values.yaml k8ssandra k8ssandra/k8ssandra
    Release "k8ssandra" has been upgraded. Happy Helming!
    NAME: k8ssandra
    LAST DEPLOYED: Fri Feb 19 12:46:49 2021
    NAMESPACE: default
    STATUS: deployed
    REVISION: 2
    ```

5. Verify that the Traefik pod is `Running`:

    ```bash
    kubectl get po -A | grep "traefik"
    traefik       traefik-55996cbb6-v6s9p                1/1     Running     0          13m
    ```

The following K8ssandra application are now available at the following persistent URLs:

* Prometheus: <http://127.0.0.1.nip.io:8080/prometheus/>
* Grafana: <http://127.0.0.1.nip.io:8080/prometheus/>
* Cassandra Reaper: <http://repair.127.0.0.1.nip.io:8080/webui>

## Access K8ssandra utilities

K8ssandra includes the following bundled and customized utilities:

### Prometheus

<http://127.0.0.1.nip.io:8080/prometheus/> is a standard metrics collection and alerting tool. For more information, see the [Prometheus](https://prometheus.io/) web site.

### Grafana

<http://127.0.0.1.nip.io:8080/grafana/login> provides a set of pre-configured dashboards displaying important K8ssandra metrics. For more information see the [Grafana](https://grafana.com/) web site.

{{% alert title="Tip" color="success" %}}
If you've followed the configuration instructions in this quick start, use `admin` for the username and `admin123` for the password.
{{% /alert %}}
### Cassandra Reaper

<http://repair.127.0.0.1.nip.io:8080/webui> provides an easy interface for managing K8ssandra cluster repairs. For details, see the [Cassandra Reaper](http://cassandra-reaper.io/) web site.

### Medusa backup and restore

K8ssandra provides a complete backup and restore solution using [Medusa](https://github.com/thelastpickle/cassandra-medusa). For detailed configuration and usage instructions, see [Backup and restore Cassandra]({{< relref "docs/topics/restore-a-backup" >}}).

## Next

* For detailed information on additional K8ssandra tasks, see [Tasks]({{< relref "docs/topics" >}}).
* For a list of frequently asked questions, see the [FAQs]({{< relref "docs/faqs" >}}).
* For detailed information on K8ssandra, see [Architecture]({{< relref "docs/architecture" >}}).
* For information on the various K8ssandra Helm charts, see [Helm chart references]({{< relref "docs/reference" >}}).
* If you'd like to contribute to K8ssandra, see [Contribution guidelines]({{< relref "docs/contribution-guidelines" >}}).
