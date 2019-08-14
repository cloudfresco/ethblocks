
SHELL := /bin/bash
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKGS = ./...

.PHONY: all build test fmt vet lint err doc

all: chk build

chk: fmt vet lint err

build: 
	@echo "Building ethblocks"	
	@go build $(PKGS)

test:

	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e 'DROP DATABASE IF EXISTS  $(ETHBLOCKS_DBNAME_TEST);'
	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e 'CREATE DATABASE $(ETHBLOCKS_DBNAME_TEST);'
	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e "GRANT ALL ON *.* TO '$(ETHBLOCKS_DBUSER_TEST)'@'$(ETHBLOCKS_DBHOST)';"
	@mysql -uroot -p$(ETHBLOCKS_DBPASSROOT) -e 'FLUSH PRIVILEGES;'
	@mysql -u$(ETHBLOCKS_DBUSER_TEST) -p$(ETHBLOCKS_DBPASS_TEST)  $(ETHBLOCKS_DBNAME_TEST) < sql/mysql/ethblocks_mysql_schema.sql

	@echo "Starting tests"
	
	@go test -v $(PKGS)

fmt:
	@echo "Running gofmt"	
	@gofmt -s -l -w $(SRC)

vet:
	@echo "Running vet"	
	@go vet $(PKGS)

linter:
	@go get -u golang.org/x/lint/golint

lint: linter
	@echo "Running lint"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done 

errcheck: 
	@go get github.com/kisielk/errcheck

err: errcheck 
	@echo "Running errcheck"
	@errcheck $(PKGS)

doc: 

