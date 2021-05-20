SHELL := /bin/bash

ORG?=k8ssandra
REG?=docker.io

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
REV=$(shell git rev-parse --short=12 HEAD)

TOOLS_IMAGE_BASE=$(REG)/$(ORG)/k8ssandra-tools
TOOLS_REV_IMAGE=$(TOOLS_IMAGE_BASE):$(REV)
TOOLS_LATEST_IMAGE=$(TOOLS_IMAGE_BASE):latest

# Image URL to use all building/pushing image targets
TOOLS_IMG ?= $(TOOLS_LATEST_IMAGE)

TESTS=all
GO_FLAGS=
ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
# Possible values: always, success, never
CLUSTER_CLEANUP=always
KIND_CLUSTER=k8ssandra-it

test: fmt vet unit-test pkg-test

unit-test:
ifeq ($(TESTS), all)
	go test $(GO_FLAGS) -test.timeout=3m ./tests/unit/... -coverprofile cover.out
else
	go test $(GO_FLAGS) -test.timeout=3m ./tests/unit/... -coverprofile cover.out -args -ginkgo.focus="$(TESTS)"
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

kind-integ-test: create-kind-cluster integ-test

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
	@echo Building ${TOOLS_REV_IMAGE}
	docker build -t ${TOOLS_REV_IMAGE} -f cmd/k8ssandra-client/Dockerfile .
	docker tag ${TOOLS_REV_IMAGE} ${TOOLS_LATEST_IMAGE}

tools-docker-push:
	docker push ${TOOLS_REV_IMAGE}
    docker push ${TOOLS_LATEST_IMAGE}
