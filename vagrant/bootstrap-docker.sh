#!/usr/bin/env bash

## Install Docker CE

echo "Bootstrapping Docker CE..."

apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg2

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key --keyring /etc/apt/trusted.gpg.d/docker.gpg add

add-apt-repository \
  "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) \
  stable"

apt-get update

apt-get install -y \
  containerd.io=1.2.13-2 \
  docker-ce=5:19.03.11~3-0~ubuntu-$(lsb_release -cs) \
  docker-ce-cli=5:19.03.11~3-0~ubuntu-$(lsb_release -cs)

cat <<EOF | sudo tee /etc/docker/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

groupadd docker

usermod -aG docker $USER

mkdir -p /etc/systemd/system/docker.service.d

systemctl daemon-reload

systemctl restart docker

systemctl enable docker

## Disable swap

swapoff -a

sed -i '/ swap / s/^/#/' /etc/fstab

echo "Completed bootstrapping Docker CE"
