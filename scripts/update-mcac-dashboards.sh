#!/usr/bin/env bash
# This is a utility script for k8ssandra developers. We need to downloads JSON versions of the MCAC Grafana dashboards into our chart's dashboards directory.
# This script:
#  1. Grabs the latest copies of each dashboard from MCAC (optionally from a specified tag or branch)
#

set -e
if [[ "$1" == "-nc" ]]; then
  NOCOLOR=''
  BOLDGREEN=''
  BOLDCYAN=''
  BOLDWHITE=''
  shift
else
  NOCOLOR='\033[0m'
  BOLDGREEN='\033[1;32m'
  BOLDCYAN='\033[1;36m'
  BOLDWHITE='\033[1;37m'
fi

MCAC_VERSION=master
if [[ -n "$1" ]]; then
  MCAC_VERSION="$1"
fi
echo -e "${BOLDCYAN}Using dashboard definitions from MCAC ${BOLDWHITE}${MCAC_VERSION}.${NOCOLOR}"

cd "$(dirname "$0")"

echo -e "${BOLDCYAN}Retrieving latest copies of the dashboards...${NOCOLOR}"
wget -q --show-progress https://raw.githubusercontent.com/datastax/metric-collector-for-apache-cassandra/${MCAC_VERSION}/dashboards/grafana/generated-dashboards/system-metrics.json -O ../charts/k8ssandra/dashboards/system-metrics.json
wget -q --show-progress https://raw.githubusercontent.com/datastax/metric-collector-for-apache-cassandra/${MCAC_VERSION}/dashboards/grafana/generated-dashboards/cassandra-condensed.json -O ../charts/k8ssandra/dashboards/cassandra-condensed.json
wget -q --show-progress https://raw.githubusercontent.com/datastax/metric-collector-for-apache-cassandra/${MCAC_VERSION}/dashboards/grafana/generated-dashboards/overview.json -O ../charts/k8ssandra/dashboards/overview.json

echo -e "${BOLDGREEN}Done!${NOCOLOR}"
