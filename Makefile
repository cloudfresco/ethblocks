SHELL := /bin/bash

APP_NAME ?= ethblocks
BIN_DIR ?= bin
COVERAGE_FILE ?= coverage.out
GOLANGCI_LINT_VERSION ?= v2.12.2
GOVULNCHECK_KNOWN ?= GO-2026-4479

# Go source files, ignoring generated/vendor-style directories.
SRC := $(shell find . -type f -name '*.go' \
	-not -path "./internal/proto/*" \
	-not -path "./internal/proto-gen/*" \
	-not -path "./vendor/*")

GO_TEST_FLAGS ?= -count=1 -race -shuffle=on
GO_TEST_COVER_FLAGS ?= -covermode=atomic -coverprofile=$(COVERAGE_FILE)
GOLANGCI_LINT ?= go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
MYSQL_ROOT_AUTH := $(if $(ETHBLOCKS_DBPASSROOT),-p$(ETHBLOCKS_DBPASSROOT),)
# Connection target for the mysql client
# Using an IP host (127.0.0.1) forces a TCP connection
# A name host (localhost) enables using a native unix socket
# CI sets these vars in the workflow's env block
MYSQL_CONN := -h$(ETHBLOCKS_DBHOST) -P$(ETHBLOCKS_DBPORT)
MYSQL_ROOT := mysql $(MYSQL_CONN) -uroot $(MYSQL_ROOT_AUTH)

.PHONY: all ci verify build test test-env-check test-db-setup test-short coverage clean fmt fmt-check lint lint-fix lint-ci vet vuln mod-verify tidy pkgupd doc

all: verify build

ci: verify build

verify: mod-verify fmt-check lint vet vuln test

build:
	@echo "Building ethblocks"
	@mkdir -p $(BIN_DIR)
	@go build ./...
	@go build -o $(BIN_DIR)/etbcli main.go

test-env-check:
	@test -n "$(ETHBLOCKS_DB)" || (echo "ETHBLOCKS_DB is required" && exit 1)
	@test -n "$(ETHBLOCKS_DBUSER_TEST)" || (echo "ETHBLOCKS_DBUSER_TEST is required" && exit 1)
	@test -n "$(ETHBLOCKS_DBPASS_TEST)" || (echo "ETHBLOCKS_DBPASS_TEST is required" && exit 1)
	@test -n "$(ETHBLOCKS_DBHOST)" || (echo "ETHBLOCKS_DBHOST is required" && exit 1)
	@test -n "$(ETHBLOCKS_DBPORT)" || (echo "ETHBLOCKS_DBPORT is required" && exit 1)
	@test -n "$(ETHBLOCKS_DBNAME_TEST)" || (echo "ETHBLOCKS_DBNAME_TEST is required" && exit 1)
	@test -n "$(ETHBLOCKS_CLIENT)" || (echo "ETHBLOCKS_CLIENT is required" && exit 1)

test-db-setup: test-env-check
	@$(MYSQL_ROOT) -e 'DROP DATABASE IF EXISTS  $(ETHBLOCKS_DBNAME_TEST);'
	@$(MYSQL_ROOT) -e 'CREATE DATABASE $(ETHBLOCKS_DBNAME_TEST);'
	@$(MYSQL_ROOT) -e "CREATE USER IF NOT EXISTS '$(ETHBLOCKS_DBUSER_TEST)'@'%' IDENTIFIED BY '$(ETHBLOCKS_DBPASS_TEST)';"
	@$(MYSQL_ROOT) -e "GRANT ALL ON $(ETHBLOCKS_DBNAME_TEST).* TO '$(ETHBLOCKS_DBUSER_TEST)'@'%';"
	@$(MYSQL_ROOT) -e 'FLUSH PRIVILEGES;'
	@mysql $(MYSQL_CONN) -u$(ETHBLOCKS_DBUSER_TEST) -p$(ETHBLOCKS_DBPASS_TEST) $(ETHBLOCKS_DBNAME_TEST) < test/sql/mysql/ethblocks_mysql_schema.sql

test: test-db-setup
	@echo "Starting tests for ethblocks"
	@go test $(GO_TEST_FLAGS) $(GO_TEST_COVER_FLAGS) -v ./...

test-short:
	@go test -short $(GO_TEST_FLAGS) ./...

coverage: test
	@go tool cover -func=$(COVERAGE_FILE) | grep total

clean:
	@rm -rf $(BIN_DIR) $(COVERAGE_FILE)

fmt:
	@echo "Running gofumpt"
	@$(GOLANGCI_LINT) fmt ./...

fmt-check:
	@echo "Checking formatting"
	@diff=$$($(GOLANGCI_LINT) fmt --diff ./...); \
	if [ -n "$$diff" ]; then \
		echo "$$diff"; \
		echo "Run 'make fmt' to apply formatting fixes."; \
		exit 1; \
	fi

lint:
	@echo "Running golangci-lint"
	@$(GOLANGCI_LINT) run ./...

lint-fix:
	@echo "Running golangci-lint with fixes"
	@$(GOLANGCI_LINT) run --fix ./...

lint-ci:
	@echo "Running golangci-lint against PR changes"
	@$(GOLANGCI_LINT) run --new-from-merge-base=origin/master --whole-files ./...

vet:
	@echo "Running go vet"
	@go vet ./...

vuln:
	@echo "Running govulncheck"
	@set +e; \
	output=$$(go run golang.org/x/vuln/cmd/govulncheck@latest ./... 2>&1); \
	status=$$?; \
	printf '%s\n' "$$output"; \
	if [ $$status -eq 0 ]; then \
		exit 0; \
	fi; \
	ids=$$(printf '%s\n' "$$output" | grep -oE 'GO-[0-9]{4}-[0-9]+' | sort -u || true); \
	unknown=$$(comm -23 <(printf '%s\n' "$$ids") <(printf '%s\n' $(GOVULNCHECK_KNOWN) | sort -u)); \
	if [ -n "$$unknown" ]; then \
		echo "govulncheck found unacknowledged vulnerabilities:"; \
		echo "$$unknown"; \
		exit $$status; \
	fi; \
	echo "govulncheck found only acknowledged vulnerabilities: $(GOVULNCHECK_KNOWN)";

mod-verify:
	@echo "Verifying module dependencies"
	@go mod download
	@go mod verify

tidy:
	@go mod tidy -v

pkgupd:
	@go get -u ./...
	@go mod tidy -v

doc:
