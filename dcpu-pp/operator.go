// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

// An Operator holds an arithmatic operator.
type Operator struct {
	*NodeBase
	Data string
}

func NewOperator(file, line, col int, value string) *Operator {
	return &Operator{
		NewNodeBase(file, line, col),
		value,
	}
}

func (n *Operator) Dump(pad string) string {
	return fmt.Sprintf("%s %T(%q)\n", n.NodeBase.Dump(pad), n, n.Data)
}
