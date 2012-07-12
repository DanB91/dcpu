// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/parser"
	"github.com/jteeuwen/dcpu/prof"
	"path"
	"strings"
)

type LabelDef struct {
	Label string
	PC    cpu.Word
}

type LabelCache []LabelDef

func (l LabelCache) IndexOfLabel(v string) int {
	for i := range l {
		if strings.EqualFold(l[i].Label, v) {
			return i
		}
	}
	return -1
}

func (l LabelCache) IndexOfAddress(v cpu.Word) int {
	for i := range l {
		if l[i].PC == v {
			return i
		}
	}
	return -1
}

var labelCache LabelCache

// GetLabel attempts to find the original label definition for a
// given instruction.
//
// We use this in places where we need to know the original name of
// a function that was called. Information like this is not retained in
// profile data, because labels are translated to physical memory addresses
// in the assembly process. They do not persist into the final program.
// We therefor have to parse the original source into an AST again and
// find it manually.
func GetLabel(p *prof.Profile, pc cpu.Word) string {
	if idx := labelCache.IndexOfAddress(pc); idx > -1 {
		return labelCache[idx].Label
	}

	file := p.Files[p.Data[pc].File]
	line := p.Data[pc].Line
	ast, err := GetAST(file)

	if err != nil {
		labelCache = append(labelCache, LabelDef{"", pc})
		return ""
	}

	node := findLabelNode(ast.Root, line, p.Data[pc].Col)
	if node == nil {
		labelCache = append(labelCache, LabelDef{"", pc})
		return ""
	}

	var name string
	lbl, ok := node.(*parser.Label)

	if !ok {
		labelCache = append(labelCache, LabelDef{"", pc})
		return ""
	} else {
		name = lbl.Data
	}

	for len(name) < 20 {
		name = name + " "
	}

	_, file = path.Split(file)
	name = fmt.Sprintf("%s %s:%d", name, file, line)

	if labelCache.IndexOfLabel(name) == -1 {
		labelCache = append(labelCache, LabelDef{name, pc})
	}

	return name
}

func findLabelNode(n parser.Node, line, col int) parser.Node {
	if n.Line() == line && n.Col() == col {
		return n
	}

	switch tt := n.(type) {
	case parser.NodeCollection:
		list := tt.Children()
		for i, v := range list {
			n = findLabelNode(v, line, col)
			if n != nil {
				// Return node preceeding this node.
				// We will assume it is a label definition.
				return list[i-1]
			}
		}
	}

	return nil
}
