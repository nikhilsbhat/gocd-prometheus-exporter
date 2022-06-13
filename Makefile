MY_ADDRESS := $(shell ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $$2}' )

.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

render.prometheus.config: ## Renders prometheus config file that could be used for running prometheus as a container
	export MY_IP_ADDRESS=${MY_ADDRESS} && envsubst <infrastructure/prometheus.yml.tmpl>infrastructure/prometheus.yml

test.setup.up: render.prometheus.config ## Brings up prometheus and GoCd setup by running container with the config specified
	@docker-compose -f infrastructure/docker-compose.infra.yaml up -d

test.setup.purge: ## Teardown prometheus and GoCd setup
	@docker-compose -f infrastructure/docker-compose.infra.yaml down

local.build:
	@go build

local.run:
	./gocd-prometheus-exporter