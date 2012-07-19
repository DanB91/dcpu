// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This package implements a DCPU code profiler.
package prof

import (
	"fmt"
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/cpu"
	"os"
	"path"
	"strings"
)

// A Profile holds timing and execution information for a single test program.
type Profile struct {
	Files     []asm.FileInfo // List of source files associated with program.
	Functions []asm.FuncInfo // List of function definitions.

	// Profile data per instruction.
	//
	// This stores a ProfileData instance for every PC value we encounter.
	// It may therefor contain multiple structures for the same opcode.
	Data []ProfileData

	fileblocks BlockList
	funcblocks BlockList
}

// New creates a new profile for the given code and debug data.
func New(code []cpu.Word, dbg *asm.DebugInfo) *Profile {
	var sym asm.SourceInfo

	p := new(Profile)
	p.Files = dbg.Files
	p.Functions = dbg.Functions
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

func (p *Profile) ListFiles() BlockList {
	if len(p.fileblocks) == 0 {
		p.fileblocks = make(BlockList, len(p.Files))

		var s, e cpu.Word

		path := os.Getenv("DCPU_PATH")

		for i := range p.fileblocks {
			s = p.Files[i].Start

			if i < len(p.fileblocks)-1 {
				e = p.Files[i+1].Start
			} else {
				e = cpu.Word(len(p.Data) - 1)
			}

			p.fileblocks[i].Data = p.Data[s:e]
			p.fileblocks[i].StartAddr = s
			p.fileblocks[i].EndAddr = e
			p.fileblocks[i].Label = strings.Replace(p.Files[i].Name, path, "$DCPU_PATH", 1)
		}
	}

	return p.fileblocks
}

// getFileName finds the filename associated with a given function.
func (p *Profile) getFileName(addr cpu.Word) string {
	files := p.ListFiles()

	for i := range files {
		if addr >= files[i].StartAddr && addr < files[i].EndAddr {
			return files[i].Label
		}
	}

	return ""
}

func (p *Profile) ListFunctions() BlockList {
	if len(p.funcblocks) == 0 {
		p.funcblocks = make(BlockList, len(p.Functions))

		var s, e cpu.Word
		var fname string

		for i := range p.funcblocks {
			s = p.Functions[i].StartAddr
			e = p.Functions[i].EndAddr

			fname = p.getFileName(s)
			_, fname = path.Split(fname)

			p.funcblocks[i].Data = p.Data[s:e]
			p.funcblocks[i].StartAddr = s
			p.funcblocks[i].EndAddr = e
			p.funcblocks[i].StartLine = p.Functions[i].StartLine
			p.funcblocks[i].EndLine = p.Functions[i].EndLine
			p.funcblocks[i].Label = fmt.Sprintf("%s (%s)",
				p.Functions[i].Name, fname)
		}
	}

	return p.funcblocks
}
