// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

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

func (n *Name) Dump(pad string) string {
	if len(n.Data) > 20 {
		return fmt.Sprintf("%s %T(%.20q...)\n", n.NodeBase.Dump(pad), n, n.Data)
	}
	return fmt.Sprintf("%s %T(%q)\n", n.NodeBase.Dump(pad), n, n.Data)
}
