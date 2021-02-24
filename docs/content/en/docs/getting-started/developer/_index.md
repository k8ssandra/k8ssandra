---
title: "K8ssandra developer quick start"
linkTitle: "Developers"
weight: 1
description: |
  Get up and coding with K8ssandra!
---

**Completion time**: **10 minutes**.

In this quick start, you'll configure your K8ssandra instance so you can:

* [Accessing K8ssandra using CQLSH]({{< relref "#access-k8ssandra-using-cqlsh" >}}): Access K8ssandra using the standard Cassandra CQLSH utility.
* [Accessing K8ssandra via Stargate]({{< relref "#access-k8ssandra-using-the-stargate-api" >}}): Access K8ssandra using the Stargate API and the GraphQL Playground.

## Set up port forwarding

In order to access Cassandra outside of the Kubernetes (K8s) cluster, if you don't have an Ingress setup as described in [Configure Ingress]({{< relref "/docs/topics/ingress" >}}), you'll need to configure port forwarding for both CQLSH and Stargate access.

Begin by getting a list of your K8ssandra K8s pods and ports:

```bash
kubectl get services
NAME                                        TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                                                 AGE
cass-operator-metrics                       ClusterIP   10.99.98.218     <none>        8383/TCP,8686/TCP                                       21h
k8ssandra-dc1-all-pods-service              ClusterIP   None             <none>        9042/TCP,8080/TCP,9103/TCP                              21h
k8ssandra-dc1-service                       ClusterIP   None             <none>        9042/TCP,9142/TCP,8080/TCP,9103/TCP,9160/TCP            21h
k8ssandra-dc1-stargate-service              ClusterIP   10.106.70.148    <none>        8080/TCP,8081/TCP,8082/TCP,8084/TCP,8085/TCP,9042/TCP   21h
k8ssandra-grafana                           ClusterIP   10.96.120.157    <none>        80/TCP                                                  21h
k8ssandra-kube-prometheus-operator          ClusterIP   10.97.21.175     <none>        443/TCP                                                 21h
k8ssandra-kube-prometheus-prometheus        ClusterIP   10.111.184.111   <none>        9090/TCP                                                21h
k8ssandra-reaper-k8ssandra-reaper-service   ClusterIP   10.104.46.103    <none>        8080/TCP                                                21h
k8ssandra-seed-service                      ClusterIP   None             <none>        <none>                                                  21h
kubernetes                                  ClusterIP   10.96.0.1        <none>        443/TCP                                                 21h
prometheus-operated                         ClusterIP   None             <none>        9090/TCP                                                2
```

In the output above, the pod of interest is:

* **k8ssandra-dc1-stargate-service**: The K8ssandra Stargate where the name is a combination of the K8ssandra cluster name you specified during the Helm install, `k8ssandra`, the datacenter name, `dc1` and the postfix, `-service`. This service listens on the ports:
  * **8080/TCP**: GraphQL interface
  * **8081/TCP**: REST authorization service for generating tokens
  * **8082/TCP**: REST interface
  * **8084/TCP**: Health check (/healthcheck, /checker/liveness, /checker/readiness)
  * **8085/TCP**: Metrics (/metrics)
  * **9042/TCP**: CQL service

Those are the ports we'll need to forward for CQLSH and Stargate access.

To configure port forwarding:

1. Open a new terminal.

1. Run the `kubectl port-forward` command:

    ```bash
    kubectl port-forward svc/k8ssandra-dc1-stargate-service 8080 8081 8082 8084 8085 9042
    Forwarding from 127.0.0.1:8080 -> 8080
    Forwarding from [::1]:8080 -> 8080
    Forwarding from 127.0.0.1:8081 -> 8081
    Forwarding from [::1]:8081 -> 8081
    Forwarding from 127.0.0.1:8082 -> 8082
    Forwarding from [::1]:8082 -> 8082
    Forwarding from 127.0.0.1:8084 -> 8084
    Forwarding from [::1]:8084 -> 8084
    Forwarding from 127.0.0.1:8085 -> 8085
    Forwarding from [::1]:8085 -> 8085
    Forwarding from 127.0.0.1:9042 -> 9042
    Forwarding from [::1]:9042 -> 9042
    ```

1. Leave the terminal running in the background.

{{% alert title="Important" color="warning" %}}
If you close the terminal for any reason, you'll shut down the port forwarding service.
{{% /alert %}}

## Access Cassandra using the Stargate API

Stargate is an open-source data gateway providing common API interfaces for backend databases. You can experiment with Stargate using [K8ssandra GraphQL Playground](http://stargate.127.0.0.1.nip.io:8080/playground).

For more detailed configuration instructions and a usage example, see [Access the Stargate API]({{< relref "docs/topics/stargate" >}}).

{{% alert title="Tip" color="success" %}}
Make a note of the K8ssandra superuser name and password for use in the K8ssandra Stargate topic.
{{% /alert %}}

For complete details on Stargate, see the [Stargate documentation](https://stargate.io/docs/stargate/1.0/quickstart/quickstart.html).

## Access Cassandra using CQLSH

To access K8ssandra using the stand alone CQLSH utility:

1. Make sure you have [Python 2.7](https://www.python.org/download/releases/2.7/) installed on your system.

1. Download CQLSH from the  [DataStax download site](https://downloads.datastax.com/#cqlsh) choosing the version for **DataStax Astra**.

1. Connect to Cassandra replacing `<k8ssandra-username>` and `<k8ssandra-password>` with the values you retrieved in [Retrieve K8ssandra superuser credentials]({{< relref "/docs/getting-started#superuser" >}}):

    ```bash
    cqlsh -u <k8ssandra-username> -p <k8ssandra-password>
    Connected to k8ssandra at 127.0.0.1:9042.
    [cqlsh 6.8.0 | Cassandra 3.11.6 | CQL spec 3.4.4 | Native protocol v4]
    Use HELP for help.
    k8ssandra-superuser@cqlsh>
   ```

1. Populate a CQL file with the following data and save as `test-data.cql`:

    ```sql
    CREATE KEYSPACE k8ssandra_test  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
    USE k8ssandra_test;
    CREATE TABLE users (email text primary key, name text, state text);
    insert into users (email, name, state) values ('alice@example.com', 'Alice Smith', 'TX');
    insert into users (email, name, state) values ('bob@example.com', 'Bob Jones', 'VA');
    insert into users (email, name, state) values ('carol@example.com', 'Carol Jackson', 'CA');
    insert into users (email, name, state) values ('david@example.com', 'David Yang', 'NV');
    ```

1. Import the test data into the Cassandra node using the CQL `SOURCE` command:

    ```bash
    SOURCE 'test-data.cql';
    ```

1. Query the data:

    ```sql
    cqlsh> SELECT * FROM k8ssandra_test.users;

     email             | name          | state
    -------------------+---------------+-------
     alice@example.com |   Alice Smith |    TX
       bob@example.com |     Bob Jones |    VA
     david@example.com |    David Yang |    NV
     carol@example.com | Carol Jackson |    CA

    (4 rows)
    ```

For complete details on Cassandra, CQL and CQLSH, see the [Apache Cassandra](https://cassandra.apache.org/) web site.
