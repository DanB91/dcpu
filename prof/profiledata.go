// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import "github.com/jteeuwen/dcpu/cpu"

// Profile data for a specific opcode.
type ProfileData struct {
	Count     uint64   // Number of times this opcode was called.
	File      int      // Original source file.
	Line      int      // Original source line.
	Col       int      // Original source column.
	Opcode    cpu.Word // Opcode this data applies to.
	A, B      cpu.Word // Arguments for this instruction.
	CycleCost uint8    // Number of cpu cycles this opcode consumed (including operands).
}

// CumulativeCount returns the cumulative cycle count for this entry.
func (p *ProfileData) CumulativeCount() uint64 {
	return p.Count * uint64(p.CycleCost)
}
