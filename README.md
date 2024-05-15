# CoAP CLI

Simple CoAP cli client written in Go.

## Usage

```bash
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
  -c, --content-format int    Content format (default 50)
  -h, --help                  help for coap-cli
  -H, --host string           Host (default "localhost")
  -k, --keep-alive uint       Send a ping after interval seconds of inactivity. If not specified (or 0), keep-alive is disabled (default).
  -m, --max-retries uint32    Max retries for keep alive (default 10)
  -O, --options num,text      Add option num with contents of text to the request. If the text begins with 0x, then the hex text (two [0-9a-f] per byte) is converted to binary data.
  -p, --port string           Port (default "5683")
  -v, --verbose               Verbose output
  -d, --data string           Data(default "") - only available for put, post and delete commands
  -o, --observe bool          Observe - only available for get command

Use "coap-cli [command] --help" for more information about a command
```

The options flag accepts a comma separated string comprising of the optionID defined by [RFC-7252](https://datatracker.ietf.org/doc/html/rfc7252) and a string or hex value. Hex values are used to set options that require numerical values e.g observe, maxAge

## Examples

```bash
coap-cli get channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -o
```

```bash
coap-cli get channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --options 6,0x00 --options 15,auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb
```

```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world"
```

```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" -H 0.0.0.0 -p 1234
```

```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -options 15,auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" -H 0.0.0.0 -p 5683
```
