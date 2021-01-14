SHELL := /bin/bash

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: fmt vet
	go test -v -test.timeout=3m ./tests/unit/... -coverprofile cover.out
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/master/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh && fetch_envtest_tools $(ENVTEST_ASSETS_DIR) && setup_envtest_env $(ENVTEST_ASSETS_DIR) && go test -v -test.timeout=3m ./pkg/... -coverprofile cover.out

integ-test:
	go test -v -test.timeout=5m ./tests/integration/... -coverprofile cover.out

fmt:
	go fmt ./pkg/...
	go fmt ./tests/...

vet:
	go vet ./pkg/...
	go vet ./tests/...

docker-build:
	docker build -f cmd/k8ssandra-client/Dockerfile .
