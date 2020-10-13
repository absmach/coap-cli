# Coap-cli
Simple CoAP cli client written in Go.


## Usage
Pre-built binary can be found here: https://github.com/mainflux/coap-cli/releases/tag/v0.2.0.
When running, please provide following format:
`coap-cli` followed by method code (`get`, `put`, `post`, `delete`) and CoAP URL. After that, you can pass following flags:

| Flag | Description                                    | Default value |
| ---- | ---------------------------------------------- | ------------- |
| o    | Observe   option - only valid with Get request | false         |
| auth | Auth option sent as URI Query                  | ""            |
| h    | Host                                           | "localhost"   |
| p    | port                                           | "5683"        |
| d    | Data to be sent in POST or PUT                 | ""            |
| cf   | Content format                                 | 50 - JSON     |
# Examples:

```bash
coap-cli get channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb -o
```
```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world"
```
```bash
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -auth=1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" -h 0.0.0.0 -p 1234
```