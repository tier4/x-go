SHELL=/bin/bash -o pipefail
.DEFAULT_GOAL := help

PACKAGE := github.com/tier4/x
GO := go
GOPATH := $(shell go env GOPATH)
GOMODCACHE=$(shell go env GOMODCACHE)
GOBIN := $(abspath .bin)
GOTEST := $(or $(GOTEST),$(GO) test)
GOIMPORTSFLAGS := -local $(PACKAGE)

export GO111MODULE := on
export PATH := $(GOBIN):${PATH}

GO_DEPENDENCIES = golang.org/x/tools/cmd/goimports@master \
				  github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  .bin/$(firstword $(subst @, ,$(notdir $1))): go.mod go.sum Makefile
		GOBIN=$(GOBIN) go install $1
endef
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))
$(call make-lint-dependency)

UNAME := $(shell uname -s)
ifeq ($(UNAME),Linux)
	OPEN := xdg-open
else ifeq ($(UNAME),Darwin)
	OPEN := open
endif

.PHONY: deps
deps: ## download dependencies
	find . -name go.mod -execdir $(GO) mod download \;
	find . -name go.mod -execdir $(GO) mod tidy \;

.PHONY: test
test: unit-test ## test all

.PHONY: unit-test
unit-test: ## unit test
	find . -name go.mod -execdir $(GOTEST) -failfast -cover -count=1 -race ./... \;

.PHONY: clean
clean: ## clean bin and cache
	rm -rf ./.bin
	$(GO) clean -testcache

.PHONY: fmt
fmt: .bin/goimports ## format source
	find . -name go.mod -execdir goimports $(GOIMPORTSFLAGS) -w . \;
	find . -name go.mod -execdir $(GO) mod tidy \;

.PHONY: golangci-lint
golangci-lint: .bin/golangci-lint ## Exec golangci-lint
	find . -name go.mod -execdir golangci-lint run ./... \;

.PHONY: lint
lint: golangci-lint ## exec lint

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -h -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

# https://postd.cc/auto-documented-makefile/
help: ## Show help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
