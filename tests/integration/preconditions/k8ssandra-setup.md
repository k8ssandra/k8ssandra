
## Integration Test Setup

### Preconditions
Ensure the following are installed in the integration test environment.
- [Helm v3+](https://helm.sh/docs/intro/install/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kind](https://kind.sigs.k8s.io/)

### Validate cluster connectivity
`kubectl cluster-info`

If cluster is not existing, install k8s cluster using Kind. 

### Create Cluster
`kind create cluster --name k8ssandra-cluster --image kindest/node:v1.18.2 --config ./k8ssandra-kind-config.yaml`

`kind get clusters`

### Validate
Invoke the cluster-info again with the expected context used in the integration tests.

`kubectl cluster-info --context kind-k8ssandra-cluster`

### Troubleshoot
To debug and diagnose any cluster problems, use:

`kubectl cluster-info dump`

Note: This will produces LOTS of details so piping output to file is recommended.

### Running Tests
It is recommended that adequate time is given for the integration tests to fully complete.  Subject to change, but recommended minimum of 5 minutes.

Example:

`go test -v -timeout=300s`

