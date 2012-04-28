// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

// An Instruction is a collection of AST nodes.
// An opcode and optional arguments.
type Instruction struct {
	*NodeBase
	Children []Node
}

func NewInstruction(file, line, col int) *Instruction {
	return &Instruction{
		NewNodeBase(file, line, col),
		nil,
	}
}
