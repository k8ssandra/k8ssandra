#!/usr/bin/env bash

if ! command -v helm-docs &> /dev/null
then
    echo "helm-docs could not be found. Please ensure it is installed and on your PATH."
    echo "https://github.com/norwoodj/helm-docs"
    exit 1
fi

cd "$(dirname "$0")/.."

# Generate docs for each chart
for directory in charts/*; do
  if [[ -d "$directory" ]]; then
    chartName="$(basename ${directory})"
    mkdir -p "docs/content/en/docs/reference/${chartName}"
    (
      set -x
      helm-docs -c "${directory}" -s file
      helm-docs -c "${directory}" -s file -t ../../docs/content/en/docs/reference/_generated.md.gotmpl -o "../../docs/content/en/docs/reference/${chartName}/_generated.md"
    )
  fi
done
