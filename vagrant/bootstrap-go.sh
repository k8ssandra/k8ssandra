#!/usr/bin/env bash

echo "Bootstrapping golang..."

echo "add-apt-repository -y ppa:longsleep/golang-backports"
add-apt-repository -y ppa:longsleep/golang-backports 

echo "apt-get update"
apt-get update

echo "apt-get install -y golang-go"
apt-get install -y golang-go

echo "export PATH=$PATH:/usr/local/go/bin"
export PATH=$PATH:/usr/local/go/bin

echo "Completed bootstrapping golang"
