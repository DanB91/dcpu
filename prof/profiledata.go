// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import "github.com/jteeuwen/dcpu/cpu"

// Cycle counts per opcode.
var opcodes = [...]uint8{
	// Basic opcodes
	0, 1, 2, 2, 2, 2, 3, 3, 3, 3, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 0, 0, 3, 3, 0, 0, 2, 2,

	// Extended opcodes.
	0, 2, 0, 0, 0, 0, 0, 0, 4, 1, 1, 3, 2, 0, 0, 0,
	2, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

// Cycle counts per operand.
var operands = [...]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1,
}

// Profile data for a specific opcode.
type ProfileData struct {
	Count uint64 // Number of times this opcode was called.

	// Number of additional cost points incurred at runtime.
	// This has to be independant of the usual instruction cost as the
	// DCPU spec employs different rules here. This is mostly relevant
	// for skipped branch instructions. They gain cost from the skipping.
	Penalty uint64

	File int      // Original source file.
	Line int      // Original source line.
	Col  int      // Original source column.
	Data cpu.Word // Encoded instruciton value.
	Size cpu.Word // Size of instruction.
}

// Copy creates a deep copy of the given data.
func (p *ProfileData) Copy() *ProfileData {
	return &ProfileData{
		Count:   p.Count,
		Penalty: p.Penalty,
		File:    p.File,
		Line:    p.Line,
		Col:     p.Col,
		Data:    p.Data,
		Size:    p.Size,
	}
}

// Cost returns the cycle cost for this entry.
func (p *ProfileData) Cost() uint8 {
	var c uint8

	op, a, b := cpu.Decode(p.Data)

	if op == cpu.EXT {
		c = opcodes[a]

		if b <= 0x1f {
			c += operands[b]
		}
	} else {
		c = opcodes[op]

		if a <= 0x1f {
			c += operands[a]
		}

		if b <= 0x1f {
			c += operands[b]
		}
	}

	return c
}

// CumulativeCost returns the cumulative cycle cost for this entry.
func (p *ProfileData) CumulativeCost() uint64 {
	return p.Count*uint64(p.Cost()) + p.Penalty
}
