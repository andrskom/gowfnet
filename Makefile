GO_IMG ?= golang:1.14
GO_CILINT ?= golangci/golangci-lint:v1.55
SRV ?= gowfnet
PROJECT_DIR ?= $(shell pwd)
COVER_PROFILE ?= cover.out
COVER_HTML_OUT ?= cover.html

GOOS ?= linux
GOARCH ?= amd64

all: lintC testC
.PHONY: all

lint:
	@echo "+ $@"
	@golangci-lint run ./...
.PHONY: lint

lintC:
	@echo "+ $@"
	@docker run --rm -i  \
		-v ${PROJECT_DIR}:/app/${SRV} \
		-v ${GOPATH}:/go \
		-w /app/${SRV} ${GO_CILINT} make lint
.PHONY: lintC

test:
	@echo "+ $@"
	@go test -v -coverprofile=${COVER_PROFILE} ./...
	@go tool cover -html=${COVER_PROFILE} -o ${COVER_HTML_OUT}
.PHONY: test

testC:
	@echo "+ $@"
	@docker run --rm -i  \
		-v ${PROJECT_DIR}:/app/${SRV} \
		-v ${GOPATH}:/go \
		-e COVER_PROFILE=${COVER_PROFILE} \
		-w /app/${SRV} \
		${GO_IMG} make test
.PHONY: testC
