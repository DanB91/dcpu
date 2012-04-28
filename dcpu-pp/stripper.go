// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

func init() {
	Register("strip", "Remove all code comments.", NewStripper)
}

// Stripper removes all code comments from the AST.
type Stripper struct{}

func NewStripper() Processor { return new(Stripper) }

func (p *Stripper) Process(ast *AST) (err error) {
	stripComments(&ast.Root.Children)
	return
}

// stripComments removes Comment nodes from the supplied list.
func stripComments(n *[]Node) {
	t := *n

loop:
	for i := range t {
		switch tt := t[i].(type) {
		case *Comment:
			copy(t[i:], t[i+1:])
			t = t[:len(t)-1]
			goto loop

		case *Block:
			stripComments(&tt.Children)

		case *Expression:
			stripComments(&tt.Children)

		case *Instruction:
			stripComments(&tt.Children)
		}
	}

	*n = t
}
