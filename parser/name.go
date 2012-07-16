// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

// An Name holds an identifier.
type Name struct {
	*NodeBase
	Data string
}

func NewName(file, line, col int, value string) *Name {
	return &Name{
		NewNodeBase(file, line, col),
		value,
	}
}

func (n *Name) Copy(file, line, col int) Node {
	return &Name{
		NodeBase: NewNodeBase(file, line, col),
		Data:     n.Data,
	}
}
