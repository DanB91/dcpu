// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

// An Comment holds a code comment.
type Comment struct {
	*NodeBase
	Data string
}

func NewComment(file, line, col int, value string) *Comment {
	return &Comment{
		NewNodeBase(file, line, col),
		value,
	}
}

func (n *Comment) Copy(file, line, col int) Node {
	return &Comment{
		NodeBase: NewNodeBase(file, line, col),
		Data:     n.Data,
	}
}
