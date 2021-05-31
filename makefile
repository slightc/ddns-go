# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BUILD_DIR=build
BINARY_NAME=ddns-go

VERSION?=0.0.0

CROSS_BUILD_DIR=$(BUILD_DIR)/$(1)-$(2)
CROSS_BUILD=CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) $(GOBUILD) -o $(call CROSS_BUILD_DIR,$(1),$(2))/$(BINARY_NAME) -v

all: build
.PHONY: all

build:
	$(GOBUILD) -o ${BUILD_DIR}/$(BINARY_NAME) -v
.PHONY: build

TAR_PACKAGE=cd ${BUILD_DIR} && tar -zcvf $(BINARY_NAME)-$(1)-$(2)-${VERSION}.tar.gz $(1)-$(2)/*

CROSS_BUILD_LINUX=$(call CROSS_BUILD,linux,$(1))\
	&& cp config.yaml.default $(call CROSS_BUILD_DIR,linux,$(1))/\
	&& cp systemd/* $(call CROSS_BUILD_DIR,linux,$(1))/\
	&& cp script/linux/* $(call CROSS_BUILD_DIR,linux,$(1))/\
	&& $(call TAR_PACKAGE,linux,$(1))

CROSS_BUILD_DARWIN=$(call CROSS_BUILD,darwin,$(1))\
	&& cp config.yaml.default $(call CROSS_BUILD_DIR,darwin,$(1))/\
	&& cp script/darwin/* $(call CROSS_BUILD_DIR,darwin,$(1))/\
	&& $(call TAR_PACKAGE,darwin,$(1))

CROSS_BUILD_WINDOWS=$(call CROSS_BUILD,windows,$(1))\
	&& cp config.yaml.default $(call CROSS_BUILD_DIR,windows,$(1))/\
	&& $(call TAR_PACKAGE,windows,$(1))

linux_arm:
	$(call CROSS_BUILD_LINUX,arm)
linux_arm64:
	$(call CROSS_BUILD_LINUX,arm64)

linux_amd64:
	$(call CROSS_BUILD_LINUX,amd64)
linux_386:
	$(call CROSS_BUILD_LINUX,386)

linux_all: linux_arm linux_arm64

darwin_amd64:
	$(call CROSS_BUILD_DARWIN,amd64)

darwin_all: darwin_amd64

windows_386:
	$(call CROSS_BUILD_WINDOWS,386)
windows_amd64:
	$(call CROSS_BUILD_WINDOWS,amd64)

windows_all: windows_386 windows_amd64

package_all: linux_all darwin_all windows_all

install:
	cp ${BUILD_DIR}/$(BINARY_NAME) /usr/bin/${BINARY_NAME}
	mkdir -p /etc/ddns-go
	if [ ! -f /etc/ddns-go/config.yaml ]; then cp ./config.yaml.default /etc/ddns-go/config.yaml; fi;
	cp ./systemd/ddns-go.service /etc/systemd/system/ddns-go.service
.PHONY: install