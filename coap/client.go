package coap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pion/dtls/v2"
	dtlsnet "github.com/pion/dtls/v2/pkg/net"
	"github.com/plgd-dev/go-coap/v3/dtls/server"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	coapNet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/net/blockwise"
	"github.com/plgd-dev/go-coap/v3/net/monitor/inactivity"
	"github.com/plgd-dev/go-coap/v3/net/responsewriter"
	"github.com/plgd-dev/go-coap/v3/options"
	"github.com/plgd-dev/go-coap/v3/udp"
	client "github.com/plgd-dev/go-coap/v3/udp/client"
)

type Client_struct struct {
	conn *client.Conn
}

var DefaultConfig = func() client.Config {
	cfg := client.DefaultConfig
	cfg.Handler = func(w *responsewriter.ResponseWriter[*client.Conn], r *pool.Message) {
		switch r.Code() {
		case codes.POST, codes.PUT, codes.GET, codes.DELETE:
			if err := w.SetResponse(codes.NotFound, message.TextPlain, nil); err != nil {
				cfg.Errors(fmt.Errorf("dtls client: cannot set response: %w", err))
			}
		}
	}
	return cfg
}()

// New returns new CoAP client connecting it to the server.
func New(addr string) (Client_struct, error) {
	c, err := udp.Dial(addr)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	return Client_struct{conn: c}, nil
}

// Dial creates a client connection to the given target.
func Dial(target string, dtlsCfg *dtls.Config, opts ...udp.Option) (*client.Conn, error) {
	cfg := DefaultConfig
	for _, o := range opts {
		o.UDPClientApply(&cfg)
	}

	c, err := cfg.Dialer.DialContext(cfg.Ctx, cfg.Net, target)
	if err != nil {
		return nil, err
	}

	conn, err := dtls.Client(dtlsnet.PacketConnFromConn(c), c.RemoteAddr(), dtlsCfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, options.WithCloseSocket())
	return Client(conn, opts...), nil
}

// Client creates client over dtls connection.
func Client(conn *dtls.Conn, opts ...udp.Option) *client.Conn {
	cfg := DefaultConfig
	for _, o := range opts {
		o.UDPClientApply(&cfg)
	}
	if cfg.Errors == nil {
		cfg.Errors = func(error) {
			// default no-op
		}
	}
	if cfg.CreateInactivityMonitor == nil {
		cfg.CreateInactivityMonitor = func() client.InactivityMonitor {
			return inactivity.NewNilMonitor[*client.Conn]()
		}
	}
	if cfg.MessagePool == nil {
		cfg.MessagePool = pool.New(0, 0)
	}
	errorsFunc := cfg.Errors
	cfg.Errors = func(err error) {
		if coapNet.IsCancelOrCloseError(err) {
			// this error was produced by cancellation context or closing connection.
			return
		}
		errorsFunc(fmt.Errorf("dtls: %v: %w", conn.RemoteAddr(), err))
	}

	createBlockWise := func(cc *client.Conn) *blockwise.BlockWise[*client.Conn] {
		return nil
	}
	if cfg.BlockwiseEnable {
		createBlockWise = func(cc *client.Conn) *blockwise.BlockWise[*client.Conn] {
			v := cc
			return blockwise.New(
				v,
				cfg.BlockwiseTransferTimeout,
				cfg.Errors,
				func(token message.Token) (*pool.Message, bool) {
					return v.GetObservationRequest(token)
				},
			)
		}
	}

	monitor := cfg.CreateInactivityMonitor()
	l := coapNet.NewConn(conn)
	session := server.NewSession(cfg.Ctx,
		l,
		cfg.MaxMessageSize,
		cfg.MTU,
		cfg.CloseSocket,
	)
	cc := client.NewConn(session,
		createBlockWise,
		monitor,
		&cfg,
	)

	cfg.PeriodicRunner(func(now time.Time) bool {
		cc.CheckExpirations(now)
		return cc.Context().Err() == nil
	})

	go func() {
		err := cc.Run()
		if err != nil {
			cfg.Errors(err)
		}
	}()

	return cc
}

// Send send a message.
func (c Client_struct) Send(path string, msgCode codes.Code, cf message.MediaType, payload io.ReadSeeker, opts ...message.Option) (*pool.Message, error) {
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
	return nil, errors.New("Invalid message code")
}

// Receive receives a message.
func (c Client_struct) Receive(path string, opts ...message.Option) (*client.Observation, error) {
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
