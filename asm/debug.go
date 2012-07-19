// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import (
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/parser"
)

type FileInfo struct {
	Name      string
	StartAddr cpu.Word
}

type FuncInfo struct {
	Name      string
	StartAddr cpu.Word
	EndAddr   cpu.Word
	StartLine int
	EndLine   int
}

// SourceInfo defines file/line/col locations in original source.
type SourceInfo struct {
	File int
	Line int
	Col  int
}

// DebugInfo will map binary instructions to original source locations.
type DebugInfo struct {
	Files         []FileInfo   // List of files used to build the original source.
	Functions     []FuncInfo   // List of function descriptors.
	SourceMapping []SourceInfo // Binary <-> Source mappings. One entry per instruction.
}

func (d *DebugInfo) getStartAddr(fileidx int) cpu.Word {
	for i := range d.SourceMapping {
		if d.SourceMapping[i].File == fileidx {
			return cpu.Word(i)
		}
	}

	return 0
}

func (d *DebugInfo) SetFileDefs(files []string) {
	var i int

	if len(files) == 0 {
		return
	}

	d.Files = make([]FileInfo, len(files))

	for i = range files {
		d.Files[i].Name = files[i]
		d.Files[i].StartAddr = d.getStartAddr(i)
	}
}

func (d *DebugInfo) SetFunctionStart(addr cpu.Word, line int, name string) {
	d.Functions = append(d.Functions, FuncInfo{
		Name:      name,
		StartAddr: addr,
		StartLine: line,
	})
}

func (d *DebugInfo) SetFunctionEnd(addr cpu.Word, line int) {
	d.Functions[len(d.Functions)-1].EndAddr = addr
	d.Functions[len(d.Functions)-1].EndLine = line
}

// Emit emits one or more debug symbols.
func (d *DebugInfo) Emit(n ...parser.Node) {
	for i := range n {
		d.SourceMapping = append(d.SourceMapping, SourceInfo{n[i].File(), n[i].Line(), n[i].Col()})
	}
}
