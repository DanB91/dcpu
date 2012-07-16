// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"github.com/jteeuwen/dcpu/cpu"
	"strconv"
)

// Number holds a numerical value.
type Number struct {
	*NodeBase
	Data string
}

// NewNumber creates and returns a new Number instance.
func NewNumber(file, line, col int, data string) *Number {
	return &Number{
		NodeBase: NewNodeBase(file, line, col),
		Data:     data,
	}
}

func (n *Number) Copy(file, line, col int) Node {
	return &Number{
		NodeBase: NewNodeBase(file, line, col),
		Data:     n.Data,
	}
}

// Parse attempts to process the node's string data as a number.
func (n *Number) Parse() (cpu.Word, error) {
	var v uint64
	var err error

	if len(n.Data) > 2 && n.Data[0] == '0' && n.Data[1] == 'b' {
		// strconv.ParseUint can't deal with 0b01010101 formatted strings.
		// So handle these manually.
		v, err = strconv.ParseUint(n.Data[2:], 2, 64)
	} else {
		// Otherwise, just let it figure out if we have octal, decimal or hex values.
		v, err = strconv.ParseUint(n.Data, 0, 64)
	}

	return cpu.Word(v), err
}
