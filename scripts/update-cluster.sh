#!/usr/bin/env bash
set -e

if [[ "$1" == "-nc" ]]; then
  NOCOLOR=''
  BOLDBLUE=''
  shift
else
  NOCOLOR='\033[0m'
  BOLDBLUE='\033[1;34m'
fi

cd "$(dirname "$0")"
echo -e "\n${BOLDBLUE}Setting up K8ssandra cluster...${NOCOLOR}"

if [[ "$#" -gt 1 ]]; then
  echo "Usage: $0 [k8ssandra-namespace=k8ssandra]"
  exit 1
elif [[ "$#" -eq 1 ]]; then
  NAME=$1
else
  NAME=k8ssandra
fi

set -x
helm upgrade ${NAME}-k8ssandra ../charts/k8ssandra  --set cassandra.cassandraLibDirVolume.storageClass=${STORAGE_CLASS:-standard} --set k8ssandra.namespace=${NAME} -n ${NAME} --create-namespace -f sample-values.yaml


