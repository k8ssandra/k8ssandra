SHELL := /bin/bash

ORG?=k8ssandra

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
REV=$(shell git rev-parse --short=12 HEAD)

TOOLS_IMAGE_BASE=$(ORG)/k8ssandra-tools
TOOLS_REV_IMAGE=$(TOOLS_IMAGE_BASE):$(REV)
TOOLS_LATEST_IMAGE=$(TOOLS_IMAGE_BASE):latest

# Image URL to use all building/pushing image targets
TOOLS_IMG ?= $(TOOLS_LATEST_IMAGE)

# Buildx params
BUILDX_PARAMS=--load

TESTS=all
GO_FLAGS=
ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
# Possible values: always, success, never
CLUSTER_CLEANUP=always
KIND_CLUSTER=k8ssandra-it

test: fmt vet unit-test pkg-test

unit-test:
ifeq ($(TESTS), all)
	go test $(GO_FLAGS) -test.timeout=5m ./tests/unit/... -coverprofile cover.out
else
	go test $(GO_FLAGS) -test.timeout=5m ./tests/unit/... -coverprofile cover.out -args -ginkgo.focus="$(TESTS)"
endif

pkg-test:
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh && fetch_envtest_tools $(ENVTEST_ASSETS_DIR) && setup_envtest_env $(ENVTEST_ASSETS_DIR) && go test $(GO_FLAGS) -test.timeout=3m ./pkg/... -coverprofile cover.out

integ-test:
ifeq ($(TESTS), all)
	CLUSTER_CLEANUP=$(CLUSTER_CLEANUP) go test $(GO_FLAGS) -test.timeout=30m ./tests/integration -run="TestFullStackScenario"
else
	CLUSTER_CLEANUP=$(CLUSTER_CLEANUP) go test $(GO_FLAGS) -test.timeout=30m ./tests/integration -run=$(TESTS)
endif

kind-integ-test: create-kind-cluster tools-docker-kind-load integ-test

create-kind-cluster:
	kind delete cluster --name $(KIND_CLUSTER)
	./tests/integration/scripts/create_kind_cluster.sh $(KIND_CLUSTER)

fmt:
	go fmt ./pkg/...
	go fmt ./tests/...

vet:
	go vet ./pkg/...
	go vet ./tests/...

tools-docker-build:
	@echo Building test version of ${TOOLS_REV_IMAGE}
	set -e ;\
	VER=$$(yq eval '.version' charts/k8ssandra/Chart.yaml) ;\
	mkdir -p build/$$VER ;\
	cp -rv charts/* build/$$VER/ ;\
	docker buildx build $(BUILDX_PARAMS) -t ${TOOLS_IMG} -f cmd/k8ssandra-client/Dockerfile . ;\
	rm -fr build/$$VER ;\

tools-docker-kind-load: tools-docker-build
	@echo Loading tools to kind
	kind load docker-image ${TOOLS_IMG} --name $(KIND_CLUSTER)

