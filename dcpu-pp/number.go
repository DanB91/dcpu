// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "strconv"

// Number holds a numerical value.
type Number struct {
	*NodeBase
	Data Word
}

// NewNumber creates and returns a new Number instance.
func NewNumber(file, line, col int, data Word) *Number {
	return &Number{
		NodeBase: NewNodeBase(file, line, col),
		Data:     data,
	}
}

// Parse attempts to process the given data as a number.
func (n *Number) Parse(data []byte) (err error) {
	var v uint64

	if len(data) > 2 && data[0] == '0' && data[1] == 'b' {
		// strconv.ParseUint can't deal with 0b01010101 formatted strings for
		// some reason. So handle these manually.
		v, err = strconv.ParseUint(string(data[2:]), 2, 64)
	} else {
		// Otherwise, just let it figure out if we have octal, decimal or hex values.
		v, err = strconv.ParseUint(string(data), 0, 64)
	}

	n.Data = Word(v)
	return
}
