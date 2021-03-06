// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

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

func (b *Instruction) Copy(file, line, col int) Node {
	nb := &Instruction{
		NodeBase: NewNodeBase(file, line, col),
		children: make([]Node, len(b.children)),
	}

	for i := range nb.children {
		nb.children[i] = b.children[i].Copy(file, line, col)
	}

	return nb
}
