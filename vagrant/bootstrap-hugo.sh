#!/usr/bin/env bash

echo "Bootstraping hugo..."

## For Ubuntu 18.04 (bionic) apt-get will install the non-extended version of hugu
## we require the extended version, so it must be built from source here

echo "mkdir ~/src"
mkdir ~/src

echo "cd ~/src/"
cd ~/src/

echo "git clone https://github.com/gohugoio/hugo.git"
git clone https://github.com/gohugoio/hugo.git

echo "cd hugo"
cd hugo

echo "go install --tags extended"
go install --tags extended

echo "Completed bootstraping hugo"
