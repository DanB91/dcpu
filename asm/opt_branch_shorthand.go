// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import dp "github.com/jteeuwen/dcpu/parser"

func init() {
	RegisterOptimization(opt_branch_shorthand)
}

// opt_branch_shorthand finds instances of "IFE 0, a" amd "IFN 0, a"
// and checks if the first operand is a short-form numeric literal.
// If so, the operands are swapped.
//
// The DCPU spec states that short-form numbers can not be encoded in
// the first operand, since its maximum value is not large enough to
// hold all allowed literals. In this case, the assembler would have
// to set A to 0x1f (next word) and store the value in a new word.
// 
// This needlessly increases the size of the program by one word.
// For the IFE and IFN instructions, we can prevent this from 
// happening by simply swapping the operands around.
//
// For other instructions this is problematic, since swapping them out
// changes the semantics of the operation. In those cases, we simply
// allow the assembler to generate the extra word.
func opt_branch_shorthand(ast *dp.AST) {
	var instr *dp.Instruction
	var argv []dp.Node
	var name *dp.Name
	var num *dp.Number
	var ok bool

	nodes := ast.Root.Children()

	for i := range nodes {
		instr, ok = nodes[i].(*dp.Instruction)
		if !ok {
			continue
		}

		argv = instr.Children()
		name = argv[0].(*dp.Name)

		if name.Data != "ife" && name.Data != "ifn" {
			continue
		}

		if len(argv) < 3 {
			continue
		}

		num, ok = argv[1].(*dp.Expression).Children()[0].(*dp.Number)

		if ok && (num.Data == 0xffff || num.Data <= 0x1e) {
			argv[1], argv[2] = argv[2], argv[1]
		}
	}
}
