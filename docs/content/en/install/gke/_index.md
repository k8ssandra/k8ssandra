---
title: "Google Kubernetes Engine"
linkTitle: "Google Kubernetes Engine"
weight: 1
description: >
  Complete production ready environment of K8ssandra on Google Kubernetes Engine (GKE).
---

[Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine) or "GKE" is a managed Kubernetes environment on the [Google Cloud Platform](https://cloud.google.com/) (GCP). GKE is a fully managed experience; it handles the management/upgrading of the Kubernetes cluster master as well as autoscaling of "nodes" through "node pool" templates.

Through GKE, your Kubernetes deployments will have first-class support for GCP IAM identities, built-in configuration of high-availability and secured clusters, as well as native access to GCP's networking features such as load balancers.

## Deployment

This guide will cover provisioning and installing the following infrastructure resources.

* 1x _Regional_ GKE cluster with instances spread across multiple Availability Zones.
* 1x Node Pool
  * 6x Kubernetes workers
    * 8 vCPUs
    * 64 GB RAM
* x Load Balancers
  * x Backend services
* x 2TB PD-SSD Volumes (provisioned automatically during installation of K8ssandra)
* 1x Google Cloud Storage bucket for backups

On this infrastructure the K8ssandra installation will consist of the following workloads.

* 3x node Apache Cassandra cluster
* 3x node Stargate deployment
* 1x node Prometheus deployment
* 1x node Grafana deployment
* 1x node Reaper deployment

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

```console
$ git clone git@github.com:k8ssandra/k8ssandra-terraform.git
Cloning into 'k8ssandra-terraform'...
remote: Enumerating objects: 273, done.
remote: Counting objects: 100% (273/273), done.
remote: Compressing objects: 100% (153/153), done.
remote: Total 273 (delta 145), reused 233 (delta 112), pack-reused 0
Receiving objects: 100% (273/273), 71.29 KiB | 1.30 MiB/s, done.
Resolving deltas: 100% (145/145), done.
$ cd k8ssandra-terraform/gcp
```

### Configure `gcloud` CLI

Ensure you have authenticated your `gcloud` client by running the following command:

```console
$ gcloud auth login
Your browser has been opened to visit:

    https://accounts.google.com/.....

You are now logged in as [kate.sandra@k8ssandra.io].
Your current project is [k8ssandra-demo].  You can change this setting by running:
  $ gcloud config set project PROJECT_ID
```

Next configure the `region`, `zone`, and `project name` configuration parameters

```console
$ gcloud config set compute/region us-central1

Updated property [compute/region].

$ gcloud config set compute/zone us-central1-c

Updated property [compute/zone].

$ gcloud config set project "k8ssandra-testing"

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

```console
$ cd env
$ terraform init
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
Acquiring state lock. This may take a few moments...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

# Output reduced for brevity

Plan: 26 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + bucket_name    = "prod-k8ssandra-storage-bucket"
  + endpoint       = (known after apply)
  + master_version = (known after apply)

```

After planning we tell terraform to `apply` the plan. This command kicks off the actual provisioning of resources for this deployment.

```console
$ terraform apply

# Output reduced for brevity

Do you want to perform these actions in workspace "my-workspace"?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

# Output reduced for brevity

Apply complete! Resources: 26 added, 0 changed, 0 destroyed.

Outputs:

bucket_name = "prod-k8ssandra-storage-bucket"
endpoint = "..."
master_version = "1.18.16-gke.502"
```

With the GKE cluster deployed you may now continue with [retrieving the kubeconfig](#retrieve-kubeconfig).

## Retrieve `kubeconfig`

After provisioning the GKE cluster we must request a copy of the `kubeconfig`. This provides the `kubectl` command with all connection information including TLS certificates and IP addresses for Kube API requests.

```console
$ gcloud container clusters get-credentials prod-k8ssandra --region us-central1 --project k8ssandra-testing
Fetching cluster endpoint and auth data.
kubeconfig entry generated for prod-k8ssandra.

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
In order to allow for backup and restore operations, we must create a service account for the Medusa operator which handles coordinating the movement of data to and from Google Cloud Storage buckets. As part of the provisioning sections a service account was generated for this purposes. Here we will retrieve the authentication JSON file for this account and push it into Kubernetes as a secret.

TODO retrieve service account credentials
TODO push service account credentials to k8s secret

### Generate `gke.values.yaml`

Here is a reference Helm `values.yaml` file with configuration options for running K8ssandra in GKE.

{{< readfilerel file="gke.values.yaml"  highlight="yaml" >}}

Note the storage class defined here, `standard-rwo`, is already created by GCP. This storage class has a `volumeBindingMode` set to `WaitForFirstConsumer`. This tells GKE to provision volumes after Kubernetes has determined which workers will be receiving the pods. This allows for the provisioning of persistent storage volumes in the same Availability Zone (AZ) as the worker.

Additionally review the `datacenters[].racks` parameters and ensure the values align with the AZs where your workers are deployed. Cassandra will strive to replicate data across rack boundaries to account for the loss of an entire rack. In our deployment this means we can tolerate the loss of an entire AZ.


### Deploy K8ssandra with Helm

With a `values.yaml` file generated which details out specific configuration overrides we can now deploy K8ssandra via Helm.

```console
$ helm install my-k8ssandra k8ssandra/k8ssandra -f gke.values.yaml
NAME: my-k8ssandra
LAST DEPLOYED: Sat Apr 24 01:15:46 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
```

## Additional Configuration

At this time there are a couple of manual post-installation steps to allow for external access to resources running within the GKE cluster.

TODO create cluster services ingress to target
TODO create ingress targeting services

## Next Steps

With a freshly provisioned cluster on GKE consider visiting the [developer]({{ relref "developer" }}) and [SRE]({{ relref "developer" }}) quickstarts for a guided experience exploring your cluster. Alternatively if you want to tear down your cluster and / or infrastructure check out the next section on cleaning up resources.

## Cleanup Resources

If this cluster is no longer needed you may optionally uninstall K8ssandra or delete all of the infrastructure

### Uninstall K8ssandra

```console
$ helm uninstall my-k8ssandra
release "my-k8ssandra" uninstalled
```

### Destroy GKE Cluster

```console
$ terraform destroy

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

Releasing state lock. This may take a few moments...

Destroy complete! Resources: 26 destroyed.
```
