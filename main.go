// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	coap "github.com/absmach/coap-cli/coap"
	"github.com/fatih/color"
	piondtls "github.com/pion/dtls/v3"
	coapmsg "github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/spf13/cobra"
)

const verboseFmt = `Date: %s
Code: %s
Type: %s
Token: %s
Message-ID: %d
`

func main() {
	req := &request{}

	rootCmd := &cobra.Command{
		Use:   "coap-cli <method> <URL> [options]",
		Short: "CLI for CoAP",
	}

	getCmd := &cobra.Command{
		Use:   "get <url>",
		Short: "Perform a GET request on a COAP resource",
		Example: "coap-cli get m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic -a 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -H localhost -p 5683 -O 17,50 -o \n" +
			"coap-cli get m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb --host localhost --port 5683 --options 17,50 --observe",
		Run: runCmd(req, codes.GET),
	}
	getCmd.Flags().BoolVarP(&req.observe, "observe", "o", false, "Observe resource")

	putCmd := &cobra.Command{
		Use:   "put <url>",
		Short: "Perform a PUT request on a COAP resource",
		Example: "coap-cli put /test -H coap.me -p 5683 -c 50 -d 'hello, world'\n" +
			"coap-cli put /test --host coap.me --port 5683 --content-format 50 --data 'hello, world'",
		Run: runCmd(req, codes.PUT),
	}
	putCmd.Flags().StringVarP(&req.data, "data", "d", "", "Data")

	postCmd := &cobra.Command{
		Use:   "post <url>",
		Short: "Perform a POST request on a COAP resource",
		Example: "coap-cli post m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic -a 1e1017e6-dee7-45b4-8a13-00e6afeb66eb  -H localhost -p 5683 -c 50 -d 'hello, world'\n" +
			"coap-cli post m/aa844fac-2f74-4ec3-8318-849b95d03bcc/c/0bb5ba61-a66e-4972-bab6-26f19962678f/subtopic  --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb --host localhost --port 5683 --content-format 50 --data 'hello, world'",
		Run: runCmd(req, codes.POST),
	}
	postCmd.Flags().StringVarP(&req.data, "data", "d", "", "Data")

	deleteCmd := &cobra.Command{
		Use:   "delete <url>",
		Short: "Perform a DELETE request on a COAP resource",
		Example: "coap-cli delete /test -H coap.me -p 5683 -c 50 -d 'hello, world' -O 17,50\n" +
			"coap-cli delete /test --host coap.me --port 5683 --content-format 50 --data 'hello, world' --options 17,50",
		Run: runCmd(req, codes.DELETE),
	}
	deleteCmd.Flags().StringVarP(&req.data, "data", "d", "", "Data")

	rootCmd.AddCommand(getCmd, putCmd, postCmd, deleteCmd)
	rootCmd.PersistentFlags().StringVarP(&req.host, "host", "H", "localhost", "Host")
	rootCmd.PersistentFlags().StringVarP(&req.port, "port", "p", "5683", "Port")
	rootCmd.PersistentFlags().StringVarP(&req.auth, "auth", "a", "", "Auth")
	rootCmd.PersistentFlags().IntVarP(&req.contentFormat, "content-format", "c", 50, "Content format")
	rootCmd.PersistentFlags().StringArrayVarP(&req.options, "options", "O", []string{}, "Add option num with contents of text to the request. If the text begins with 0x, then the hex text (two [0-9a-f] per byte) is converted to binary data.")
	rootCmd.PersistentFlags().Uint64VarP(&req.keepAlive, "keep-alive", "k", 0, "Send a ping after interval seconds of inactivity. If not specified (or 0), keep-alive is disabled (default).")
	rootCmd.PersistentFlags().Uint32VarP(&req.maxRetries, "max-retries", "m", 10, "Max retries for keep alive")
	rootCmd.PersistentFlags().BoolVarP(&req.verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&req.certFile, "cert-file", "C", "", "Client certificate file")
	rootCmd.PersistentFlags().StringVarP(&req.keyFile, "key-file", "K", "", "Client key file")
	rootCmd.PersistentFlags().StringVarP(&req.clientCAFile, "ca-file", "A", "", "Client CA file")
	rootCmd.PersistentFlags().BoolVarP(&req.useDTLS, "use DTLS", "s", false, "Use DTLS")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func printMsg(m *pool.Message, verbose bool) {
	if m != nil && verbose {
		fmt.Printf(verboseFmt,
			time.Now().Format(time.RFC1123),
			m.Code(),
			m.Type(),
			m.Token(),
			m.MessageID())
		cf, err := m.ContentFormat()
		if err == nil {
			fmt.Printf("Content-Format: %s \n", cf.String())
		}
		bs, err := m.BodySize()
		if err == nil {
			fmt.Printf("Content-Length: %d\n", bs)
		}
	}
	body, err := m.ReadBody()
	if err != nil {
		log.Fatalf("failed to read body %v", err)
	}
	if len(body) > 0 {
		fmt.Printf("\n%s\n", string(body))
	}
}

func makeRequest(req *request, args []string) {
	dtlsConfig, err := req.createDTLSConfig()
	if err != nil {
		log.Fatalf("Error creating DTLS config: %v", err)
	}
	client, err := coap.NewClient(req.host+":"+req.port, req.keepAlive, req.maxRetries, dtlsConfig)
	if err != nil {
		log.Fatalf("Error coap creating client: %v", err)
	}

	var opts coapmsg.Options
	for _, optString := range req.options {
		opt := strings.Split(optString, ",")
		if len(opt) < 2 {
			log.Fatal("Invalid option format")
		}
		optId, err := strconv.ParseUint(opt[0], 10, 16)
		if err != nil {
			log.Fatal("Error parsing option id")
		}
		if strings.HasPrefix(opt[1], "0x") {
			op := strings.TrimPrefix(opt[1], "0x")
			optValue, err := hex.DecodeString(op)
			if err != nil {
				log.Fatal("Invalid option value ", err.Error())
			}
			opts = append(opts, coapmsg.Option{ID: coapmsg.OptionID(optId), Value: optValue})
		} else {
			opts = append(opts, coapmsg.Option{ID: coapmsg.OptionID(optId), Value: []byte(opt[1])})
		}
	}
	if req.auth != "" {
		opts = append(opts, coapmsg.Option{ID: coapmsg.URIQuery, Value: []byte("auth=" + req.auth)})
	}
	if opts.HasOption(coapmsg.Observe) {
		if value, _ := opts.GetBytes(coapmsg.Observe); len(value) == 1 && value[0] == 0 && !req.observe {
			req.observe = true
		}
	}

	switch req.code {
	case codes.GET:
		switch {
		case req.observe:
			obs, err := client.Receive(args[0], req.verbose, opts...)
			if err != nil {
				log.Fatalf("Error observing resource: %v", err)
			}
			errs := make(chan error, 1)
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT)
				errs <- fmt.Errorf("%s", <-c)
			}()

			err = <-errs
			if err := obs.Cancel(context.Background(), opts...); err != nil {
				log.Fatalf("Error cancelling observation: %v", err)
			}
			log.Fatalf("Observation terminated: %v", err)
		default:
			res, err := client.Send(args[0], req.code, coapmsg.MediaType(req.contentFormat), nil, opts...)
			if err != nil {
				log.Fatalf("Error sending message: %v", err)
			}
			printMsg(res, req.verbose)
		}
	default:
		pld := strings.NewReader(req.data)
		res, err := client.Send(args[0], req.code, coapmsg.MediaType(req.contentFormat), pld, opts...)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		printMsg(res, req.verbose)
	}
}

func runCmd(req *request, code codes.Code) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintf(os.Stdout, color.YellowString("\nusage: %s\n\n"), cmd.Use)
			return
		}
		req.code = code
		makeRequest(req, args)
	}
}

type request struct {
	code          codes.Code
	host          string
	port          string
	contentFormat int
	auth          string
	observe       bool
	data          string
	options       []string
	keepAlive     uint64
	verbose       bool
	maxRetries    uint32
	certFile      string
	keyFile       string
	clientCAFile  string
	useDTLS       bool
}

func (r *request) createDTLSConfig() (*piondtls.Config, error) {
	if !r.useDTLS {
		return nil, nil
	}
	dc := &piondtls.Config{}
	if r.certFile != "" && r.keyFile != "" {
		cert, err := tls.LoadX509KeyPair(r.certFile, r.keyFile)
		if err != nil {
			return nil, errors.Join(errors.New("failed to load certificates"), err)
		}
		dc.Certificates = []tls.Certificate{cert}
	}
	rootCA, err := loadCertFile(r.clientCAFile)
	if err != nil {
		return nil, errors.Join(errors.New("failed to load Client CA"), err)
	}
	if len(rootCA) > 0 {
		if dc.RootCAs == nil {
			dc.RootCAs = x509.NewCertPool()
		}
		if !dc.RootCAs.AppendCertsFromPEM(rootCA) {
			return nil, errors.New("failed to append root ca tls.Config")
		}
	}
	return dc, nil
}

func loadCertFile(certFile string) ([]byte, error) {
	if certFile != "" {
		return os.ReadFile(certFile)
	}
	return []byte{}, nil
}
