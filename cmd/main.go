// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/hex"
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
	coapmsg "github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/spf13/cobra"
)

var (
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
)

const verboseFmt = `Date: %s
Code: %s
Type: %s
Token: %s
Message-ID: %d
`

func main() {
	rootCmd := &cobra.Command{
		Use:   "coap-cli <method> <URL> [options]",
		Short: "CLI for CoAP",
	}

	getCmd := &cobra.Command{
		Use:   "get <url>",
		Short: "Perform a GET request on a COAP resource",
		Example: "coap-cli get channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -a 1e1017e6-dee7-45b4-8a13-00e6afeb66eb -H localhost -p 5683 -O 17,50 -o \n" +
			"coap-cli get channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb --host localhost --port 5683 --options 17,50 --observe",
		Run: runCmd(codes.GET),
	}
	getCmd.Flags().BoolVarP(&observe, "observe", "o", false, "Observe resource")

	putCmd := &cobra.Command{
		Use:   "put <url>",
		Short: "Perform a PUT request on a COAP resource",
		Example: "coap-cli put /test -H coap.me -p 5683 -c 50 -d 'hello, world'\n" +
			"coap-cli put /test --host coap.me --port 5683 --content-format 50 --data 'hello, world'",
		Run: runCmd(codes.PUT),
	}
	putCmd.Flags().StringVarP(&data, "data", "d", "", "Data")

	postCmd := &cobra.Command{
		Use:   "post <url>",
		Short: "Perform a POST request on a COAP resource",
		Example: "coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic -a 1e1017e6-dee7-45b4-8a13-00e6afeb66eb  -H localhost -p 5683 -c 50 -d 'hello, world'\n" +
			"coap-cli post channels/0bb5ba61-a66e-4972-bab6-26f19962678f/messages/subtopic  --auth 1e1017e6-dee7-45b4-8a13-00e6afeb66eb --host localhost --port 5683 --content-format 50 --data 'hello, world'",
		Run: runCmd(codes.POST),
	}
	postCmd.Flags().StringVarP(&data, "data", "d", "", "Data")

	deleteCmd := &cobra.Command{
		Use:   "delete <url>",
		Short: "Perform a DELETE request on a COAP resource",
		Example: "coap-cli delete /test -H coap.me -p 5683 -c 50 -d 'hello, world' -O 17,50\n" +
			"coap-cli delete /test --host coap.me --port 5683 --content-format 50 --data 'hello, world' --options 17,50",
		Run: runCmd(codes.DELETE),
	}
	deleteCmd.Flags().StringVarP(&data, "data", "d", "", "Data")

	rootCmd.AddCommand(getCmd, putCmd, postCmd, deleteCmd)
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "localhost", "Host")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "5683", "Port")
	rootCmd.PersistentFlags().StringVarP(&auth, "auth", "a", "", "Auth")
	rootCmd.PersistentFlags().IntVarP(&contentFormat, "content-format", "c", 50, "Content format")
	rootCmd.PersistentFlags().StringArrayVarP(&options, "options", "O", []string{}, "Add option num with contents of text to the request. If the text begins with 0x, then the hex text (two [0-9a-f] per byte) is converted to binary data.")
	rootCmd.PersistentFlags().Uint64VarP(&keepAlive, "keep-alive", "k", 0, "Send a ping after interval seconds of inactivity. If not specified (or 0), keep-alive is disabled (default).")
	rootCmd.PersistentFlags().Uint32VarP(&maxRetries, "max-retries", "m", 10, "Max retries for keep alive")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

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

func makeRequest(code codes.Code, args []string) {
	client, err := coap.NewClient(host+":"+port, keepAlive, maxRetries)
	if err != nil {
		log.Fatalf("Error coap creating client: %v", err)
	}

	var opts coapmsg.Options
	for _, optString := range options {
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
	if auth != "" {
		opts = append(opts, coapmsg.Option{ID: coapmsg.URIQuery, Value: []byte("auth=" + auth)})
	}
	if opts.HasOption(coapmsg.Observe) {
		if value, _ := opts.GetBytes(coapmsg.Observe); len(value) == 1 && value[0] == 0 && !observe {
			observe = true
		}
	}

	switch code {
	case codes.GET:
		switch {
		case observe:
			obs, err := client.Receive(args[0], verbose, opts...)
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
			res, err := client.Send(args[0], code, coapmsg.MediaType(contentFormat), nil, opts...)
			if err != nil {
				log.Fatalf("Error sending message: %v", err)
			}
			printMsg(res, verbose)
		}
	default:
		pld := strings.NewReader(data)
		res, err := client.Send(args[0], code, coapmsg.MediaType(contentFormat), pld, opts...)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		printMsg(res, verbose)
	}
}

func runCmd(code codes.Code) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintf(os.Stdout, color.YellowString("\nusage: %s\n\n"), cmd.Use)
			return
		}
		makeRequest(code, args)
	}
}
