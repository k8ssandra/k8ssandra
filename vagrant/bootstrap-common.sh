#!/usr/bin/env bash

## Common bootstraping

echo "Bootstraping common components..."

echo "apt-get update"
apt-get update

echo "apt-get install -y git"
apt-get install -y git-all

echo "Completed bootstraping common components"
