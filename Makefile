
SHELL := /bin/bash
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKGS = ./...
.PHONY: all deps build test fmt vet lint err doc

all: fmt vet lint err build

deps: 
	@dep ensure

build: 
	@go build $(PKGS)

test: deps build
	@go test -v $(PKGS)

fmt:
	@gofmt -s -l -w $(SRC)

vet:
	@go vet $(PKGS)

linter:
	@go get -u github.com/golang/lint/golint

lint: linter
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done 

errcheck: 
	@go get github.com/kisielk/errcheck

err: errcheck 
	@errcheck $(PKGS)

doc: 

