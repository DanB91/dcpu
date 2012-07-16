// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package util

import "github.com/jteeuwen/dcpu/parser"

// FindLabels finds Label all nodes.
func FindLabels(n []parser.Node, l *[]*parser.Label) {
	for i := range n {
		switch tt := n[i].(type) {
		case parser.NodeCollection:
			FindLabels(tt.Children(), l)

		case *parser.Label:
			*l = append(*l, tt)
		}
	}
}

// FindReferences finds all Label references.
func FindReferences(n []parser.Node, l *[]*parser.Name) {
	for i := range n {
		switch tt := n[i].(type) {
		case parser.NodeCollection:
			FindReferences(tt.Children(), l)

		case *parser.Name:
			if parser.IsRegister(tt.Data) || parser.IsInstruction(tt.Data) {
				break
			}

			*l = append(*l, tt)
		}
	}
}

// FindFunctions finds all function definitions.
func FindFunctions(n []parser.Node, l *[]*parser.Function) {
	for i := range n {
		switch tt := n[i].(type) {
		case *parser.Function:
			*l = append(*l, tt)

		case parser.NodeCollection:
			FindFunctions(tt.Children(), l)
		}
	}
}

// FindConstants finds all constant names.
func FindConstants(n []parser.Node, l *[]*parser.Name) {
	for i := range n {
		switch tt := n[i].(type) {
		case *parser.Instruction:
			list := tt.Children()
			name := list[0].(*parser.Name)

			if name.Data != "equ" {
				break
			}

			expr := list[1].(*parser.Expression)
			name = expr.Children()[0].(*parser.Name)
			*l = append(*l, name)

		case parser.NodeCollection:
			FindConstants(tt.Children(), l)
		}
	}
}

func containsString(in []string, data string) bool {
	for i := range in {
		if in[i] == data {
			return true
		}
	}
	return false
}

func containsLabel(in []*parser.Label, v string) bool {
	for i := range in {
		if in[i].Data == v {
			return true
		}
	}
	return false
}

func containsName(in []*parser.Name, v string) bool {
	for i := range in {
		if in[i].Data == v {
			return true
		}
	}
	return false
}

func containsFunction(in []*parser.Function, v string) bool {
	var name *parser.Name
	var list []parser.Node
	var lbl *parser.Label
	var ok bool

	for i := range in {
		list = in[i].Children()

		if name, ok = list[0].(*parser.Name); ok {
			if name.Data == v {
				return true
			}
			continue
		}

		if lbl, ok = list[0].(*parser.Label); ok {
			if lbl.Data == v {
				return true
			}
		}
	}
	return false
}
