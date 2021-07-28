---
title: "Private registries"
linkTitle: "Private registries"
weight: 1
description: "Optional steps to use private registries with K8ssandra and its deployed components."
---

Starting in the 1.3.0 release, K8ssandra supports the use of private registries. This feature is designed for those users who have limited or restricted access to public registries such as [Docker Hub](https://hub.docker.com). Site Reliability Engineers and developers can use the examples in this topic (and linked related topics) to provide access to all the images required for their K8ssandra environment.

## Introduction

All images used by K8ssandra are pulled from Docker Hub, a publicly accessible registry. There are situations where some users need the ability to pull images from a different registry. Examples:

* Consider ["air gapped"](https://en.wikipedia.org/wiki/Air_gap_(networking)) deployments, where there is no access to the public Internet. This scenario means images cannot be pulled from Docker Hub or any other external registry. Images must be pulled from a registry on the user's internal network.

* A third-party organization may provide its own registry for its users where images are built from source, and go through additional security audits. Again, relying on the default registry used by K8ssandra -- Docker Hub -- is insufficient.

## Image coordinates

K8ssandra 1.3.0 adds the ability to fully specify coordinates for every image that it deploys. To better understand what we mean by image coordinates, let's look closer.

### Format of image coordinates

Consider this example:

```
docker.io/k8ssandra/cass-management-api:3.11.10-v0.1.26
```

* `docker.io` is the registry.

* `k8ssandra` is the repository.

* `cass-management-api` is the image.

* `3.11.10-v0.1.26` is the image tag.

Given this, the complete **format** of coordinates looks like this:

```
registry/repository/image:tag
```

K8ssandra allows you to configure each of these parts for each image that it deploys.

### Chart properties YAML

For each container that K8ssandra deploys, there are chart properties like this:

```yaml
<container-name>:
  image:
    registry:
    repository:
    tag:
    pullPolicy: 
```

Example YAML:

```yaml
cassandra:
  serviceAccount: cassandra
  image:
    registry: myregistry
    repository: myrepo/cass-management-api
    tag: 3.11.10-v0.1.26
    pullPolicy: IfNotPresent
```

{{% alert title="Tip" color="success" %}}
While `pullPolicy` is shown above, it is not part of the image coordinates. Possible values for `pullPolicy` are:

* `IfNotPresent` (default)
* `Always` (default when the image tag is `latest`)
* `Never`

See [ImagePullPolicy](https://helm.sh/docs/chart_best_practices/pods/#imagepullpolicy) in the Helm v3 documentation.
{{% /alert %}}

### List of chart properties

Here is a list of chart properties that correspond to K8ssandra deployed containers having an `image` property:

| Chart property / container name                          | Summary                               |
| -------------------------------------------------------- | ------------------------------------- |
| `cassandra`                                              | Runs the Management API and Cassandra&reg; |
| `cassandra.configBuilder`                                | An init container that provisions Cassandra configuration files such as cassandra.yaml |
| `cassandra.jmxCredentialsConfig` | An init container that configures authentication for remote JMX access |
| `cassandra.loggingSidecar` | A container that tails Cassandra's system.log file |
| `stargate` | Runs Stargate, which is an API for interacting with Cassandra |
| `stargate.waitForCassandra` | An init container that waits until all the Cassandra pods are ready |
| `reaper` | Run Reaper, which provides repair functionality for Cassandra data |
| `medusa` | Runs the medusa container that performs backups, and the `medusa-restore` init container that performs restores |
| `cleaner` | A Helm pre-delete hook that deletes the CassandraDatacenter |
| `client` | Used in the `crd-upgrader` pre-upgrade to apply Custom Resource Definition (CRD) updates |
| `cass-operator` | Cass Operator |
| `reaper-operator` | Reaper Operator |
| `medusa-operator` | Medusa Operator |
| `kube-prometheus-stack.prometheus-operator` | Prometheus Operator |
| `kube-prometheus-stack.prometheus-operator.prometheusSpec` | Prometheus to collect Cassandra cluster, OS, and node metrics |
| `grafana` | Grafana to visualize the metrics collected by Prometheus |

{{% alert title="Tip" color="success" %}}
Operators are software extensions to Kubernetes that make use of custom resources to manage applications and their components. Thus, for example, "Reaper Operator" deploys and configures Reaper. "Reaper" itself manages the actual Cassandra repair operations. Similarly, "Prometheus Operator" deploys and configures Prometheus. "Prometheus" itself manages the actual collection of relevant OS / Cassandra metrics. "Medusa Operator" configures and orchestrates the backup and restore operations. "Medusa" itself runs the container that performs backups of Cassandra data. 
{{% /alert %}}

### Cassandra container defaults

By default, K8ssandra uses a version mapping to decide on the image to use for Cassandra. K8ssandra looks up the image coordinates based on the value of the `cassandra.version` property. (Starting in K8ssandra 1.3.0, the default `cassandra.version` is `4.0.0`.)

If you choose to set the `cassandra.image` property, the version mapping won't be used. You will be responsible for ensuring that a valid image is used. 

## Image Pull Secrets and Service Accounts

Private registries require authentication to pull images. Pods need credentials in order to pull images. 
This is accomplished with **service accounts** and **image pull secrets**. Every pod has a service account, and image pull secrets can be assigned to a service account. 

For related details, see these topics in the Kubernetes documentation:

* [Pull an Image from a Private Registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/).

* [Add ImagePullSecrets to a service account ](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account)
 for details on configuring and creating a service account with image pull secrets.

### Properties for configuring service accounts

The following table lists the relevant properties for configuring service accounts:

| Pod       | Chart property                     | Summary                               |
| ----------|----------------------------------- | ------------------------------------- |
| Cassandra | `cassandra.serviceAccount`         | Specifies the name of an existing service account. If not specified the default service account is used. |
| Stargate  | `stargate.serviceAccount`          | Specifies the name of an existing service account. If not specified the default service account is used. |
| Reaper    | `reaper.serviceAccount`            | Specifies the name of an existing service account. If not specified the default service account is used. |
| Prometheus Operator | `kube-prometheus-stack.prometheusOperator.serviceAccount.name` | Specifies the name of the service account that will be created. The `kube-prometheus-stack` chart creates this service account by default. |
| Prometheus | `kube-prometheus-stack.prometheus.serviceAccountName` | Specifies the name of the service account that will be created. The `kube-prometheus-stack` chart creates this service account by default. |
| Reaper Operator | `reaper-operator.serviceAccount` | Configure the name and image pull secrets for the service account that will be created. The `reaper-operator chart` creates this service account by default. |
| Medusa Operator | `medusa-operator.serviceAccount.name` | Configure the name and image pull secrets for the service account that will be created. The `k8ssandra` chart creates this service account by default. 
| Cass Operator | `cass-operator.serviceAccount.name` | Configure the name and image pull secrets for the service account that will be created. The `cass-operator` chart creates this service account by default. |
| Cleaner | `cleaner.serviceAccount` | Specifies the name of the service account that will be created. If specified, the service account must already exist. |
| CRD Updater | `client.serviceAccount` | Specifies the name of the service account that will be created. If specified, the service account must already exist. |


### Configure Image Pull Secrets

There are some differences in the way K8ssandra handles the configuration of service accounts, depending on the component. Let's look 
at each component to see how to configure image pull secrets.

#### Cassandra, Stargate, Reaper

For each of these components -- Cassandra, Stargate, Reaper -- K8ssandra does not create or configure the service account. For details on configuring and creating a service account with image pull secrets for these components, see [Configure Service Accounts for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account) in the Kubernetes documentation. 

#### Prometheus Operator

Configuration for its service account:

```yaml
global:
  imagePullSecrets:
    - myregistrykey
kube-prometheus-stack:
  prometheusOperator:
      serviceAccount:
     name: prometheus-operator
```

#### Prometheus

The configuration for the Prometheus service account is the same as for Prometheus Operator, but each requires its own service account:

```yaml
global:
  imagePullSecrets:
    - myregistrykey
kube-prometheus-stack:
  prometheus:
    serviceAccount:
      name: prometheus-operator
```

#### Reaper Operator    

Configuration for its service account:

```yaml
reaper-operator:
  imagePullSecrets:
  - myregistrykey
  serviceAccount:
    name: reaper-operator
```

#### Medusa Operator

Configuration for its service account:

```yaml
medusa-operator:
  imagePullSecrets:
  - myregistrykey
  serviceAccount:
    name: medusa-operator
```

#### Cass Operator

Configuration for its service account:

```yaml
cass-operator:
  imagePullSecrets:
  - myregistrykey
  serviceAccount:
    name: cass-operator
```

#### cleaner

Configuration for its service account:

```yaml
cleaner:
  serviceAccount: cleaner
```

#### CRD Updater

Configuration for its service account:

```yaml
client:
  serviceAccount: cleaner
```

{{% alert title="Tip" color="success" %}}
You must create this service account for CRD Updater separately.
{{% /alert %}}


## Image property examples

Let's look at two examples that demonstrates how to work with the image properties.

### Custom stargate example with the default registry

Suppose we are working on a fork of [Stargate](https://github.com/stargate/stargate) and we want to deploy it in K8ssandra. We also want to use a custom image for the `waitForCassandra` init container. We could create and apply an overrides values YAML file with properties like this:

``` yaml
stargate:
  image:    
    repository: example-user/stargate
    tag: 1.0.29-dev
  waitForCassandra:
    image:
      repository: example-user/stargate-init
      tag: 0.1.0
```

Notice that the property names (`stargate` and `waitForCassandra`) in the example correspond directly to properties in the [list of chart properties]({{< relref "#list-of-chart-properties" >}}) table above. 

Because the `registry` property was not specified, the default (`docker.io`) will be used.

### Complete example with a private registry

Here is a complete example of an overrides file that is configured to use a private registry. See the comments embedded in the YAML.

```yaml
# We have to specify global image pull secrets for kube-prometheus-stack
global:
  imagePullSecrets:
    - myregistrykey

cassandra:
  # This service account needs to be configured and created prior to installing
  # the chart.
  serviceAccount: cassandra
  image:
    registry: myregistry
    repository: myrepo/cass-management-api
    # We have to specify the tag since we are not using the version/image mapping
    tag: 3.11.10-v0.1.26

  configBuilder:
    image:
      registry: myregistry
      repository: myrepo/cass-config-builder

  jmxCredentialsConfig:
    image:
      registry: myregistry
      repository: myrepo/busybox

  loggingSidecar:
    image:
      registry: myregistry
      repository: myrepo/system-logger

stargate:
  # This service account needs to be configured and created prior to installing
  # the chart.
  serviceAccount: stargate
  image:
    registry: myregistry
    repository: myrepo/stargate-3_11
    # We have to specify the tag since we are not using the image mapping
    tag: 1.0.29

  waitForCassandra:
    image:
      registry: myregistry
      repository: myrepo/alpine

reaper:
  # This service account needs to be configured and created prior to installing
  # the chart.
  serviceAccount: reaper
  image:
    registry: myregistry
    repository: myrepo/cassandra-reaper

medusa:
  # We don't configure the service account here because Medusa
  # runs in the Cassandra pod.
  enabled: true
  image:
    registry: myregistry
    repository: myrepo/medusa

cleaner:
  serviceAccount: cleaner
  image:
    registry: myregistry
    repository: myrepo/k8ssandra-tools

client:
  serviceAccount: client
  image:
    registry: myregistry
    repository: myrepo/k8ssandra-tools

cass-operator:
  imagePullSecrets:
  - myregistrykey
  serviceAccount:
    name: cass-operator
  image:
    registry: myregistry
    repository: myrepo/cass-operator

reaper-operator:
  imagePullSecrets:
  - myregistrykey
  serviceAccount:
    name: reaper-operator
  image:
    registry: myregistry
    repository: myrepo/reaper-operator

medusa-operator:
  imagePullSecrets:
  - myregistrykey
  serviceAccount:
    name: medusa-operator
  image:
    registry: myregistry
    repository: myrepo/medusa-operator

kube-prometheus-stack:
  prometheusOperator:
    serviceAccount:
      name: prometheus-operator
    image:
      repository: myregistry/myrepo/prometheus-operator
  prometheus:
    serviceAccount:
      name: prometheus-operator
    image:
      repository: myregistry/myrepo/prometheus
```
## Related issues and PR

From the K8ssandra GitHub repo, here are the related issues and pull request for your reference:

* https://github.com/k8ssandra/k8ssandra/pull/901
* https://github.com/k8ssandra/k8ssandra/issues/420
* https://github.com/k8ssandra/k8ssandra/issues/839
* https://github.com/k8ssandra/k8ssandra/issues/840

## Next steps

Explore other K8ssandra [tasks]({{< relref "/tasks" >}}).

See the [Helm Chart]({{< relref "/reference/helm-charts" >}}) reference topics.
