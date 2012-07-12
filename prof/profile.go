// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This package implements a DCPU code profiler.
package prof

import (
	"fmt"
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/cpu"
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

	// List of function definitions. The address at which they start and end.
	funcs FuncList
}

// New creates a new profile for the given code and debug data.
func New(code []cpu.Word, dbg *asm.DebugInfo) *Profile {
	p := new(Profile)
	p.Files = dbg.Files
	p.Data = make([]ProfileData, len(code))

	for pc := range code {
		sym := dbg.SourceMapping[pc]

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

var fp = fmt.Printf

// findFuncAddresses finds all unique function address that have been called.
//
// XXX(jimt): This makes one important assumption; address references in a
// JSR instruction can not be null (0). This is done to prevent the routine
// from interpreting data sections as code and mistaking some data for a valid
// jump. This is rather flaky, but I know of no better solution at this point.
func (p *Profile) findFuncAddresses() []cpu.Word {
	var list []cpu.Word
	var pc, op, a, b, addr cpu.Word

	for pc = 0; pc < cpu.Word(len(p.Data)); pc += p.Data[pc].Size {
		op, a, b = cpu.Decode(p.Data[pc].Data)

		if op != cpu.EXT || a != cpu.JSR || b == 0 {
			continue
		}

		if b == 0x1f {
			addr = p.Data[pc+1].Data
		} else {
			addr = b - 0x21
		}

		if !listContains(list, addr) {
			list = append(list, addr)
		}
	}

	return list
}

// isBranch checks if the target address is a branching instruction.
func (p *Profile) isBranch(pc cpu.Word) bool {
	op, _, _ := cpu.Decode(p.Data[pc].Data)
	return op >= cpu.IFB && op <= cpu.IFU
}

// indexFunctions finds all function definitions and computes their
// start and end addresses.
//
// XXX(jimt): This makes some assumptions about code layout which may not
// necessarily be valid. Notably the one described in Profile.findFuncAddresses()
// Additionally, it assumes functions will always return by means of 
// a `SET PC, POP` instruction.
func (p *Profile) indexFunctions() {
	addresslist := p.findFuncAddresses()
	if len(addresslist) == 0 {
		return
	}

	p.funcs = make(FuncList, len(addresslist))

	var pc, op, a, b, size cpu.Word

	cap := cpu.Word(len(p.Data))
	for i := range p.funcs {
		pc = addresslist[i]
		p.funcs[i].Addr = pc
		size = p.Data[pc].Size
		pc += size

		// Find end address for this function.
		for ; pc < cap; pc += size {
			op, a, b = cpu.Decode(p.Data[pc].Data)

			// Do we have a 'return' statement (set pc, pop)?
			if op == cpu.SET && a == 0x1c && b == 0x18 {
				// Check if it is part of a branching instruction.
				// If not, consider it the end address of the function.
				if !p.isBranch(pc - size) {
					p.funcs[i].Data = p.Data[p.funcs[i].Addr : pc+1]
					break
				}
			}

			size = p.Data[pc].Size
		}
	}
}

// Functions yields the address ranges of all known functions.
func (p *Profile) Functions() []FuncDef {
	if len(p.funcs) == 0 {
		// Functions have not been indexed yet.
		p.indexFunctions()
	}
	return p.funcs
}
