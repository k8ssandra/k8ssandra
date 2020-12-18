#!/usr/bin/env bash
# This is a utility script for k8ssandra developers. We need to embed Helm-template-compatible versions of the MCAC Grafana dashboards into our chart.
# Due to a syntax conflict (both Helm and Grafana use the double-curly-brace notation), we can't just use the MCAC versions as-is.
# This script takes one or more unmodified MCAC dashboard YAMLs as arguments. For each one, it:
#  1. Extracts and unescapes the JSON from the GrafanaDashboard YAML.
#  2. Uses SED to apply a regex replacement
#  3. Re-escapes the JSON
#  4. Outputs a new YAML file containing the updated JSON.
#
# The output file will match the name of the input file, replacing the file extension (presumably ".yaml") with "-helm-template.yaml".
# For example, overview.dashboard.yaml becomes overview.dashboard-helm-template.yaml.
#
# This script depends on the presence of both `jq` and `yq`. If they're missing, you'll get helpful output.
# If you're a Mac user with Homebrew installed, just run `brew install jq python-yq`.

if [[ "$1" == "-nc" ]]; then
  NOCOLOR=''
  RED=''
  GREEN=''
  CYAN=''
  shift
else
  NOCOLOR='\033[0m'
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  CYAN='\033[0;36m'
fi

if ! command -v yq > /dev/null ; then
  echo -e "${RED}The 'yq' utility is required. Please install: https://kislyuk.github.io/yq/${NOCOLOR}"
  exit 1
fi

if ! command -v jq > /dev/null ; then
  echo -e "${RED}The 'jq' utility is required. Please install: https://stedolan.github.io/jq/download/${NOCOLOR}"
  exit 1
fi

# take one or more YAML files as input
if [[ $# -lt 1 ]]; then
  echo -e "${RED}Expected at least one argument (one or more YAML files defining GrafanaDashboards)${NOCOLOR}"
  exit 2
fi

while [[ -n "$1" ]]; do
  INFILE="$1"
  INFILE_NAME="$(basename ${INFILE})"
  INFILE_PREFIX="${INFILE_NAME%.*}"
  OUTFILE="${INFILE_PREFIX}-helm-template.yaml"

  echo -ne "${CYAN}${INFILE_NAME}${NOCOLOR}"

  # extract unescaped JSON from YAML file
  DASHBOARD_JSON="$(yq -r '.spec.json' ${INFILE})"

  echo -ne "${CYAN}.${NOCOLOR}"

  # run file through sed to escape the "{{xxxx}}" into "{{`{{xxxx}}`}}"
  DASHBOARD_TEMPLATE_JSON="$(sed -E 's/"([^"]*)({{[^"]+}})([^"]*)"/"{{`\1\2\3`}}"/g' <<< "${DASHBOARD_JSON}")"

  echo -ne "${CYAN}.${NOCOLOR}"

  # re-escape the JSON
  ESCAPED_JSON="$(jq '.|@json' <<< "${DASHBOARD_TEMPLATE_JSON}")"

  echo -ne "${CYAN}.${NOCOLOR}"

  # generate new YAML with updated JSON and write to the outfile
  yq -Y ".spec.json|=${ESCAPED_JSON}" "${INFILE}" > "${OUTFILE}"

  echo -e " -> ${GREEN}${OUTFILE}${NOCOLOR}"
  shift
done
