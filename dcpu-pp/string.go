// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

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

func (n *String) Dump(pad string) string {
	if len(n.Data) > 20 {
		return fmt.Sprintf("%s %T(%.20q...)\n", n.NodeBase.Dump(pad), n, n.Data)
	}
	return fmt.Sprintf("%s %T(%q)\n", n.NodeBase.Dump(pad), n, n.Data)
}
