// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

// An Instruction is a collection of AST nodes.
// An opcode and optional arguments.
type Instruction struct {
	*NodeBase
	children []Node
}

func NewInstruction(file, line, col int) *Instruction {
	return &Instruction{
		NewNodeBase(file, line, col),
		nil,
	}
}

func (i *Instruction) Children() []Node     { return i.children }
func (i *Instruction) SetChildren(n []Node) { i.children = n }
