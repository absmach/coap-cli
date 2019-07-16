package main

import (
	"errors"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"

	coap "github.com/dusanb94/coap-cli/coap"

	gocoap "github.com/dustin/go-coap"
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

func main() {
	code, err := parseCode(strings.ToUpper(os.Args[1]))
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	addr := os.Args[2]
	os.Args = os.Args[2:]

	log.Println("CODE:", code)
	c := flag.Bool("C", false, "CONFIRMABLE")
	n := flag.Bool("NC", false, "NON-CONFIRMABLE")
	a := flag.Bool("ACK", false, "ACKNOWLEDGEMENT")
	r := flag.Bool("RST", false, "RESET")

	o := flag.Bool("O", false, "Observe")

	d := flag.String("d", "", "Message data")
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

	res, err := client.Send(t, code, 12, nil, []byte(*d), opts)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	log.Println(res)
	if *o {
		msgs := make(chan *gocoap.Message)
		go func() {
			for {
				msg, err := client.Receive()
				if err != nil {
					log.Fatal("ERROR RECEIVING: ", err)
					return
				}
				msgs <- msg
			}
		}()
		for {
			select {
			case m := <-msgs:
				log.Println(string(m.Payload))
			}
		}
	}
}
