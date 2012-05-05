// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package spc2000

import (
	"github.com/jteeuwen/dcpu/cpu"
)

// Keyboard - Generic hardware Keyboard.
type SPC2000 struct {
	int  cpu.IntFunc
	time uint64
	unit cpu.Word
}

// New creates and initializes a new device instance.
func New(f cpu.IntFunc) cpu.Device {
	spc := new(SPC2000)
	spc.int = f
	spc.unit = Days
	return spc
}

func (spc *SPC2000) Manufacturer() uint32 { return 0x1c6c8b36 }
func (spc *SPC2000) Id() uint32           { return 0x40e41d9d }
func (spc *SPC2000) Revision() uint16     { return 0x005e }

func (spc *SPC2000) Handler(s *cpu.Storage) {
	switch s.A {
	case GetStatus:
		s.C, s.B = spc.status()

	case SetSkipAmount:
		spc.time = uint64(s.Mem[s.B+0])<<48 | uint64(s.Mem[s.B+1])<<32 |
			uint64(s.Mem[s.B+2])<<16 | uint64(s.Mem[s.B+3])

	case TriggerSleepCycle:
		s.C, s.B = spc.status()

		// This is contrary to the spec, which is inconsistent.
		// It notes that GET_STATUS sets C to 1 if the device is ready.
		//
		// It then notes that TRIGGER_DEVICE should actually trigger
		// the device when C is 0.
		if s.C == 1 {
			go spc.trigger()
		}

	case SetSkipUnit:
		spc.unit = s.B
	}
}

// status should return the current device status.
//
// The implementation requires game mechanics that this emulator
// does not supply. Notably access to sensors that can detect certain
// physical world properties which determine if this chamber can be
// triggered or not. As such, we just supply the 'UnknownError' status.
func (spc *SPC2000) status() (ready, error cpu.Word) {
	return 0, UnknownError
}

// trigger triggers the sleep chamber.
//
// The implementation requires game mechanics that this emulator
// does not supply. Notably access to sensors that can detect certain
// physical world properties which determine if this chamber can be
// triggered or not.
func (spc *SPC2000) trigger() {

}
