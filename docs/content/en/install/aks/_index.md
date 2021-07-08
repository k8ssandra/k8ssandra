---
title: "Install K8ssandra on AKS"
linkTitle: "Azure AKS"
weight: 2
description: >
  Complete **production** ready environment of K8ssandra on Azure Kubernetes Service (AKS).
---

[Azure Elastic Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/) or "AKS" is a managed Kubernetes service that makes it easy for you to run Kubernetes on Azure. AKS offers serverless Kubernetes, an integrated continuous integration and continuous delivery (CI/CD) experience, and enterprise-grade security and governance.

{{% alert title="Tip" color="success" %}}
Also available in followup topics are post-install steps and role-based considerations for [developers]({{< relref "/quickstarts/developer">}}) or [site reliability engineers]({{< relref "/quickstarts/site-reliability-engineer">}}) (SREs).
{{% /alert %}}

Note that at this time backup and restore support for Azure Blob Storage is still in progress. Follow [k8ssandra#685](https://github.com/k8ssandra/k8ssandra/issues/685) for progress in this area.

## Deployment

This topic covers provisioning and installing the following infrastructure resources.

* AKS cluster
* AKS default node pool
* Managed Identity
* Storage Account
* Storage container
* Virtual Network(Vnet)
* Subnets
* Network Security Group
* NAT Gateway
* Public IP
* Route Table
* Route Table association

On this infrastructure the K8ssandra installation will consist of the following workloads.

* 3x instance Apache Cassandra cluster
* 3x instance Stargate deployment
* 1x instance Prometheus deployment
* 1x instance Grafana deployment
* 1x instance Reaper deployment

Feel free to update the parameters used during this guide to match your target deployment. This should be considered a minimum for production workloads.

{{% alert title="Tip" color="primary" %}}
If you already have an AKS cluster provisioned skip to the [installation procedures](#install-k8ssandra).
{{% /alert %}}

## Terraform

As a convenience we provide reference [Terraform](https://www.terraform.io/) modules for orchestrating the provisioning of cloud resources necessary to run K8ssandra.

### Tools

| Tool | Version | 
|------|---------|
| [Terraform](https://www.terraform.io/downloads.html) | 0.14 |
| [Azurerm provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest) | ~>3.0 |
| [Helm](https://helm.sh/) | 3 |
| [AZ CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) | 2.22.1 |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.17.17 |
| [Python](https://www.python.org/downloads/) | 3 |

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
cd k8ssandra-terraform/azure
```

### Configure `az` CLI

Ensure you have authenticated your `az` client by running the following command:

```console
az login
```

**Output**:

```console
The default web browser has been opened at https://login.microsoftonline.com/common/oauth2/authorize. Please continue the login in the web browser. If no web browser is available or if the web browser fails to open, use device code flow with `az login --use-device-code`.
You have logged in. Now let us find all the subscriptions to which you have access...
[
  {
    "cloudName": "AzureCloud",
    "homeTenantId":  "00000000-0000-0000-0000-000000000000",
    "id": "00000000-0000-0000-0000-000000000000",
    "isDefault": true,
    "managedByTenants": [],
    "name": "k8ssandra",
    "state": "Enabled",
    "tenantId": "00000000-0000-0000-0000-000000000000",
    "user": {
      "name": "kate.sandra@k8ssandra.io",
      "type": "user"
    }
  }
]
```

### Setup Environment Variables

These values will be used to define where infrastructure is provisioned along with the naming of resources.

```bash
export TF_VAR_environment=prod
export TF_VAR_name=k8ssandra
export TF_VAR_region=eastus
```

### Provision Infrastructure

We begin this process by initializing our environment. To start we run `terraform init` which handles pulling down any plugins required and configures the backend.

```console
cd env
terraform init
```

**Output**:

```console
Initializing modules...
- aks in ../modules/aks
- iam in ../modules/iam
- storage in ../modules/storage
- vnet in ../modules/vnet

Initializing the backend...

Initializing provider plugins...
- Reusing previous version of hashicorp/azurerm from the dependency lock file
- Installing hashicorp/azurerm v2.49.0...
- Installed hashicorp/azurerm v2.49.0 (self-signed, key ID 34365D9472D7468F)

# Output reduced for brevity

Terraform has been successfully initialized!
```

With the workspace configured we now instruct terraform to `plan` the required changes to our infrastructure (in this case creation).

```console
terraform plan
```

**Output**:

```console
Terraform used the selected providers to generate the following execution plan. Resource actions
are indicated with the following symbols:
  + create

Terraform will perform the following actions:

# Output reduced for brevity

Plan: 21 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + aks_fqdn           = (known after apply)
  + aks_id             = (known after apply)
  + connect_cluster    = "az aks get-credentials --resource-group prod-k8ssandra-resource-group --name prod-k8ssandra-aks-cluster"
  + resource_group     = "prod-k8ssandra-resource-group"
  + storage_account_id = (known after apply)
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

Apply complete! Resources: 21 added, 0 changed, 0 destroyed.

Outputs:

aks_fqdn = "prod-k8ssandra-00000000.hcp.eastus.azmk8s.io"
aks_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/prod-k8ssandra-resource-group/providers/Microsoft.ContainerService/managedClusters/prod-k8ssandra-aks-cluster"
connect_cluster = "az aks get-credentials --resource-group prod-k8ssandra-resource-group --name prod-k8ssandra-aks-cluster"
resource_group = "prod-k8ssandra-resource-group"
storage_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-k8ssandra-resource-group/providers/Microsoft.Storage/storageAccounts/prodk8ssandrastorage"
```

With the AKS cluster deployed you may now continue with [retrieving the kubeconfig](#retrieve-kubeconfig).

## Retrieve `kubeconfig`

After provisioning the AKS cluster we must request a copy of the `kubeconfig`. This provides the `kubectl` command with all connection information including TLS certificates and IP addresses for Kube API requests.

```console
az aks get-credentials --resource-group prod-k8ssandra-resource-group --name prod-k8ssandra-aks-cluster
```

**Output**:

```bash
Merged "prod-k8ssandra-aks-cluster" as current context in /home/bradfordcp/.kube/config
```

```bash
kubectl cluster-info
```

**Output**:

```bash
Kubernetes control plane is running at https://prod-k8ssandra-00000000.hcp.eastus.azmk8s.io:443
CoreDNS is running at https://prod-k8ssandra-00000000.hcp.eastus.azmk8s.io:443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
Metrics-server is running at https://prod-k8ssandra-00000000.hcp.eastus.azmk8s.io:443/api/v1/namespaces/kube-system/services/https:metrics-server:/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

```

```bash
kubectl version
```

**Output**:

```bash
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.1", GitCommit:"5e58841cce77d4bc13713ad2b91fa0d961e69192", GitTreeState:"clean", BuildDate:"2021-05-12T14:18:45Z", GoVersion:"go1.16.4", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.9", GitCommit:"6c90dbd9d6bb1ae8a4c0b0778752be06873e7c55", GitTreeState:"clean", BuildDate:"2021-03-22T23:02:49Z", GoVersion:"go1.15.8", Compiler:"gc", Platform:"linux/amd64"}
WARNING: version difference between client (1.21) and server (1.19) exceeds the supported minor version skew of +/-1
```

## Install K8ssandra

With all of the infrastructure provisioned we can now focus on installing K8ssandra. This will require configuring a service account for the backup and restore service, creating a set of Helm variable overrides, and setting up AKS specific ingress configurations.

### Create Backup / Restore Service Account Secrets

In order to allow for backup and restore operations, we must provide a storage account to the Medusa operator which
handles coordinating the movement of data to and from Azure Storage blobs. As part of the provisioning sections, a
storage account was generated for these purposes. Here we will generate a credentials file for this account and push it
into Kubernetes as a secret.

Inspect the output of `terraform output` to retrieve the information we need: the resource group, and the storage
account name. The resource group is displayed as part of the output; the storage account name to use is the last part of
the `storage_account_id` entry.

In our reference implementation, the resource group is `prod-k8ssandra-resource-group` and the storage account name is
`prodk8ssandrastorage`.

If in doubt, you can retrieve all the available storage accounts and corresponding resource groups with:

```console
az storage account list --query '[].{StorageAccountName:name,ResourceGroup:resourceGroup}' -o table
``` 

Now we are going to retrieve one of the access keys for our target storage account, and use it generate a credentials
file in your local machine called `credentials.json`:

```console
az storage account keys list \
   --account-name prodk8ssandrastorage \
   --query "[0].value|{storage_account:'prodk8ssandrastorage',key:@}" > credentials.json
```

The generated file should look like this:

```bash
{
    "storage_account": "prodk8ssandrastorage",
    "key": "<ACCESS KEY>"
}
```

We can now push this file to Kubernetes as a secret with `kubectl`:

```bash
kubectl create secret generic prod-k8ssandra-medusa-key \
    --from-file=medusa_azure_credentials.json=./credentials.json
```

**Output**:

```bash
secret/prod-k8ssandra-medusa-key created
```

{{% alert title="Important" color="primary" %}} The name of the JSON credentials file within the secret MUST be
`medusa_azure_credentials.json`. _Any_ other value will result in Medusa not finding the secret and backups failing. {{%
/alert %}}

This secret, `prod-k8ssandra-medusa-key`, can now be referenced in our K8ssandra configuration to allow for backing up
data to Azure with Medusa.

### Generate `aks.values.yaml`

Here is a reference Helm `values.yaml` file with configuration options for running K8ssandra in AKS.

{{< readfilerel file="aks.values.yaml"  highlight="yaml" >}}

{{% alert title="Important" color="primary" %}}
Take note of the comments in this file. If you have changed the name of your secret, are deploying in a different region, or have tweaked any other values it is imperative that you update this file before proceeding.
{{% /alert %}}

### Deploy K8ssandra with Helm

With a `values.yaml` file generated which details out specific configuration overrides we can now deploy K8ssandra via Helm.

```bash
helm install prod-k8ssandra k8ssandra/k8ssandra -f aks.values.yaml
```

**Output**:

```bash
NAME: prod-k8ssandra
LAST DEPLOYED: Fri May 21 16:17:33 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
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
    o4NlBIWcO2KzKuBPU9mB
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

### Destroy AKS Cluster

```console
terraform destroy
```

**Output**:

```bash
# Output omitted for brevity

Plan: 0 to add, 0 to change, 21 to destroy.

Changes to Outputs:
  - aks_fqdn           = "prod-k8ssandra-0000000.hcp.eastus.azmk8s.io" -> null
  - aks_id             = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/prod-k8ssandra-resource-group/providers/Microsoft.ContainerService/managedClusters/prod-k8ssandra-aks-cluster" -> null
  - connect_cluster    = "az aks get-credentials --resource-group prod-k8ssandra-resource-group --name prod-k8ssandra-aks-cluster" -> null
  - resource_group     = "prod-k8ssandra-resource-group" -> null
  - storage_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-k8ssandra-resource-group/providers/Microsoft.Storage/storageAccounts/prodk8ssandrastorage" -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes


# Output omitted for brevity

Destroy complete! Resources: 21 destroyed.
```

## Next steps

With a freshly provisioned cluster on AKS, consider visiting the [developer]({{< relref "/quickstarts/developer" >}}) and [Site Reliability Engineer]({{< relref "/quickstarts/site-reliability-engineer" >}}) quickstarts for a guided experience exploring your cluster. 

Alternatively, if you want to tear down your AKS cluster and / or infrastructure, refer to the section above that covers [cleaning up resources]({{< relref "#cleanup-resources" >}}).

