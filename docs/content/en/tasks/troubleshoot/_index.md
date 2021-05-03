---
title: "Troubleshoot K8ssandra"
linkTitle: "Troubleshoot"
description: "Troubleshooting tips for K8ssandra users."
---

Find troubleshooting tips in these sections:

* [Check quotas]({{< relref "#check-quotas" >}})
* [Bucket region or name for backups is misconfigured]({{< relref "#bucket-region-or-name-for-backups-is-misconfigured" >}}) 
* [Incorrect credentials are configured for backups]({{< relref "#incorrect-credentials-are-configured-for-backups" >}}) 
* [Check dependencies in Helm charts]({{< relref "#check-dependencies-in-helm-charts" >}}) 
* [Collect useful information]({{< relref "#collect-useful-information" >}})

## Check quotas

In some cases, pods can become "unhealthy" and the root cause may be an insufficient quota. You can check quotas in the cloud provider's UI. For example, in the Google Cloud Platform (GCP) console, check for any unhealthy pods in your GKE project. Then in the IAM &amp; Admin section of the GCP console, navigate to Quotas. Look for any reported issues with backend services:

![Backend service quota error](gcp-quota-example1.png)

From the GCP &gt; IAM &amp; Admin &gt; Quotas display:

1. Select the row for the service name that is reporting a quota issue
2. Click **All Quotas** from the Details column
3. Check the box for the affected quota, and click **Edit Quota**. 
4. The dialog indicates: "Enter a new quota limit. Your request will be sent to your service provider for approval." Examine the displayed current value and set a new value. 
5. Enter a brief request description and click **Next**.
6. Verify your contact information, and click **Submit Request**.

Notice how in the following example the Backend services quota is set to '5', and we're changing it to '50'. For the K8ssandra deployments (Stargate, cass-operator, Reaper, Medusa, and so on), actually `10` might be a sufficient quota.

![Quota UI showing change in Backend service quota from 5 to 50](gcp-quota-example2.png)


## Bucket region or name for backups is misconfigured

Among the operators installed by K8ssandra is Medusa, which provides backup and restore for Cassandra data.

If the storage object's name or region used by an Amazon S3 bucket does not match the values expected by Medusa, an error is written to the Medusa section of the logs.  Example:

```bash
kubectl logs demo-dc1-default-sts-0 -c medusa
.
.
.
File "/usr/local/lib/python3.6/dist-packages/libcloud/storage/drivers/s3.py", line 143, in parse_error driver=S3StorageDriver)
libcloud.common.types.LibcloudError: <LibcloudError in <class 'libcloud.storage.drivers.s3.S3StorageDriver'> 
'This bucket is located in a different region. Please use the correct driver. Bucket region "us-east-2", used region "us-east-1".'>
As a result of the region mismatch, the Medusa container within the <cluster-name>-dc1-default-sts-0 pod fails to start. While other pods launched by the K8ssandra install may start successfully, the <cluster-name>-dc1-default-sts-0 pod will not due to the Medusa error.
```

Separately in Amazon AWS, confirm that you know the correct region and name to use for your bucket. Example:

![Amazon S3 Bucket Overview shows the name of the region](amazon-s3-bucket-overview.png)

Declare the appropriate name and region in a values YAML. For example, create a file named `my-backup-restore-values.yaml`. Notice below the `storage_properties` setting for the region `us-east-1`, which matches the region configured and shown in the Amazon S3 user interface:

```yaml
size: 3
backupRestore: 
  medusa:
    enabled: true
    bucketName: jsmart-k8ssandra-bucket2
    bucketSecret: medusa-bucket-key
    storage: s3
    storage_properties:
      region: us-east-1
```

Also make sure the bucketName matches: `jsmart-k8ssandra-bucket2`, in this example.

For example, relreferring again to the S3 UI, confirm the bucket name:

![Confirm the bucket name as shown in the Amazon S3 UI](amazon-s3-confirm-bucket-name.png)

Then for a new or existing K8ssandra installation, reference the values file. 

New install:

```bash
helm install demo k8ssandra/k8ssandra -f my-backup-restore-values.yaml
```

Upgrade:

```bash
helm upgrade demo k8ssandra/k8ssandra -f my-backup-restore-values.yaml
```

{{% alert title="Tip" color="success" %}}
If you're using Google Cloud Storage for your backups, you do not need to include the region setting in a values YAML. 
{{% /alert %}}


## Incorrect credentials are configured for backups

If the Medusa log reports an authentication error, check that you provided the correct credentials. For example, with Amazon S3 buckets, check the credentials in the configured `aws_access_key_id` and `aws_secret_access_key` settings. 

For example, `my-medusa-bucket-key.yaml` contains:

```yaml
apiVersion: v1
kind: Secret
metadata:
 name: medusa-bucket-key
type: Opaque
stringData:
# Note that this currently has to be set to medusa_s3_credentials!
medusa_s3_credentials: |-
  [default]
  aws_access_key_id = FakeValues99ESPW3ALMEZ6U
  aws_secret_access_key = FakeValues99cl9bqJFVA3iFUm+yqVe08HxhXFE/ilK
``` 

If your IT group manages S3 credentials, contact IT to get the correct values.

Before installing or upgrading K8ssandra, and before starting a backup, apply the Medusa bucket values to your Kubernetes environment. Example:

```bash
kubectl apply -f my-medusa-bucket-key.yaml
```

**Output:**

```bash
 secret/medusa-bucket-key configured
```

## Check dependencies in Helm charts

You may experience a `missing in charts/ directory` error message.  

If so, you can utilize a K8ssandra script: 

[update-helm-deps.sh](https://github.com/k8ssandra/k8ssandra/blob/main/scripts/update-helm-deps.sh)

This script assists with updating dependencies for each chart in an appropriate order.  

Be sure to run this script so the `./charts` folder is properly located.

## Collect useful information

Suppose you have an error after editing a K8ssandra configuration, or you want to inspect some things as you learn.  There are some useful commands that come in handy when needing to dig a bit deeper. The following examples assume you are using a `k8ssandra` namespace, but this can be adjusted as needed.

Issue the following `kubectl` command to view the `Management-api` logs.  Replace *cassandra-pod* with an actual pod instance name:

```bash
kubectl logs *cassandra-pod* -c cassandra -n k8ssandra
```

Issue the following `kubectl` command to view the `Cassandra` logs.  Replace *cassandra-pod* with an actual pod instance name:

```bash
kubectl logs *cassandra-pod* -c server-system-logger -n k8ssandra
```

Issue the following `kubectl` command to view `Medusa` logs.  Replace *cassandra-pod* with an actual pod instance name:

```bash
kubectl logs *cassandra-pod* -c medusa -n k8ssandra
```

Issue the following `kubectl` command to describe the `CassandraDatacenter` resource.  This provides a wealth of information about the resource, which includes `aged events` that assist when trying to troubleshoot an issue:

```bash
kubectl describe cassandradatacenter/dc1 -n k8ssandra
```

Gather container specific information for a pod.

 First, list out the pods scoped to the K8ssandra namespace and instance with a target release:

```bash
kubectl get pods -l app.kubernetes.io/instance=*release-name* -n k8ssandra
```

{{% alert title="Note" color="success" %}}
If you don't know the release name, look it up with:
```bash
helm list -n k8ssandra
```
{{% /alert %}}

Next, targeting a specific pod, filter out `container` specific information. Replace the name of the pod with the pod of interest:

```bash
kubectl describe pod/*pod-name* -n k8ssandra | grep container -C 3
```

A slight variation: list out pods having the label for a `cassandra` cluster:

```bash
kubectl get pods -l cassandra.datastax.com/cluster=*release-name* -n k8ssandra
```

Now, using a pod-name returned, describe all the details:

```bash
kubectl describe pod/*pod-name* -n k8ssandra
```

## Next steps

Explore other K8ssandra [tasks]({{< relref "/tasks" >}}).

See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Helm charts, and a glossary. 
