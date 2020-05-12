GO_IMG ?= golang:1.14
GO_CILINT ?= golangci/golangci-lint:v1.26
SRV ?= gowfnet
PROJECT_DIR ?= $(shell pwd)
COVER_PROFILE ?=

GOOS ?= linux
GOARCH ?= amd64

all: lintC testC
.PHONY: all

lint:
	@echo "+ $@"
	@golangci-lint run --enable-all --disable testpackage ./...
.PHONY: lint

lintC:
	@echo "+ $@"
	@docker run --rm -i  \
		-v ${PROJECT_DIR}:/app/${SRV} \
		-w /app/${SRV} ${GO_CILINT} make lint
.PHONY: lintC

test:
	@echo "+ $@"
	@go test -v -coverprofile=${COVER_PROFILE}  ./...
.PHONY: test

testC:
	@echo "+ $@"
	@docker run --rm -i  \
		-v ${PROJECT_DIR}:/app/${SRV} \
		-e COVER_PROFILE=${COVER_PROFILE} \
		-w /app/${SRV} \
		${GO_IMG} make test
.PHONY: testС
