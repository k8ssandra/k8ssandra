---
title: "Amazon Elastic Kubernetes Service"
linkTitle: "Amazon EKS"
weight: 2
description: >
  Complete **production** ready environment of K8ssandra on Amazon Elastic Kubernetes Service (EKS).
---

Amazon [Elastic Kubernetes Service](https://aws.amazon.com/eks/features/) or "EKS" is a managed Kubernetes service that makes it easy for you to run Kubernetes on AWS and on-premises. EKS is certified Kubernetes conformant, so existing applications that run on upstream Kubernetes are compatible with EKS. AWS automatically manages the availability and scalability of the Kubernetes control plane nodes responsible scheduling containers, managing the availability of applications, storing cluster data, and other key tasks.

{{% alert title="Tip" color="success" %}}
Also available in followup topics are post-install steps and role-based considerations for [developers]({{< relref "/quickstarts/developer">}}) or [site reliability engineers]({{< relref "/quickstarts/site-reliability-engineer">}}) (SREs).
{{% /alert %}}

## Deployment

This guide will cover provisioning and installing the following infrastructure resources.

* 1x Virtual Private Cloud
* 10x Subnets
* 3x Security Groups (& Rules)
* 1x NAT Gateway
* 1x Internet Gateway
* 3x Elastic IP
* 6x Route Table
* 4x Route Table Association
* 1x EKS cluster with instances spread across multiple Availability Zones.
* 1x EKS Node Group
  * 6x Kubernetes workers
    * 8 vCPUs
    * 64 GB RAM
* 3x 2TB EBS Volumes (provisioned automatically during installation of K8ssandra)
* 1x Amazon S3 bucket for K8ssandra Medusa backups

On this infrastructure the K8ssandra installation will consist of the following workloads.

* 3x node Apache Cassandra cluster
* 3x node Stargate deployment
* 1x node Prometheus deployment
* 1x node Grafana deployment
* 1x node Reaper deployment

Feel free to update the parameters used during this guide to match your target deployment. The following should be considered a minimum for production workloads.

## Terraform

As a convenience we provide reference [Terraform](https://www.terraform.io/) modules for orchestrating the provisioning of cloud resources necessary to run K8ssandra.

### Tools

| Tool | Version | 
|------|---------|
| [Terraform](https://www.terraform.io/downloads.html) | 0.14 |
| [Terraform EKS provider](https://learn.hashicorp.com/collections/terraform/aws-get-started) | ~>N.n |
| [Helm](https://helm.sh/) | 3 |
| [Amazon AWS SDK](https://aws.amazon.com/cli/)  | 2.2.0 |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.17.17 |
| [Python](https://www.python.org/) | 3 |
| [aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html) | 0.5.0 |

### Checkout the `k8ssandra-terraform` project

Each of our reference deployment may be found in the GitHub [k8ssandra/k8ssandra-terraform](https://github.com/k8ssandra/k8ssandra-terraform) project. Download the latest release or current `main` branch locally.

```bash
git clone git@github.com:k8ssandra/k8ssandra-terraform.git
```

**Output**:

```bash
Cloning into 'k8ssandra-terraform'...
remote: Enumerating objects: 273, done.
remote: Counting objects: 100% (273/273), done.
remote: Compressing objects: 100% (153/153), done.
remote: Total 273 (delta 145), reused 233 (delta 112), pack-reused 0
Receiving objects: 100% (273/273), 71.29 KiB | 1.30 MiB/s, done.
Resolving deltas: 100% (145/145), done.
```

```bash
cd k8ssandra-terraform/aws
```

### Configure `aws` CLI

Ensure you have authenticated your `aws` client by running the following command:

```console
$ aws configure
AWS Access Key ID [None]: ....
AWS Secret Access Key [None]: ....
Default region name [None]: us-east-1
Default output format [None]: 
```

### Setup Environment Variables

```bash
export TF_VAR_environment=prod
export TF_VAR_name=k8ssandra
export TF_VAR_region=us-east-1
```

### Provision Infrastructure

We begin this process by initializing our environment and configuring a workspace. To start we run `terraform init` which handles pulling down any plugins required and configures the backend.

```console
cd env
terraform init
```

**Output**:

```console
Initializing modules...
- eks in ../modules/eks
- iam in ../modules/iam
- s3 in ../modules/s3
- vpc in ../modules/vpc

Initializing the backend...

Successfully configured the backend "s3"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/aws versions matching "~> 3.0"...
- Installing hashicorp/aws v3.37.0...
- Installed hashicorp/aws v3.37.0 (self-signed, key ID 34365D9472D7468F)

# Output reduced for brevity

Terraform has been successfully initialized!
```

Now we configure a workspace to hold our terraform state information.

```console
terraform workspace new my-workspace
```

**Output**:

```bash
Created and switched to workspace "my-workspace"!

You're now on a new, empty workspace. Workspaces isolate their state,
so if you run "terraform plan" Terraform will not see any existing state
for this configuration.
```

With the workspace configured we now instruct terraform to `plan` the required changes to our infrastructure (in this case creation).

```console
terraform plan
```

**Output**:

```console
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

# Output reduced for brevity

Plan: 50 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + bucket_id        = (known after apply)
  + cluster_Endpoint = (known after apply)
  + cluster_name     = (known after apply)
  + cluster_version  = "1.19"
```

After planning we tell terraform to `apply` the plan. This command kicks off the actual provisioning of resources for this deployment.

```console
terraform apply
```

**Output**:

```console
# Output reduced for brevity

Plan: 50 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + bucket_id        = (known after apply)
  + cluster_Endpoint = (known after apply)
  + cluster_name     = (known after apply)
  + cluster_version  = "1.19"

Do you want to perform these actions in workspace "my-workspace"?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

# Output reduced for brevity

Apply complete! Resources: 50 added, 0 changed, 0 destroyed.

Outputs:

bucket_id = "prod-k8ssandra-s3-bucket"
cluster_Endpoint = "https://....us-east-1.eks.amazonaws.com"
cluster_name = "prod-k8ssandra-eks-cluster"
cluster_version = "1.19"
```

With the EKS cluster deployed you may now continue with [installing K8ssandra](#install-k8ssandra).

## Validate Kubernetes Cluster Connectivity

After provisioning the EKS cluster with `terraform apply` the local Kubeconfig will automatically be updated with the appropriate entries. Let's test this connectivity with `kubectl`.

```console
kubectl cluster-info
```

**Output**:

```console
Kubernetes control plane is running at https://.....
CoreDNS is running at https://..../api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

```console
kubectl version
```

**Output**:

```console
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.0", GitCommit:"cb303e613a121a29364f75cc67d3d580833a7479", GitTreeState:"clean", BuildDate:"2021-04-08T16:31:21Z", GoVersion:"go1.16.1", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"19+", GitVersion:"v1.19.8-eks-96780e", GitCommit:"96780e1b30acbf0a52c38b6030d7853e575bcdf3", GitTreeState:"clean", BuildDate:"2021-03-10T21:32:29Z", GoVersion:"go1.15.8", Compiler:"gc", Platform:"linux/amd64"}
WARNING: version difference between client (1.21) and server (1.19) exceeds the supported minor version skew of +/-1
```

## Install K8ssandra

With all of the infrastructure provisioned we can now focus on installing K8ssandra. This will require configuring a service account for the backup and restore service, creating a set of Helm variable overrides, and setting up EKS specific ingress configurations.

### Create Backup / Restore Service Account Secrets
As part of deploying infrastructure with Terraform an IAM policy is created allowing the EKS cluster workers to access S3 for backup and restore operations. At this time as part of deploying Medusa we _must_ provide a secret for the pods to successfully get scheduled. In this case we will create an empty secret to bypass this limitation until [k8ssandra/k8ssandra#712](https://github.com/k8ssandra/k8ssandra/issues/712) is resolved.

```console
kubectl create secret generic prod-k8ssandra-medusa-key
```

**Output**:

```console
secret/prod-k8ssandra-medusa-key created
```

### Generate `eks.values.yaml`

Here is a reference Helm `values.yaml` file with configuration options for running K8ssandra in EKS.

{{< readfilerel file="eks.values.yaml"  highlight="yaml" >}}

{{% alert title="Important" color="primary" %}}
Take note of the comments in this file. If you have changed the name of your secret, are deploying in a different region, or have tweaked any other values it is imperative that you update this file before proceeding.
{{% /alert %}}

### Deploy K8ssandra with Helm

With a `values.yaml` file generated which details out specific configuration overrides we can now deploy K8ssandra via Helm.

```console
helm install prod-k8ssandra k8ssandra/k8ssandra -f eks.values.yaml
```

### Retrieve K8ssandra superuser credentials {#superuser}

You'll need the K8ssandra superuser name and password in order to access Cassandra utilities and do things like generate a Stargate access token.

{{% alert title="Tip" color="success" %}}
In `kubectl get secret` commands, be sure to prepend the environment name. In this topic's examples, we have used `prod-k8ssandra`. Notice how it's prepended in the examples below. Also, save the displayed superuser name and the generated password for your environment. You will need the credentials when following the 
[Quickstart for developers]({{< relref "/quickstarts/developer" >}}) or [Quickstart for Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer" >}}) post-install steps.
{{% /alert %}}

To retrieve K8ssandra superuser credentials:

1. Retrieve the K8ssandra superuser name:

    ```bash
    kubectl get secret prod-k8ssandra-superuser -o jsonpath="{.data.username}" | base64 --decode ; echo
    ```

    **Output**:

    ```bash
    prod-k8ssandra-superuser
    ```

1. Retrieve the K8ssandra superuser password:

    ```bash
    kubectl get secret prod-k8ssandra-superuser -o jsonpath="{.data.password}" | base64 --decode ; echo
    ```

    **Output**:

    ```bash
    PGo8kROUgAJOa8vhjQrE49Lgruw7s32HCPyVvcfVmmACW8oUhfoO9A
    ```

## Cleanup Resources

If this cluster is no longer needed you may optionally uninstall K8ssandra or delete all of the infrastructure.

### Uninstall K8ssandra

```console
$ helm uninstall prod-k8ssandra
release "prod-k8ssandra" uninstalled
```

### Destroy EKS Cluster

```console
terraform destroy
```

**Output**:

```console
# Output omitted for brevity

Plan: 0 to add, 0 to change, 50 to destroy.

Do you really want to destroy all resources in workspace "my-workspace"?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

# Output omitted for brevity

Destroy complete! Resources: 50 destroyed.
```

## Next steps

With a freshly provisioned cluster on Amazon EKS, consider visiting the [developer]({{< relref "/quickstarts/developer" >}}) and [site reliability engineer]({{< relref "/quickstarts/site-reliability-engineer" >}}) quickstarts for a guided experience exploring your cluster. 

Alternatively, if you want to tear down your Amazon EKS cluster and / or infrastructure, refer to the sections above that cover [cleaning up resources]({{< relref "#cleanup-resources" >}}).
