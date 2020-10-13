all:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags "-s -w" -o build/coap-cli-linux cmd/main.go
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -ldflags "-s -w" -o build/coap-cli-darwin cmd/main.go
	CGO_ENABLED=0 GOOS=windows go build -mod=vendor -ldflags "-s -w" -o build/coap-cli-windows cmd/main.go