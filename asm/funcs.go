// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import "github.com/jteeuwen/dcpu/parser"

// parseFunctions resolves function constants and it injects the
// necessary function prolog and epilog code.
//
// What the latter means, is that we analyze the registers being
// used. If any of them are registers that are guaranteed to
// be preserved across function calls, then we explicitly inject stack
// push/pop instructions to save their state.
//
// If there are any `return` instructions, we replace them with `set pc, $$`. 
// Where `$$` is a label we generate which points to the function epilog.
func parseFunctions(ast *parser.AST) (err error) {
	list := ast.Functions()

	for i := range list {
		if err = parseFuncConst(ast, list[i]); err != nil {
			return
		}

		fixFunctionReturns(ast, list[i])
		injectFunctionCode(ast, list[i])
	}

	return
}

// fixFunctionReturns finds 'return' instructions and replaces them
// with appropriate `set pc, $__<name>_epilog` versions.
func fixFunctionReturns(ast *parser.AST, f *parser.Function) {
	var instr *parser.Instruction
	var expr *parser.Expression
	var file, line, col int
	var code []parser.Node
	var ok bool

	name := f.Children()[0].(*parser.Name)
	exitLabel := "$__" + name.Data + "_epilog"
	list := f.Children()[1:]

	for i := range list {
		if instr, ok = list[i].(*parser.Instruction); !ok {
			continue
		}

		code = instr.Children()
		name = code[0].(*parser.Name)
		if name.Data != "return" {
			continue
		}

		file, line, col = instr.File(), instr.Line(), instr.Col()
		name.Data = "set"
		code = append(code, nil, nil)

		expr = parser.NewExpression(file, line, col)
		expr.SetChildren([]parser.Node{
			parser.NewName(file, line, col, "pc"),
		})
		code[1] = expr

		expr = parser.NewExpression(file, line, col)
		expr.SetChildren([]parser.Node{
			parser.NewName(file, line, col, exitLabel),
		})
		code[2] = expr

		instr.SetChildren(code)
	}
}

// injectFunctionCode injects function prologs and epilogs.
func injectFunctionCode(ast *parser.AST, f *parser.Function) {
	var tmp, code []parser.Node
	var instr *parser.Instruction
	var expr *parser.Expression
	var idx int

	name := f.Children()[0].(*parser.Name)
	exitLabel := "$__" + name.Data + "_epilog"
	file, line, col := name.File(), name.Line(), name.Col()

	// Find list of referenced protected registers.
	var regs []string
	findProtectedRegisters(f.Children()[1:], &regs)

	// Inject prolog and epilog code.
	//
	// For each protected register we found, we add appropriate
	// stack push/pop instructions to preserve their state.
	tmp = make([]parser.Node, len(f.Children())+(len(regs)*2)+2)
	tmp[idx] = parser.NewLabel(file, line, col, name.Data)
	idx++

	// set push, $reg
	for n := range regs {
		instr = parser.NewInstruction(file, line, col)
		code = make([]parser.Node, 3)
		code[0] = parser.NewName(file, line, col, "set")

		expr = parser.NewExpression(file, line, col)
		expr.SetChildren([]parser.Node{
			parser.NewName(file, line, col, "push"),
		})
		code[1] = expr

		expr = parser.NewExpression(file, line, col)
		expr.SetChildren([]parser.Node{
			parser.NewName(file, line, col, regs[n]),
		})
		code[2] = expr

		instr.SetChildren(code)
		tmp[idx] = instr
		idx++
	}

	// Regular code goes here.
	copy(tmp[idx:], f.Children()[1:])
	idx += len(f.Children()) - 1

	// Label denoting the start of the epilog.
	line, col = tmp[idx-1].Line()+1, 1
	tmp[idx] = parser.NewLabel(file, line, col, exitLabel)
	idx++

	// set $reg, pop
	for n := range regs {
		instr = parser.NewInstruction(file, line, col)
		code = make([]parser.Node, 3)
		code[0] = parser.NewName(file, line, col, "set")

		expr = parser.NewExpression(file, line, col)
		expr.SetChildren([]parser.Node{
			parser.NewName(file, line, col, regs[n]),
		})
		code[1] = expr

		expr = parser.NewExpression(file, line, col)
		expr.SetChildren([]parser.Node{
			parser.NewName(file, line, col, "pop"),
		})
		code[2] = expr

		instr.SetChildren(code)
		tmp[idx] = instr
		idx++
	}

	// set pc, pop
	instr = parser.NewInstruction(file, line, col)
	code = make([]parser.Node, 3)
	code[0] = parser.NewName(file, line, col, "set")

	expr = parser.NewExpression(file, line, col)
	expr.SetChildren([]parser.Node{
		parser.NewName(file, line, col, "pc"),
	})
	code[1] = expr

	expr = parser.NewExpression(file, line, col)
	expr.SetChildren([]parser.Node{
		parser.NewName(file, line, col, "pop"),
	})
	code[2] = expr

	instr.SetChildren(code)
	tmp[idx] = instr

	f.SetChildren(tmp)
}

func parseFuncConst(ast *parser.AST, f *parser.Function) (err error) {
	list, err := parseConstants(ast, f.Children())
	if err != nil {
		return
	}

	f.SetChildren(list)
	return
}

// hasProtectedRegister checks if the given list contains the specified name.
func hasProtectedRegister(l []string, v string) bool {
	for i := range l {
		if l[i] == v {
			return true
		}
	}

	return false
}

// findProtectedRegisters recursively finds references to protected registers.
func findProtectedRegisters(in []parser.Node, regs *[]string) {
	for i := range in {
		switch tt := in[i].(type) {
		case parser.NodeCollection:
			findProtectedRegisters(tt.Children(), regs)

		case *parser.Name:
			switch tt.Data {
			case "x", "y", "z", "i", "j":
				if !hasProtectedRegister(*regs, tt.Data) {
					*regs = append(*regs, tt.Data)
				}
			}
		}
	}

	return
}
