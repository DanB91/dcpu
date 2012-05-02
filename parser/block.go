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
