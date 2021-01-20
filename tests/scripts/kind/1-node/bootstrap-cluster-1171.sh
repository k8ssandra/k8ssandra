#!/usr/bin/env bash

## Cluster environment bootstrapping

echo "Bootstrapping initial cluster environment..."

echo "kind delete cluster --name k8ssandra"
kind delete cluster --name k8ssandra

echo "kind create cluster --config kind-config.yaml"
kind create cluster --image "kindest/node:v1.17.11"  --config kind.config.yaml

echo "kubectl config use-context kind-k8ssandra"
kubectl config use-context kind-k8ssandra

echo "Bootstrapping initial cluster environment..."