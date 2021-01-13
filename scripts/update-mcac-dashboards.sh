#!/usr/bin/env bash
# This is a utility script for k8ssandra developers. We need to embed Helm-template-compatible versions of the MCAC Grafana dashboards into our chart.
# Due to a syntax conflict (both Helm and Grafana use the double-curly-brace notation), we can't just use the MCAC versions as-is.
# This script:
#  1. grabs the latest copies of each dashboard from MCAC (optionally from a specified tag or branch)
#  2. updates the "Overview" dashboard to include a plugin specifier for grafana-polystat-panel v1.2.2
#  3. passes each dashboard through another script which adds the necessary escaping
#  4. copies the results into the k8ssandra chart's templates directory.
#  5. cleans up after itself
#
# This script depends on the presence of both `jq` and `yq`. If they're missing, you'll get helpful output.
# If you're a Mac user with Homebrew installed, just run `brew install jq python-yq`.


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

if ! command -v yq > /dev/null ; then
  echo -e "${RED}The 'yq' utility is required. Please install: https://kislyuk.github.io/yq/${NOCOLOR}"
  exit 1
fi

if ! command -v jq > /dev/null ; then
  echo -e "${RED}The 'jq' utility is required. Please install: https://stedolan.github.io/jq/download/${NOCOLOR}"
  exit 1
fi

MCAC_VERSION=master
if [[ -n "$1" ]]; then
  MCAC_VERSION="$1"
fi
echo -e "${BOLDCYAN}Using dashboard definitions from MCAC ${BOLDWHITE}${MCAC_VERSION}.${NOCOLOR}"

cd "$(dirname "$0")"

mkdir temp
cd temp

echo -e "${BOLDCYAN}Retrieving latest copies of the dashboards...${NOCOLOR}"
wget -q --show-progress https://raw.githubusercontent.com/datastax/metric-collector-for-apache-cassandra/${MCAC_VERSION}/dashboards/k8s-build/generated/grafana/system-metrics.dashboard.yaml
wget -q --show-progress https://raw.githubusercontent.com/datastax/metric-collector-for-apache-cassandra/${MCAC_VERSION}/dashboards/k8s-build/generated/grafana/cassandra-condensed.dashboard.yaml
wget -q --show-progress https://raw.githubusercontent.com/datastax/metric-collector-for-apache-cassandra/${MCAC_VERSION}/dashboards/k8s-build/generated/grafana/overview.dashboard.yaml

echo -e "${BOLDCYAN}Adding plugin dependencies...${NOCOLOR}"
yq -Y '.spec.plugins=[{"name":"grafana-polystat-panel", "version":"1.2.2"}]' overview.dashboard.yaml > overview-with-plugin.dashboard.yaml
rm -vf overview.dashboard.yaml

echo -e "${BOLDCYAN}Templatizing...${NOCOLOR}"
../templatize-dashboard.sh *.dashboard.yaml

echo -e "${BOLDCYAN}Copying templatized dashboards to chart...${NOCOLOR}"
cp -v *.dashboard-helm-template.yaml ../../../k8ssandra/charts/k8ssandra/templates/grafana/dashboards/

echo -e "${BOLDCYAN}Cleaning up...${NOCOLOR}"
cd ..
rm -rf temp

echo -e "${BOLDGREEN}Done!${NOCOLOR}"
