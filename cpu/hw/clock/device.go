// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package clock

import (
	"github.com/jteeuwen/dcpu/cpu"
	"time"
)

// Known interrupt messages.
const (
	SetInterval = iota
	GetTicks
	SetInterruptId
)

// Clock - Generic hardware clock.
type Clock struct {
	int     cpu.IntFunc // Interrupt function we can call on the CPU.
	ticker *time.Ticker
	ticks  cpu.Word
	id     cpu.Word
}

// New creates and initializes a new device instance.
func New(f cpu.IntFunc) cpu.Device {
	c := new(Clock)
	c.int = f
	c.ticker = time.NewTicker(time.Duration(1e9) / 60)
	go c.poll()
	return c
}

func (c *Clock) Manufacturer() uint32 { return 0x0 }
func (c *Clock) Id() uint32           { return 0x12d0b402 }
func (c *Clock) Revision() uint16     { return 0x1 }

func (c *Clock) Handler(s *cpu.Storage) {
	switch s.A {
	case SetInterval:
		c.ticker.Stop()
		c.ticks = 0

		if s.B > 0 {
			c.ticker = time.NewTicker(
				time.Duration(1e9) / (60 / time.Duration(s.B)))
		}

	case GetTicks:
		s.C = c.ticks

	case SetInterruptId:
		c.id = s.B
	}
}

// poll waits for ticks and optionally sends interrupt messages.
func (c *Clock) poll() {
	for {
		select {
		case <-c.ticker.C:
			c.ticks++

			if c.id > 0 {
				c.int(c.id)
			}
		}
	}
}
