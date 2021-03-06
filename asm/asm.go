// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// DCPU Assembler package.
package asm

import (
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/parser"
)

// assembler holds some assembler state.
type assembler struct {
	ast    *parser.AST               // Source AST.
	code   []cpu.Word                // Final program.
	labels map[string]cpu.Word       // Map of defined labels with their address.
	refs   map[cpu.Word]*parser.Name // Indices into `code` holding unresolved label references.
	debug  *DebugInfo                // Maps binary instructions to original source locations.
}

// Assemble takes the given AST and attempts to assemble it into a compiled program.
//
// It returns either an error, or the program along with debug symbols.
func Assemble(ast *parser.AST) (prog []cpu.Word, dbg *DebugInfo, err error) {
	var asm assembler
	asm.ast = ast
	asm.labels = make(map[string]cpu.Word)
	asm.refs = make(map[cpu.Word]*parser.Name)
	asm.debug = new(DebugInfo)

	// Process function definitions.
	// This also processes function-local constants.
	if err = parseFunctions(ast); err != nil {
		return
	}

	// Process global constants.
	list, err := parseConstants(ast, ast.Root.Children())
	if err != nil {
		return
	}

	ast.Root.SetChildren(list)

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

		return nil, nil, NewBuildError(
			ast.Files[v.File()], v.Line(), v.Col(),
			"Unknown label reference %q.", v.Data)
	}

	asm.debug.SetFileDefs(ast.Files)

	prog = asm.code
	dbg = asm.debug
	return
}

// buildNodes compiles the given ast root nodes
func (a *assembler) buildNodes(nodes []parser.Node) (err error) {
	for i := range nodes {
		switch tt := nodes[i].(type) {
		case *parser.Comment:
			/* ignore */

		case *parser.Label:
			a.labels[tt.Data] = cpu.Word(len(a.code))

		case *parser.Function:
			err = a.buildFunction(tt)

		case *parser.Instruction:
			err = a.buildInstruction(tt.Children())

		default:
			err = NewBuildError(
				a.ast.Files[tt.File()], tt.Line(), tt.Col(),
				"Unexpected node %T. Want Comment, Label, Function or Instruction.", tt,
			)
		}

		if err != nil {
			return
		}
	}

	return
}

// buildFunction compiles the given function.
func (a *assembler) buildFunction(f *parser.Function) (err error) {
	nodes := f.Children()
	name := nodes[0].(*parser.Label)

	a.debug.SetFunctionStart(cpu.Word(len(a.code)), f.Line(), name.Data)

	for i := range nodes {
		switch tt := nodes[i].(type) {
		case *parser.Comment:
			/* ignore */

		case *parser.Label:
			a.labels[tt.Data] = cpu.Word(len(a.code))

		case *parser.Instruction:
			err = a.buildInstruction(tt.Children())

		default:
			err = NewBuildError(
				a.ast.Files[tt.File()], tt.Line(), tt.Col(),
				"Unexpected node %T. Want Comment, Label, Instruction.", tt,
			)
		}

		if err != nil {
			return
		}
	}

	a.debug.SetFunctionEnd(cpu.Word(len(a.code)), nodes[len(nodes)-1].Line())
	return
}

// buildInstruction compiles the given instruction.
func (a *assembler) buildInstruction(nodes []parser.Node) (err error) {
	name := nodes[0].(*parser.Name)

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
	var symbols []parser.Node

	symbols = append(symbols, name)

	switch op.argc {
	case 2:
		va, err = a.buildOperand(&argv, &symbols, nodes[1].(*parser.Expression).Children()[0], true)
		if err != nil {
			return
		}

		vb, err = a.buildOperand(&argv, &symbols, nodes[2].(*parser.Expression).Children()[0], false)

	case 1:
		va, err = a.buildOperand(&argv, &symbols, nodes[1].(*parser.Expression).Children()[0], false)
	}

	if err != nil {
		return
	}

	if op.ext {
		a.code = append(a.code, cpu.Encode(cpu.EXT, op.code, va))
	} else {
		a.code = append(a.code, cpu.Encode(op.code, va, vb))
	}

	a.debug.Emit(symbols...)
	a.code = append(a.code, argv...)
	return
}

// buildOperand compiles the given instruction operand.
//
// The `first` parameter determines if we are parsing the A or B parameter
// in something like 'set A, B'. This makes a difference when encoding
// small literal numbers.
func (a *assembler) buildOperand(argv *[]cpu.Word, symbols *[]parser.Node, node parser.Node, first bool) (val cpu.Word, err error) {
	switch tt := node.(type) {
	case *parser.Name:
		if reg, ok := registers[tt.Data]; ok {
			return reg, nil
		}

		if addr, ok := a.labels[tt.Data]; ok {
			if !first && (addr == 0xffff || addr <= 0x1e) {
				return addr + 0x21, nil
			}

			*symbols = append(*symbols, tt)
			*argv = append(*argv, addr)
			return 0x1f, nil
		}

		a.refs[cpu.Word(len(a.code)+1+len(*argv))] = tt
		*symbols = append(*symbols, tt)
		*argv = append(*argv, 0)
		return 0x1f, nil

	case parser.NumericNode:
		num, err := tt.Parse()
		if err != nil {
			return 0, err
		}

		if !first && (num == 0xffff || num <= 0x1e) {
			return num + 0x21, nil
		}

		*symbols = append(*symbols, tt)
		*argv = append(*argv, num)
		return 0x1f, nil

	case *parser.Block:
		return a.buildBlock(argv, symbols, tt)

	default:
		return 0, NewBuildError(
			a.ast.Files[tt.File()], tt.Line(), tt.Col(),
			"Unexpected node %T. Want Name, Number or Block.", tt,
		)
	}

	return
}

// buildBlock builds a block expression.
func (a *assembler) buildBlock(argv *[]cpu.Word, symbols *[]parser.Node, b *parser.Block) (val cpu.Word, err error) {
	nodes := b.Children()

	switch len(nodes) {
	case 1:
		switch tt := nodes[0].(type) {
		case *parser.Name:
			if reg, ok := registers[tt.Data]; ok {
				return reg + 0x08, nil
			}

			if addr, ok := a.labels[tt.Data]; ok {
				*symbols = append(*symbols, tt)
				*argv = append(*argv, addr)
				return 0x1e, nil
			}

			a.refs[cpu.Word(len(a.code)+1+len(*argv))] = tt
			*symbols = append(*symbols, tt)
			*argv = append(*argv, 0)
			return 0x1e, nil

		case parser.NumericNode:
			num, err := tt.Parse()
			if err != nil {
				return 0, err
			}

			*symbols = append(*symbols, tt)
			*argv = append(*argv, num)
			return 0x1e, nil

		default:
			return 0, NewBuildError(
				a.ast.Files[tt.File()], tt.Line(), tt.Col(),
				"Unexpected node %T. Want Name, Number or Block.", tt,
			)
		}

	case 3:
		var va, vb cpu.Word

		if va, err = a.buildBlockOperand(argv, symbols, nodes[0]); err != nil {
			return
		}

		if vb, err = a.buildBlockOperand(argv, symbols, nodes[2]); err != nil {
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

func (a *assembler) buildBlockOperand(argv *[]cpu.Word, symbols *[]parser.Node, node parser.Node) (cpu.Word, error) {
	switch tt := node.(type) {
	case *parser.Name:
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
			*symbols = append(*symbols, tt)
			*argv = append(*argv, addr)
			return 0, nil
		}

		a.refs[cpu.Word(len(a.code)+1+len(*argv))] = tt
		*symbols = append(*symbols, tt)
		*argv = append(*argv, 0)
		return 0, nil

	case parser.NumericNode:
		num, err := tt.Parse()
		if err != nil {
			return 0, err
		}

		*symbols = append(*symbols, tt)
		*argv = append(*argv, num)
		return 0, nil
	}

	return 0, NewBuildError(
		a.ast.Files[node.File()], node.Line(), node.Col(),
		"Unexpected node %T. Want Name or Number.", node,
	)
}

// buildData compiles the given data section
func (a *assembler) buildData(nodes []parser.Node) (err error) {
	var r rune
	nodes = nodes[1:] // Skip 'dat' instruction.

	for i := range nodes {
		expr, ok := nodes[i].(*parser.Expression)
		if !ok {
			continue
		}

		switch tt := expr.Children()[0].(type) {
		case *parser.String:
			for _, r = range tt.Data {
				a.debug.Emit(tt)
				a.code = append(a.code, cpu.Word(r))
			}

		case parser.NumericNode:
			num, err := tt.Parse()
			if err != nil {
				return err
			}

			a.debug.Emit(tt)
			a.code = append(a.code, num)
		}
	}

	return
}
