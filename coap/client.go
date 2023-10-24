package coap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	certutil "github.com/mainflux/coap-cli/certutil"

	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/plgd-dev/go-coap/v3/udp/client"
)

// Client represents CoAP client.
type Client struct {
	conn *client.Conn
}

// Observation interface
type NewObservation interface {
	Cancel(ctx context.Context, opts ...message.Option) error
	Canceled() bool
}

// New returns new CoAP client connecting it to the server.
func New(addr string, certPath string) (Client, error) {
	switch {
	case certPath != "":
		config, err := certutil.CreateClientConfig(context.Background(), certPath)
		if err != nil {
			log.Fatalln(err)
		}
		co, err := dtls.Dial(addr, config)
		if err != nil {
			log.Fatalf("Error dialing: %v", err)
		}
		return Client{conn: co}, err
	default:
		c, err := udp.Dial(addr)
		if err != nil {
			log.Fatalf("Error dialing: %v", err)
		}
		return Client{conn: c}, nil

	}
}

// Send send a message.
func (c Client) Send(path string, msgCode codes.Code, cf message.MediaType, payload io.ReadSeeker, opts ...message.Option) (*pool.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	switch msgCode {
	case codes.GET:
		return c.conn.Get(ctx, path, opts...)
	case codes.POST:
		return c.conn.Post(ctx, path, cf, payload, opts...)
	case codes.PUT:
		return c.conn.Put(ctx, path, cf, payload, opts...)
	case codes.DELETE:
		return c.conn.Delete(ctx, path, opts...)
	}
	return nil, errors.New("INVALID MESSAGE CODE")
}

// Receive receives a message.
func (c Client) Receive(path string, opts ...message.Option) (NewObservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return c.conn.Observe(ctx, path, func(res *pool.Message) {
		fmt.Printf("\nRECEIVED OBSERVE: %v\n", res)
		body, err := res.ReadBody()
		if err != nil {
			fmt.Println("Error reading message body: ", err)
			return
		}
		if len(body) > 0 {
			fmt.Println("Payload: ", string(body))
		}
	}, opts...)

}
