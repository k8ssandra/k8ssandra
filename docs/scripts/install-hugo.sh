#!/bin/bash

if [ ! -d "bin" ]
then
  mkdir bin
fi

if [ ! -f "bin/hugo" ]
then
  curl -L "https://github.com/gohugoio/hugo/releases/download/v0.76.5/hugo_extended_0.76.5_Linux-64bit.tar.gz" | tar -xvzpf - -C bin hugo
fi

bin/hugo version
