

OS ?= $(shell uname -s)
GO ?= go
GOGET ?= $(GO) get -u
GOBIN ?= $(go env GOPATH)/bin

# Build (default target)
GOBUILD ?= gox
BUILD_OPTIONS ?= -mod=readonly -output="build/terraform-provider-maas_${CIRCLE_TAG:=}_{{.OS}}_{{.Arch}}"
build: install_gox
	$(GOBUILD) $(BUILD_OPTIONS)
.PHONY: build
.DEFAULT_GOAL := build

install_gox:
	command -v gox || go get github.com/mitchellh/gox
.PHONY: install_gox

# Lint (https://github.com/golangci/golangci-lint)
LINTER_OPTIONS ?= run# Arguments to golangci-lint
LINTER_BINARY ?= golangci-lint# Name of the binary of golangci-lint
LINTER_VERSION ?= 1.29.0# Version of golangci-lint to use in CI

lint: install_lint
	$(LINTER_BINARY) $(LINTER_OPTIONS)
.PHONY: lint

install_lint:
ifneq (1,$(shell $(LINTER_BINARY) version 2>&1 | grep -c $(LINTER_VERSION)))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v$(LINTER_VERSION)
endif
.PHONY: install_lint

# Test (go test)
GOTEST ?= $(GO) test
TEST_OPTS ?= -race -bench -mod=readonly# Perform any benchmarks and enable the race detector
test:
	$(GOTEST) $(TEST_OPTS) ./...
.PHONY: test

# Check everything
check: lint test
.PHONY: check

# Use Go Modules
depend: # Fetch modules with `go mod download`
	go mod download
.PHONY: depend

update_depend: # Update (I think) modules with `go mod tidy`
	go mod tidy
.PHONY: update_depend

# Hook up the local environment
local: .git/hooks/pre-commit install_lint depend ## Set up a local development environment
.PHONY: local

.git/hooks/pre-commit:
	@echo "Installing pre-commit hook"
	echo "make check" > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
