# Copyright (c) Abstract Machines
# SPDX-License-Identifier: Apache-2.0

INSTALL_DIR ?= /usr/local/bin
BUILD_DIR ?= build
BUILD_FLAGS ?= -ldflags "-s -w"
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: all linux darwin windows install

all: linux darwin windows

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-linux-amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-linux-arm64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-linux-riscv64 main.go

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-darwin-amd64 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-darwin-arm64 main.go

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-windows-amd64 main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-windows-arm64 main.go

install:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli main.go
	cp $(BUILD_DIR)/coap-cli $(INSTALL_DIR)/coap-cli || { echo "Installation failed"; exit 1; }

clean:
	rm -rf $(BUILD_DIR)/*
