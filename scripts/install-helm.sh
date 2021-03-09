#!/bin/bash

VERSION=3.4.0
if [[ -n $1 ]]; then
  VERSION=$1
  shift
fi

if [ ! -d "bin" ]
then
  mkdir bin
fi

if [ ! -f "bin/helm" ]
then
  curl -L https://get.helm.sh/helm-v$VERSION-linux-amd64.tar.gz | tar -xvzf - --strip-components 1 --directory bin/ linux-amd64/helm
  chmod +x bin/helm
fi

bin/helm version
