.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: help
help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Runs all the tests
	@./test.sh

generate-mocks: ## Regenerates all mocks with mockery
	go generate ./...

.PHONY: update-dependencies
update-dependencies: ## Updates all dependencies
	@echo "Updating Go dependencies"
	@cat go.mod | grep -E "^\t" | grep -v "// indirect" | cut -f 2 | cut -d ' ' -f 1 | xargs -n 1 -t go get -d -u
	@go mod vendor
	@go mod tidy
