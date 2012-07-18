// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"log"
	"os"
	"sync/atomic"
	"time"
)

// StateTracker tracks request intervals.
//
// It ensures we shut down the server when the client has
// not been sending us any keep-alive requests in a timely
// fashion.
//
// This is the only reliable way for us to know if the client
// closed the website.
type StateTracker struct {
	lastRequest int64
	timeout     int64
}

// NewStateTracker creates a new state tracker with the given timeout
// value in seconds.
func NewStateTracker(timeout uint) *StateTracker {
	s := new(StateTracker)
	s.timeout = int64(timeout) * int64(time.Second)
	s.lastRequest = time.Now().UnixNano()
	return s
}

// Ping is called on every single request.
// It atomically updates the last request time.
func (s *StateTracker) Ping() {
	atomic.StoreInt64(&s.lastRequest, time.Now().UnixNano())
}

// Poll runs in the background and compares the current
// time with the stored last request time. If we exceed
// the given timeout, this will radically shut down the server.
func (s *StateTracker) Poll() {
	for {
		select {
		case t := <-time.After(time.Second):
			a := t.UnixNano()
			b := atomic.LoadInt64(&s.lastRequest)

			if a-b >= s.timeout {
				log.Printf("Idle for %d second(s). Shutting down.",
					(a-b)/int64(time.Second))
				os.Exit(0)
			}
		}
	}
}
