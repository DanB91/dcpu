// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import "github.com/jteeuwen/dcpu/cpu"

// Profile data for a specific opcode.
type ProfileData struct {
	Count  uint64   // Number of times this opcode was called.
	File   int      // Original source file.
	Line   int      // Original source line.
	Col    int      // Original source column.
	Opcode cpu.Word // Opcode this data applies to.
	A, B   cpu.Word // Arguments for this instruction.
}

// Cost returns the cycle cost for this entry.
func (p *ProfileData) Cost() uint8 {
	var c uint8

	if p.Opcode == cpu.EXT {
		c = opcodes[p.A]

		if p.B <= 0x1f {
			c += operands[p.B]
		}
	} else {
		c = opcodes[p.Opcode]

		if p.A <= 0x1f {
			c += operands[p.A]
		}

		if p.B <= 0x1f {
			c += operands[p.B]
		}
	}

	return c
}

// CumulativeCount returns the cumulative cycle cost for this entry.
func (p *ProfileData) CumulativeCost() uint64 {
	return p.Count * uint64(p.Cost())
}
