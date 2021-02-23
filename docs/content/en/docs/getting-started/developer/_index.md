---
title: "K8ssandra developer quick start"
linkTitle: "Developers"
weight: 1
description: |
  Kick the tires and take K8ssandra for a spin!
---

* [Accessing K8ssandra using CQLSH]({{< relref "#access-k8ssandra-using-cqlsh" >}}): Accessing K8ssandra using the standard Cassandra CQLSH utility.
* [Accessing K8ssandra via Stargate]({{< relref "#access-k8ssandra-using-the-stargate-api" >}}): Accessing K8ssandra using the Stargate API and the GraphQL Playground.

## Access Cassandra using CQLSH

Now that K8ssandra is installed and running as expected, we can interact with the actual Cassandra node, `k8ssandra-dc1-default-sts-0`, using CQLSH. We'll do this using the `kubectl exec` utility which doesn't require any fancy Ingress configuration.

Let's prepare some data, copy it to the Cassandra node, and then run a query using an interactive CQLSH session:

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

2. Copy `test-data.cql` to the `/tmp/` directory of the Cassandra node using `kubectl cp`:

    ```bash
    kubectl cp ./test-data.cql k8ssandra-dc1-default-sts-0:/tmp/ -c cassandra
    ```

3. Import the test data into the Cassandra node using `cqlsh -f`, providing the superuser name and password from the previous section:

    ```bash
    kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- cqlsh -u k8ssandra-superuser -p PGo8kROUgAJOa8vhjQrE49Lgruw7s32HCPyVvcfVmmACW8oUhfoO9A -f /tmp/test-data.cql
    ```

4. Open an interactive `cqlsh` session on the node, providing the superuser name and password from the previous section:

    ```bash
    kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- cqlsh -u k8ssandra-superuser -p PGo8kROUgAJOa8vhjQrE49Lgruw7s32HCPyVvcfVmmACW8oUhfoO9A
    Connected to k8ssandra at 127.0.0.1:9042.
    [cqlsh 5.0.1 | Cassandra 3.11.7 | CQL spec 3.4.4 | Native protocol v4]
    Use HELP for help.
    cqlsh>
   ```

5. Query the data:

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

For complete details on Cassandra, see the [Apache Cassandra](https://cassandra.apache.org/) web site.

## Access Cassandra using the Stargate API

Stargate is an open-source data gateway providing common API interfaces for backend databases. You can experiment with Stargate using the [K8ssandra GraphQL Playground](http://stargate.127.0.0.1.nip.io:8080/playground) 

For more detailed configuration instructions and a usage example, see [Access the Stargate API]({{< relref "docs/topics/stargate" >}}).

{{% alert title="Tip" color="success" %}}
Make a note of the K8ssandra superuser name and password for use in the K8ssandra Stargate topic.
{{% /alert %}}

For complete details on Stargate, see the [Stargate documentation](https://stargate.io/docs/stargate/1.0/quickstart/quickstart.html).
