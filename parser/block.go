// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

// An Block is a collection of AST nodes.
type Block struct {
	*NodeBase
	children []Node
}

func NewBlock(file, line, col int) *Block {
	return &Block{
		NewNodeBase(file, line, col),
		nil,
	}
}

func (b *Block) Children() []Node     { return b.children }
func (b *Block) SetChildren(n []Node) { b.children = n }

func (b *Block) Copy(file, line, col int) Node {
	nb := &Block{
		NodeBase: NewNodeBase(file, line, col),
		children: make([]Node, len(b.children)),
	}

	for i := range nb.children {
		nb.children[i] = b.children[i].Copy(file, line, col)
	}

	return nb
}
