#!/usr/bin/env bash

## Cluster bootstrapping

echo "Bootstrapping cluster..."

echo "helm install k8ssandra k8ssandra/k8ssandra --set k8ssandra.cassandraVersion=\"3.11.7\" -f k8ssandra.values.yaml"
helm install k8ssandra k8ssandra/k8ssandra --set k8ssandra.cassandraVersion="3.11.7" -f k8ssandra.values.yaml

echo "Completed bootstrapping cluster..."

watch kubectl get pods
