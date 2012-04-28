// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

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
