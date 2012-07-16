// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This package implements a DCPU code profiler.
package prof

import (
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/cpu"
	"path"
)

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

// A Profile holds timing and execution information for a single test program.
type Profile struct {
	Files []string // List of source files associated with program.

	// Profile data per instruction.
	//
	// This stores a ProfileData instance for every PC value we encounter.
	// It may therefor contain multiple structures for the same opcode.
	Data []ProfileData

	funcs BlockList // List of function definitions.
	files BlockList // List of file definitions.
}

// New creates a new profile for the given code and debug data.
func New(code []cpu.Word, dbg *asm.DebugInfo) *Profile {
	var sym *asm.SourceInfo

	p := new(Profile)
	p.Files = dbg.Files
	p.Data = make([]ProfileData, len(code))

	for pc := range code {
		sym = dbg.SourceMapping[pc]

		pd := p.Data[pc]
		pd.Data = code[pc]
		pd.File = sym.File
		pd.Line = sym.Line
		pd.Col = sym.Col
		p.Data[pc] = pd
	}

	p.getInstructionSizes()
	return p
}

// Update updates the information for each instruction as it is executed.
func (p *Profile) Update(pc cpu.Word, s *cpu.Storage) {
	p.Data[pc].Count++
}

// UpdateCost alters the cumulative cost of a given instruction where necessary.
// This can happen when a branching instruction failed its check and had to
// be skipped. This increases its cost by an amount dependant on how many
// instructions where skipped. Nested branches can make this amount increase
// considerably.
func (p *Profile) UpdateCost(pc, cost cpu.Word) {
	p.Data[pc].Penalty += uint64(cost)
}

// getInstructionSizes computes and stores the size of each instruction.
// The size is the number of words the instruction occupies.
func (p *Profile) getInstructionSizes() {
	var pc, size cpu.Word

	for pc = 0; pc < cpu.Word(len(p.Data)); pc += size {
		size = cpu.Sizeof(cpu.Decode(p.Data[pc].Data))
		p.Data[pc].Size = size
	}
}

// countOpcodes counts the number of times the given opcode is used.
func (p *Profile) countOpcodes(opcode cpu.Word, extended bool) int {
	var c int
	var pc, op, a cpu.Word

	for pc = 0; pc < cpu.Word(len(p.Data)); pc += p.Data[pc].Size {
		op, a, _ = cpu.Decode(p.Data[pc].Data)

		if extended {
			if op == cpu.EXT && a == opcode {
				c++
			}
		} else if op == opcode {
			c++
		}
	}

	return c
}

func listContains(l []cpu.Word, w cpu.Word) bool {
	for i := range l {
		if l[i] == w {
			return true
		}
	}

	return false
}

// indexFunctions finds all function definitions and computes their
// start and end addresses.
func (p *Profile) indexFunctions() {

}

// indexFiles finds the start and end addresses of all code restricted
// to a specific source file.
func (p *Profile) indexFiles() {
	var start, pc cpu.Word
	var filename string
	var file int

	size := cpu.Word(len(p.Data))
	for pc = start; pc < size; pc++ {
		if p.Data[pc].File != file || pc == size-1 {
			_, filename = path.Split(p.Files[file])

			p.files = append(p.files, Block{
				Data:  p.Data[start:pc],
				Addr:  start,
				Label: filename,
			})

			start = pc
			file = p.Data[pc].File
		}
	}
}

// ListFiles yields the address ranges of all source files.
func (p *Profile) ListFiles() []Block {
	if len(p.files) == 0 {
		p.indexFiles()
	}
	return p.files
}

// ListFunctions yields the address ranges of all known functions.
func (p *Profile) ListFunctions() []Block {
	if len(p.funcs) == 0 {
		p.indexFunctions()
	}
	return p.funcs
}
