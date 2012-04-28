// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

// An Label holds a label name.
type Label struct {
	*NodeBase
	Data string
}

func NewLabel(file, line, col int, value string) *Label {
	return &Label{
		NewNodeBase(file, line, col),
		value,
	}
}

func (n *Label) Dump(pad string) string {
	if len(n.Data) > 20 {
		return fmt.Sprintf("%s %T(%.20q...)\n", n.NodeBase.Dump(pad), n, n.Data)
	}
	return fmt.Sprintf("%s %T(%q)\n", n.NodeBase.Dump(pad), n, n.Data)
}
