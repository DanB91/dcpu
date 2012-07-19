// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
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
// a given timeout, this will send a termination signal
// down the returned channel.
func (s *StateTracker) Poll() <-chan struct{} {
	c := make(chan struct{})

	go func() {
		defer close(c)
		tick := time.Tick(time.Second)

		for {
			select {
			case t := <-tick:
				a := t.UnixNano()
				b := atomic.LoadInt64(&s.lastRequest)

				if a-b >= s.timeout {
					c <- struct{}{}
					return
				}
			}
		}
	}()

	return c
}
