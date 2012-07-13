// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package util

import "github.com/jteeuwen/dcpu/parser"

// Instructions and registers according to spec v1.7

var (
	instructions = []string{
		"set", "add", "sub", "mul", "mli", "div", "dvi", "mod", "mdi", "and",
		"bor", "xor", "shr", "asr", "shl", "ifb", "ifc", "ife", "ifn", "ifg",
		"ifa", "ifl", "ifu", "adx", "sbx", "sti", "std", "jsr", "int", "iag",
		"ias", "rfi", "iaq", "hwn", "hwq", "hwi", "dat", "panic", "exit",
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

// containsName returns true if the given list contains the supplied Name.
func containsName(r []*parser.Name, data string) bool {
	for i := range r {
		if r[i].Data == data {
			return true
		}
	}
	return false
}

// FindLabels recursively finds Label nodes.
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

// FindReferences recursively finds Label references.
func FindReferences(n []parser.Node, l *[]*parser.Name) {
	for i := range n {
		switch tt := n[i].(type) {
		case parser.NodeCollection:
			FindReferences(tt.Children(), l)

		case *parser.Name:
			if isRegister(tt.Data) || isInstruction(tt.Data) {
				continue
			}

			*l = append(*l, tt)
		}
	}
}

// StripDuplicateNames removes duplicate entries from the given Name list.
func StripDuplicateNames(r []*parser.Name) []*parser.Name {
	l := make([]*parser.Name, 0, len(r))

	for i := range r {
		if !containsName(l, r[i].Data) {
			l = append(l, r[i])
		}
	}

	return l
}
