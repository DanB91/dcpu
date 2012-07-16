// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

// String holds a string value.
type String struct {
	*NodeBase
	Data string
}

// NewString creates and returns a new String instance.
func NewString(file, line, col int, data string) *String {
	return &String{
		NewNodeBase(file, line, col),
		data,
	}
}

func (n *String) Copy(file, line, col int) Node {
	return &String{
		NodeBase: NewNodeBase(file, line, col),
		Data:     n.Data,
	}
}
