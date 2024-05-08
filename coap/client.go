// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package coap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
)

var errInvalidMsgCode = errors.New("message can be GET, POST, PUT or DELETE")

// Client represents CoAP client.
type Client struct {
	conn *client.ClientConn
}

// New returns new CoAP client connecting it to the server.
func New(addr string) (Client, error) {
	c, err := udp.Dial(addr)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	return Client{conn: c}, nil
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
	default:
		return nil, errInvalidMsgCode
	}
}

// Receive receives a message.
func (c Client) Receive(path string, opts ...message.Option) (*client.Observation, error) {
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
