# Copyright (c) Abstract Machines
# SPDX-License-Identifier: Apache-2.0

INSTALL_DIR=/usr/local/bin
BUILD_DIR=build
BUILD_FLAGS=-ldflags "-s -w"

.PHONY: all linux darwin windows install install-linux

all: linux darwin windows

linux:
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-linux cmd/main.go

darwin:
	CGO_ENABLED=0 GOOS=darwin go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-darwin cmd/main.go

windows:
	CGO_ENABLED=0 GOOS=windows go build $(BUILD_FLAGS) -o $(BUILD_DIR)/coap-cli-windows cmd/main.go
	
install: install-linux 

install-linux:
	@cp $(BUILD_DIR)/coap-cli-linux $(INSTALL_DIR)/coap-cli || { echo "Installation failed"; exit 1; }

clean:
	rm -rf $(BUILD_DIR)/*
