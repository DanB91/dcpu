// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/parser"
	"github.com/jteeuwen/dcpu/prof"
	"path"
)

var labelCache map[cpu.Word]string

func init() {
	labelCache = make(map[cpu.Word]string)
}

// getLabel attempts to find the original label definition for a
// given instruction.
//
// We use this in places where we need to know the original name of
// a function that was called. Information like this is not retained in
// profile data, because labels are translated to physical memory addresses
// in the assembly process. They do not persist into the final program.
// We therefor have to parse the original source into an AST again and
// find it manually.
func getLabel(p *prof.Profile, pc cpu.Word) string {
	if label, ok := labelCache[pc]; ok {
		return label
	}

	file := p.Files[p.Data[pc].File]
	line := p.Data[pc].Line
	ast, err := GetAST(file)

	if err != nil {
		labelCache[pc] = ""
		return ""
	}

	node := findLabelNode(ast.Root, line, p.Data[pc].Col)
	if node == nil {
		labelCache[pc] = ""
		return ""
	}

	var name string
	lbl, ok := node.(*parser.Label)

	if !ok {
		name = fmt.Sprintf("%04x", pc)
	} else {
		name = lbl.Data
	}

	for len(name) < 20 {
		name = name + " "
	}

	_, file = path.Split(file)

	labelCache[pc] = fmt.Sprintf("%s %s:%d", name, file, line)

	return labelCache[pc]
}

func findLabelNode(n parser.Node, line, col int) parser.Node {
	if n.Line() == line && n.Col() == col {
		return n
	}

	switch tt := n.(type) {
	case parser.NodeCollection:
		list := tt.Children()
		for i, v := range list {
			if n = findLabelNode(v, line, col); n != nil {
				if i == 0 {
					return n
				}

				// Return node preceeding this node.
				// We will assume it is a label definition.
				return list[i-1]
			}
		}
	}

	return nil
}
