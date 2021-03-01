SHELL := /bin/bash

ORG?=k8ssandra
REG?=docker.io

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
REV=$(shell git rev-parse --short=12 HEAD)

CLEANER_IMAGE_BASE=$(REG)/$(ORG)/k8ssandra-cleaner
CLEANER_REV_IMAGE=$(CLEANER_IMAGE_BASE):$(REV)
CLEANER_LATEST_IMAGE=$(CLEANER_IMAGE_BASE):latest

# Image URL to use all building/pushing image targets
CLEANER_IMG ?= $(CLEANER_LATEST_IMAGE)

TESTS=all
GO_FLAGS=
ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: fmt vet unit-test pkg-test

unit-test:
ifeq ($(TESTS), all)
	go test $(GO_FLAGS) -test.timeout=3m ./tests/unit/... -coverprofile cover.out
else
	go test $(GO_FLAGS) -test.timeout=3m ./tests/unit/... -coverprofile cover.out -args -ginkgo.focus="$(TESTS)"
endif

pkg-test:
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/master/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh && fetch_envtest_tools $(ENVTEST_ASSETS_DIR) && setup_envtest_env $(ENVTEST_ASSETS_DIR) && go test $(GO_FLAGS) -test.timeout=3m ./pkg/... -coverprofile cover.out

integ-test:
	go test $(GO_FLAGS) -test.timeout=5m ./tests/integration/... -coverprofile cover.out

fmt:
	go fmt ./pkg/...
	go fmt ./tests/...

vet:
	go vet ./pkg/...
	go vet ./tests/...

cleaner-docker-build:
	@echo Building ${CLEANER_REV_IMAGE}
	docker build -t ${CLEANER_REV_IMAGE} -f cmd/k8ssandra-client/Dockerfile .
	docker tag ${CLEANER_REV_IMAGE} ${CLEANER_LATEST_IMAGE}

cleaner-docker-push:
	docker push ${CLEANER_REV_IMAGE}
    docker push ${CLEANER_LATEST_IMAGE}
