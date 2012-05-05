// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package keyboard

import (
	"github.com/jteeuwen/dcpu/cpu"
)

const (
	ClearBuffer = iota
	GetNextKey
	GetKeyState
	SetInterruptId
)

// Keyboard - Generic hardware Keyboard.
type Keyboard struct {
	buf  []cpu.Word
	keys []uint8
	int  cpu.IntFunc
	id   cpu.Word
}

// New creates and initializes a new device instance.
func New(f cpu.IntFunc) cpu.Device {
	k := new(Keyboard)
	k.int = f
	go k.poll()
	return k
}

func (k *Keyboard) Manufacturer() uint32 { return 0x0 }
func (k *Keyboard) Id() uint32           { return 0x30cf7406 }
func (k *Keyboard) Revision() uint16     { return 0x1 }

func (k *Keyboard) Handler(s *cpu.Storage) {
	switch s.A {
	case ClearBuffer:
		k.buf = k.buf[:0]

	case GetNextKey:
		s.C = 0
		if sz := len(k.buf); sz > 0 {
			s.C = k.buf[0]
			copy(k.buf, k.buf[1:])
			k.buf = k.buf[:sz-1]
		}

	case GetKeyState:
		s.C = 0
		if int(s.B) < len(k.keys) {
			s.C = cpu.Word(k.keys[s.B] & 1)
		}

	case SetInterruptId:
		k.id = s.B
	}
}

func (k *Keyboard) poll() {
	// TODO: Implement polling for input using something
	// like GLFW, Termbox, SDL, etc.
}
