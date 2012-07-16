// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"errors"
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"unicode/utf8"
)

// Char holds a character literal.
type Char struct {
	*NodeBase
	Data string
}

// NewChar creates and returns a new Char instance.
func NewChar(file, line, col int, data string) *Char {
	return &Char{
		NodeBase: NewNodeBase(file, line, col),
		Data:     data,
	}
}

func (n *Char) Copy(file, line, col int) Node {
	return &Char{
		NodeBase: NewNodeBase(file, line, col),
		Data:     n.Data,
	}
}

// Parse attempts to process the node's string data as a number.
func (n *Char) Parse() (cpu.Word, error) {
	r, size := utf8.DecodeRuneInString(n.Data)

	// If the encoding is invalid, utf8.DecodeRune yields (RuneError, 1).
	// This constitutes an impossible result for correct UTF-8.
	if r == utf8.RuneError && size == 1 {
		return 0, errors.New(fmt.Sprintf("Invalid utf8 character literal: %s", n.Data))
	}

	return cpu.Word(r), nil
}
