package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	log     io.Writer
	conn    net.Conn
	cancel  context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer, log io.Writer, cancel context.CancelFunc) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		log:     log,
		cancel:  cancel,
	}
}

func (c *Client) Connect() error {
	if c.in == nil {
		return fmt.Errorf("incomming is absent")
	}
	if c.out == nil {
		return fmt.Errorf("outcomming is absent")
	}

	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("can't connection: %w", err)
	}

	c.conn = conn

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Send() error {
	defer c.cancel()

	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		return fmt.Errorf("can't send msg: %w", err)
	}

	return nil
}

func (c *Client) Receive() error {
	defer c.cancel()

	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return fmt.Errorf("can't recive msg: %w", err)
	}

	return nil
}
