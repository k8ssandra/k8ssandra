#!/usr/bin/env bash

echo "Extracting cassandra username..."

echo "kubectl get secret k8ssandra-cluster-superuser -o jsonpath=\"{.data.username}\" | base64 --decode"
kubectl get secret k8ssandra-cluster-superuser -o jsonpath="{.data.username}" | base64 --decode
