#!/usr/bin/env bash

if ! command -v helm &> /dev/null
then
    echo "helm could not be found. Please ensure it is installed and on your PATH."
    echo "https://github.com/helm/helm/"
    exit 1
fi

cd "$(dirname "$0")/.."

# Update dependencies for each chart, order is important!
CHARTS=("k8ssandra-common" "backup" "cass-operator" "medusa-operator" "reaper-operator" "restore" "k8ssandra")
for CHART in ${CHARTS[@]}; do
  if [[ -d "charts/$CHART" ]]; then
    helm dep update charts/$CHART
  else
    echo "Error fetching dependency for $CHART, directory not found"
    exit 1
  fi
done
