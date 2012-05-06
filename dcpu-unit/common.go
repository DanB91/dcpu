// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import dp "github.com/jteeuwen/dcpu/parser"

// Instructions and registers according to spec v1.7
//
// With additional 'test' and 'exit' instructions for this tool.

var (
	instructions = []string{
		"set", "add", "sub", "mul", "mli", "div", "dvi", "mod", "mdi", "and",
		"bor", "xor", "shr", "asr", "shl", "ifb", "ifc", "ife", "ifn", "ifg",
		"ifa", "ifl", "ifu", "adx", "sbx", "sti", "std", "jsr", "int", "iag",
		"ias", "rfi", "iaq", "hwn", "hwq", "hwi", "dat", "test", "exit",
	}

	regs = []string{
		"a", "b", "c", "x", "y", "z", "i", "j", "ia",
		"ex", "peek", "push", "pop", "pc", "sp",
	}
)

// isRegister returns true if the given value qualifies as
// a valid register name.
func isRegister(v string) bool {
	for i := range regs {
		if regs[i] == v {
			return true
		}
	}
	return false
}

// isInstruction returns true if the given value qualifies as
// a valid instruction name.
func isInstruction(v string) bool {
	for i := range instructions {
		if instructions[i] == v {
			return true
		}
	}
	return false
}

// findLabels recursively finds Label nodes.
func findLabels(n []dp.Node, l *[]*dp.Label) {
	for i := range n {
		switch tt := n[i].(type) {
		case dp.NodeCollection:
			findLabels(tt.Children(), l)

		case *dp.Label:
			*l = append(*l, tt)
		}
	}
}

// findReferences recursively finds Label references.
func findReferences(n []dp.Node, l *[]*dp.Name) {
	for i := range n {
		switch tt := n[i].(type) {
		case dp.NodeCollection:
			findReferences(tt.Children(), l)

		case *dp.Name:
			if isRegister(tt.Data) || isInstruction(tt.Data) {
				continue
			}

			*l = append(*l, tt)
		}
	}
}

// stripDuplicateNames removes duplicate entries from the given Name list.
func stripDuplicateNames(r []*dp.Name) []*dp.Name {
	l := make([]*dp.Name, 0, len(r))

	for i := range r {
		if !containsName(l, r[i].Data) {
			l = append(l, r[i])
		}
	}

	return l
}

// containsName returns true if the given list contains the supplied Name.
func containsName(r []*dp.Name, data string) bool {
	for i := range r {
		if r[i].Data == data {
			return true
		}
	}
	return false
}
