// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import "github.com/jteeuwen/dcpu/parser"

// parseConstants finds constant definitions and references to them.
// It replaces the references with the nodes denoting the value of
// the respective constants.
func parseConstants(ast *parser.AST, list []parser.Node) (out []parser.Node, err error) {
	var instr *parser.Instruction
	var expr *parser.Expression
	var name *parser.Name
	var ok bool

	consts := make(map[string][]parser.Node)

	for i := 0; i < len(list); i++ {
		if instr, ok = list[i].(*parser.Instruction); !ok {
			continue
		}

		name = instr.Children()[0].(*parser.Name)
		if name.Data != "equ" {
			continue
		}

		expr = instr.Children()[1].(*parser.Expression)
		name = expr.Children()[0].(*parser.Name)

		if _, ok = consts[name.Data]; ok {
			return nil, NewBuildError(
				ast.Files[instr.File()], instr.Line(), instr.Col(),
				"Duplicate constant %q.", name.Data)
		}

		expr = instr.Children()[2].(*parser.Expression)
		consts[name.Data] = expr.Children()

		// Remove constant node from AST.
		copy(list[i:], list[i+1:])
		list = list[:len(list)-1]
		i--
	}

	for k, v := range consts {
		list = replaceConstantRef(list, k, v)
	}

	return list, nil
}

// definition finds and replaces a single constant reference.
func replaceConstantRef(in []parser.Node, key string, value []parser.Node) []parser.Node {
	vsize := len(value)

	for i := 0; i < len(in); i++ {
		switch tt := in[i].(type) {
		case parser.NodeCollection:
			tt.SetChildren(replaceConstantRef(tt.Children(), key, value))

		case *parser.Name:
			if tt.Data != key {
				continue
			}

			switch {
			case vsize == 1:
				in[i] = value[0]

			case vsize == 0:
				copy(in[i:], in[i+1:])
				in = in[:len(in)-1]

			default:
				tmp := make([]parser.Node, (len(in)-1)+vsize)
				copy(tmp, in[:i])
				copy(tmp[i:], value)
				copy(tmp[i+vsize:], in[i+1:])
				in = tmp
			}
		}
	}

	return in
}
