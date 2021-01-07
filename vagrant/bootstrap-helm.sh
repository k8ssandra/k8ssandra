#!/usr/bin/env bash

echo "Bootstraping helm..."

echo "curl https://baltocdn.com/helm/signing.asc | apt-key add -"
curl https://baltocdn.com/helm/signing.asc | apt-key add -

echo "apt-get install -y apt-transport-https"
apt-get install -y apt-transport-https

echo "echo \"deb https://baltocdn.com/helm/stable/debian/ all main\" | tee /etc/apt/sources.list.d/helm-stable-debian.list"
echo "deb https://baltocdn.com/helm/stable/debian/ all main" | tee /etc/apt/sources.list.d/helm-stable-debian.list

echo "apt-get update"
apt-get update

echo "apt-get install helm"
apt-get install -y helm

echo "Completed bootstraping helm"
