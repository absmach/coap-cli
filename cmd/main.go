package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	coap "github.com/mainflux/coap-cli/coap"
	"github.com/plgd-dev/go-coap/v2/message"
	coapmsg "github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
)

const (
	get    = "GET"
	put    = "PUT"
	post   = "POST"
	delete = "DELETE"
)

func parseCode(code string) (codes.Code, error) {
	switch code {
	case get:
		return codes.GET, nil
	case put:
		return codes.PUT, nil
	case post:
		return codes.POST, nil
	case delete:
		return codes.DELETE, nil
	}
	return 0, errors.New("Message can be GET, POST, PUT or DELETE")
}

func printMsg(m *pool.Message) {
	if m != nil {
		log.Printf("\nMESSAGE:\n %v", m)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Message code must be GET, PUT, POST or DELETE")
	}
	code, err := parseCode(strings.ToUpper(os.Args[1]))
	if err != nil {
		log.Fatal("error: ", err)
	}
	if len(os.Args) < 3 {
		log.Fatal("Please enter valid CoAP URL")
	}
	path := os.Args[2]
	os.Args = os.Args[2:]

	o := flag.Bool("o", false, "Observe")
	h := flag.String("h", "localhost", "Host")
	p := flag.String("p", "5683", "Port")
	// Default type is JSON.
	cf := flag.Int("q", 50, "Content format")
	d := flag.String("d", "", "Message data")
	a := flag.String("auth", "", "Auth token")

	flag.Parse()

	client, err := coap.New(*h + ":" + *p)
	if err != nil {
		log.Fatal("Error creating client: ", err)
	}
	var opts coapmsg.Options
	if a != nil {
		opts = append(opts, coapmsg.Option{ID: coapmsg.URIQuery, Value: []byte(fmt.Sprintf("auth=%s", *a))})
	}

	if o == nil {
		pld := strings.NewReader(*d)

		res, err := client.Send(path, code, message.MediaType(*cf), pld, opts...)
		if err != nil {
			log.Fatal("Error sending message: ", err)
		}
		printMsg(res)
		return
	}
	if code != codes.GET {
		log.Fatal("Only GET requests accept observe option.")
	}
	obs, err := client.Receive(path, opts...)
	if err != nil {
		log.Fatal("Error observing resource: ", err)
	}
	errs := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	obs.Cancel(context.Background())
	log.Fatal("Observation terminated: ", err)
}
