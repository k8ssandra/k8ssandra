#!/usr/bin/env bash
# This is a utility script for opening the Grafana dashboard for a k8ssandra cluster. It:
#  1. Ensures that Grafana is ready (note that there is no timeout; if you've made a mistake, this script will wait forever.)
#  2. Retrieves the username and password from the Kubernetes secret. (Make sure your kubectl is configured with sufficient access.)
#  3. Forwards a local port to the Grafana service in Kubernetes.
#  4. Launches your web browser (on Mac) or just tells you where to go.

if [[ "$1" == "-nc" ]]; then
  NOCOLOR=''
  RED=''
  BOLDRED=''
  GREEN=''
  BOLDGREEN=''
  BOLDBLUE=''
  CYAN=''
  BOLDWHITE=''
  shift
else
  NOCOLOR='\033[0m'
  RED='\033[0;31m'
  BOLDRED='\033[1;31m'
  GREEN='\033[0;32m'
  BOLDGREEN='\033[1;32m'
  BOLDBLUE='\033[1;34m'
  CYAN='\033[0;36m'
  BOLDWHITE='\033[1;37m'
fi

NAMESPACE=default
if [[ -n $1 ]]; then
  NAMESPACE=$1
  shift
fi
echo -e "${BOLDCYAN}Using namespace ${BOLDWHITE}${NAMESPACE}.${NOCOLOR}"

echo -e "\n${BOLDBLUE}Waiting for Grafana pod to be ready...${NOCOLOR}"
until kubectl wait --for=condition=ready -n ${NAMESPACE} pod -l app=grafana > /dev/null 2>&1; do sleep 1; echo -ne "${BOLDBLUE}.${NOCOLOR}"; done

GRAFANA_URL="http://127.0.0.1:3000/"
GRAFANA_CREDS_JSON=$(kubectl get secret -n ${NAMESPACE} grafana-admin-credentials -o json)
GRAFANA_USER=$(jq -r ".data.\"GF_SECURITY_ADMIN_USER\"" <<< "${GRAFANA_CREDS_JSON}" | base64 -d)
GRAFANA_PASSWORD=$(jq -r ".data.\"GF_SECURITY_ADMIN_PASSWORD\"" <<< "${GRAFANA_CREDS_JSON}" | base64 -d)
echo -e "\n  ${BOLDWHITE}Grafana URL: ${NOCOLOR} ${GRAFANA_URL}"
echo -e "  ${BOLDWHITE}Grafana User:${NOCOLOR} ${GRAFANA_USER}"
echo -e "  ${BOLDWHITE}Grafana Password:${NOCOLOR} ${GRAFANA_PASSWORD}"

kill $(cat grafana-port-forward.pid) > /dev/null 2>&1

echo -e "\n${BOLDBLUE}Forwarding local port 9080 to grafana...${NOCOLOR}"
kubectl port-forward -n ${NAMESPACE} service/grafana-service 3000 &
PORT_FORWARD_PID=$!
echo "${PORT_FORWARD_PID}" > grafana-port-forward.pid
sleep 0.5

if kill -0 ${PORT_FORWARD_PID} > /dev/null 2>&1; then

  echo -e "\n${BOLDBLUE}Waiting for Grafana service to be ready on port 3000...${NOCOLOR}"
  until nc -zv localhost 3000 > /dev/null 2>&1; do sleep 1; echo -ne "${BOLDBLUE}.${NOCOLOR}"; done

  if command -v open > /dev/null ; then
    echo -e "\n${BOLDBLUE}Launching in browser.${NOCOLOR}"
    open "${GRAFANA_URL}"
  else
    echo -e "\n${BOLDBLUE}Grafana is ready. Navigate your browser to: ${BOLDWHITE}${GRAFANA_URL}${NOCOLOR}"
  fi

  echo "When you're done, release the port using: kill ${PORT_FORWARD_PID}"
else
  echo -e "${BOLDRED}Failed to create proxy.${NOCOLOR}"
fi
