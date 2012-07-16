// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

// A Function is a pseudo construct that represents a block of code
// and is treated by various of our tools as a single code unit.
type Function struct {
	*NodeBase
	children []Node
}

func NewFunction(file, line, col int) *Function {
	return &Function{
		NewNodeBase(file, line, col),
		nil,
	}
}

func (f *Function) Children() []Node     { return f.children }
func (f *Function) SetChildren(n []Node) { f.children = n }

func (b *Function) Copy(file, line, col int) Node {
	nb := &Function{
		NodeBase: NewNodeBase(file, line, col),
		children: make([]Node, len(b.children)),
	}

	for i := range nb.children {
		nb.children[i] = b.children[i].Copy(file, line, col)
	}

	return nb
}
