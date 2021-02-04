#!/bin/bash

VERSION=0.76.5
if [[ -n $1 ]]; then
  VERSION=$1
  shift
fi

if [ ! -d "bin" ]
then
  mkdir bin
fi

if [ ! -f "bin/hugo" ]
then
  curl -L "https://github.com/gohugoio/hugo/releases/download/v${VERSION}/hugo_extended_${VERSION}_Linux-64bit.tar.gz" | tar -xvzpf - -C bin hugo
fi

bin/hugo version
