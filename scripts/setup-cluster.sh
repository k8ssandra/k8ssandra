#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

if [[ "$#" -gt 3 ]]; then
  echo "Usage: $0 [k8ssandra-namespace=k8ssandra] <name> [size=1]"
  exit 1
elif [[ "$#" -eq 3 ]]; then
  # all three args given
  NAMESPACE=$1
  CLUSTERNAME=$2
  SIZE=$3
elif [[ "$#" -eq 2 ]]; then
  # two args given -- is it namespace+name or name+size?
  case "$2" in
    ''[0-9]*)
      # second arg is numeric; must be name+size
      NAMESPACE=k8ssandra
      CLUSTERNAME=$1
      SIZE=$2
      ;;
    *)
      NAMESPACE=$1
      CLUSTERNAME=$2
      SIZE=1
      ;;
  esac
elif [[ "$#" -eq 1 ]]; then
  # one arg given
  NAMESPACE=k8ssandra
  CLUSTERNAME=$1
  SIZE=1
else
  # no args given
  echo "Usage: $0 [k8ssandra-namespace=k8ssandra] <name> [size=1]"
  exit 1
fi

set -x
helm install ${CLUSTERNAME} ../charts/k8ssandra-cluster --set size=${SIZE} --set clusterName=${CLUSTERNAME} --set k8ssandra.namespace=${NAMESPACE} -n ${CLUSTERNAME} --create-namespace


