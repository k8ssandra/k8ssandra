#!/usr/bin/env bash

set -x

# Generate docs for each chart
ls charts | while read c; do
  if [[ -d charts/$c ]]; then
    helm-docs -c charts/$c -s file

    mkdir -p docs/content/en/docs/reference/$c
    helm-docs -c charts/$c -s file -t ../../docs/content/en/docs/reference/_generated.md.gotmpl -o ../../docs/content/en/docs/reference/$c/_generated.md
  fi
done
