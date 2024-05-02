# CoAP CLI

Simple CoAP cli client written in Go.

## Usage

coap-cli [`command`]

#### Available Commands:

- `completion` Generate the autocompletion script for the specified shell
- `delete` Perform a DELETE request on a COAP resource
- `get` Perform a GET request on a COAP resource
- `help` Help about any command
- `post` Perform a POST request on a COAP resource
- `put` Perform a PUT request on a COAP resource

#### Flags:

- `-a`, `--auth` string Auth (default "")
- `-c`, `--content-format` int Content format (default 50)
- `-h`, `--help` help for coap-cli
- `-H`, `--host` string Host (default "localhost")
- `-p`, `--port` string Port (default "5683")

Use "coap-cli [command] --help" for more information about a command.

## Examples:

```bash
coap-cli get channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -o
```

```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world"
```

```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" -H 0.0.0.0 -p 1234
```
