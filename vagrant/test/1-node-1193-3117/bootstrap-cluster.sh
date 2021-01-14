#!/usr/bin/env bash

## Cluster bootstrapping

echo "Bootstrapping cluster..."

echo "sudo kind create cluster --config kind-config.yaml"
sudo kind create cluster --config kind.config.yaml

echo "sudo chown vagrant ~/.kube/config"
sudo chown vagrant ~/.kube/config

echo "sudo chgrp vagrant ~/.kube/config"
sudo chgrp vagrant ~/.kube/config

echo "sudo kubectl config use-context kind-kind"
sudo kubectl config use-context kind-k8ssandra-cluster-1-1193-3117

echo "sudo helm repo add k8ssandra https://helm.k8ssandra.io"
sudo helm repo add k8ssandra https://helm.k8ssandra.io

echo "sudo helm repo add traefik https://helm.traefik.io/traefik"
sudo helm repo add traefik https://helm.traefik.io/traefik

echo "sudo helm repo update"
sudo helm repo update

echo "sudo helm install traefik traefik/traefik -n traefik --create-namespace -f traefik.values.yaml"
sudo helm install traefik traefik/traefik -n traefik --create-namespace -f ../common/traefik.values.yaml

echo "sudo kubectl apply -f k8ssandra-stargate.ingress.yaml"
sudo kubectl apply -f ../common/k8ssandra-stargate.ingress.yaml

echo "sudo helm install k8ssandra k8ssandra/k8ssandra"
sudo helm install k8ssandra k8ssandra/k8ssandra

echo "sudo helm install k8ssandra-cluster k8ssandra/k8ssandra-cluster -f k8ssandra-cluster.yaml"
sudo helm install k8ssandra-cluster k8ssandra/k8ssandra-cluster -f k8ssandra-cluster.values.yaml

echo "Completed bootstrapping cluster..."