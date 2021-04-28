---
title: "Google Kubernetes Engine"
linkTitle: "Google GKE"
weight: 2
description: >
  Complete **production** ready environment of K8ssandra on Google Kubernetes Engine (GKE).
---

[Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) or "GKE" is a managed Kubernetes environment on the [Google Cloud Platform](https://cloud.google.com/) (GCP). GKE is a fully managed experience; it handles the management/upgrading of the Kubernetes cluster master as well as autoscaling of "nodes" through "node pool" templates.

Through GKE, your Kubernetes deployments will have first-class support for GCP IAM identities, built-in configuration of high-availability and secured clusters, as well as native access to GCP's networking features such as load balancers.

{{% alert title="Tip" color="success" %}}
Also available in followup topics are post-install steps and role-based considerations for [developers]({{< relref "/quickstarts/developer">}}) or [Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer">}}) (SREs).
{{% /alert %}}

## Deployment

This topic covers provisioning and installing the following infrastructure resources.

* 1x Google Compute Network (Virtual Private Cloud, or VPC)
* TODOx Subnet
* TODOx Router
* TODOx Compute Router NAT
* 1x _Regional_ GKE cluster with instances spread across multiple Availability Zones.
* 1x Node Pool
  * 6x Kubernetes workers
    * 8 vCPUs
    * 64 GB RAM
* 2x Load Balancers
* 3x 2TB PD-SSD Volumes (provisioned automatically during installation of K8ssandra)
* 1x Google Cloud Storage bucket for backups
* 1x Google Storage Bucket IAM member

On this infrastructure the K8ssandra installation will consist of the following workloads.

* 3x instance Apache Cassandra cluster
* 3x instance Stargate deployment
* 1x instance Prometheus deployment
* 1x instance Grafana deployment
* 1x instance Reaper deployment

Feel free to update the parameters used during this guide to match your target deployment. The following should be considered a minimum for production workloads.

{{% alert title="Quotas" color="primary" %}}
This installation slightly exceeds the default quotas provided within a new project. Consider requesting the following quota requests to allow for the provisioning of this installation.

TODO identify quota limits that need to be updated
{{% /alert %}}

## Terraform

As a convenience we provide reference [Terraform](https://www.terraform.io/) modules for orchestrating the provisioning of cloud resources necessary to run K8ssandra.

### Tools

| Tool | Version | 
|------|---------|
| [Terraform](https://www.terraform.io/downloads.html) | 0.14 |
| [Terraform GCP provider](https://registry.terraform.io/providers/hashicorp/google/latest) | ~>3.0 |
| [Helm](https://helm.sh/) | 3 |
| [Google Cloud SDK](https://cloud.google.com/sdk)  | 333.0.0 |
| - bq | 2.0.65 |
| - core | 2021.03.19 |
| - gsutil | 4.60 |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.17.17 |

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
cd k8ssandra-terraform/gcp
```

### Configure `gcloud` CLI

Ensure you have authenticated your `gcloud` client by running the following command:

```console
gcloud auth login
```

**Output**:

```console
Your browser has been opened to visit:

    https://accounts.google.com/.....

You are now logged in as [kate.sandra@k8ssandra.io].
Your current project is [k8ssandra-demo].  You can change this setting by running:
  $ gcloud config set project PROJECT_ID
```

Next configure the `region`, `zone`, and `project name` configuration parameters.

Set the region:

```console
gcloud config set compute/region us-central1
```

**Output**:

```console
Updated property [compute/region].
```

Set the zone:

```console
gcloud config set compute/zone us-central1-c
```

**Output**:

```console
Updated property [compute/zone].
```

Set the project:

```console
gcloud config set project "k8ssandra-testing"
```

**Output**:

```console
Updated property [core/project].
```

### Setup Environment Variables

These values will be used to define where infrastructure is provisioned along with the naming of resources.

```bash
export TF_VAR_environment=prod
export TF_VAR_name=k8ssandra
export TF_VAR_project_id=k8ssandra-testing
export TF_VAR_region=us-central1
```

{{% alert title="Limits" color="primary" %}}
GCP limits the total length of resource names. If your deployment fails to plan try reducing the number of characters in the `environment` and `name` parameters.
{{% /alert %}}

### Provision Infrastructure

We begin this process by initializing our environment and configuring a workspace. To start we run `terraform init` which handles pulling down any plugins required and configures the backend.

```bash
cd env
terraform init
```

**Output**:

```bash
Initializing modules...

Initializing the backend...

Successfully configured the backend "gcs"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/google versions matching "~> 3.0"...
- Finding latest version of hashicorp/google-beta...
- Installing hashicorp/google v3.65.0...
- Installed hashicorp/google v3.65.0 (signed by HashiCorp)
- Installing hashicorp/google-beta v3.65.0...
- Installed hashicorp/google-beta v3.65.0 (signed by HashiCorp)

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

```bash
Acquiring state lock. This may take a few moments...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

# Output reduced for brevity

Plan: 26 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + bucket_name     = "prod-k8ssandra-storage-bucket"
  + endpoint        = (known after apply)
  + master_version  = (known after apply)
  + service_account = (known after apply)
```

After planning we tell terraform to `apply` the plan. This command kicks off the actual provisioning of resources for this deployment.

```console
terraform apply
```

**Output**:

```bash
# Output reduced for brevity

Do you want to perform these actions in workspace "my-workspace"?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

# Output reduced for brevity

Apply complete! Resources: 26 added, 0 changed, 0 destroyed.

Outputs:

bucket_name = "prod-k8ssandra-storage-bucket"
endpoint = "......"
master_version = "1.18.16-gke.502"
service_account = "prod-k8ssandra-sa@k8ssandra-testing.iam.gserviceaccount.com"
```

With the GKE cluster deployed you may now continue with [retrieving the kubeconfig](#retrieve-kubeconfig).

## Retrieve `kubeconfig`

After provisioning the GKE cluster we must request a copy of the `kubeconfig`. This provides the `kubectl` command with all connection information including TLS certificates and IP addresses for Kube API requests.

```console
gcloud container clusters get-credentials prod-k8ssandra --region us-central1 --project k8ssandra-testing
```

**Output**:

```bash
Fetching cluster endpoint and auth data.
kubeconfig entry generated for prod-k8ssandra.
```

```bash
kubectl cluster-info
```

**Output**:

```bash
Kubernetes control plane is running at https://.....
GLBCDefaultBackend is running at https://...../api/v1/namespaces/kube-system/services/default-http-backend:http/proxy
KubeDNS is running at https://...../api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
Metrics-server is running at https://...../api/v1/namespaces/kube-system/services/https:metrics-server:/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

```bash
kubectl version
```

**Output**:

```bash
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.0", GitCommit:"cb303e613a121a29364f75cc67d3d580833a7479", GitTreeState:"clean", BuildDate:"2021-04-08T16:31:21Z", GoVersion:"go1.16.1", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"18+", GitVersion:"v1.18.16-gke.502", GitCommit:"a2a88ab32201dca596d0cdb116bbba3f765ebd36", GitTreeState:"clean", BuildDate:"2021-03-08T22:06:24Z", GoVersion:"go1.13.15b4", Compiler:"gc", Platform:"linux/amd64"}
WARNING: version difference between client (1.21) and server (1.18) exceeds the supported minor version skew of +/-1
```

## Install K8ssandra

With all of the infrastructure provisioned we can now focus on installing K8ssandra. This will require configuring a service account for the backup and restore service, creating a set of Helm variable overrides, and setting up GKE specific ingress configurations.

### Create Backup / Restore Service Account Secrets

In order to allow for backup and restore operations, we must create a service account for the Medusa operator which handles coordinating the movement of data to and from Google Cloud Storage (GCS) buckets. As part of the provisioning sections a service account was generated for this purposes. Here we will retrieve the authentication JSON file for this account and push it into Kubernetes as a secret.

Looking at the output of `terraform plan` and `terraform apply` we can see the name of the service account which has been provisioned. Here we use the `gcloud` command line tools to retrieve keys for use by Medusa. In our reference implementation this value is `prod-k8ssandra-sa@k8ssandra-testing.iam.gserviceaccount.com`.

```console
gcloud iam service-accounts keys create medusa.key.json --iam-account=prod-k8ssandra-sa@k8ssandra-testing.iam.gserviceaccount.com
```

**Output**:

```bash
created key [3e5b6e4a02936b20f6ae39bffe7d28f870c94fe6] of type [json] as [medusa.key.json] for [prod-k8ssandra-sa@k8ssandra-testing.iam.gserviceaccount.com]
```

With the key file on our local machine we can now push this file to Kubernetes as a secret with `kubectl`.

```bash
kubectl create secret generic prod-k8ssandra-medusa-key --from-file=medusa_gcp_key.json=./medusa.key.json 
```

**Output**:

```bash
secret/prod-k8ssandra-medusa-key created
```

{{% alert title="Important" color="primary" %}}
The name of the JSON key file within the secret MUST be `medusa_gcp_key.json`. _Any_ other value will result in Medusa not finding the secret and backups failing.
{{% /alert %}}

This secret, `prod-k8ssandra-medusa-key`, can now be referenced in our K8ssandra configuration to allow for backing up data to GCS with Medusa.

### Generate `gke.values.yaml`

Here is a reference Helm `values.yaml` file with configuration options for running K8ssandra in GKE.

{{< readfilerel file="gke.values.yaml"  highlight="yaml" >}}

{{% alert title="Important" color="primary" %}}
Take note of the comments in this file. If you have changed the name of your secret, are deploying in a different region, or have tweaked any other values it is imperative that you update this file before proceeding.
{{% /alert %}}

### Deploy K8ssandra with Helm

With a `values.yaml` file generated which details out specific configuration overrides we can now deploy K8ssandra via Helm.

```bash
helm install prod-k8ssandra k8ssandra/k8ssandra -f gke.values.yaml
```

**Output**:

```bash
NAME: prod-k8ssandra
LAST DEPLOYED: Sat Apr 24 01:15:46 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
```

### Retrieve K8ssandra superuser credentials {#superuser}

You'll need the K8ssandra superuser name and password in order to access Cassandra utilities and do things like generate a Stargate access token.

To retrieve K8ssandra superuser credentials:

1. Retrieve the K8ssandra superuser name:

    ```bash
    kubectl get secret k8ssandra-superuser -o jsonpath="{.data.username}" | base64 --decode ; echo
    ```

    **Output**:

    ```bash
    k8ssandra-superuser
    ```

1. Retrieve the K8ssandra superuser password:

    ```bash
    kubectl get secret k8ssandra-superuser -o jsonpath="{.data.password}" | base64 --decode ; echo
    ```

    **Output**:

    ```bash
    PGo8kROUgAJOa8vhjQrE49Lgruw7s32HCPyVvcfVmmACW8oUhfoO9A
    ```

{{% alert title="Tip" color="success" %}}
Save the superuser name and the generated password for your environment. You will need the credentials when following the 
[Quickstart for developers]({{< relref "/quickstarts/developer" >}}) or [Quickstart for Site Reliability Engineers]({{< relref "/quickstarts/site-reliability-engineer" >}}) post-install steps.
{{% /alert %}}

## Additional Configuration

At this time there are a couple of manual post-installation steps to allow for external access to resources running within the GKE cluster.

TODO create cluster services ingress to target
TODO create ingress targeting services

## Cleanup Resources

If this cluster is no longer needed you may optionally uninstall K8ssandra or delete all of the infrastructure.

### Uninstall K8ssandra

```console
helm uninstall prod-k8ssandra
```

**Output**:

```bash
release "prod-k8ssandra" uninstalled
```

### Destroy GKE Cluster

```console
terraform destroy
```

**Output**:

```bash
# Output omitted for brevity

Plan: 0 to add, 0 to change, 26 to destroy.

Changes to Outputs:
  - bucket_name    = "prod-k8ssandra-storage-bucket" -> null
  - endpoint       = "....." -> null
  - master_version = "1.18.16-gke.502" -> null

Do you really want to destroy all resources in workspace "my-workspace"?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

# Output omitted for brevity

Destroy complete! Resources: 26 destroyed.
```

## Next steps

With a freshly provisioned cluster on GKE, consider visiting the [developer]({{< relref "/quickstarts/developer" >}}) and [Site Reliability Engineer]({{< relref "/quickstarts/site-reliability-engineer" >}}) quickstarts for a guided experience exploring your cluster. 

Alternatively, if you want to tear down your GKE cluster and / or infrastructure, refer to the section above that covers [cleaning up resources]({{< relref "#cleanup-resources" >}}).
