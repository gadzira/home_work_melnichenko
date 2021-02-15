package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	var t string
	flag.StringVar(&t, "timeout", "10s", "Duration of connections")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("Mandatory arguments \"host\" and \"port\" not define")
	}

	targetSource := flag.Arg(0) + ":" + flag.Arg(1)
	timeout, err := time.ParseDuration(t)
	if err != nil {
		log.Fatalf("E! Can't parse the timeout duration: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	tn := NewTelnetClient(targetSource, timeout, os.Stdin, os.Stdout, os.Stderr, cancel)

	err = tn.Connect()
	if err != nil {
		log.Fatalf("E! Can't connect: %s", err)
	}
	defer tn.Close()

	go func() {
		err := tn.Receive()
		if err != nil {
			log.Fatalf("E! Can't receieve: %s", err)

			return
		}
	}()

	go func() {
		err := tn.Send()
		if err != nil {
			log.Fatalf("E! Can't send: %s", err)

			return
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
		close(sigint)
	}
}
