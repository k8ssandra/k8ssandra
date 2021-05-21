#!/bin/bash

# Download latest Fossa CLI distribution

curl -H 'Cache-Control: no-cache' https://raw.githubusercontent.com/fossas/fossa-cli/master/install.sh | bash

## Init the project (if necessary)

fossa init

## Analyze the project
## This requires an environment varaible called FOSSA_API_KEY

fossa analyze