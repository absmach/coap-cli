# Coap-cli
Simple CoAP cli client written in Go.


## Usage
Pre-built binary can be found here: https://github.com/mainflux/coap-cli/releases/tag/v0.1.0.
When running, please provide following format:
`coap-cli` followed by method code (`get`, `put`, `post`, `delete`) and CoAP URL. After that, you can pass following flags:

| Flag  | Description                    | Default value                          |
|-------|--------------------------------|----------------------------------------|
| ACK   | Acknowledgement                | false                                  |
| C     | Confirmable                    | false                                  |
| NC    | Non-Confirmable                | false                                  |
| O     | Observe                        | false                                  |
| d     | Data to be sent in POST or PUT | ""                                     |
| id    | Message ID                     | 0                                      |
| token | Token                          | Byte array of empty string: [49 50 51] |
# Examples:

```bash
coap-cli get coap://localhost/channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic\?authorization=1e1017e6-dee7-45b4-8a13-00e6afeb66eb -O
```
```bash
coap-cli post coap://localhost/channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic\?authorization=1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world"
```