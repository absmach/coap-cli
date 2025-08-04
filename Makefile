# Copyright (c) Abstract Machines
# SPDX-License-Identifier: Apache-2.0

all:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o build/coap-cli-linux main.go
	CGO_ENABLED=0 GOOS=darwin go build -ldflags "-s -w" -o build/coap-cli-darwin main.go
	CGO_ENABLED=0 GOOS=windows go build -ldflags "-s -w" -o build/coap-cli-windows main.go
