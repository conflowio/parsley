.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: help
help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Runs all the tests
	@./test.sh

generate-mocks: ## Regenerates all mocks with mockery
	cd ast && mockery --all --case=underscore
	cd parser && mockery --all --case=underscore
	cd reader && mockery --all --case=underscore
