---
title: "Install K8ssandra on GKE"
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

## Minimum deployment

This topic covers provisioning the following infrastructure resources as a minimum for production. See the [next section]({{< relref "#infrastructure-and-cassandra-recommendations" >}}) for additional considerations discovered during performance benchmarks.

* 1x Google Compute Network (Virtual Private Cloud, or VPC)
* 1x Subnet
* 1x Router
* 1x Compute Router NAT
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

Feel free to update the parameters used during this guide to match your target deployment. This should be considered a minimum for production workloads.

{{% alert title="Quotas" color="primary" %}}
This installation slightly exceeds the default quotas provided within a new project. Consider requesting the following quota requests to allow for the provisioning of this installation:

* 6x Backend services

For more, see this [troubleshooting tip]({{< relref "tasks/troubleshoot/#check-quotas" >}}).
{{% /alert %}}

## Infrastructure and Cassandra recommendations

While the section above includes infrastructure settings for **minimum** production workloads, performance benchmarks reveal a wider range of recommendations that are important to consider. The performance benchmark report, available in this [detailed blog post](https://k8ssandra.io/blog/articles/k8ssandra-performance-benchmarks-on-cloud-managed-kubernetes/), compared the throughput and latency between:

* The baseline performance of a Cassandra cluster running on AWS EC2 instances -- a common setup for enterprises operating Cassandra clusters
* The performance of K8ssandra running on Amazon EKS, Google GCP GKE, and Microsoft Azure AKS. 

It's important to note the following additional GCP infrastructure recommendations and observations from the benchmark:

* 8 to 16 vCPUs 
  * n2 instances: Intel Cascade Lake series
* 32 GB to 128 GB RAM
* 2 to 4 TB of disk space
  * In the benchmark, we used 3.4TB volumes to achieve enough power to match high performance production requirements.
* 5k to 10k IOPS (observed)

For the disk performance, the benchmark used [Cassandra inspired fio profiles](https://github.com/ibspoof/cassandra-fio) that attempt to emulate Leveled Compaction Strategy and Size Tiered Compaction Strategy behaviors. In the report, throughput and latency results were measured for three types of GCP disks:

* GCP GKE n2-highmem-8, pd-ssd disk (`premium-rwo`)
* GCP GKE n2-highmem-8, pd-balanced disk (`standard-rwo`)
* GCP GKE e2-highmem-8, pd-ssd disk (`premium-rwo`)

{{% alert title="Disk recommendations" color="success" %}}
The e2 instance class in GCP is the default in GKE, but proved to be under-powered, compared to other cloud vendors' latest generation instances, with approximately 20% lower [sysbench](https://github.com/akopytov/sysbench) score. The n2 instance class offered much better performance with a score that was approximately 10% higher than the benchmark's baseline AWS r5 instances. The pd-balanced (standard-rwo) disks provided enough performance for this specific benchmarks; thus we didn’t need to upgrade to the slightly more expensive pd-ssd (premium-rwo) disks.
{{% /alert %}}

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

First, these steps assume you already have a GCP account and project. If not, see [Creating and managing projects](https://cloud.google.com/resource-manager/docs/creating-managing-projects) in the Google documentation.

Next, ensure you have the prerequisite tools installed (or subsequent versions), including the Terraform binary. Links below go to download resources:

| Tool | Version | 
|------|---------|
| [Terraform](https://www.terraform.io/downloads.html) | 1.0.0 |
| [Terraform GCP provider](https://registry.terraform.io/providers/hashicorp/google/latest) (installed by `terraform init`, as explained below) | ~>3.0 |
| [Helm](https://helm.sh/) | 3 |
| [Google Cloud SDK](https://cloud.google.com/sdk)  | 347.0.0 |
| - bq | 2.0.65 |
| - core | 2021.03.19 |
| - gsutil | 4.60 |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.17.17 |

#### Install the Terraform binary

If you haven't already, install Terraform. Refer to the helpful Terraform installation video on this [hashicorp.com page](https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/gcp-get-started). Follow the instructions for your OS type, then return here.  

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


#### Install the Google Cloud SDK

If you haven't already, [install the Google Cloud SDK](https://cloud.google.com/sdk/docs/install) on your client where Terraform was installed. 

Example:

```bash
cd
curl -LO https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-347.0.0-linux-x86_64.tar.gz
tar -xzvf google-cloud-sdk-347.0.0-linux-x86_64.tar.gz
./google-cloud-sdk/install.sh
```

If you agree, take the defaults when prompted.

Initialize Google Cloud:

```bash
./google-cloud-sdk/bin/gcloud init
```

To associate your account with Google Cloud and authenticate, follow the prompts in your Terminal and a browser. Then return to this page.

```bash
source ~/.bashrc
```

Now you can use the SDK to install additional components.

#### Install kubectl

If you haven't already, install `kubectl`. You'll use `kubectl` commands to interact with your K8ssandra resources. 

Get `kubectl` from the Google Cloud SDK, and verify the version:

```console
gcloud components install kubectl
```

Verify the kubectl install:

```bash
which kubectl
```

**Output**:

```bash
/Users/thats.me/google-cloud-sdk/bin/kubectl
```

#### Install helm v3

If you haven't already, install Helm v3. On the [Helm site](https://helm.sh/docs/), notice the latest version near the top-right menu. For example: 3.6.2. 

Example:

```bash
cd
curl -LO https://get.helm.sh/helm-v3.6.2-linux-amd64.tar.gz
tar -zxvf helm-v3.6.2-linux-amd64.tar.gz
sudo mv linux-amd64/helm /usr/local/bin/helm
```

Verify the install:

```bash
which helm
```

**Output**:

```bash
/usr/local/bin/helm
```

### Checkout the `k8ssandra-terraform` GitHub project

Each of our reference deployment may be found in the GitHub [k8ssandra/k8ssandra-terraform](https://github.com/k8ssandra/k8ssandra-terraform) project. Download the latest release or current `main` branch locally.

```bash
mkdir ~/github
cd github
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
cd ~/github/k8ssandra-terraform/gcp
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
gcloud config set project "your-gcp-project-id"
```

**Output**:

```console
Updated property [core/project].
```

### Setup Environment Variables

These values will be used to define where infrastructure is provisioned along with the naming of resources. 

{{% alert title="Tip" color="success" %}}
Be sure to specify your GCP project-id in the `TF_VAR_project_id` variable. 
{{% /alert %}}

```bash
export TF_VAR_environment=prod
export TF_VAR_name=k8ssandra
export TF_VAR_project_id=your-gcp-project-id
export TF_VAR_region=us-central1
```

{{% alert title="Limits" color="primary" %}}
GCP limits the total length of resource names. If your deployment fails to plan try reducing the number of characters in the `environment` and `name` parameters.
{{% /alert %}}

### Provision Infrastructure

We begin this process by initializing our environment and configuring a workspace. To start we run `terraform init` which handles pulling down any plugins required and configures the backend.

```bash
cd ~/github/k8ssandra-terraform/gcp/env
terraform init
```

**Output**:

```bash
Initializing modules...

Initializing the backend...

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

With the workspace configured we now instruct terraform to `plan` the required changes to our infrastructure (in this case creation).

```console
terraform plan
```

{{% alert title="Tip" color="warning" %}}
In certain circumstances, Google Cloud may return the following message after you enter `terraform plan`:

```bash
Error: Attempted to load application default credentials 
since neither `credentials` nor `access_token` was set 
in the provider block.  No credentials loaded.
```

If encountered, run this command to use your gcloud credentials:

`gcloud auth application-default login`

Follow the prompts to authenticate. Then run the `terraform plan` command again.
{{% /alert %}}

**Abbreviated output during successful terraform plan**:

```console
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

Do you want to perform these actions?
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

With the GKE cluster deployed you may now continue with the next step, [retrieving the kubeconfig](#retrieve-kubeconfig).

## Retrieve `kubeconfig`

After provisioning the GKE cluster we must request a copy of the `kubeconfig`. This provides the `kubectl` command with all connection information including TLS certificates and IP addresses for Kube API requests.

In the commands, specify your GCP project-id instead of `k8ssandra-testing` - it's the project-id from the reference implementation (as an example).  

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

With all of the infrastructure provisioned we can now focus on installing K8ssandra. This will require configuring a service account for the backup and restore service (Medusa), creating a set of Helm variable overrides, and setting up GKE specific ingress configurations.

### Create Backup / Restore Service Account Secrets

In order to allow for backup and restore operations, we must create a service account for the Medusa operator which handles coordinating the movement of data to and from Google Cloud Storage (GCS) buckets. As part of the provisioning sections a service account was generated for this purposes. Here we will retrieve the authentication JSON file for this account and push it into Kubernetes as a secret.

Looking at the output of `terraform plan` and `terraform apply` we can see the name of the service account which has been provisioned. Here we use `terraform output` to retrieve keys for use by Medusa. In our reference implementation this value is `prod-k8ssandra-sa@k8ssandra-testing.iam.gserviceaccount.com`.

Because Medusa will expect a JSON file to pass in the secret, we'll use `-raw` to avoid a Terraform output issue: 

```console
terraform output -raw service_account_key > medusa_gcp_key
```

With the key file on our local machine we can now push this file to Kubernetes as a secret with `kubectl`.  

```bash
kubectl create secret generic prod-k8ssandra-medusa-key --from-file=medusa_gcp_key.json=medusa_gcp_key
```

**Output**:

```bash
secret/prod-k8ssandra-medusa-key created
```

{{% alert title="Important" color="primary" %}}
The name of the JSON key file within the secret MUST be `medusa_gcp_key.json`. _Any_ other value will result in Medusa not finding the secret and backups failing. The Medusa container within each Cassandra pod will not start. 

To ensure that doesn't happen, notice in the `kubectl create secret...` command above how we pushed the content of the generated **raw** `medusa_gcp_key` into the `medusa_gcp_key.json` file that Medusa expects to find in the deployed secret.

{{% /alert %}}

This secret, `prod-k8ssandra-medusa-key`, can now be referenced in our K8ssandra configuration to allow for backing up data to GCS with Medusa.

### Generate `gke.values.yaml`

Here is a reference Helm `values.yaml` file with configuration options for running K8ssandra in GKE.

{{< readfilerel file="gke.values.yaml"  highlight="yaml" >}}

{{% alert title="Important" color="primary" %}}
Take note of the comments in this file. If you have changed the name of your secret, are deploying in a different region, or have tweaked any other values it is imperative that you update this file before proceeding.
{{% /alert %}}

### Deploy K8ssandra with Helm

With a `values.yaml` file generated, which details specific configuration overrides, we can now deploy K8ssandra via Helm.

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

Now install K8ssandra and specify the `gke.values.yaml` file that you customized in a prior step:

```bash
helm install prod-k8ssandra k8ssandra/k8ssandra -f gke.values.yaml
```

**Output**:

```bash
NAME: prod-k8ssandra
LAST DEPLOYED: Tue Jul  6 18:41:56 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
```

After a few minutes, you can start checking the deployed pods. Note that the Stargate pods are configured to wait until after all the Cassandra Stateful Set pods have been initialized and are running successfully. During the first few minutes of the K8ssandra deployment, you may notice that the Stargate pods are in an `Init` phase. Here's the view of all pods after 9+ minutes:

```bash
kubectl get pods
```

**Output**:

```bash
NAME                                                  READY   STATUS    RESTARTS   AGE
prod-k8ssandra-cass-operator-56446cc654-zh4tg         1/1     Running   0          9m52s
prod-k8ssandra-dc1-stargate-68cb9d56d6-g2xfp          1/1     Running   4          9m52s
prod-k8ssandra-dc1-stargate-68cb9d56d6-kv6dp          1/1     Running   3          9m52s
prod-k8ssandra-dc1-stargate-68cb9d56d6-w8bmh          1/1     Running   5          9m52s
prod-k8ssandra-dc1-us-central1-a-sts-0                3/3     Running   0          9m37s
prod-k8ssandra-dc1-us-central1-c-sts-0                3/3     Running   0          9m37s
prod-k8ssandra-dc1-us-central1-f-sts-0                3/3     Running   0          9m37s
prod-k8ssandra-grafana-5f8d54d5fc-jv2p6               2/2     Running   0          9m52s
prod-k8ssandra-kube-promet-operator-fc975b8f4-rdsbd   1/1     Running   0          9m51s
prod-k8ssandra-medusa-operator-b9577db9-wxq4c         1/1     Running   0          9m52s
prod-k8ssandra-reaper-747c48d7c6-p9zn9                1/1     Running   0          5m41s
prod-k8ssandra-reaper-operator-d9599c75f-l8ccr        1/1     Running   0          9m52s
prometheus-prod-k8ssandra-kube-promet-prometheus-0    2/2     Running   1          9m49s
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
