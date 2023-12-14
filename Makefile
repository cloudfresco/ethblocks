SHELL := /bin/bash


# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./internal/proto-gen/*")
#SRCPROTO = $(shell find . -type f -name '*.proto'")
MFILE = cmd/main.go
EXEC = cmd/ethblocks
PKGS = $(go list ./... | grep -v /proto/ | grep -v /proto-gen/)

.PHONY: all chk lint build test clean fmt gocritic staticcheck errcheck revive golangcilint tidy pkgupd doc

all: chk buildp

chk: goimports fmt gocritic staticcheck errcheck

rev: revive

lint: golangcilint
 
build: 
	@echo "Building ethblocks"	
	@go build ./...

test:

	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e 'DROP DATABASE IF EXISTS  $(ETHBLOCKS_DBNAME_TEST);'
	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e 'CREATE DATABASE $(ETHBLOCKS_DBNAME_TEST);'
	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e "GRANT ALL ON *.* TO '$(ETHBLOCKS_DBUSER_TEST)'@'$(ETHBLOCKS_DBHOST)';"
	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e 'FLUSH PRIVILEGES;'
	@mysql -u$(ETHBLOCKS_DBUSER_TEST) -p$(ETHBLOCKS_DBPASS_TEST)  $(ETHBLOCKS_DBNAME_TEST) < test/sql/mysql/ethblocks_mysql_schema.sql

	@echo "Starting tests for ethblocks"
	
	@go test -v ./...

clean:

goimports:
	@echo "Running goimports"		
	@goimports -l -w $(SRC)

fmt:
	@echo "Running gofumpt"
	@gofumpt -l -w .
	@echo "Running gofmt"		
	@gofmt -s -l -w $(SRC)

gocritic:
	@echo "Running gocritic"
	@gocritic check $(SRC)

staticcheck:
	@echo "Running staticcheck"
	@staticcheck ./...

errcheck:
	@echo "Running errcheck"
	@errcheck ./...

revive:
	@echo "Running revive"
	@revive $(SRC)

golangcilint:
	@echo "Running golangci-lint"
	@golangci-lint run

tidy:
	go mod tidy -v -e

pkgupd:
	go get -u ./...
	go mod tidy -v -e

doc: 

