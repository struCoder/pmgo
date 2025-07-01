# PMGO Makefile
.PHONY: build clean test test-coverage lint run install dev docs

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=pmgo
BINARY_UNIX=$(BINARY_NAME)_unix

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD)

# Linker flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/pmgo

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v ./cmd/pmgo

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

run:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/pmgo
	./$(BINARY_NAME)

install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./cmd/pmgo

dev:
	air -c .air.toml

docs:
	swag init -g cmd/pmgo/main.go -o ./docs

deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Cross compilation
build-all: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/$(BINARY_NAME)-linux-amd64 ./cmd/pmgo
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/$(BINARY_NAME)-darwin-amd64 ./cmd/pmgo
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/$(BINARY_NAME)-windows-amd64.exe ./cmd/pmgo

# Docker
docker-build:
	docker build -t pmgo:$(VERSION) .

docker-run:
	docker run -p 9876:9876 -p 8080:8080 pmgo:$(VERSION)

release: build-all
	@echo "Building release $(VERSION)"
	@mkdir -p build
	@echo "$(VERSION)" > build/VERSION