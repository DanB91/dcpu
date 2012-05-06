// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import (
	"github.com/jteeuwen/dcpu/cpu"
	dp "github.com/jteeuwen/dcpu/parser"
	"os"
)

// assembler holds some assembler state.
type assembler struct {
	ast    *dp.AST               // Source AST.
	code   []cpu.Word            // Final program.
	labels map[string]cpu.Word   // Map of defined labels with their address.
	refs   map[cpu.Word]*dp.Name // Indices into `code` holding unresolved label references.
}

func encode(a, b, c cpu.Word) cpu.Word { return a | (b << 5) | (c << 10) }

// Assemble takes the given AST and attempts to assemble
// it into a compiled program.
func Assemble(ast *dp.AST) (prog []cpu.Word, err error) {
	var asm assembler
	asm.ast = ast
	asm.labels = make(map[string]cpu.Word)
	asm.refs = make(map[cpu.Word]*dp.Name)

	dp.WriteAst(os.Stdout, ast)
	if err = asm.buildNodes(ast.Root.Children()); err != nil {
		return
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

	if len(nodes)-1 != op.argc {
		return NewBuildError(
			a.ast.Files[name.File()], name.Line(), name.Col(),
			"Invalid argument count for %q. Want %d", name.Data, op.argc,
		)
	}

	var va, vb cpu.Word

	switch op.argc {
	case 2:
		vb, err = a.buildOperand(nodes[2].(*dp.Expression).Children()[0])
		if err != nil {
			return
		}

		fallthrough
	case 1:
		va, err = a.buildOperand(nodes[1].(*dp.Expression).Children()[0])
		if err != nil {
			return
		}
	}

	if op.ext {
		a.code = append(a.code, encode(cpu.EXT, op.code, va))
	} else {
		a.code = append(a.code, encode(op.code, va, vb))
	}

	return
}

// buildOperand compiles the given instruction operand.
func (a *assembler) buildOperand(node dp.Node) (val cpu.Word, err error) {
	switch tt := node.(type) {
	case *dp.Name:
		return a.buildName(tt)

	case *dp.Number:
		return a.buildNumber(tt)

	case *dp.Block:

	default:
		return 0, NewBuildError(
			a.ast.Files[tt.File()], tt.Line(), tt.Col(),
			"Unexpected node %T. Want Name, Number or Block.", tt,
		)
	}

	return
}

// buildName builds a name operand.
func (a *assembler) buildName(name *dp.Name) (val cpu.Word, err error) {
	if reg, ok := registers[name.Data]; ok {
		return reg, nil
	}

	if addr, ok := a.labels[name.Data]; ok {
		return addr, nil
	}

	// Undefined label. Save it for later and reserve
	// a new word for the address we will be resolving.
	a.refs[cpu.Word(len(a.code))] = name
	a.code = append(a.code, 0)
	return 0, nil
}

// buildNumber builds a numerical operand.
func (a *assembler) buildNumber(num *dp.Number) (val cpu.Word, err error) {
	if num.Data == 0xffff || num.Data <= 0x1e {
		return num.Data + 0x20, nil
	}

	return 0, nil
}
