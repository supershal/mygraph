export GO111MODULE=on
PACKAGES=$(shell go list ./... | grep -v '/vendor')
MAIN_GO_SERVER="./cmd/server"
MAIN_GO_CLIENT="./cmd/client"
SERVER_BINARY="graphserver"
CLIENT_BINARY="graphclient"
GOOS?=darwin

GOMODFLAG?=-mod=vendor
CGO_ENABLED?=0
GOARCH?=amd64


.PHONY: test build install lint setup clean all
default: all

all: setup deps fmt lint vet test bin

setup: ${GOPATH}/bin/golint
	
${GOPATH}/bin/golint:
	go get -u golang.org/x/lint/golint

deps:
	go mod download
	go mod vendor -v

test:
	@echo "tests..."
	@go test -v -timeout=2m ${PACKAGES}

build:
	@echo "building..."
	@go build ${PACKAGES}

bin:
	@echo "building and installing..."
	@GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build --installsuffix cgo  -o ${SERVER_BINARY} ${MAIN_GO_SERVER}
	@GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build --installsuffix cgo  -o ${CLIENT_BINARY} ${MAIN_GO_CLIENT}

lint:
	@echo "linting..."
	@${GOPATH}/bin/golint ${PACKAGES}


clean:
	go clean -i -x -r


vet:
	@echo "govet..."
	@go vet ${PACKAGES}

fmt:
	@go fmt ${PACKAGES}