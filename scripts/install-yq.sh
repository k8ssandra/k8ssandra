#!/bin/bash

VERSION=3.4.1
if [[ -n $1 ]]; then
  VERSION=$1
  shift
fi

if [ ! -d "bin" ]
then
  mkdir bin
fi

if [ ! -f "bin/yq" ]
then
  curl -L https://github.com/mikefarah/yq/releases/download/$VERSION/yq_linux_amd64 > bin/yq
  chmod +x bin/yq
fi

bin/yq --version
