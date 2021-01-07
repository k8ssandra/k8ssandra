#!/usr/bin/env bash

echo "Finalizing bootstraping"

## Disable firewall

echo "ufw disable"
ufw disable

echo "apt-get upgrade -y"
apt-get upgrade -y

echo "Completed finalizing bootstraping"
