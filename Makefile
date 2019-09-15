GO_IMG ?= golang:1.12.9
GO_CILINT ?= golangci/golangci-lint:v1.18
SRV ?= gowfnet
PROJECT = github.com/andrskom/${SRV}
PROJECT_DIR ?= $(shell pwd)
DOCKER_BUILD_EXTRA_OPTS ?= -v ~/.netrc:/root/.netrc
COVER_PROFILE ?=

GOOS ?= linux
GOARCH ?= amd64

all: vendorC lintC testC
.PHONY: all

vendor:
	@echo "+ $@"
	@GO111MODULE='on' go mod tidy
	@GO111MODULE='on' go mod vendor
.PHONY: vendor

vendorC:
	@echo "+ $@"
	@docker run --rm -i  \
		${DOCKER_BUILD_EXTRA_OPTS} \
		-v ${PROJECT_DIR}:/go/src/${PROJECT} \
		-v ${GOPATH}/pkg/mod:/go/pkg/mod \
		-e GO111MODULE='on' \
		-w /go/src/${PROJECT} ${GO_IMG} make vendor
.PHONY: vendorC

lint:
	@echo "+ $@"
	@golangci-lint run --enable-all --skip-dirs vendor ./...
.PHONY: lint

lintC:
	@echo "+ $@"
	@docker run --rm -i  \
		${DOCKER_BUILD_EXTRA_OPTS} \
		-v ${PROJECT_DIR}:/go/src/${PROJECT} \
		-w /go/src/${PROJECT} ${GO_CILINT} make lint
.PHONY: lintC

test:
	@echo "+ $@"
	@go test -v -coverprofile=${COVER_PROFILE}  ./...
.PHONY: test

testC:
	@echo "+ $@"
	@docker run --rm -i  \
		${DOCKER_BUILD_EXTRA_OPTS} \
		-v ${PROJECT_DIR}:/go/src/${PROJECT} \
		-e COVER_PROFILE=${COVER_PROFILE} \
		-w /go/src/${PROJECT} ${GO_IMG} make test
.PHONY: testС

clean:
	@echo "+ $@"
	@docker run --rm -i \
		-v ${PROJECT_DIR}:/go/src/${PROJECT} \
		-w /go/src/${PROJECT} ${GO_IMG} rm -rf vendor
.PHONY: clean

cleanC:
	@echo "+ $@ ${GOOS}"
	@docker run --rm -i  \
		${DOCKER_BUILD_EXTRA_OPTS} \
		-v ${PROJECT_DIR}:/go/src/${PROJECT} \
		-e CGO_ENABLED=0 \
		-e GOOS=${GOOS} \
		-e GOARCH=${GOARCH} \
		-w /go/src/${PROJECT} ${GO_IMG} rm -rf ./bin
.PHONY: cleanC
