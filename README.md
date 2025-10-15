# CoAP CLI

Simple CoAP cli client written in Go.

## Installation

### Linux

```bash
# x86_64
curl -sL https://github.com/absmach/coap-cli/releases/download/v0.4.1/coap-cli-linux-amd64 -o coap-cli && chmod +x coap-cli
# arm64
curl -sL https://github.com/absmach/coap-cli/releases/download/v0.4.1/coap-cli-linux-arm64 -o coap-cli && chmod +x coap-cli
# riscv64
curl -sL https://github.com/absmach/coap-cli/releases/download/v0.4.1/coap-cli-linux-riscv64 -o coap-cli && chmod +x coap-cli
```

### MacOS

```bash
# x86_64
curl -sL https://github.com/absmach/coap-cli/releases/download/v0.4.1/coap-cli-darwin-amd64 -o coap-cli && chmod +x coap-cli
# arm64
curl -sL https://github.com/absmach/coap-cli/releases/download/v0.4.1/coap-cli-darwin-arm64 -o coap-cli && chmod +x coap-cli
```

### Windows

```bash
# x86_64
curl -sL https://github.com/absmach/coap-cli/releases/download/v0.4.1/coap-cli-windows-amd64 -o coap-cli.exe
# arm64
curl -sL https://github.com/absmach/coap-cli/releases/download/v0.4.1/coap-cli-windows-arm64 -o coap-cli.exe
```

### Build from source

Make sure you have Go and Make installed.

```bash
git clone https://github.com/absmach/coap-cli.git
cd coap-cli
make install # or INSTALL_DIR=~/go/bin make install
```

## Usage

```bash
CLI for CoAP

Usage:
  coap-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Perform a DELETE request on a COAP resource
  get         Perform a GET request on a COAP resource
  help        Help about any command
  post        Perform a POST request on a COAP resource
  put         Perform a PUT request on a COAP resource

Flags:
  -a, --auth string           Auth
  -A, --ca-file string        Client CA file
  -C, --cert-file string      Client certificate file
  -c, --content-format int    Content format (default 50)
  -h, --help                  help for coap-cli
  -H, --host string           Host (default "localhost")
  -k, --keep-alive uint       Send a ping after interval seconds of inactivity. If not specified (or 0), keep-alive is disabled (default).
  -K, --key-file string       Client key file
  -m, --max-retries uint32    Max retries for keep alive (default 10)
  -O, --options stringArray   Add option num with contents of text to the request. If the text begins with 0x, then the hex text (two [0-9a-f] per byte) is converted to binary data.
  -p, --port string           Port (default "5683")
  -v, --verbose               Verbose output

Use "coap-cli [command] --help" for more information about a command.
```

The options flag accepts a comma separated string comprising of the optionID defined by [RFC-7252](https://datatracker.ietf.org/doc/html/rfc7252) and a string or hex value. Hex values are used to set options that require numerical values e.g observe, maxAge

## Examples

```bash
coap-cli get m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -o
```

```bash
coap-cli get m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic --options 6,0x00 --options 15,auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb
```

```bash
coap-cli get m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic --options 6,0x00 --options 15,auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb --ca-file ssl/certs/ca.crt --cert-file ssl/certs/client.crt --key-file ssl/certs/client.key
```

```bash
coap-cli post m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world"
```

```bash
coap-cli post m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" -H 0.0.0.0 -p 1234
```

```bash
coap-cli post m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic --options 15,auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" -H 0.0.0.0 -p 5683
```

```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" --ca-file ssl/certs/ca.crt --cert-file ssl/certs/client.crt --key-file ssl/certs/client.key
```
