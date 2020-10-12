package main

import (
	"errors"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"

	gocoap "github.com/dustin/go-coap"
	coap "github.com/mainflux/coap-cli/coap"
)

const (
	get    = "GET"
	put    = "PUT"
	post   = "POST"
	delete = "DELETE"
)

func parseCode(code string) (gocoap.COAPCode, error) {
	switch code {
	case get:
		return gocoap.GET, nil
	case put:
		return gocoap.PUT, nil
	case post:
		return gocoap.POST, nil
	case delete:
		return gocoap.DELETE, nil
	}
	return 0, errors.New("Message can be GET, POST, PUT or DELETE")
}

func checkType(c, n, a, r *bool) (gocoap.COAPType, error) {
	arr := []bool{*c, *n, *a, *r}
	var counter int
	for _, v := range arr {
		if v {
			counter++
		}
	}
	if counter > 1 {
		return 0, errors.New("invalid message type")
	}
	switch {
	case *c:
		return gocoap.Confirmable, nil
	case *n:
		return gocoap.NonConfirmable, nil
	case *a:
		return gocoap.Acknowledgement, nil
	case *r:
		return gocoap.Reset, nil
	}
	return gocoap.Confirmable, nil
}

func printMsg(m *gocoap.Message) {
	if m != nil {
		log.Printf("\nMESSAGE:\nType: %d\nCode: %s\nMessageID: %d\nToken: %s\nPayload: %s\n",
			m.Type, m.Code.String(), m.MessageID, m.Token, m.Payload)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Message code must be GET, PUT, POST or DELETE")
	}
	code, err := parseCode(strings.ToUpper(os.Args[1]))
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	if len(os.Args) < 3 {
		log.Fatal("Please enter valid CoAP URL")
	}
	addr := os.Args[2]
	os.Args = os.Args[2:]

	c := flag.Bool("C", false, "Confirmable")
	n := flag.Bool("NC", false, "Non-confirmable")
	a := flag.Bool("ACK", false, "Acknowledgement")
	r := flag.Bool("RST", false, "Reset")
	o := flag.Bool("O", false, "Observe")
	// Default type is JSON.
	cf := flag.Int("CF", 50, "Content format")
	d := flag.String("d", "", "Message data")
	id := flag.Uint("id", 0, "Message ID")
	token := flag.String("token", "", "Message data")

	flag.Parse()

	t, err := checkType(c, n, a, r)
	if err != nil {
		log.Fatal("ERR TYPE: ", err)
	}
	address, err := url.Parse(addr)
	if err != nil {
		log.Fatal("ERR PARSING ADDR:", err)
	}
	client, err := coap.New(address)
	if err != nil {
		log.Fatal("ERROR CREATING CLIENT: ", err)
	}

	opts := coap.ParseOptions(address)
	if *o {
		opts = append(opts, coap.Option{
			ID:    gocoap.Observe,
			Value: 0,
		})
	}
	if *cf != 0 {
		opts = append(opts, coap.Option{
			ID:    gocoap.ContentFormat,
			Value: *cf,
		})
	}

	res, err := client.Send(t, code, uint16(*id), []byte(*token), []byte(*d), opts)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	printMsg(res)
	if res == nil {
		os.Exit(0)
	}
	switch res.Code {
	case gocoap.Forbidden, gocoap.BadRequest, gocoap.InternalServerError, gocoap.NotFound:
		log.Fatalf("Response code: %s", res.Code)
	}
	if *o {
		if code != gocoap.GET {
			log.Fatalln("Can observe non GET requests.")
		}
		msgs := make(chan *gocoap.Message)
		go func() {
			for {
				msg, err := client.Receive()
				if err != nil {
					log.Fatal("ERROR RECEIVING: ", err)
				}
				msgs <- msg
			}
		}()
		for {
			select {
			case m := <-msgs:
				printMsg(m)
			}
		}
	}
}
