#!/usr/bin/env bash

echo "Bootstraping kind..."

## Install kind

echo "curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.9.0/kind-linux-amd64"
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.9.0/kind-linux-amd64

echo "chmod +x ./kind"
chmod +x ./kind

echo "mv ./kind /usr/local/bin/kind"
mv ./kind /usr/local/bin/kind

echo "Completed bootstraping kind"
