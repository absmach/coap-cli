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

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/mux"
	"github.com/plgd-dev/go-coap/v3/options"
	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/plgd-dev/go-coap/v3/udp/client"
)

var (
	errInvalidMsgCode = errors.New("message can be GET, POST, PUT or DELETE")
	errDialFailed     = errors.New("failed to dial the connection")
)

const verboseFmt = `Date: %s
Code: %s
Type: %s
Token: %s
Message-ID: %d
Content-Length: %d
`

// Client represents CoAP client.
type Client struct {
	conn *client.Conn
}

// NewClient returns new CoAP client connecting it to the server.
func NewClient(addr string, keepAlive uint64, maxRetries uint32) (Client, error) {
	var dialOptions []udp.Option
	if keepAlive > 0 {
		dialOptions = append(dialOptions, options.WithKeepAlive(maxRetries, time.Duration(keepAlive)*time.Second, onInactive))
	}
	c, err := udp.Dial(addr, dialOptions...)
	if err != nil {
		return Client{}, errors.Join(errDialFailed, err)
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
func (c Client) Receive(path string, verbose bool, opts ...message.Option) (mux.Observation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return c.conn.Observe(ctx, path, func(res *pool.Message) {
		body, err := res.ReadBody()
		if err != nil {
			fmt.Println("Error reading message body: ", err)
			return
		}
		bodySize, err := res.BodySize()
		if err != nil {
			fmt.Println("Error getting body size: ", err)
			return
		}
		if bodySize == 0 {
			fmt.Println("Received observe")
		}
		switch verbose {
		case true:
			fmt.Printf(verboseFmt,
				time.Now().Format(time.RFC1123),
				res.Code(),
				res.Type(),
				res.Token(),
				res.MessageID(),
				bodySize)
			if len(body) > 0 {
				fmt.Printf("Payload: %s\n\n", string(body))
			}
		case false:
			if len(body) > 0 {
				fmt.Printf("Payload: %s\n", string(body))
			}
		}
	}, opts...)
}

func onInactive(cc *client.Conn) {
	if err := cc.Ping(cc.Context()); err != nil {
		log.Fatalf("Error pinging: %v", err)
	}
}
