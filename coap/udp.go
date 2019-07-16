package coap

import (
	"errors"
	"net"
	"net/url"
	"time"

	gocoap "github.com/dustin/go-coap"
)

const (
	network   = "udp"
	maxPktLen = 65536
	defPort   = ":5683"
)

var errInvalidScheme = errors.New("Invalid porotcol scheme")

type conn struct {
	conn *net.UDPConn
	buf  []byte
}

func parseAddr(addr *url.URL) (string, error) {
	if addr.Scheme != "coap" {
		return "", errInvalidScheme
	}
	var a, p string
	if addr.Port() == "" {
		p = defPort
	}
	a = addr.Host
	return a + p, nil
}

// Dial connects a CoAP client.
func dial(addr string) (*conn, error) {
	uaddr, err := net.ResolveUDPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	s, err := net.DialUDP(network, nil, uaddr)
	if err != nil {
		return nil, err
	}

	return &conn{s, make([]byte, maxPktLen)}, nil
}

// Send a message.  Get a response if there is one.
func (c *conn) send(req gocoap.Message) (*gocoap.Message, error) {
	err := transmit(c.conn, nil, req)

	if err != nil {
		return nil, err
	}

	if !req.IsConfirmable() {
		return nil, nil
	}

	rv, err := receive(c.conn, c.buf)
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

// Receive a message.
func (c *conn) receive() (*gocoap.Message, error) {
	rv, err := receive(c.conn, c.buf)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// Transmit a message.
func transmit(l *net.UDPConn, a *net.UDPAddr, m gocoap.Message) error {
	d, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	if a == nil {
		_, err = l.Write(d)
	} else {
		_, err = l.WriteTo(d, a)
	}
	return err
}

// Receive a message.
func receive(conn *net.UDPConn, buf []byte) (gocoap.Message, error) {
	conn.SetReadDeadline(time.Now().Add(time.Minute * 5))

	nr, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return gocoap.Message{}, err
	}
	return gocoap.ParseMessage(buf[:nr])
}
