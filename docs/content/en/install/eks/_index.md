---
title: "Install K8ssandra on EKS"
linkTitle: "Amazon EKS"
weight: 2
description: >
  Complete **production** ready environment of K8ssandra on Amazon Elastic Kubernetes Service (EKS).
---

Amazon [Elastic Kubernetes Service](https://aws.amazon.com/eks/features/) or "EKS" is a managed Kubernetes service that makes it easy for you to run Kubernetes on AWS and on-premises. EKS is certified Kubernetes conformant, so existing applications that run on upstream Kubernetes are compatible with EKS. AWS automatically manages the availability and scalability of the Kubernetes control plane nodes responsible scheduling containers, managing the availability of applications, storing cluster data, and other key tasks.

{{% alert title="Tip" color="success" %}}
Also available in followup topics are post-install steps and role-based considerations for [developers]({{< relref "/quickstarts/developer">}}) or [site reliability engineers]({{< relref "/quickstarts/site-reliability-engineer">}}) (SREs).
{{% /alert %}}

## Minimum deployment

This topic covers provisioning the following infrastructure resources as a minimum for production. See the [next section]({{< relref "#infrastructure-and-cassandra-recommendations" >}}) for additional considerations discovered during performance benchmarks.

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

Feel free to update the parameters used during this guide to match your target deployment. This should be considered a minimum for production workloads.

## Infrastructure and Cassandra recommendations

While the section above includes infrastructure settings for **minimum** production workloads, performance benchmarks reveal a wider range of recommendations that are important to consider. The performance benchmark report, available in this [detailed blog post](https://k8ssandra.io/blog/articles/k8ssandra-performance-benchmarks-on-cloud-managed-kubernetes/), compared the throughput and latency between:

* The baseline performance of a Cassandra cluster running on AWS EC2 instances -- a common setup for enterprises operating Cassandra clusters
* The performance of K8ssandra running on Amazon EKS, Google GCP GKE, and Microsoft Azure AKS. 

It's important to note the following additional AWS infrastructure settings and observations from the benchmark:

* 8 to 16 vCPUs 
  * r5 instances: Intel Xeon Platinum 8000 series. Cassandra workloads are mostly CPU bound and the core speed made a difference in the throughput benchmarks.
* 32 GB to 128 GB RAM (we used 64 GB RAM during the benchmark)
* 2 to 4 TB of disk space
  * In the benchmark, we used 1x 3.4 TB EBS gp2 volume
* 10k IOPS (observed)

For the disk performance, the benchmark used [Cassandra inspired fio profiles](https://github.com/ibspoof/cassandra-fio) that attempt to emulate Leveled Compaction Strategy and Size Tiered Compaction Strategy behaviors.  

Regarding the Cassandra version and settings:

* The benchmark used Cassandra 4.0-beta4.
* Cassandra default settings were applied with the exception of garbage collection (GC) settings. This used G1GC with 31GB of heap size, along with a few GC related JVM flags:

  ```
  -XX:+UseG1GC
  -XX:G1RSetUpdatingPauseTimePercent=5
  -XX:MaxGCPauseMillis=300
  -XX:InitiatingHeapOccupancyPercent=70 -Xms31G -Xmx31G
  ```

To summarize the findings, running Cassandra in Kubernetes using K8ssandra didn't introduce any notable performance impacts in throughput or latency, all while K8ssandra simplified the deployment steps. See the [blog post](https://k8ssandra.io/blog/articles/k8ssandra-performance-benchmarks-on-cloud-managed-kubernetes/) for more detailed settings, results, and the measures taken to ensure fair production comparisons.

## Terraform

As a convenience we provide reference [Terraform](https://www.terraform.io/) modules for orchestrating the provisioning of cloud resources necessary to run K8ssandra.

### Prerequisite tools

First, these steps assume you already have an AWS account. If not, see [Create and activate a new AWS account](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/) in the AWS documentation. 

Next, ensure you have the prerequisite tools installed (or subsequent versions), including the Terraform binary. Links below go to download resources:


| Tool | Version | 
|------|---------|
| [Terraform](https://www.terraform.io/downloads.html) | 1.0.0 or higher |
| [Terraform EKS provider](https://learn.hashicorp.com/collections/terraform/aws-get-started) | ~>N.n |
| [Helm](https://helm.sh/) | 3 |
| [Amazon AWS SDK](https://aws.amazon.com/cli/)  | 2.2.0 |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.17.17 |
| [Python](https://www.python.org/) | 3 |
| [aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html) | 0.5.0 |

#### Install the Terraform binary

If you haven't already, install Terraform. Refer to the helpful Terraform installation video on this [hashicorp.com page](https://learn.hashicorp.com/tutorials/terraform/infrastructure-as-code?in=terraform/aws-get-started). Follow the instructions for your OS type, then return here.  

Terraform install example for Ubuntu Linux:

```bash
sudo apt-get update && sudo apt-get install -y gnupg software-properties-common curl
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install terraform
```

Verify the installation:

```bash
terraform version
```

**Output**:

```bash
Terraform v1.0.0
on darwin_amd64

Your version of Terraform is out of date! The latest version
is 1.0.1. You can update by downloading from https://www.terraform.io/downloads.html
```

#### Set up the AWS CLI v2

If you haven't already, set up the AWS CLI v2. The steps assume you already have an AWS account. 

Follow the instructions in [Installing, updating, and uninstalling the AWS CLI version 2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html).

Next, FYI, refer to [Installing aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html). As explained on that page, 
if you're running the AWS CLI version 1.16.156 or later, you don't need to install the authenticator.

#### Install kubectl

If you haven't already, install `kubectl`. You'll use `kubectl` commands to interact with your K8ssandra resources. 

One option to get `kubectl` is described in this AWS topic, [Installing kubectl](https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html). See the OS-specific examples. Here's an example on Linux and the 1.17 Kubernetes version:

```bash
curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.17.12/2020-11-02/bin/linux/amd64/kubectl
```

Verify the `kubectl` install:

```bash
kubectl version --short --client
```

**Output**:

```bash
Client Version: v1.17.12
```

#### Install helm v3

If you haven't already, install Helm v3. See this EKS topic, [Using Helm with Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/helm.html). Note the  prerequisite: before you can install Helm charts on your Amazon EKS cluster, you must configure `kubectl` to work for Amazon EKS. If you have not already done this, see [Create a kubeconfig for Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html).

Once you've completed the prerequisites, see the [Helm install] steps for your OS. Here's a Linux example:

```bash
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh
```

#### Install Python v3

If you haven't already, install Python version 3.x for your OS. See the [python downloads](https://www.python.org/downloads/) page.


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

Set up the following environment variables for Terraform's use. Be sure to specify the region you're using in AWS.

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

If you haven't already, add the latest K8ssandra repo:

```bash
helm repo add k8ssandra https://helm.k8ssandra.io
```

**Output**:

```bash
"k8ssandra" has been added to your repositories
```

To ensure you have the latest from all your repos:

```bash
helm repo update
```

**Output**:

```bash
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "k8ssandra" chart repository
Update Complete. ⎈Happy Helming!⎈
```

Now install K8ssandra and specify the `eks.values.yaml` file that you customized in a prior step:

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
