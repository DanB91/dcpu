// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

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

func (n *Comment) Dump(pad string) string {
	if len(n.Data) > 20 {
		return fmt.Sprintf("%s %T(%.20q...)\n", n.NodeBase.Dump(pad), n, n.Data)
	}
	return fmt.Sprintf("%s %T(%q)\n", n.NodeBase.Dump(pad), n, n.Data)
}
