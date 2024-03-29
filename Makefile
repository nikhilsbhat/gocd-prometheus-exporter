GOFMT_FILES?=$(shell find . -not -path "./vendor/*" -type f -name '*.go')
MY_ADDRESS := $(shell ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $$2}' )
BUILD_ENVIRONMENT?=${ENVIRONMENT}
GOVERSION?=$(shell go version | awk '{printf $$3}')
APP_DIR?=$(shell git rev-parse --show-toplevel)
HELM_CHART_REPO=oci://ghcr.io/nikhilsbhat/charts

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

render.prometheus.config: ## Renders prometheus config file that could be used for running prometheus as a container
	export MY_IP_ADDRESS=${MY_ADDRESS} && envsubst <infrastructure/prometheus.yml.tmpl>infrastructure/prometheus.yml

test.setup.up: render.prometheus.config ## Brings up prometheus and GoCd setup by running container with the config specified
	@docker-compose -f infrastructure/docker-compose.infra.yaml up -d

test.setup.purge: ## Teardown prometheus and GoCd setup
	@docker-compose -f infrastructure/docker-compose.infra.yaml down

local.fmt: ## Lints all the go code in the application.
	@gofmt -w $(GOFMT_FILES)
	$(GOBIN)/gofumpt -l -w $(GOFMT_FILES)
	$(GOBIN)/goimports -w $(GOFMT_FILES)
	$(GOBIN)/gci write $(GOFMT_FILES) --skip-generated

local.check: local.fmt ## Loads all the dependencies to vendor directory
	@go mod vendor
	@go mod tidy

local.build: local.check ## Generates the artifact with the help of 'go build'
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} goreleaser build --rm-dist

publish: local.check ## Builds and publishes the app
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} PLUGIN_PATH=${APP_DIR} goreleaser release --rm-dist

mock.publish: local.check ## Builds and mocks app release
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} PLUGIN_PATH=${APP_DIR} goreleaser release --skip-publish --rm-dist

lint: ## Lint's application for errors, it is a linters aggregator (https://github.com/golangci/golangci-lint).
	if [ -z "${DEV}" ]; then golangci-lint run --color always ; else docker run --rm -v $(APP_DIR):/app -w /app golangci/golangci-lint:v1.46.2-alpine golangci-lint run --color always ; fi

report: ## Publishes the go-report of the appliction (uses go-reportcard)
	if [ -z "${DEV}" ]; then goreportcard -v ; else docker run --rm -v $(APP_DIR):/app -w /app basnik/goreportcard-cli:latest goreportcard-cli -v ; fi

generate.document: ## generates cli documents using 'github.com/nikhilsbhat/urfavecli-docgen'.
	@go generate github.com/nikhilsbhat/gocd-prometheus-exporter/docs
	@helm-docs -c charts/gocd-prometheus-exporter

test: ## Runs test cases
	@go test ./... -mod=vendor -coverprofile cover.out && go tool cover -html=cover.out -o cover.html && open cover.html

docker.login: ## Should login to ghcr docker registry.
	@echo "${GITHUB_TOKEN}" | docker login ghcr.io -u nikshilsbhat --password-stdin

helm.package: ## Packages the helm chart to make it ready for publishing.
	@rm -rf gocd-prometheus-exporter-*
	@helm package charts/gocd-prometheus-exporter/ -d charts/

helm.registry.login: ## Logins to ghcr oci helm registry.
	echo "${GITHUB_TOKEN}" | helm registry login ghcr.io/nikshilsbhat --username nikhilsbhat --password-stdin

helm.publish: helm.package helm.registry.login ## Should publish helm chart to oci based ghcr helm registries.
	@helm push charts/*.tgz ${HELM_CHART_REPO}
