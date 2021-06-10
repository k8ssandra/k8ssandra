
# Developer Scripts
The scripts described are for developers to customize and use with `k8ssandra`.

## Install/Setup

### [setup-cluster.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/setup-cluster.sh)

Performs Helm installation of the `k8ssandra` cluster, using default or customized values for: `size`, `clusterName`,
& `namespace`. The namespace is created if not already existing.

Usage: setup-cluster [_k8ssandra-namespace=k8ssandra_] <name> [_size=1_]

### [setup-k8ssandra.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/setup-k8ssandra.sh)

Performs Helm installation of the `k8ssandra` cluster, using default or customized values as specified in
the `k8ssandra` chart Values file. The namespace is created if not already existing.

## Update

### [update-cluster.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/update-cluster.sh)

Performs Helm upgrade of release `{NAME}-k8ssandra` for the `k8ssandra` chart.

The standard `storageClass` for `cassandra.cassandraLibDirVolume` can be customized
using a `STORAGE_CLASS` environment variable.  Other values can also be defined
in a `sample-values.yaml` file. The namespace is created if not already existing.

### [update-mcac-dashboards.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/update-mcac-dashboards.sh)

Performs a download of `MCAC Grafana` dashboards into the k8ssandra `dashboards` directory.  
The latest copies (master by default, unless a branch or tag is specified) of each dashboard from MCAC.


## Delete

### [delete-crds.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/delete-crds.sh)

Performs cleanup of `k8ssandra` custom resource definitions (CRDs).   Includes cleanup for `Traefik`.  Easily
customizable to target more or fewer CRDs as needed.

## Document Management

### [install-helm-docs.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/install-helm-docs.sh)
Performs installation of `Helm` documentation from the helm-docs [releases](https://github.com/norwoodj/helm-docs/releases) GitHub.

### [install-hugo.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/install-hugo.sh)
Performs installation of `Hugo` (a static site generator) from the [releases](https://github.com/gohugoio/hugo/releases/) github site.

### [install-npm-dependencies.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/install-npm-dependencies.sh)
Performs a directory change to the `docs` directory and issues an `npm install` command.

## Misc.

### [open-grafana.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/open-grafana.sh)

Performs launching of the k8ssandra cluster's Grafana dashboard.  Activities performed by this script include:
* Ensures Grafana is active.
* Retrieves username and password from the Kubernetes secret.
* Forwards a local port to the Grafana service in Kubernetes.
* Launches your web browser (on Mac) or provides navigation information with other platforms.

### [helm-version.py](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/helm-version.py)

Performs update of the Helm chart versions in `k8ssandra` project.

### [generate-helm-docs.sh](https://github.com/k8ssandra/k8ssandra/tree/main/scripts/generate-helm-docs.sh)
Generates Helm docs folder for each chart (e.g. docs/content/en/docs/reference).


