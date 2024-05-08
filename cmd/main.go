// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

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

	coap "github.com/absmach/coap-cli/coap"
	"github.com/plgd-dev/go-coap/v2/message"
	coapmsg "github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
)

const (
	helpCmd = `Use "coap-cli --help" for help.`
	helpMsg = `
Usage: coap-cli <method> <URL> [options]
mathod: get, put, post or delete
-o      observe   option - only valid with GET request (default: false)
-auth   auth option sent as URI Query                  (default: "")
-h      host                                           (default: "localhost")
-p      port                                           (default: "5683")
-d      data to be sent in POST or PUT                 (default: "")
-cf     content format                                 (default: 50 - JSON format))

Examples:
coap-cli get channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -o
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world"
coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -d "hello world" -h 0.0.0.0 -p 1234
`
)

var (
	errCreateClient     = errors.New("failed to create client")
	errSendMessage      = errors.New("failed to send message")
	errInvalidObsOpt    = errors.New("invalid observe option")
	errFailedObserve    = errors.New("failed to observe resource")
	errTerminatedObs    = errors.New("observation terminated")
	errCancelObs        = errors.New("failed to cancel observation")
	errCodeNotSupported = errors.New("message can be GET, POST, PUT or DELETE")
)

type request struct {
	code codes.Code
	path string
	host *string
	port *string
	cf   *int
	data *string
	auth *string
	obs  *bool
}

func parseCode(code string) (codes.Code, error) {
	ret, err := codes.ToCode(code)
	if err != nil {
		return 0, err
	}
	switch ret {
	case codes.GET, codes.POST, codes.PUT, codes.DELETE:
		return ret, nil
	default:
		return 0, errCodeNotSupported
	}
}

func printMsg(m *pool.Message) {
	if m != nil {
		log.Printf("\nMESSAGE:\n %v", m)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(helpCmd)
	}
	help := strings.ToLower(os.Args[1])
	if help == "-h" || help == "--help" {
		log.Print(helpMsg)
		return
	}
	req := request{}
	var err error
	req.code, err = parseCode(strings.ToUpper(os.Args[1]))
	if err != nil {
		log.Fatalf("Can't read request code: %s\n%s", err, helpCmd)
	}

	if len(os.Args) < 3 {
		log.Fatalf("CoAP URL must not be empty.\n%s", helpCmd)
	}
	req.path = os.Args[2]
	if strings.HasPrefix(req.path, "-") {
		log.Fatalf("Please enter a valid CoAP URL.\n%s", helpCmd)
	}

	os.Args = os.Args[2:]
	req.obs = flag.Bool("o", false, "Observe")
	req.host = flag.String("h", "localhost", "Host")
	req.port = flag.String("p", "5683", "Port")
	// Default type is JSON.
	req.cf = flag.Int("cf", 50, "Content format")
	req.data = flag.String("d", "", "Message data")
	req.auth = flag.String("auth", "", "Auth token")
	flag.Parse()

	if err = makeRequest(req); err != nil {
		log.Fatal(err)
	}
}

func makeRequest(req request) error {
	client, err := coap.New(*req.host + ":" + *req.port)
	if err != nil {
		return errors.Join(errCreateClient, err)
	}
	var opts coapmsg.Options
	if req.auth != nil {
		opts = append(opts, coapmsg.Option{ID: coapmsg.URIQuery, Value: []byte(fmt.Sprintf("auth=%s", *req.auth))})
	}

	if req.obs == nil || (!*req.obs) {
		pld := strings.NewReader(*req.data)

		res, err := client.Send(req.path, req.code, message.MediaType(*req.cf), pld, opts...)
		if err != nil {
			return errors.Join(errSendMessage, err)
		}
		printMsg(res)
		return nil
	}
	if req.code != codes.GET {
		return errInvalidObsOpt
	}
	obs, err := client.Receive(req.path, opts...)
	if err != nil {
		return errors.Join(errFailedObserve, err)
	}

	errs := make(chan error, 1)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT)

		sig := <-sigChan
		errs <- fmt.Errorf("%v", sig)
	}()

	err = <-errs
	if err != nil {
		return errors.Join(errTerminatedObs, err)
	}
	if err := obs.Cancel(context.Background()); err != nil {
		return errors.Join(errCancelObs, err)
	}
	return nil
}
