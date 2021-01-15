#!/usr/bin/env bash

echo "Extracting cassandra password..."

echo "kubectl get secret k8ssandra-superuser -o jsonpath=\"{.data.password}\" | base64 --decode"
kubectl get secret k8ssandra-superuser -o jsonpath="{.data.password}" | base64 --decode