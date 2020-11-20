# K8ssandra Helm Charts
K8ssandra is installed and configured through Helm charts.

Prerequisites
-------------
In your local environment the following tools are required for provisioning a K8ssandra cluster.
- [Helm 3+](https://helm.sh/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

As K8ssandra deploys on a Kubernetes cluster one must be available to target for installation. This may be a local version running on your development machine, on-premises self-hosted environment, or managed cloud offering. To that end the cluster must be up and available to your kubectl command.
```
# Validate cluster connectivity
kubectl cluster-info
```

If you do not have a Kubernetes cluster available consider one of the following local versions that run within Docker or a virtual machine.

- [K3D](https://k3d.io/)
- [Kind](https://kind.sigs.k8s.io/)
- [OpenShift CodeReady Containers](https://developers.redhat.com/products/codeready-containers/overview)


Configure Helm Repository
-------------------------

K8ssandra is delivered as a collection of Helm Charts. In order to leverage these charts we have provided a k8ssandra Helm Repository for easy installation.

Also add the Traefik Ingress repo - you’ll need its resources to access services from outside the Kubernetes cluster.

```
helm repo add k8ssandra https://helm.k8ssandra.io/
helm repo add traefik https://helm.traefik.io/traefik
helm repo update
```

Alternatively, you may download the individual charts directly from the project’s [releases](https://github.com/k8ssandra/k8ssandra/releases) page

Install K8ssandra
-----------------


From a packaging perspective, K8ssandra is composed of a number of helm charts. The `k8ssandra-tools` chart handles the installation of operators and custom resources. The `k8ssandra-cluster` chart (which you can uniquely name) is focused on provisioning cluster instances. This loose coupling allows for separate lifecycles of components with an easy procedure - submitting just two `helm install` commands.

```
# Install shared dependencies / tooling
helm install k8ssandra-tools k8ssandra/k8ssandra

# Provision a K8ssandra cluster named "k8ssandra-cluster-a" 
helm install k8ssandra-cluster-a k8ssandra/k8ssandra-cluster  

```


In later steps, you can upgrade your k8ssandra-cluster via `helm upgrade` commands, for example to access services from outside Kubernetes via a Traefik Ingress controller.
