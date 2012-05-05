// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package hmd2043

import (
	"github.com/jteeuwen/dcpu/cpu"
)

// Media represents a single media device that can be
// plugged into the HMD2043 drive.
type Media interface {
	WordsPerSector() cpu.Word
	SectorCount() cpu.Word
	WriteLocked() bool
}

// Implements the HMD2043 disk drive controller.
type HMD2043 struct {
	int     cpu.IntFunc // Interrupt function we can call on the CPU.
	media   Media       // Media we are currently working on.
	id      cpu.Word    // Interrupt message we send to the CPU.
	flags   cpu.Word    // Device flags.
	lastint cpu.Word    // Last interrupt type we raised.
}

// New creates and initializes a new device instance.
func New(f cpu.IntFunc) cpu.Device {
	h := new(HMD2043)
	h.int = f
	h.flags = 0
	return h
}

func (h *HMD2043) Manufacturer() uint32 { return 0x21544948 }
func (h *HMD2043) Id() uint32           { return 0x74fa4cae }
func (h *HMD2043) Revision() uint16     { return 0x07c2 }

// Insert loads new media into the drive.
// This fails silently when media is already present.
//
// When the device flag MediaStatusInterrupt is set, this
// will trigger an interrupt.
func (h *HMD2043) Insert(m Media) {
	if h.media != nil {
		return
	}

	h.media = m

	if h.flags&MediaStatusInterrupt != 0 {
		h.int(h.id)
	}
}

// Eject unloads existing media from the drive.
// If no media is present, this fails silently.
//
// When the device flag MediaStatusInterrupt is set, this
// will trigger an interrupt.
func (h *HMD2043) Eject(m Media) {
	if h.media == nil {
		return
	}

	h.media = nil

	if h.flags&MediaStatusInterrupt != 0 {
		h.int(h.id)
	}
}

func (h *HMD2043) Handler(s *cpu.Storage) {
	switch s.A {
	case QueryMediaPresent:
		if h.media == nil {
			s.A = ErrorNoMedia
			s.B = 0
			return
		}

		s.A, s.B = ErrorNone, 0

		if isSupported(h.media) {
			s.B = 1
		}

	case QueryMediaParams:
		if h.media == nil {
			s.A, s.B, s.C, s.X = ErrorNoMedia, 0, 0, 0
			return
		}

		s.A = ErrorNone
		s.B = h.media.WordsPerSector()
		s.C = h.media.SectorCount()
		s.X = 0

		if h.media.WriteLocked() {
			s.X = 1
		}

	case QueryDeviceFlags:
		s.A, s.B = ErrorNone, h.flags

	case UpdateDeviceFlags:
		s.A, h.flags = ErrorNone, s.B

	case QueryInterruptType:

	case SetInterruptId:
		s.A, h.id = ErrorNone, s.B

	case ReadSectors:

	case WriteSectors:

	case QueryMediaQuality:

	}
}

// Returns true if the given media is supported by our drive.
//
// TODO: Find some metric to determine of the media is OK or not.
// The HMU1440 spec defines no manufacturer or device ids.
func isSupported(m Media) bool {
	return true
}
