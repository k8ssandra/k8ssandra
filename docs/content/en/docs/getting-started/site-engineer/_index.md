---
title: "K8ssandra site engineer quick start"
linkTitle: "Site engineers"
weight: 2
description: |
  Kick the tires and take K8ssandra for a spin!
---

**Completion time**: **10 minutes**.

* [Upgrading K8ssandra with Ingress support]({{< relref "#configure-ingress" >}}): Enabling access to K8ssandra from outside the K8s cluster via Traefik.
* [Accessing K8ssandra utilities]({{< relref "#access-k8ssandra-utilities" >}}): Accessing useful utilities like the Cassandra Reaper repair tool and Grafana metrics reporting.
* [Starting and stopping K8ssandra]({{< relref "#cassandra-operations" >}}): Cleanly stopping and restarting the K8ssandra pod.

## Configure Ingress

Right now, you can only interact with your Cassandra node within the K8s cluster using `kubectl` commands which is a pretty severe limitation. To enable external access for applications as well as enable access to K8ssandra features like Grafana dashboards, the Reaper Repair UI, and the Stargate GraphQL playground, you'll need to configure Ingress. In this section we'll upgrade your existing K8ssandra installation to support Ingress using a Traefik Helm chart.

{{% alert title="Tip" color="success" %}}
Kubernetes Ingress configuration is a complicated topic. Really all that is happening in this upgrade is that a set of external URIs and ports are being mapped to internal Kubernetes resources. It's not necessary to understand all of the configuration particulars in order to successfully complete this section.
{{% /alert %}}

To upgrade K8ssandra:

1. Copy the following YAML into a new file named `k8ssandra-traefik.values.yaml`:

    ```yaml
    cassandra:
      version: "3.11.7"
      auth:
        enabled: false
      clusterName: k8ssandra
      datacenters:
      - name: dc1
        size: 1
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

1. Install Traefik:

    ```bash
    helm install -f traefik.values.yaml traefik traefik/traefik -n traefik --create-namespace
    NAME: traefik
    LAST DEPLOYED: Fri Feb 19 12:40:36 2021
    NAMESPACE: traefik
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None
    ```

1. Upgrade the K8ssandra pod with Traefik support:

    ```bash
    helm upgrade -f k8ssandra-traefik.values.yaml k8ssandra k8ssandra/k8ssandra
    Release "k8ssandra" has been upgraded. Happy Helming!
    NAME: k8ssandra
    LAST DEPLOYED: Fri Feb 19 12:46:49 2021
    NAMESPACE: default
    STATUS: deployed
    REVISION: 2
    ```

1. Verify that the Traefik pod is `Running`:

    ```bash
    kubectl get po -A | grep "traefik"
    traefik       traefik-55996cbb6-v6s9p                1/1     Running     0          13m
    ```

Now you're ready to access the full suite of K8ssandra supporting applications.

{{% alert title="Tip" color="success" %}}
If you've followed the configuration instructions in this quick start, the URLs in the sections below should take you directly to the associated utilities.
{{% /alert %}}

## Access K8ssandra utilities

K8ssandra includes the following bundled and customized utilities:

### Cassandra Reaper

[Cassandra Reaper](<"http://repair.127.0.0.1.nip.io:8080/webui">) provides an easy interface for managing K8ssandra cluster repairs. For details, see the [Cassandra Reaper](http://cassandra-reaper.io/) web site.

### Prometheus

[Prometheus](http://127.0.0.1.nip.io:8080/prometheus/) is a standard metrics collection and alerting tool. For more information, see the [Prometheus](https://prometheus.io/) web site.

### Grafana

[Grafana](http://127.0.0.1.nip.io:8080/grafana/login) provides a set of pre-configured dashboards displaying important K8ssandra metrics. For more information see the [Grafana](https://grafana.com/) web site.

{{% alert title="Tip" color="success" %}}
If you've followed the configuration instructions in this quick start, use `admin` for the username and `admin123` for the password.
{{% /alert %}}

### Medusa backup and restore

K8ssandra provides a complete backup and restore solution using [Medusa](https://github.com/thelastpickle/cassandra-medusa). For detailed configuration and usage instructions, see [Backup and restore Cassandra]({{< relref "docs/topics/restore-a-backup" >}}).

## Stopping and starting Cassandra {#cassandra-operations}

Before shutting down your Kubernetes cluster, you'll want to make sure you cleanly shut down your Cassandra datacenters. You can do that using the `kubectl patch` command and setting the `spec:stopped` property to either `true` (stopped) or `false` (running).

### Shut down Cassandra

To shut down a Cassandra datacenter:

```bash
kubectl patch cassdc <datacenter-name> --type merge -p '{"spec":{"stopped":true}}'
```

Example:

```bash
kubectl patch cassdc dc1 --type merge -p '{"spec":{"stopped":true}}'
cassandradatacenter.cassandra.datastax.com/dc1 patched
```

### Start up Cassandra

To start up a Cassandra datacenter

```bash
kubectl patch cassdc <datacenter-name> --type merge -p '{"spec":{"stopped":false}}'
```

Example:

```bash
kubectl patch cassdc dc1 --type merge -p '{"spec":{"stopped":false}}'
cassandradatacenter.cassandra.datastax.com/dc1 patched
```

## Nodetool fragment

1. Run `nodetool status`, using the Cassandra node name `k8ssandra-dc1-default-sts-0`, and passing the superuser name and password. Verify that the node is in the state `UN` or Up Normal:

    ```bash
    kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- nodetool -u k8ssandra-superuser -pw 6WH3pp8scIOvyJqGc0m_Ubb-Ft07lmKkYb3Ye_hXpc6TiscaQtwcuA status
    Datacenter: dc1
    ===============
    Status=Up/Down
    |/ State=Normal/Leaving/Joining/Moving
    --  Address      Load       Owns    Host ID                               Token                                    Rack
    UN  10.244.1.12  215.3 KiB  ?       75e52e51-edc9-49f8-84f6-f044999ac130  -1080085985719557225                     default

    Note: Non-system keyspaces don't have the same replication settings, effective ownership information is meaningless
    ```
