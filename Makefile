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
ENVTEST=$(shell pwd)/testbin/setup-envtest
ENVTEST_K8S_VERSION=1.21
# Possible values: always, success, never
CLUSTER_CLEANUP=always
KIND_CLUSTER=k8ssandra-it

test: fmt vet unit-test pkg-test

unit-test:
ifeq ($(TESTS), all)
	go test $(GO_FLAGS) -test.timeout=10m ./tests/unit/... -coverprofile cover.out
else
	go test $(GO_FLAGS) -test.timeout=10m ./tests/unit/... -coverprofile cover.out -args -ginkgo.focus="$(TESTS)"
endif

.PHONY: envtest
envtest: ## Download envtest-setup locally if necessary.
	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

.PHONY: pkg-test
pkg-test: envtest
	export KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --use-env -p path )" && go test $(GO_FLAGS) -test.timeout=3m ./pkg/... -coverprofile cover.out

integ-test:
ifeq ($(TESTS), all)
	CLUSTER_CLEANUP=$(CLUSTER_CLEANUP) go test $(GO_FLAGS) -test.timeout=30m ./tests/integration -run="TestFullStackScenario"
else
	CLUSTER_CLEANUP=$(CLUSTER_CLEANUP) go test $(GO_FLAGS) -test.timeout=30m ./tests/integration -run=$(TESTS)
endif

.PHONY: cert-manager
cert-manager: ## Install cert-manager to the cluster
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.7.1/cert-manager.yaml
	kubectl wait --for=condition=Established crd certificates.cert-manager.io
	kubectl rollout status deployment cert-manager-webhook -n cert-manager

kind-integ-test: create-kind-cluster cert-manager tools-docker-kind-load integ-test 

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


# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/testbin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef