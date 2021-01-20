#!/usr/bin/env bash

## k8ssandra environment bootstrapping

echo "Bootstrapping k8ssandra environment..."

echo "helm repo add k8ssandra https://helm.k8ssandra.io"
helm repo add k8ssandra https://helm.k8ssandra.io

echo "helm repo add traefik https://helm.traefik.io/traefik"
helm repo add traefik https://helm.traefik.io/traefik

echo "helm repo update"
helm repo update

echo "helm install traefik traefik/traefik -n traefik --create-namespace -f traefik.values.yaml"
helm install traefik traefik/traefik -n traefik --create-namespace -f traefik.values.yaml

echo "kubectl apply -f stargate.ingress.yaml"
kubectl apply -f stargate.ingress.yaml

echo "Completed bootstrapping k8ssandra environment."