---
title: "Troubleshoot K8ssandra"
linkTitle: "Troubleshoot K8ssandra"
description: "Troubleshooting tips for K8ssandra users."
---

## Check quotas in the cloud provider's UI

In some cases, pods can become "unhealthy" and the root cause may be an insufficient quota. For example, in the Google Cloud Platform (GCP) console, check for any unhealthy pods in your GKE project. Then in the IAM &amp; Admin section of the GCP console, navigate to Quotas. Look for any reported issues with backend services:

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

## Amazon S3 bucket's region or name is misconfigured for backups

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

Then for a new or existing K8ssandra installation, relreference the values file. 

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

## Medusa backup/restore connection to Amazon S3 bucket fails - check credentials

If the Medusa log reports an authentication error, check that you provided the correct S3 credentials in the `aws_access_key_id` and `aws_secret_access_key` settings. 

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

## Next

See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Helm charts, a glossary, and cheat sheets.  

