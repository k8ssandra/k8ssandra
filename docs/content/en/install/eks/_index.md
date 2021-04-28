---
title: "Amazon Elastic Kubernetes Service"
linkTitle: "Amazon EKS"
weight: 3
description: >
  Complete **production** ready environment of K8ssandra on Amazon Elastic Kubernetes Service (EKS).
---

Amazon [Elastic Kubernetes Service](https://aws.amazon.com/eks/features/) or "EKS" is a managed Kubernetes service that makes it easy for you to run Kubernetes on AWS and on-premises. EKS is certified Kubernetes conformant, so existing applications that run on upstream Kubernetes are compatible with EKS. This cloud provider automatically manages the availability and scalability of the Kubernetes control plane nodes responsible scheduling containers, managing the availability of applications, storing cluster data, and other key tasks.

{{% alert title="Tip" color="success" %}}
Also available in followup topics are post-install steps and role-based considerations for [developers]({{< relref "/quickstarts/developer">}}) or [Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer">}}) (SREs).
{{% /alert %}}

## Deployment

This guide will cover provisioning and installing the following infrastructure resources.

* 1x _Regional_ EKS cluster with instances spread across multiple Availability Zones.
* 1x Node Pool
  * 6x Kubernetes workers
    * 8 vCPUs
    * 64 GB RAM
* x Load Balancers
  * x Backend services
* x 2TB PD-SSD Volumes (provisioned automatically during installation of K8ssandra)
* 1x Amazon S3 bucket for K8ssandra Medusa backups

On this infrastructure the K8ssandra installation will consist of the following workloads.

* 3x node Apache Cassandra cluster
* 3x node Stargate deployment
* 1x node Prometheus deployment
* 1x node Grafana deployment
* 1x node Reaper deployment

Feel free to update the parameters used during this guide to match your target deployment. The following should be considered a minimum for production workloads.

## Terraform

As a convenience we provide reference [Terraform](https://www.terraform.io/) modules for orchestrating the provisioning of cloud resources necessary to run K8ssandra. If you do not have Terraform available, or prefer to create cloud resources manually, skip over to the [Manual Provisioning](#manual-provisioning) section of this topic.

### Tools

| Tool | Version | 
|------|---------|
| [Terraform](https://www.terraform.io/downloads.html) | 0.14 |
| [Terraform EKS provider](https://learn.hashicorp.com/collections/terraform/aws-get-started) | ~>N.n |
| [Helm](https://helm.sh/) | 3 |
| [Amazon AWS SDK](https://aws.amazon.com/getting-started/tools-sdks/)  | N.n.n |
|   bq | N.n.n |
|   core | yyyy.mm.dd |
|   util | N.nn |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.17.17 |

### Checkout the `k8ssandra-terraform` project

Each of our reference deployment may be found in the GitHub [k8ssandra/k8ssandra-terraform](https://github.com/k8ssandra/k8ssandra-terraform) project. Download the latest release or current `main` branch locally.

```console
$ git clone git@github.com:k8ssandra/k8ssandra-terraform.git
Cloning into 'k8ssandra-terraform'...
remote: Enumerating objects: 273, done.
remote: Counting objects: 100% (273/273), done.
remote: Compressing objects: 100% (153/153), done.
remote: Total 273 (delta 145), reused 233 (delta 112), pack-reused 0
Receiving objects: 100% (273/273), 71.29 KiB | 1.30 MiB/s, done.
Resolving deltas: 100% (145/145), done.
$ cd k8ssandra-terraform/aws
```

### Configure `...` CLI

Ensure you have authenticated your `...` client by running the following command:

```console
$ ... auth login
Your browser has been opened to visit:

https://aws.amazon.com/free/.....

You are now logged in as [kate.sandra@amazon.com].
Your current project is [k8ssandra-demo].  You can change this setting by running:
  $ ... config set project PROJECT_ID
```

Next configure the `region`, `zone`, and `project name` configuration parameters

```console
$ ... config set compute/region us-central1

Updated property [compute/region].

$ ... config set compute/zone us-central1-c

Updated property [compute/zone].

$ ... config set project "k8ssandra-testing"

Updated property [core/project].
```

### Setup Environment Variables

```bash
export TF_VAR_environment=production
export TF_VAR_name=k8ssandra
export TF_VAR_project_id=k8ssandra-testing
export TF_VAR_region=us-central1
```

### Provision Infrastructure

We begin this process by initializing our environment and configuring a workspace. To start we run `terraform init` which handles pulling down any plugins required and configures the backend.

```console
$ cd env
$ terraform init
Initializing modules...

Initializing the backend...

Successfully configured the backend "aws"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/amazon versions matching "~> N.n"...
- Finding latest version of hashicorp/amazon-beta...
- Installing hashicorp/amazon vN.nn.n...
- Installed hashicorp/amazon vN.nn.n (signed by HashiCorp)
- Installing hashicorp/amazon-beta vN.nn.n...
- Installed hashicorp/amazon-beta vN.nn.nn (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

Now we configure a workspace to hold our terraform state information.

```console
$ terraform workspace new my-workspace
Created and switched to workspace "my-workspace"!

You're now on a new, empty workspace. Workspaces isolate their state,
so if you run "terraform plan" Terraform will not see any existing state
for this configuration.

$ terraform workspace select my-workspace
```

With the workspace configured we now instruct terraform to `plan` the required changes to our infrastructure (in this case creation).

```console
$ terraform plan
╷
│ Error: "account_id" ("production-k8ssandra-service-account") doesn't match regexp "^[a-z](?:[-a-z0-9]{4,28}[a-z0-9])$"
│ 
│   on ../modules/iam/main.tf line 17, in resource "amazon_service_account" "service_account":
│   17:   account_id   = format("%s-service-account", var.name)
```

After planning we tell terraform to `apply` the plan. This command kicks off the actual provisioning of resources for this deployment.

```console
$ terraform apply
...
```

With the EKS cluster deployed you may now continue with [installing K8ssandra](#install-k8ssandra). The next section covers the manual provisioning of resources which Terraform has handled for you.

## Manual Provisioning

### Tools

| Tool | Version | 
|------|---------|
| [Helm](https://helm.sh/) | 3 |
| [Amazon AWS SDK](https://aws.amazon.com/getting-started/tools-sdks/)  | Nnn.n.n |
|   bq | N.n.nn |
|   core | yyyy.mm.dd |
|   util | N.nn |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.17.17 |

### Create Service Account

### Create a GKE Cluster

TODO develop steps and provide screenshots

### Create a GCS Bucket

TODO develop steps and provide screenshots
TODO include adding permissions for service account to read / write objects

### Configure `aws` CLI

Ensure you have authenticated your `aws` client by running the following command:

```console
# Authenticate
$ ... auth login
Your browser has been opened to visit:

    https://aws.amazon.com/account/.....

You are now logged in as [kate.sandra@k8ssandra.io].
Your current project is [k8ssandra-demo].  You can change this setting by running:
  $ ... config set project PROJECT_ID

# Set default application login credentials
$ ... auth application-default login
```

Next configure the `region`, `zone`, and `project name` configuration parameters

```console
$ ... config set compute/region us-central1

Updated property [compute/region].

$ ... config set compute/zone us-central1-c

Updated property [compute/zone].

$ ... config set project "k8ssandra-testing"

Updated property [core/project].
```

## Retrieve `kubeconfig`

After provisioning the EKS cluster we must request a copy of the `kubeconfig`. This provides the `kubectl` command with all connection information including TLS certificates and connection information for Kube API requests.

```console
$ ... container clusters get-credentials dev-k8ssandra --region us-central1 --project k8ssandra-testing
Fetching cluster endpoint and auth data.
kubeconfig entry generated for dev-k8ssandra.

$ kubectl cluster-info
Kubernetes control plane is running at https://.....
GLBCDefaultBackend is running at https://...../api/v1/namespaces/kube-system/services/default-http-backend:http/proxy
KubeDNS is running at https://...../api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
Metrics-server is running at https://...../api/v1/namespaces/kube-system/services/https:metrics-server:/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.0", GitCommit:"cb303e613a121a29364f75cc67d3d580833a7479", GitTreeState:"clean", BuildDate:"2021-04-08T16:31:21Z", GoVersion:"go1.16.1", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"18+", GitVersion:"v1.18.16-gke.502", GitCommit:"a2a88ab32201dca596d0cdb116bbba3f765ebd36", GitTreeState:"clean", BuildDate:"2021-03-08T22:06:24Z", GoVersion:"go1.13.15b4", Compiler:"gc", Platform:"linux/amd64"}
WARNING: version difference between client (1.21) and server (1.18) exceeds the supported minor version skew of +/-1
```

## Install K8ssandra

With all of the infrastructure provisioned we can now focus on installing K8ssandra. This will require configuring a service account for the backup and restore service, creating a set of Helm variable overrides, and setting up GKE specific ingress configurations.

### Create Backup / Restore Service Account Secrets
In order to allow for backup and restore operations, we must create a service account for the Medusa operator which handles coordinating the movement of data to and from Amazon S3 buckets. As part of the provisioning sections a service account was generated for this purposes. Here we will retrieve the authentication JSON file for this account and push it into Kubernetes as a secret.

TODO retrieve service account credentials
TODO push service account credentials to k8s secret

### Generate `eks.values.yaml`

```yaml
values.yaml file
```

### Deploy K8ssandra with Helm

With a `values.yaml` file generated which details out specific configuration overrides we can now deploy K8ssandra via Helm.

```console
$ helm install my-k8ssandra k8ssandra/k8ssandra -f eks.values.yaml
```

### Additional Configuration

At this time there are a couple of manual post-installation steps to allow for external access to resources running within the EKS cluster.

TODO create cluster services ingress to target
TODO create ingress targeting services
