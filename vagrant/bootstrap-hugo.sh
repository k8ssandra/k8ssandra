#!/usr/bin/env bash

echo "Bootstrapping hugo..."

## For Ubuntu 18.04 (bionic) apt-get will install the non-extended version of hugu
## we require the extended version, so it must be built from source here

echo "mkdir ~/tmp/hugo"
mkdir ~/tmp/hugo

echo "cd ~/tmp/hugo"
cd ~/tmp/hugo/

echo "wget https://github.com/gohugoio/hugo/releases/download/v0.80.0/hugo_extended_0.80.0_Linux-64bit.tar.gz"
wget https://github.com/gohugoio/hugo/releases/download/v0.80.0/hugo_extended_0.80.0_Linux-64bit.tar.gz

echo "tar -xf hugo_extended_0.80.0_Linux-64bit.tar.gz"
tar -xf hugo_extended_0.80.0_Linux-64bit.tar.gz

echo "mv hugo /usr/local/bin/hugo"
mv hugo /usr/local/bin/hugo

echo "rm -rf ~/tmp/hugo/"
rm -rf ~/tmp/hugo/

echo "Completed bootstrapping hugo"
