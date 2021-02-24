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
  curl -L https://github.com/mikefarah/yq/releases/download/$YQ_VERSION/yq_linux_amd64 > $HOME/bin/yq
  chmod +x $HOME/bin/yq
fi

bin/yq --version
