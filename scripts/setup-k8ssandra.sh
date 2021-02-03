#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

helm install k8ssandra ../charts/k8ssandra -n k8ssandra --create-namespace



