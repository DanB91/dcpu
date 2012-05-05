// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import (
	"github.com/jteeuwen/dcpu/cpu"
	dp "github.com/jteeuwen/dcpu/parser"
)

// assembler holds some assembler state.
type assembler struct {
	ast    *dp.AST             // Source AST.
	code   []cpu.Word          // Final program.
	labels map[string]cpu.Word // Map of defined labels with their address.
}

func encode(a, b, c cpu.Word) cpu.Word { return a | (b << 5) | (c << 10) }

// Assemble takes the given AST and attempts to assemble
// it into a compiled program.
func Assemble(ast *dp.AST) (prog []cpu.Word, err error) {
	var asm assembler
	asm.ast = ast
	asm.labels = make(map[string]cpu.Word)

	if err = asm.buildNodes(ast.Root.Children()); err != nil {
		return
	}

	for k, v := range asm.labels {
		println(k, v)
	}

	prog = asm.code
	return
}

// buildNodes compiles the given ast root nodes
func (a *assembler) buildNodes(nodes []dp.Node) (err error) {
	for i := range nodes {
		switch tt := nodes[i].(type) {
		case *dp.Comment:
			/* ignore */

		case *dp.Label:
			a.labels[tt.Data] = cpu.Word(len(a.code))

		case *dp.Instruction:
			err = a.buildInstruction(tt.Children())

		default:
			err = NewBuildError(
				a.ast.Files[tt.File()], tt.Line(), tt.Col(),
				"Unexpected node %T. Want Comment, Label or Instruction.", tt,
			)
		}

		if err != nil {
			return
		}
	}

	return
}

// buildInstruction compiles the given instruction.
func (a *assembler) buildInstruction(nodes []dp.Node) (err error) {
	name := nodes[0].(*dp.Name)
	op, ok := opcodes[name.Data]

	if !ok {
		return NewBuildError(
			a.ast.Files[name.File()], name.Line(), name.Col(),
			"Unknown instruction: %s", name.Data,
		)
	}

	if len(nodes) == 1 {
		// No arguments -- There are only two instructions with
		// this signature. Both are non-standard, extended instructions.
		a.code = append(a.code, encode(cpu.EXT, op.code, 0x21))
		return
	}

	return
}
