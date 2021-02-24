#!/usr/bin/env bash

if ! command -v helm &> /dev/null
then
    echo "helm could not be found. Please ensure it is installed and on your PATH."
    echo "https://github.com/helm/helm/"
    exit 1
fi

cd "$(dirname "$0")/.."

# Generate docs for each chart
for directory in charts/*; do
  if [[ -d "$directory" ]]; then
    helm dep update $directory
  fi
done
