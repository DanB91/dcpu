// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import dp "github.com/jteeuwen/dcpu/parser"

func init() {
	RegisterPreProcessor("shorthand",
		"Fixes 'IF[E|N] 0, $R' to use one word less.",
		NewOptShorthand, true)
}

// OptShorthand finds instances of "IFE B, A" and "IFN B, A"
// and checks if the first operand is a short-form numeric literal.
// If so, the operands are swapped.
//
// The DCPU spec states that short-form numbers can not be encoded in
// the first operand, since its maximum value is not large enough to
// hold all allowed literals. In this case, the assembler would have
// to set B to 0x1f (next word) and store the value in a new word.
// 
// This needlessly increases the size of the program by one word.
// For the IFE and IFN instructions, we can prevent this from 
// happening by simply swapping the operands around.
//
// For other instructions this is problematic, since swapping them out
// changes the semantics of the operation. In those cases, we simply
// allow the assembler to generate the extra word.
type OptShorthand struct{}

func NewOptShorthand() PreProcessor { return new(OptShorthand) }

func (*OptShorthand) Process(ast *dp.AST) (err error) {
	nodes := ast.Root.Children()

	for i := range nodes {
		instr, ok := nodes[i].(*dp.Instruction)
		if !ok {
			continue
		}

		argv := instr.Children()
		name := argv[0].(*dp.Name)

		if name.Data != "ife" && name.Data != "ifn" {
			continue
		}

		if len(argv) < 3 {
			continue
		}

		num, ok := argv[1].(*dp.Expression).Children()[0].(dp.NumericNode)
		if !ok {
			continue
		}

		word, err := num.Parse()
		if err != nil {
			return err
		}

		if ok && (word == 0xffff || word <= 0x1e) {
			argv[1], argv[2] = argv[2], argv[1]
		}
	}

	return
}
