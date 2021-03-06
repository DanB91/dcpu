// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import "fmt"

// Represents a parse error.
type ParseError struct {
	File string
	Msg  string
	Line int
	Col  int
}

// NewParseError creates a new parse error from the given values.
func NewParseError(file string, line, col int, f string, argv ...interface{}) *ParseError {
	return &ParseError{file, fmt.Sprintf(f, argv...), line, col}
}

// Error returns a string representation of this error.
func (e *ParseError) Error() string {
	return fmt.Sprintf("%s:%d:%d %s", e.File, e.Line, e.Col, e.Msg)
}
