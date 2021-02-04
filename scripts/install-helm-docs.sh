#!/bin/bash

VERSION=1.5.0
if [[ -n $1 ]]; then
  VERSION=$1
  shift
fi

if [ ! -d "bin" ]
then
  mkdir bin
fi

if [ ! -f "bin/helm-docs" ]
then
  curl -L "https://github.com/norwoodj/helm-docs/releases/download/v${VERSION}/helm-docs_${VERSION}_Linux_x86_64.tar.gz" | tar -xvzpf - -C bin helm-docs
  chmod +x bin/helm-docs
fi

bin/helm-docs --version
