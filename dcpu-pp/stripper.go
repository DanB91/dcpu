// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import dp "github.com/jteeuwen/dcpu/parser"

func init() {
	Register("strip", "Remove all code comments.", NewStripper)
}

// Stripper removes all code comments from the AST.
type Stripper struct{}

func NewStripper() Processor { return new(Stripper) }

func (p *Stripper) Process(ast *dp.AST) (err error) {
	stripComments(ast.Root)
	return
}

// stripComments removes Comment nodes from the supplied list.
func stripComments(n dp.NodeCollection) {
	list := n.Children()

loop:
	for i := range list {
		switch tt := list[i].(type) {
		case *dp.Comment:
			copy(list[i:], list[i+1:])
			list = list[:len(list)-1]
			goto loop

		case dp.NodeCollection:
			stripComments(tt)
		}
	}

	n.SetChildren(list)
}
