SHELL := /bin/bash

# Current Operator version
VERSION ?= 0.0.1
# Default bundle image tag
BUNDLE_IMG ?= controller-bundle:$(VERSION)
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Image URL to use all building/pushing image targets
IMG ?= stakater/slack-operator:v$(VERSION)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,crdVersions=v1,preserveUnknownFields=false"

# GOLANGCI_LINT env
GOLANGCI_LINT = _output/tools/golangci-lint
GOLANGCI_LINT_CACHE = $(PWD)/_output/golangci-lint-cache
GOLANGCI_LINT_VERSION = v1.24

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager


ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet ## Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.7.2/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out


# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests install
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build -t ${IMG} .

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v4@v4.5.4)

# go-install-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: manifests kustomize ## Generate bundle manifests and metadata, then validate generated files.
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

# Options for "packagemanifests".
ifneq ($(origin FROM_VERSION), undefined)
PKG_FROM_VERSION := --from-version=$(FROM_VERSION)
endif
ifneq ($(origin CHANNEL), undefined)
PKG_CHANNELS := --channel=$(CHANNEL)
endif
ifeq ($(IS_CHANNEL_DEFAULT), 1)
PKG_IS_DEFAULT_CHANNEL := --default-channel
endif
PKG_MAN_OPTS ?= $(FROM_VERSION) $(PKG_CHANNELS) $(PKG_IS_DEFAULT_CHANNEL)

# Generate package manifests.
packagemanifests: manifests
	operator-sdk generate kustomize manifests -q
	$(KUSTOMIZE) build config/manifests | operator-sdk generate packagemanifests -q --version $(VERSION) $(PKG_MAN_OPTS)

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

verify-fmt:
	./hack/verify-gofmt.sh

$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(dir $@) v1.36.0

verify-golangci-lint: $(GOLANGCI_LINT)
	GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) $(GOLANGCI_LINT) run --timeout=300s ./...

verify: verify-fmt verify-golangci-lint

bump-chart-operator:
	sed -i "s/^version:.*/version:  $(VERSION)/" charts/slack-operator/Chart.yaml
	sed -i "s/^appVersion:.*/appVersion:  $(VERSION)/" charts/slack-operator/Chart.yaml
	sed -i "s/tag:.*/tag:  v$(VERSION)/" charts/slack-operator/values.yaml

# Bump Chart
bump-chart: bump-chart-operator 

generate-crds: controller-gen
	$(CONTROLLER_GEN) crd paths="./..." output:crd:artifacts:config=charts/slack-operator/crds