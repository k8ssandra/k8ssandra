#!/usr/bin/env bash

echo "Bootstraping docker permissions..."

echo "sudo usermod -aG docker $USER"
sudo usermod -aG docker $USER

echo "Rebooting VM..."
sudo reboot