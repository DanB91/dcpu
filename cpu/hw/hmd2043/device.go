// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package hmd2043

import (
	"github.com/jteeuwen/dcpu/cpu"
)

// Media represents a single media device that can be
// plugged into the HMD2043 drive.
type Media interface {
	cpu.Device
	WordsPerSector() cpu.Word
	SectorCount() cpu.Word
	WriteLocked() bool
}

// Implements the HMD2043 disk drive controller.
type HMD2043 struct {
	f     cpu.IntFunc
	id    cpu.Word
	media Media
}

// New creates and initializes a new device instance.
func New(f cpu.IntFunc) cpu.Device {
	h := new(HMD2043)
	h.f = f
	return h
}

func (h *HMD2043) Manufacturer() uint32 { return 0x21544948 }
func (h *HMD2043) Id() uint32           { return 0x74fa4cae }
func (h *HMD2043) Revision() uint16     { return 0x07c2 }

// Load loads new media into the drive.
// This simply swaps it out with any existing media.
func (h *HMD2043) Load(m Media) { h.media = m }

func (h *HMD2043) Handler(s *cpu.Storage) {
	switch s.A {
	case QueryMediaPresent:
		s.A, s.B = ErrorNone, 0

		if h.media == nil {
			s.A = ErrorNoMedia
			return
		}

		if isSupported(h.media) {
			s.B = 1
		}

	case QueryMediaParams:
		s.X = 0

		if h.media == nil {
			s.A = ErrorNoMedia
			s.B = 0
			s.C = 0
			return
		}

		s.B = h.media.WordsPerSector()
		s.C = h.media.SectorCount()

		if h.media.WriteLocked() {
			s.X = 1
		}

	case QueryDeviceFlags:

	case UpdateDeviceFlags:

	case QueryInterruptType:

	case SetInterruptId:

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
