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
		log.Fatal(helpCmd)
	}
	help := strings.ToLower(os.Args[1])
	if help == "-h" || help == "--help" {
		log.Println(helpMsg)
		os.Exit(0)
	}

	code, err := parseCode(strings.ToUpper(os.Args[1]))
	if err != nil {
		log.Fatalf("Can't read request code: %s\n%s", err, helpCmd)
	}

	if len(os.Args) < 3 {
		log.Fatalf("CoAP URL must not be empty.\n%s", helpCmd)
	}
	path := os.Args[2]
	if strings.HasPrefix(path, "-") {
		log.Fatalf("Please enter a valid CoAP URL.\n%s", helpCmd)
	}

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

	if o == nil || (!*o) {
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
