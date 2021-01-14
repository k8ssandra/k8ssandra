#!/usr/bin/env bash

echo "Finalizing bootstrapping"

## Disable firewall

echo "ufw disable"
ufw disable

## Add docker permissions to current user - prevent sudo requirements

#echo "sudo usermod -aG docker $USER"
#sudo usermod -aG docker $USER

reboot

echo "Completed finalizing bootstrapping"
