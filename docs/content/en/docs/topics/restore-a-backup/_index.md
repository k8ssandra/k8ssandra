---
title: "Backup and Restore"
linkTitle: "Backup and Restore"
weight: 3
description: K8ssandra provides backup/restore via Medusa
---

This topic walks you through the steps to backup and restore Cassandra data running in a Kubernetes cluster.

## Tools

* K8ssandra Helm chart, which we'll extend with `backupRestore` Medusa buckets for Amazon S3 integration
* Sample files in GitHub:
  * [medusa-bucket-key.yaml](medusa-bucket-key.yaml) to create a secret with credentials for AWS S3 buckets
  * [backup-restore-values.yaml](backup-restore-values.yaml) to enable Medusa (backup/restore service) and set related minimal values
  * [test_data.cql](test_data.cql) to populate a Cassandra keyspace and table with data

## Prerequisites

* A Kubernetes environment
* Storage for the backups - see below
* [Helm](https://helm.sh/), a packaging manager for Kubernetes
* An edited version of [medusa-bucket-key.yaml](./medusa-bucket-key.yaml), as noted below

All other prerequisites are handled by the installed tools listed above. The sample files are checked into GitHub.

## Steps

### Verify you've met the prereqs

You will need storage for the backups. This topic shows the use of AWS S3 buckets.

* If you'll use AWS S3, before proceeding with the configuration described below, verify that you know the `aws_access_key_id` and `aws_secret_access_key` values. Or  contact your IT team if they manage those assets. You'll provide those details in an edited version of the [medusa-bucket-key.yaml](medusa-bucket-key.yaml) file. For information about the S3 setup steps, see this helpful [readme](https://github.com/thelastpickle/cassandra-medusa/blob/master/docs/aws_s3_setup.md).  

* Add and update the following repo, which has in one chart all the settings for K8ssandra plus the backup/restore settings:

```
helm repo add k8ssandra https://helm.k8ssandra.io/
helm repo update
```

```
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "k8ssandra" chart repository
Update Complete. ⎈Happy Helming!⎈
```

### Create secret for read/write access to an S3 bucket

Before installing the k8ssandra cluster, we need to supply credentials so that Medusa has read/write to an AWS S3 bucket, which is where the backup will be stored.  Currently, Medusa supports local, Amazon S3, Google Cloud Storage, and Azure buckets. Currently, K8ssandra supports S3. 

**Note:** See [AWS S3 setup](https://github.com/thelastpickle/cassandra-medusa/blob/master/docs/aws_s3_setup.md) on the Medusa wiki for more details for configuring S3.

To do this, start by creating a secret with the credentials for the S3 bucket.

The [medusa-bucket-key.yaml](medusa-bucket-key.yaml) sample in GitHub contains:

```
apiVersion: v1
kind: Secret
metadata:
 name: medusa-bucket-key
type: Opaque
stringData:
 # Note that this currently has to be set to medusa_s3_credentials!
 medusa_s3_credentials: |-
   [default]
   aws_access_key_id = my_access_key
   aws_secret_access_key = my_secret_key
```
   
**Make a copy** of [medusa-bucket-key.yaml](medusa-bucket-key.yaml), and then replace `my_access_key` and `my_secret_key` with your S3 values. 

In the YAML, notice the `stringData` property value: `medusa_s3_credentials`. The secret gets mounted to this location; this is where Medusa expects to get the AWS credentials.

Apply the YAML to your Kubernetes environment. In this example, assume that you had copied medusa-bucket-key.yaml to my-medusa-bucket-key.yaml:

```
kubectl apply -f my-medusa-bucket-key.yaml
secret/medusa-bucket-key configured
```

### Create or update the k8ssandra cluster

Install the `k8ssandra` chart with the following properties. You can reference the provided [backup-restore-values.yaml](backup-restore-values.yaml) file. It contains:

```
size: 3
backupRestore: 
  medusa:
    enabled: true
    bucketName: k8ssandra-bucket-dev
    bucketSecret: medusa-bucket-key
    storage: s3
```

The chart's entries relate to a Kubernetes Secret, which contains the object store credentials. Specifically, the `bucketSecret` property specifies the name of a secret that should contain an AWS access key. As described in the [Medusa documentation](https://github.com/thelastpickle/cassandra-medusa/blob/master/docs/aws_s3_setup.md), the AWS account with which the key is associated should have the permissions that are required for Medusa to access the S3 bucket.

Example for a new k8ssandra installation:

`helm install k8ssandra k8ssandra/k8ssandra -f backup-restore-values.yaml`

Example for an existing k8ssandra installation:

`helm upgrade k8ssandra k8ssandra/k8ssandra -f backup-restore-values.yaml`

Allow a few minutes for the pods to start and proceed to a Ready state; check the pod status periodically:

```
kubectl get pods                              
NAME                                                         READY   STATUS    RESTARTS   AGE
cass-operator-86d4dc45cd-8p7cq                               1/1     Running   0          98s
k8ssandra-tools-kube-prome-operator-6bcdf668d4-b2r6v         1/1     Running   0          98s
.
.
.
```

Backup and restore operations are enabled by default. In the example YAML, `bucketName` corresponds to the name of the S3 bucket: `K8ssanda-bucket-dev`.  The `bucketSecret` corresponds to the secret credentials.

The `k8ssandra` Helm chart includes the Grafana Operator. Notice that `k8ssandra` adds a number of properties in the `cassdc` datacenter.  

`kubectl get cassdc dc1 -o yaml`

In the output, see the `podTemplateSpec` property; two containers were added for Medusa.  Here’s the entry for the GRPC backup service:

`    name: medusa`

Here’s the entry for the restore’s init container. K8ssandra looks for an environment variable to be set, which would indicate when to perform a restore operation.

`    name: medusa-restore`

After a few minutes, once the pods have started, check the status:

```
kubectl get cassdc dc1 -o yaml`
.
.
.
status:
  cassandraOperatorProgress: Ready
  conditions:
  ...
  - lastTransitionTime: "2021-02-03T17:04:52Z"
    message: ""
    reason: ""
    status: "True"
    type: Ready
  ...
```

### Add test data

Now let’s create some test data.  The [test_data.cql](test_data.cql) sample file in GitHub contains:

```
CREATE KEYSPACE medusa_test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
USE medusa_test;
CREATE TABLE users (email text primary key, name text, state text);
insert into users (email, name, state) values ('john@gamil.com', 'John Smith', 'NC');
insert into users (email, name, state) values ('joe@gamil.com', 'Joe Jones', 'VA');
insert into users (email, name, state) values ('sue@help.com', 'Sue Sas', 'CA');
insert into users (email, name, state) values ('tom@yes.com', 'Tom and Jerry', 'NV');
```

Copy the cql file to the k8ssandra container (pod) :

`kubectl cp test_data.cql k8ssandra-dc1-default-sts-0:/tmp -c cassandra`

Add the data to the Cassandra database:

`kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- cqlsh -f /tmp/test_data.cql`

Exec open cqlsh:

`kubectl exec -it k8ssandra-dc1-default-sts-0 -c cassandra -- cqlsh`

```
Connected to k8ssandra at 127.0.0.1:9042.
[cqlsh 5.0.1 | Cassandra 3.11.7 | CQL spec 3.4.4 | Native protocol v4]
Use HELP for help.
cqlsh> use medusa_test;
cqlsh:medusa_test> select * from medusa_test.users;

 email          | name          | state
----------------+---------------+-------
 john@gamil.com |    John Smith |    NC
  joe@gamil.com |     Joe Jones |    VA
   sue@help.com |       Sue Sas |    CA
    tom@yes.com | Tom and Jerry |    NV

(4 rows)
```

Exit out of CQLSH:

`cqlsh:medusa_test> exit`

Review the current charts that are in use, so far:

`helm list`

```
NAME               	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART                  	APP VERSION
k8ssandra          	default  	1       	2021-02-03 04:17:23.107265 -0700 MST	deployed	k8ssandra-0.38.0        3.11.7  
```

Also get the deployment status, so far:

`kubectl get deployment`
```
NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
cass-operator                              1/1     1            1           159m
k8ssandra-grafana-operator-k8ssandra       1/1     1            1           159m
k8ssandra-kube-prometheus-stack-operator   1/1     1            1           159m
k8ssandra-reaper-k8ssandra                 1/1     1            1           156m
k8ssandra-reaper-operator-k8ssandra        1/1     1            1           159m
grafana-deployment                         1/1     1            1           158m
```

The output above shows the addition of medusa-test-medusa-operator-k8ssandra pod. 

### Create the backup

Now create a backup using a `test` chart:

`helm install test charts/backup --set name=test,cassandraDatacenter.name=dc1`

```
kubectl get cassandrabackup
NAME       AGE
test       17s
```

Examine the YAML:

`kubectl get cassandrabackup test -o yaml`

The Status section in the YAML shows the backup operation’s start and finish timestamps.

### Amazon S3 buckets

Let's look at the resources in the Amazon S3 dashboard. 

S3 maintains the `backup_index` bucket so it only has to store a single copy of an SSTable across backups.  S3 stores pointers in the index to the SSTables. That implementation avoids a large amount of storage.  For example: 

![Amazon S3 with Medusa buckets](s3K8ssandraMedusaBuckets.png)

### Restore data from the backup

Consider the case where an unexpected event occurred, such as an authorized user accidentally entering cqlsh `TRUNCATE` commands that wiped out data in Cassandra. You can restore data from the backup. For example:

`helm install restore-test ./restore --set name=helm-test,backup.name=test,cassandraDatacenter.name=dc1`

Examine the YAML:

`kubectl get cassandrarestore helm-test -o yaml`

The output shows the restore operation’s start time and that the `cassandraDatacenter` is being recreated.

You can also examine the in-progress logs:

`kubectl logs cassandra-dc1-default-sts-0 -c medusa-restore`

### Launch cqlsh again and verify the restore

Exec into cqlsh and select the data again, to verify the restore operation.

```
kubectl exec -it k8ssandra-dc1-default-stc-0 -c cassandra -cqlsh

Connected to k8ssandra at 127.0.0.1:9042.
[cqlsh 5.0.1 | Cassandra 3.11.7 | CQL spec 3.4.4 | Native protocol v4]
Use HELP for help.
cqlsh> use medusa_test;
cqlsh:medusa_test> select * from medusa_test.users;

 email          | name          | state
----------------+---------------+-------
 john@gamil.com |    John Smith |    NC
  joe@gamil.com |     Joe Jones |    VA
   sue@help.com |       Sue Sas |    CA
    tom@yes.com | Tom and Jerry |    NV

(4 rows)
```

You can look again at the cassandrarestore helm-test YAML for the start and ending timestamps:

`kubectl get cassadrarestore helm-test -o yaml`

![Log output from restore operation](k8ssanda-restore-start-end-timestamps-example.png)

## Next

Learn how to use the Repair Web Interface (Reaper).
