---
title: "Backup and restore with Azure Storage"
linkTitle: "Azure Blob Storage"
no_list: true
weight: 3
description: Use Medusa to backup and restore Apache Cassandra® data in Kubernetes to Azure Storage.
---

Medusa is a Cassandra backup and restore tool. It's packaged with K8ssandra and supports a variety of backends,
including Azure Storage.

## Introduction

Azure Storage is a RESTful online file storage service for storing and accessing data on Microsoft Azure / AKS
infrastructure.

For details about storing blobs on Azure, read the [Azure Blob Storage
documentation](https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blobs-introduction).

## Create a storage account for backups

In order to perform backups in Azure, Medusa needs to use a storage account with appropriate permissions.

Using the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/), run the following command to create the
`medusabackups` storage account:

```bash
az storage account create \
  --name medusabackups \
  --resource-group storage-resource-group \
  --location eastus \
  --sku Standard_RAGRS \
  --kind StorageV2
```

See [this article](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create) for details on storage
account creation on Azure.

## Create the storage account credentials file

List the access keys that were generated for the storage account:

```bash
az storage account keys list \
    --account-name medusabackups \
    --resource-group storage-resource-group
```

Create a file on your local machine called `credentials.json` with the following contents:

```bash
{
    "storage_account": "medusabackups",
    "key": "<YOUR_KEY>"
}
```

`YOUR_KEY` can be any of the generated keys listed with the command above. You can automate the creation of this file
with the following command:

```bash
az storage account keys list \
    --account-name medusabackups \
    --resource-group storage-resource-group \
    --query "[0].value|{storage_account:'medusabackups',key:@}" > credentials.json
```

Some accounts need to use a non-standard host name to contact the Azure Blob Storage REST API; for example, the host
name to use for Azure US Government accounts is `<storageAccount>.blob.core.usgovcloudapi.net`. If you are in this case,
you will need to define two other fields in the credentials file, `host` and `connection_string`:

```bash
{
    "storage_account": "medusabackups",
    "key": "<YOUR_KEY>",
    "host": "<YOUR_HOST>",
    "connection_string": "<YOUR_CONNECTION_STRING>"
}
```

The connection string to use can be found with the following command:

```bash
az storage account show-connection-string \
    --name medusabackups \
    --resource-group storage-resource-group
```

## Store the credentials file as a Kubernetes secret

Now, using `kubectl`, push the `credentials.json` file to your Kubernetes cluster as a secret under the `k8ssandra`
namespace:

```bash
kubectl create secret generic medusa-bucket-key \
    --from-file=medusa_azure_credentials.json=./credentials.json \
    -n k8ssandra
```

The secret itself can be named anything, but it must contain one single entry named `medusa_azure_credentials.json`.
**Any other name would result in Medusa not finding the credentials file.**

Check that the secret is correct:

```bash
kubectl describe secret medusa-bucket-key -n k8ssandra
```

Expected output:

```
Name:         medusa-bucket-key
Namespace:    k8ssandra
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
medusa_azure_credentials.json:  142 bytes
```
## Create an Azure storage container 

Next, create the Azure container that will store Medusa backups for your K8ssandra cluster. The following command will
create a container named `k8ssandra-backups` inside the `medusabackups` storage account:

```bash
az storage container create \
    --name k8ssandra-backups \
    --account-name medusabackups \
    --account-key "<YOUR_KEY>" \
    --resource-group storage-resource-group
```

Expected output:

```json
{
  "created": true
}
```

## Configure Medusa to use the right Azure storage account and container

Lastly, create or update your K8ssandra Helm values file to instruct Medusa to use the Azure storage account and
container that we just created:

```yaml
# k8ssandra-medusa-azure.yaml
# other settings omitted for brevity
medusa:
  enabled: true
  multiTenant: true
  storage: azure_blobs
  bucketName: k8ssandra-backups
  storageSecret: medusa-bucket-key
```

Important:

* `medusa.storage` must be `azure_blobs`;
* `medusa.bucketName` must be the name of the target Azure Blob container (`k8ssandra-backups` in our example);
* `medusa.storageSecret` must be the name of the secret created above, containing the storage account name and access 
  key.

## Deploy K8ssandra

Now that you have Azure Blob Storage set up, install K8ssandra and enable Medusa backups on Azure by providing the
Helm values file created in the previous step:

```bash
helm install k8ssandra k8ssandra/k8ssandra -f ./k8ssandra-medusa-azure.yaml -n k8ssandra
```

Now wait for the Cassandra cluster to be ready by using the following `wait` command:

```bash
kubectl wait --for=condition=Ready cassandradatacenter/dc1 --timeout=900s -n k8ssandra
```

When ready, you should now see a list of pods. Example:

```bash
kubectl get pods -n k8ssandra
```

**Output:**

```bash
NAME                                                  READY   STATUS      RESTARTS   AGE
k8ssandra-cass-operator-547845459-dwg68               1/1     Running     0          6m36s
k8ssandra-dc1-default-sts-0                           3/3     Running     0          5m56s
k8ssandra-dc1-stargate-776f88f945-p9twg               0/1     Running     0          6m36s
k8ssandra-grafana-75b9cb64cc-kndtc                    2/2     Running     0          6m36s
k8ssandra-kube-prometheus-operator-5bdd97c666-qz5vv   1/1     Running     0          6m36s
k8ssandra-medusa-operator-d766d5b66-wjt7j             1/1     Running     0          6m36s
k8ssandra-reaper-5f9bbfc989-j59xk                     1/1     Running     0          2m48s
k8ssandra-reaper-operator-858cd89bdd-7gfjj            1/1     Running     0          6m36s
k8ssandra-reaper-schema-4gshj                         0/1     Completed   0          3m3s
prometheus-k8ssandra-kube-prometheus-prometheus-0     2/2     Running     1          6m32s
```

## Create some data and back it up

Next, let's define some sample data in Cassandra by creating a `test_data.cql` file that contains DDL and DML statements:

```cql
CREATE KEYSPACE medusa_test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
USE medusa_test;
CREATE TABLE users (email text primary key, name text, state text);
INSERT INTO users (email, name, state) VALUES ('alice@example.com', 'Alice Smith', 'TX');
INSERT INTO users (email, name, state) VALUES ('bob@example.com', 'Bob Jones', 'VA');
INSERT INTO users (email, name, state) VALUES ('carol@example.com', 'Carol Jackson', 'CA');
INSERT INTO users (email, name, state) VALUES ('david@example.com', 'David Yang', 'NV');
```

Copy the CQL file into the Cassandra pod; that is, the StatefulSet one, which contains `-sts-` in its name:

```bash
kubectl cp test_data.cql k8ssandra-dc1-default-sts-0:/tmp -c cassandra -n k8ssandra
```

Let’s now run the uploaded CQL script and check that you can read the data. 

```bash
K8S_PWD=$(kubectl get secret k8ssandra-superuser -o jsonpath="{.data.password}" -n k8ssandra | base64 --decode)

kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- \
    cqlsh -u k8ssandra-superuser -p $K8S_PWD -f /tmp/test_data.cql

kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- \
    cqlsh -u k8ssandra-superuser -p $K8S_PWD -e "SELECT * FROM medusa_test.users"
```

**Output:**

```bash
 email             | name          | state
-------------------+---------------+-------
 alice@example.com |   Alice Smith |    TX
   bob@example.com |     Bob Jones |    VA
 david@example.com |    David Yang |    NV
 carol@example.com | Carol Jackson |    CA

(4 rows)
```

Now backup this data, and check that files get created in your Azure Storage container. 

```bash
helm install my-backup k8ssandra/backup -n k8ssandra \
    --set name=backup1 \
    --set cassandraDatacenter.name=dc1
```

Because the backup operation is asynchronous, you can monitor its completion by running the following command:

```bash
kubectl get cassandrabackup backup1 -n k8ssandra -o jsonpath={.status.finishTime}
```

As long as this command doesn’t output a date and time, you know that the backup is still running. With the amount of
data present, this should complete quickly.

Once the backup is done, check that Medusa created files in the `k8ssandra-backups` container: 

```bash
az storage blob list \
    --container-name k8ssandra-backups \
    --account-name medusabackups \
    --account-key "<YOUR_KEY>" \
    --query "[].{name:name}" \
    --output tsv
```

## Deleting the data and restoring the backup

Now delete the data by truncating the table, and check that the table is empty.

```bash
kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- \
    cqlsh -u k8ssandra-superuser -p $K8S_PWD -e "TRUNCATE medusa_test.users"

kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- \
     cqlsh -u k8ssandra-superuser -p $K8S_PWD -e "SELECT * FROM medusa_test.users"
```

**Output:**

```cql
 email | name | state
-------+------+-------

(0 rows)
```

Now restore the backup taken previously:

```bash
helm install restore-test k8ssandra/restore -n k8ssandra \
    --set name=restore-backup1 \
    --set backup.name=backup1 \
    --set cassandraDatacenter.name=dc1 \
```

The restore operation will take a little longer because it requires K8ssandra to stop the StatefulSet pod and perform
the restore as part of the init containers, before the Cassandra container can start. You can monitor progress using
this command:

```bash
watch -d kubectl get cassandrarestore restore-backup1 -o jsonpath={.status} -n k8ssandra
```

The restore operation is fully completed once the `finishTime` value appears in the output. Example:

```bash
{"finishTime":"2021-03-30T13:58:36Z","restoreKey":"83977399-44dd-4752-b4c4-407273f0339e","startTime":"2021-03-30T13:55:35Z"}
```

Verify that you can now read the restored data from the previously truncated table:

```bash
kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- \
    cqlsh -u k8ssandra-superuser -p $K8S_PWD -e "SELECT * FROM medusa_test.users"
```

**Output:**

```bash
 email             | name          | state
-------------------+---------------+-------
 alice@example.com |   Alice Smith |    TX
   bob@example.com |     Bob Jones |    VA
 david@example.com |    David Yang |    NV
 carol@example.com | Carol Jackson |    CA

(4 rows)
```

Success! You’ve successfully restored your lost data in just a few commands.

## Next steps

See the following reference topics:

* [Medusa Operator Helm Chart]({{< relref "/reference/helm-charts/medusa-operator" >}})
* [Backup Helm Chart]({{< relref "/reference/helm-charts/backup" >}})
* [Restore Helm Chart]({{< relref "/reference/helm-charts/restore" >}})
