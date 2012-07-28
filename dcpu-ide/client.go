// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
	"time"
)

type ClientHandler func(*Client)

const (
	ClientTimeout = time.Minute * 3
	SessionName   = "dcpu-ide"
)

// Client holds client data and a connection.
type Client struct {
	onClose []ClientHandler
	conn    *websocket.Conn
}

// NewClient creates a new client for the given connection.
func NewClient(conn *websocket.Conn) *Client {
	c := new(Client)
	c.conn = conn
	c.conn.SetDeadline(time.Now().Add(ClientTimeout))
	return c
}

// OnClose adds a new onclose handler for this client.
// All handlers are fired when the client connection is closed down
// for whatever reason.
func (c *Client) OnClose(f ClientHandler) {
	c.onClose = append(c.onClose, f)
}

// Close closes the client connection and cleans up resources.
func (c *Client) Close() (err error) {
	for _, f := range c.onClose {
		f(c)
	}

	c.onClose = nil

	if c.conn != nil {
		err = c.conn.Close()
		c.conn = nil
	}

	return
}

// Write sends the given data to the remote client.
func (c *Client) Write(argv ...byte) error {
	if c.conn == nil {
		return io.EOF
	}
	return websocket.Message.Send(c.conn, argv)
}

// Poll polls the connection for incoming data and dispatches
// it to the api handlers.
//
// It also ensures the connection deadline is moved forward on
// every message, so we do not disconnect prematurely.
func (c *Client) Poll() {
	var err error
	var msg []byte

	defer c.Close()

	for {
		err = websocket.Message.Receive(c.conn, &msg)
		if err != nil {
			if err != io.EOF {
				log.Printf("client.poll: %v", err)
			}
			break
		}

		if len(msg) == 0 {
			continue
		}

		if tracker != nil {
			tracker.Ping()
		}

		// Push the connection deadline forward.
		// Otherwise we will timeout when it isn't necessary.
		c.conn.SetDeadline(time.Now().Add(ClientTimeout))

		if msg[0] != ApiPing {
			apiDispatch(c, msg)
		}
	}
}
