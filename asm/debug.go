// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import (
	"bytes"
	"fmt"
	dp "github.com/jteeuwen/dcpu/parser"
)

// SourceInfo defines file/line/col locations in original source.
type SourceInfo struct {
	File int
	Line int
	Col  int
}

// DebugInfo will map binary instructions to original source locations.
type DebugInfo struct {
	Files []string      // List of files used to build the original source.
	Data  []*SourceInfo // Binary <-> Source mappings.
}

func NewDebugInfo(files []string) *DebugInfo {
	d := new(DebugInfo)
	d.Files = files
	return d
}

// Emit emits one or more debug symbols.
func (d *DebugInfo) Emit(n ...dp.Node) {
	for i := range n {
		d.Data = append(d.Data, &SourceInfo{n[i].File(), n[i].Line(), n[i].Col()})
	}
}

func (d *DebugInfo) String() string {
	var b bytes.Buffer

	for k, v := range d.Data {
		fmt.Fprintf(&b, "%04x: %s:%d:%d\n",
			k, d.Files[v.File], v.Line, v.Col)
	}

	return b.String()
}
