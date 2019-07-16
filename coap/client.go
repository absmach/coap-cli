package coap

import (
	"net/url"

	gocoap "github.com/dustin/go-coap"
)

// Client represents CoAP client.
type Client struct {
	conn *conn
}

// New returns new CoAP client connecting it to the server.
func New(addr *url.URL) (Client, error) {
	address, err := parseAddr(addr)
	if err != nil {
		return Client{}, err
	}
	c, err := dial(address)
	if err != nil {
		return Client{}, err
	}
	return Client{conn: c}, nil
}

// Send send a message.
func (c Client) Send(msgType gocoap.COAPType, msgCode gocoap.COAPCode, msgID uint16, token, payload []byte, opts []Option) (*gocoap.Message, error) {
	msg := gocoap.Message{
		Type:      msgType,
		Code:      msgCode,
		MessageID: msgID,
		Token:     token,
		Payload:   payload,
	}

	for _, o := range opts {
		msg.SetOption(o.ID, o.Value)
	}
	return c.conn.send(msg)
}

// Receive receives a message.
func (c Client) Receive() (*gocoap.Message, error) {
	return c.conn.receive()
}
