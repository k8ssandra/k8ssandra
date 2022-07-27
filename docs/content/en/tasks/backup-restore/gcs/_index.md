---
title: "Backup and restore with Google Cloud Storage"
linkTitle: "Google Cloud Storage"
toc_hide: true
no_list: true
weight: 3
description: Use Medusa to backup and restore Apache Cassandra® data in Kubernetes to Google Cloud Storage (GCS).
---

**Note:** The information in this topic has not been verified yet for use with K8ssandra Operator.  

Medusa is a Cassandra backup and restore tool. It's packaged with K8ssandra and supports a variety of backends, including GCS for storage in Google Kubernetes Engine (GKE) environments.

## Introduction

Google Cloud Storage (GCS) is a RESTful online file storage web service for storing and accessing data on Google Cloud Platform (GCP) / GKE infrastructure. The service combines the performance and scalability of Google's cloud with advanced security and sharing capabilities.

For details about GCS, see the [Google Cloud Storage documentation](https://cloud.google.com/storage).

## Create a role for backups

In order to perform backups in GCS, Medusa needs to use a service account with appropriate permissions. See this [permissions setup](https://github.com/spotify/cassandra-medusa/blob/master/docs/permissions-setup.md) article.

Using the [Google Cloud SDK](https://cloud.google.com/sdk/install), run the following command to create the `MedusaStorageRole` (set the `$GCP_PROJECT` env variable appropriately):  

```bash
gcloud iam roles create MedusaStorageRole \
        --project ${GCP_PROJECT} \
        --stage GA \
        --title MedusaStorageRole \
        --description "Custom role for Medusa for accessing GCS safely" \
        --permissions storage.buckets.get,storage.buckets.getIamPolicy,storage.objects.create,storage.objects.delete,storage.objects.get,storage.objects.getIamPolicy,storage.objects.list
```

## Create a GCS bucket

Create a bucket for each Cassandra cluster, using the following command line (set the env variables appropriately):

```bash
gsutil mb -p ${GCP_PROJECT} -c regional -l ${LOCATION} ${BUCKET_URL}
```

## Create a service account and download its keys

Medusa will require a `credentials.json` file with the informations and keys for a service account with the appropriate role in order to interact with the bucket.

Create the service account (if it doesn't exist yet):

```bash
gcloud --project ${GCP_PROJECT} iam service-accounts create ${SERVICE_ACCOUNT_NAME} --display-name ${SERVICE_ACCOUNT_NAME}
``` 

## Configure the service account with the role

Once the service account has been created, and considering [jq](https://stedolan.github.io/jq/) is installed, run the following command to add the `MedusaStorageRole` to it, for our backup bucket:

```bash
gsutil iam set <(gsutil iam get ${BUCKET_URL} | jq ".bindings += [{\"members\":[\"serviceAccount:${SERVICE_ACCOUNT_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com\"],\"role\":\"projects/${GCP_PROJECT}/roles/MedusaStorageRole\"}]") ${BUCKET_URL}
```

## Configure Medusa

Generate a json key file called `credentials.json`, for the service account:

```bash
gcloud --project ${GCP_PROJECT} iam service-accounts keys create credentials.json --iam-account=${SERVICE_ACCOUNT_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com
```

Place this file on all Cassandra nodes running medusa under `/etc/medusa` and set the rights appropriately so that only users running Medusa can read/modify it.
Set the `key_file` value in the `[storage]` section of `/etc/medusa/medusa.ini` to the credentials file:  

```ini
bucket_name = my_gcs_bucket
key_file = /etc/medusa/credentials.json
```

Medusa should now be able to access the bucket and perform all required operations, as explained below.

## Deploy K8ssandra

Now that you have GCS set up, install K8ssandra, which includes Medusa:

```bash
helm install k8ssandra k8ssandra/k8ssandra -n k8ssandra
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
CREATE KEYSPACE medusa_test  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
USE medusa_test;
CREATE TABLE users (email text primary key, name text, state text);
insert into users (email, name, state) values ('alice@example.com', 'Alice Smith', 'TX');
insert into users (email, name, state) values ('bob@example.com', 'Bob Jones', 'VA');
insert into users (email, name, state) values ('carol@example.com', 'Carol Jackson', 'CA');
insert into users (email, name, state) values ('david@example.com', 'David Yang', 'NV');
```

Copy the CQL file into the Cassandra pod; that is, the StatefulSet one, which contains `-sts-` in its name:

```bash
kubectl cp test_data.cql k8ssandra-dc1-default-sts-0:/tmp -n k8ssandra -c cassandra
```

Now extract the password to access Cassandra with the `k8ssandra-superuser`. (The password is different for each installation unless it is explicitly set at install time.)

```bash
kubectl get secret k8ssandra-superuser -n k8ssandra -o jsonpath="{.data.password}" | base64 --decode ; echo
```

**Output:**

```bash
XHsZ943WBg5RPNhVAT8x
```

{{% alert title="Tip" color="success" %}}
The password above is an example. The value will be different for your environment. In the subsequent examples, don’t forget to replace the sample password with the one extracted in your environment.
{{% /alert %}}

Let’s now run the uploaded cql script and check that you can read the data. 

```bash
kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- cqlsh -u k8ssandra-superuser -p XHsZ943WBg5RPNhVAT8x -f /tmp/test_data.cql
```

```bash
kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- cqlsh -u k8ssandra-superuser -p XHsZ943WBg5RPNhVAT8x -e "SELECT * FROM medusa_test.users"
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

Now backup this data, and check that files get created in your GCS bucket. 

```bash
helm install my-backup k8ssandra/backup -n k8ssandra --set name=backup1,cassandraDatacenter.name=dc1
```

Because the backup operation is asynchronous, you can monitor its completion by running the following command:

```bash
kubectl get cassandrabackup backup1 -n k8ssandra -o jsonpath={.status.finishTime}
```

As long as this command doesn’t output a date and time, you know that the backup is still running. With the amount of data present and the fact that you’re using a locally accessible backend, this should complete quickly.

Now refresh the GCS UI and you should see some files in the `k8ssandra-medusa` bucket.

In the GCS UI, you should see an index folder, which is the Medusa backup index, and another folder that is specific to each Cassandra node in the cluster. 

## Deleting the data and restoring the backup

Now delete the data by truncating the table, and check that the table is empty.

```bash
kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- cqlsh -u k8ssandra-superuser -p XHsZ943WBg5RPNhVAT8x -e "TRUNCATE medusa_test.users"
```

```bash
 kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- cqlsh -u k8ssandra-superuser -p XHsZ943WBg5RPNhVAT8x -e "SELECT * FROM medusa_test.users"
```

**Output:**

```cql
 email | name | state
-------+------+-------

(0 rows)
```

Now restore the backup taken previously:

```bash
helm install restore-test k8ssandra/restore --set name=restore-backup1,backup.name=backup1,cassandraDatacenter.name=dc1 -n k8ssandra
```

The restore operation will take a little longer because it requires K8ssandra to stop the StatefulSet pod and perform the restore as part of the init containers, before the Cassandra container can start. You can monitor progress using this command:

```bash
watch -d kubectl get cassandrarestore restore-backup1 -o jsonpath={.status} -n k8ssandra
```

The restore operation is fully completed once the finishTime value appears in the output. Example:

```bash
{"finishTime":"2021-03-30T13:58:36Z","restoreKey":"83977399-44dd-4752-b4c4-407273f0339e","startTime":"2021-03-30T13:55:35Z"}
```

Verify that you can now read the restored data from the previously truncated table:

```bash
kubectl exec -it k8ssandra-dc1-default-sts-0 -n k8ssandra -c cassandra -- cqlsh -u k8ssandra-superuser -p XHsZ943WBg5RPNhVAT8x -e "SELECT * FROM medusa_test.users"
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

## Next steps

See the [Custom Resource Definition (CRD) reference]({{< relref "/reference/crd" >}}) topics.

