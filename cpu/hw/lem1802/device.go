// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package lem1802

import (
	"github.com/jteeuwen/dcpu/cpu"
)

// Interrupt messages supported by this device.
const (
	MemMapScreen = iota
	MemMapFont
	MemMapPalette
	SetBorderColor
	MemDumpFont
	MemDumpPalette
)

const (
	ScreenSize  = 384 // Screen size in words.
	FontSize    = 256 // Font size in words.
	PaletteSize = 16  // Size of palette in words.
)

// LEM1802 - Low Energy Monitor.
// http://dcpu.com/highnerd/lem1802.txt
type Lem1802 struct {
	f       cpu.IntFunc
	buffer  []cpu.Word
	font    []cpu.Word
	palette []cpu.Word
	border  cpu.Word
}

// New creates and initializes a new device instance.
func New(f cpu.IntFunc) cpu.Device {
	return &Lem1802{
		f:       f,
		font:    DefaultFont(),
		palette: DefaultPalette,
		border:  0,
	}
}

func (d *Lem1802) Manufacturer() uint32 { return 0x1c6c8b36 }
func (d *Lem1802) Id() uint32           { return 0x7349f615 }
func (d *Lem1802) Revision() uint16     { return 0x1802 }

func (d *Lem1802) Handler(s *cpu.Storage) {
	switch s.A {
	case MemMapScreen:
		if s.B == 0 {
			d.buffer = nil
		} else {
			d.buffer = s.Mem[s.B : s.B+ScreenSize]
		}

	case MemMapFont:
		if s.B == 0 {
			d.font = DefaultFont()
		} else {
			d.font = s.Mem[s.B : s.B+FontSize]
		}

	case MemMapPalette:
		if s.B == 0 {
			d.palette = DefaultPalette
		} else {
			d.palette = s.Mem[s.B : s.B+PaletteSize]
		}

	case SetBorderColor:
		d.border = s.B & 0xf

	case MemDumpFont:
		copy(s.Mem[s.B:], DefaultFont())

	case MemDumpPalette:
		copy(s.Mem[s.B:], DefaultPalette)

	}
}

// decode decodes character/colour values from the given word.
func decode(w cpu.Word) (ch, blink, fg, bg cpu.Word) {
	return w & 0x7f, (w >> 7) & 1, (w >> 8) & 0xf, (w >> 12) & 0xf
}
