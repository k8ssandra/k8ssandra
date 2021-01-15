#!/usr/bin/env bash

## Cluster bootstrapping

echo "Bootstrapping cluster..."

echo "kind create cluster --config kind-config.yaml"
kind create cluster --config kind.config.yaml

echo "kubectl config use-context kind-k8ssandra-cluster-1-1193-3117"
kubectl config use-context kind-k8ssandra-cluster-1-1193-3117

echo "helm repo add k8ssandra https://helm.k8ssandra.io"
helm repo add k8ssandra https://helm.k8ssandra.io

echo "helm repo add traefik https://helm.traefik.io/traefik"
helm repo add traefik https://helm.traefik.io/traefik

echo "helm repo update"
helm repo update

echo "helm install traefik traefik/traefik -n traefik --create-namespace -f ../common/traefik.values.yaml"
helm install traefik traefik/traefik -n traefik --create-namespace -f ../common/traefik.values.yaml

echo "kubectl apply -f ../common/k8ssandra-stargate.ingress.yaml"
kubectl apply -f ../common/k8ssandra-stargate.ingress.yaml

echo "helm install k8ssandra k8ssandra/k8ssandra -f k8ssandra.values.yaml"
helm install k8ssandra k8ssandra/k8ssandra -f k8ssandra.values.yaml

echo "Completed bootstrapping cluster..."