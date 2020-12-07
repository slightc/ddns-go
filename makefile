# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BUILD_DIR=build
BINARY_NAME=ddns-go

all: build
.PHONY: all

build:
	$(GOBUILD) -o ${BUILD_DIR}/$(BINARY_NAME) -v
.PHONY: build

install:
	cp ${BUILD_DIR}/$(BINARY_NAME) /usr/bin/${BINARY_NAME}
	mkdir -p /etc/ddns-go
	cp ./config.yaml.default /etc/ddns-go/config.yaml
	cp ./systemd/ddns-go.service /etc/systemd/systemd/ddns-go.service
.PHONY: install