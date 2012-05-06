// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import (
	"github.com/jteeuwen/dcpu/cpu"
	dp "github.com/jteeuwen/dcpu/parser"
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

	// Compile program.
	if err = asm.buildNodes(ast.Root.Children()); err != nil {
		return
	}

	// Fix unresolved label references.
	for k, v := range asm.refs {
		addr, ok := asm.labels[v.Data]
		if ok {
			asm.code[k] = addr
			continue
		}

		return nil, NewBuildError(
			ast.Files[v.File()], v.Line(), v.Col(),
			"Unknown label reference %q.", v.Data)
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

	if name.Data == "dat" {
		return a.buildData(nodes)
	}

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
	var argv []cpu.Word

	switch op.argc {
	case 2:
		vb, err = a.buildOperand(&argv, nodes[2].(*dp.Expression).Children()[0])
		if err != nil {
			return
		}

		fallthrough
	case 1:
		va, err = a.buildOperand(&argv, nodes[1].(*dp.Expression).Children()[0])
		if err != nil {
			return
		}
	}

	if op.ext {
		a.code = append(a.code, encode(cpu.EXT, op.code, va))
	} else {
		a.code = append(a.code, encode(op.code, va, vb))
	}

	a.code = append(a.code, argv...)
	return
}

// buildOperand compiles the given instruction operand.
func (a *assembler) buildOperand(argv *[]cpu.Word, node dp.Node) (val cpu.Word, err error) {
	switch tt := node.(type) {
	case *dp.Name:
		if reg, ok := registers[tt.Data]; ok {
			return reg, nil
		}

		if addr, ok := a.labels[tt.Data]; ok {
			if addr == 0xffff || addr <= 0x1e {
				return addr + 0x21, nil
			}

			*argv = append(*argv, addr)
			return 0x1f, nil
		}

		a.refs[cpu.Word(len(a.code)+1+len(*argv))] = tt
		*argv = append(*argv, 0)
		return 0x1f, nil

	case *dp.Number:
		if tt.Data == 0xffff || tt.Data <= 0x1e {
			return tt.Data + 0x21, nil
		}

		*argv = append(*argv, tt.Data)
		return 0x1f, nil

	case *dp.Block:
		return a.buildBlock(argv, tt)

	default:
		return 0, NewBuildError(
			a.ast.Files[tt.File()], tt.Line(), tt.Col(),
			"Unexpected node %T. Want Name, Number or Block.", tt,
		)
	}

	return
}

// buildBlock builds a block expression.
func (a *assembler) buildBlock(argv *[]cpu.Word, b *dp.Block) (val cpu.Word, err error) {
	nodes := b.Children()

	switch len(nodes) {
	case 1:
		switch tt := nodes[0].(type) {
		case *dp.Name:
			if reg, ok := registers[tt.Data]; ok {
				return reg + 0x08, nil
			}

			if addr, ok := a.labels[tt.Data]; ok {
				*argv = append(*argv, addr)
				return 0x1e, nil
			}

			a.refs[cpu.Word(len(a.code)+1+len(*argv))] = tt
			*argv = append(*argv, 0)
			return 0x1e, nil

		case *dp.Number:
			*argv = append(*argv, tt.Data)
			return 0x1e, nil

		default:
			return 0, NewBuildError(
				a.ast.Files[tt.File()], tt.Line(), tt.Col(),
				"Unexpected node %T. Want Name, Number or Block.", tt,
			)
		}

	case 3:
		var va, vb cpu.Word
		if va, err = a.buildBlockOperand(argv, nodes[0]); err != nil {
			return
		}

		if vb, err = a.buildBlockOperand(argv, nodes[2]); err != nil {
			return
		}

		if va != 0 {
			return va, nil
		}

		return vb, nil

	default:
		err = NewBuildError(
			a.ast.Files[b.File()], b.Line(), b.Col(),
			"Unexpected node count. Want 1 or 3. Got %d", len(nodes))
	}

	return
}

func (a *assembler) buildBlockOperand(argv *[]cpu.Word, node dp.Node) (cpu.Word, error) {
	switch tt := node.(type) {
	case *dp.Name:
		if reg, ok := registers[tt.Data]; ok {
			if reg <= 0x7 {
				return reg + 0x10, nil
			}

			if reg == 0x1b {
				return 0x1a, nil
			}

			return 0, NewBuildError(
				a.ast.Files[tt.File()], tt.Line(), tt.Col(),
				"Illegal use of register %q.", tt.Data,
			)
		}

		if addr, ok := a.labels[tt.Data]; ok {
			*argv = append(*argv, addr)
			return 0, nil
		}

		a.refs[cpu.Word(len(a.code)+1+len(*argv))] = tt
		*argv = append(*argv, 0)
		return 0, nil

	case *dp.Number:
		*argv = append(*argv, tt.Data)
		return 0, nil
	}

	return 0, NewBuildError(
		a.ast.Files[node.File()], node.Line(), node.Col(),
		"Unexpected node %T. Want Name or Number.", node,
	)
}

// buildData compiles the given data section
func (a *assembler) buildData(nodes []dp.Node) (err error) {
	var r rune
	nodes = nodes[1:] // Skip 'dat' instruction.

	for i := range nodes {
		expr, ok := nodes[i].(*dp.Expression)

		if !ok {
			continue
		}

		switch tt := expr.Children()[0].(type) {
		case *dp.String:
			for _, r = range tt.Data {
				a.code = append(a.code, cpu.Word(r))
			}

		case *dp.Number:
			a.code = append(a.code, tt.Data)
		}
	}

	return
}
