#!/usr/bin/env bash

echo "Bootstraping kubectl..."

## Install kubectl
echo "apt-get update && apt-get install -y apt-transport-https gnupg2 curl"
apt-get update && apt-get install -y apt-transport-https gnupg2 curl

echo "curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -"
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

echo "echo \"deb https://apt.kubernetes.io/ kubernetes-xenial main\" | tee -a /etc/apt/sources.list.d/kubernetes.list"
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list

echo "apt-get update"
apt-get update

echo "apt-get install -y kubectl"
apt-get install -y kubectl

echo "Completed bootstraping kubectl"