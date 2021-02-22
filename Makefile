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

cleaner-docker-build:
	@echo Building ${CLEANER_REV_IMAGE}
	docker build -t ${CLEANER_REV_IMAGE} -f cmd/k8ssandra-client/Dockerfile .
	docker tag ${CLEANER_REV_IMAGE} ${CLEANER_LATEST_IMAGE}

cleaner-docker-push:
	docker push ${CLEANER_REV_IMAGE}
    docker push ${CLEANER_LATEST_IMAGE}
