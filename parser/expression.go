// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

// An Expression is a collection of AST nodes.
type Expression struct {
	*NodeBase
	children []Node
}

func NewExpression(file, line, col int) *Expression {
	return &Expression{
		NewNodeBase(file, line, col),
		nil,
	}
}

func (e *Expression) Children() []Node     { return e.children }
func (e *Expression) SetChildren(n []Node) { e.children = n }

func (b *Expression) Copy(file, line, col int) Node {
	nb := &Expression{
		NodeBase: NewNodeBase(file, line, col),
		children: make([]Node, len(b.children)),
	}

	for i := range nb.children {
		nb.children[i] = b.children[i].Copy(file, line, col)
	}

	return nb
}
