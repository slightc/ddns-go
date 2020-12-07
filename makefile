# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=build/ddns-go

all: build
.PHONY: all

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: build