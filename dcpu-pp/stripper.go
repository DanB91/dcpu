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
	return
}
