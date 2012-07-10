// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import dp "github.com/jteeuwen/dcpu/parser"

// SourceInfo defines file/line/col locations in original source.
type SourceInfo struct {
	File int
	Line int
	Col  int
}

// DebugInfo will map binary instructions to original source locations.
type DebugInfo struct {
	Files         []string      // List of files used to build the original source.
	SourceMapping []*SourceInfo // Binary <-> Source mappings. One entry per instruction.
}

func NewDebugInfo(files []string) *DebugInfo {
	d := new(DebugInfo)
	d.Files = files
	return d
}

// Emit emits one or more debug symbols.
func (d *DebugInfo) Emit(n ...dp.Node) {
	for i := range n {
		d.SourceMapping = append(d.SourceMapping, &SourceInfo{n[i].File(), n[i].Line(), n[i].Col()})
	}
}
