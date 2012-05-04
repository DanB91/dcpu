// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package lem1802

import "github.com/jteeuwen/dcpu/cpu"

// DefaultPalette defines the default colour palette for the LEM1802 monitor.
var DefaultPalette = []cpu.Word{
	0x0000, // black
	0x0007, // dark blue
	0x0070, // dark green
	0x0077, // dark teal
	0x0700, // dark red
	0x0707, // dark purple
	0x0770, // dark yellow
	0x0555, // dark gray
	0x0aaa, // light gray
	0x000f, // blue
	0x00f0, // green
	0x00ff, // teal
	0x0f00, // red
	0x0f0f, // purple
	0x0ffF, // yellow
	0x0fff, // white
}
