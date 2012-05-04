// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package keyboard

import (
	"github.com/jteeuwen/dcpu/cpu"
)

// Keyboard - Generic hardware Keyboard.
type Keyboard struct {
	f cpu.IntFunc
}

// New creates and initializes a new device instance.
func New(f cpu.IntFunc) cpu.Device {
	return &Keyboard{
		f: f,
	}
}

func (d *Keyboard) Manufacturer() uint32   { return 0x0 }
func (d *Keyboard) Id() uint32             { return 0x30cf7406 }
func (d *Keyboard) Revision() uint16       { return 0x1 }
func (d *Keyboard) Handler(s *cpu.Storage) {}
