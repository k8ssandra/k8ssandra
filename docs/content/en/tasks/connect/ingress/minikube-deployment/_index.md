---
title: "Deploy K8ssandra and Traefik with Minikube"
linkTitle: "Minikube Deployment"
toc_hide: true
description: "Deploy Apache Cassandra® on Kubernetes in a local Minikube cluster with Traefik ingress installed and configured."
---

{{< tbs >}}

Minikube is a popular tool for testing k8s clusters locally. When deploying to Minikube, we use the native Kubernetes `port-forward` command from `kubectl` in order to access the services. 

### 1. Install Minikube and start the cluster

Installation instructions for Minikube are [here](https://minikube.sigs.k8s.io/docs/start/).

Once installed, running `minikube start` will bring up a local cluster.

### 2. Install Traefik via Helm

Traefik can be installed via Helm in conjunction with the below Traefik values file - 

### [`Traefik Values`](traefik-values.yaml)

The `traefik.values.yaml` file is [here](traefik-values.yaml).
 
{{< readfilerel file="traefik-values.yaml"  highlight="yaml" >}}

It can be applied to the cluster using the below commands. 

```
$ helm repo add traefik https://helm.traefik.io/traefik
$ helm repo update
$ helm install traefik traefik/traefik -n traefik --create-namespace -f traefik.values.yaml
NAME: traefik
LAST DEPLOYED: Thu Nov 12 16:59:40 2020
NAMESPACE: traefik
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

### 3. Port forward to the Traefik dashboard

Assuming you have used the above configuration, Traefik will have been installed into a new namespace called `traefik`. The Traefik dashboard can now be accessed using `kubectl port-forward --namespace traefik services/traefik 9000:9000`

Minikube offers several options to access a service from the host machine (including via a NodePort on the Minikube host). Additional information can be found in the Minikube [documentation](https://minikube.sigs.k8s.io/docs/handbook/accessing/). Note that it is sometimes preferable to access services via `port-forward` as this is closer to the way that services are accessed on a remote cluster.

### Note - these settings are not suitable for production.

The steps above will create a NodePort service, which will serve the Traefik dashboard to clients outside the k8s cluster. This is not suitable for deployment to cloud platforms (EKS, AKS, GKE etc.) as it will make the dashboard publicly available on any external IPs attached to the node. (Albeit only if access is provided by any security group or firewall on the cluster.) For a cloud-ready configuration, ensure the Traefik dashboard and other internal services are not visible from outside the cluster (except via `port-forward` as above).

## Next steps

Feel free to explore the other [Traefik ingress]({{< relref "/tasks/connect/ingress/" >}}) topics. Also see the additional K8ssandra [tasks]({{< relref "tasks" >}}).

