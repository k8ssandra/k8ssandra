# K8ssandra Helm Charts
K8ssandra is installed and configured through Helm charts.

[Helm 3](https://helm.sh/) must be installed to use the charts.

The charts are not yet deployed to a Helm repo. For now you need use the charts directly from this Git repo.

# Getting Started
First you need to deploy the k8ssandra stack:

```
$ helm install k8ssandra charts/k8ssandra
```

Next, create a cluster:

```
$ helm install k8ssandra-cluster charts/k8ssandra-cluster
```