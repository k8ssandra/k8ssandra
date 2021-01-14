#!/usr/bin/env bash

## Common bootstrapping

echo "Bootstrapping common components..."

echo "apt-get update"
apt-get update

echo "apt-get install -y git"
apt-get install -y git-all

echo "Completed bootstrapping common components"
