---
title: "Quickstart for developers"
linkTitle: "Developers"
weight: 1
description: "Get up and coding with K8ssandra by exposing access to Stargate and CQL APIs!"
---

**Completion time**: **10 minutes**.

{{% alert title="Tip" color="success" %}}
Be sure to first complete one of the [K8ssandra install]({{< relref "/install" >}}) options (locally or on a cloud provider) before performing these post-install steps. 
{{% /alert %}}

In this quickstart for developers, we'll cover:

* [Setting up port forwarding]({{< relref "#set-up-port-forwarding" >}}) to access Stargate services and CQLSH outside your Kubernetes (K8s) cluster.
* [Accessing Cassandra using Stargate]({{< relref "#access-cassandra-using-the-stargate-apis" >}}) by creating an access token, and using Stargates's REST, GraphQL and document interfaces.
* [Accessing Cassandra using CQLSH]({{< relref "#access-cassandra-using-cqlsh" >}}) including some basic CQL commands.

## Set up port forwarding

In order to access Apache Cassandra® outside of the K8s cluster, you'll need to utilize port forwarding unless ingress is [configured]({{< relref "/tasks/connect/ingress" >}}).

Begin by getting a list of your K8ssandra K8s services and ports:

```bash
kubectl get services
```

**Output**:

```bash
NAME                                   TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)                                                 AGE
cass-operator-metrics                  ClusterIP   10.80.3.92     <none>        8383/TCP,8686/TCP                                       21h
k8ssandra-dc1-all-pods-service         ClusterIP   None           <none>        9042/TCP,8080/TCP,9103/TCP                              21h
k8ssandra-dc1-service                  ClusterIP   None           <none>        9042/TCP,9142/TCP,8080/TCP,9103/TCP,9160/TCP            21h
k8ssandra-dc1-stargate-service         ClusterIP   10.80.13.197   <none>        8080/TCP,8081/TCP,8082/TCP,8084/TCP,8085/TCP,9042/TCP   21h
k8ssandra-grafana                      ClusterIP   10.80.7.168    <none>        80/TCP                                                  21h
k8ssandra-kube-prometheus-operator     ClusterIP   10.80.8.109    <none>        443/TCP                                                 21h
k8ssandra-kube-prometheus-prometheus   ClusterIP   10.80.2.44     <none>        9090/TCP                                                21h
k8ssandra-reaper-reaper-service        ClusterIP   10.80.5.77     <none>        8080/TCP                                                21h
k8ssandra-seed-service                 ClusterIP   None           <none>        <none>                                                  21h
kubernetes                             ClusterIP   10.80.0.1      <none>        443/TCP                                                 23h
prometheus-operated                    ClusterIP   None           <none>        9090/TCP                                                21h
```

In the output above, the service of interest is:

* **k8ssandra-dc1-stargate-service**: The K8ssandra Stargate service where the name is a combination of the K8ssandra cluster name you specified during the Helm install, `k8ssandra`, the datacenter name, `dc1` and the postfix, `-service`. This service listens on the ports:
  * **8080/TCP**: GraphQL interface
  * **8081/TCP**: REST authorization service for generating tokens
  * **8082/TCP**: REST interface
  * **9042/TCP**: CQL service

Those are the ports we'll need to forward for CQLSH and Stargate access.

To configure port forwarding:

1. Open a new terminal.

2. Run the `kubectl port-forward` command in the background:

    ```bash
    kubectl port-forward svc/k8ssandra-dc1-stargate-service 8080 8081 8082 9042 &
    ```

    **Output**:

    ```bash
    [1] 80940

    ~/
    Forwarding from 127.0.0.1:8080 -> 8080
    Forwarding from [::1]:8080 -> 8080
    Forwarding from 127.0.0.1:8081 -> 8081
    Forwarding from [::1]:8081 -> 8081
    Forwarding from 127.0.0.1:8082 -> 8082
    Forwarding from [::1]:8082 -> 8082
    ```

### Terminate port forwarding

To terminate the port forwarding service:

1. Get the process ID:

    ```bash
    jobs -l
    ```

    **Output**:

    ```bash
    [1]  + 80940 running    kubectl port-forward svc/k8ssandra-dc1-stargate-service 8080 8081 8082
    ```

1. Kill the process

    ```bash
    kill 80940
    ```

    **Output**:

    ```bash
    [1]  + terminated  kubectl port-forward svc/k8ssandra-dc1-stargate-service 8080 8081 8082
    ```

{{% alert title="Tip" color="success" %}}
Exiting the terminal instance will terminate the port forwarding service.
{{% /alert %}}

## Access Cassandra using the Stargate APIs

[Stargate](https://stargate.io/) provides APIs, data types and access methods that bring new capabilities to existing databases. Currently Stargate adds Document, REST and GraphQL APIs for CRUD access to data stored in Apache Cassandra® and there are many more APIs coming soon. Separating compute and storage also has benefits for maximizing resource consumption in cloud environments. When using Stargate with Cassandra, you can offload the request coordination overhead from your storage instances onto Stargate instances which has shown latency improvements in preliminary testing.

To access K8ssandra using Stargate:

1. Generate a Stargate access token replacing `<k8ssandra-username>` and `<k8ssandra-password>` with the values you retrieved in [Retrieve K8ssandra superuser credentials]({{< relref "/install/local#superuser" >}}):

    ```bash
    curl -L -X POST 'http://localhost:8081/v1/auth' -H 'Content-Type: application/json' --data-raw '{"username": "<k8ssandra-username>", "password": "<k8ssandra-password>"}'
    ```

    **Output**:

    ```json
    {"authToken":"<access-token>"}
    ```

1. Use `<access-token>` to populate the `x-cassandra-token` header for all Stargate requests.

Once you've got the access token, take a look at the following Stargate access options:

* [Access Document Data API]({{< relref "develop#access-document-data-api" >}})
* [Access REST Data API]({{< relref "develop#access-rest-data-api" >}})
* [Access GraphQL Data API]({{< relref "develop#access-graphql-data-api" >}})

You can access the following interfaces to make development easier as well:

* Stargate swagger UI: <http://127.0.0.1:8082/swagger-ui>
* GraphQL Playground: <http://127.0.0.1:8080/playground>

For complete details on Stargate, see the [Stargate documentation](https://stargate.io/docs/stargate/1.0/quickstart/quickstart.html).

## Access Cassandra using CQLSH

If you're familiar with Cassandra, then you're familiar with CQLSH. You can download a full-featured [stand alone CQLSH utility](https://docs.datastax.com/en/dse/6.8/cql/cql/cql_using/startCqlshStandalone.html) from Datastax and use that to interact with K8ssandra as if you were in a native Cassandra environment.

To access K8ssandra using the stand alone CQLSH utility:

1. Make sure you have [Python 2.7](https://www.python.org/download/releases/2.7/) installed on your system.

1. Download CQLSH from the  [DataStax download site](https://downloads.datastax.com/#cqlsh) choosing the version for **DataStax Astra**.

1. Connect to Cassandra replacing `<k8ssandra-username>` and `<k8ssandra-password>` with the values you retrieved in [Retrieve K8ssandra superuser credentials]({{< relref "/install/local#superuser" >}}):

    ```bash
    cqlsh -u <k8ssandra-username> -p <k8ssandra-password>
    ```

    **Output**:

    ```bash
    Connected to k8ssandra at 127.0.0.1:9042.
    [cqlsh 6.8.0 | Cassandra 3.11.6 | CQL spec 3.4.4 | Native protocol v4]
    Use HELP for help.
    k8ssandra-superuser@cqlsh>
   ```

1. Create a new keyspace, `k8ssandra_test`, using [CREATE KEYSPACE](https://docs.datastax.com/en/cql-oss/3.x/cql/cql_reference/cqlCreateKeyspace.html):

    ```sql
    CREATE KEYSPACE k8ssandra_test  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
    ```

1. Switch to the new keyspace using [USE](https://docs.datastax.com/en/cql-oss/3.x/cql/cql_reference/cqlUse.html):

    ```sql
    USE k8ssandra_test;
    ```

1. Create a new table, `users` using [CREATE TABLE](https://docs.datastax.com/en/cql-oss/3.x/cql/cql_reference/cqlCreateTable.html#cqlCreateTable)

    ```sql
    CREATE TABLE users (email text primary key, name text, state text);
    ```

1. Insert some sample data into the new table using [INSERT](https://docs.datastax.com/en/cql-oss/3.x/cql/cql_reference/cqlInsert.html)

    ```sql
    INSERT INTO users (email, name, state) values ('alice@example.com', 'Alice Smith', 'TX');
    INSERT INTO users (email, name, state) values ('bob@example.com', 'Bob Jones', 'VA');
    INSERT INTO users (email, name, state) values ('carol@example.com', 'Carol Jackson', 'CA');
    INSERT INTO users (email, name, state) values ('david@example.com', 'David Yang', 'NV');
    ```

1. Query the data using [SELECT](https://docs.datastax.com/en/cql-oss/3.x/cql/cql_reference/cqlSelect.html) and validate the return results:

    ```sql
    cqlsh> SELECT * FROM k8ssandra_test.users;
    ```

    **Output**:

    ```sql
     email             | name          | state
    -------------------+---------------+-------
     alice@example.com |   Alice Smith |    TX
       bob@example.com |     Bob Jones |    VA
     david@example.com |    David Yang |    NV
     carol@example.com | Carol Jackson |    CA

    (4 rows)
    ```

1. When you're done, exit CQLSH using `QUIT`:

    ```sql
    cqlsh> QUIT;
    ```

For complete details on Cassandra, CQL and CQLSH, see the [Apache Cassandra](https://cassandra.apache.org/) web site.

## Next steps

* [FAQs]({{< relref "faqs" >}}): If you're new to K8ssandra, these FAQs are for you. 
* [Components]({{< relref "components" >}}): Dig in to each deployed component of the K8ssandra stack and see how it communicates with the others.
* [Tasks]({{< relref "tasks" >}}): Need to get something done? Check out the Tasks topics for a helpful collection of outcome-based solutions.
* [Reference]({{< relref "reference" >}}): Explore the K8ssandra configuration interface (Helm charts), the available options, and a Glossary.

We encourage developers to actively participate in the [K8ssandra community](https://k8ssandra.io/community/).
