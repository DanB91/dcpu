// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"log"
)

type ApiHandler func(*Client, []byte)

var handlers map[byte]ApiHandler

func Register(id byte, ah ApiHandler) {
	if handlers == nil {
		handlers = make(map[byte]ApiHandler)
	}
	handlers[id] = ah
}

// apiDispatch processes an incoming socket message and
// delagates it to the appropriate handler.
func apiDispatch(c *Client, in []byte) {
	id := in[0]
	in = in[1:]

	if s, ok := MessageStrings[id]; ok {
		log.Printf("%s [i] %s", c.conn.Request().RemoteAddr, s)
	}

	if ah, ok := handlers[id]; ok {
		ah(c, in)
	}
}
