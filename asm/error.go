// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import "fmt"

// Represents an assembler/build error.
type BuildError struct {
	File string
	Msg  string
	Line int
	Col  int
}

// NewBuildError creates a new assembler/build error from the given values.
func NewBuildError(file string, line, col int, f string, argv ...interface{}) *BuildError {
	return &BuildError{file, fmt.Sprintf(f, argv...), line, col}
}

// Error returns a string representation of this error.
func (e *BuildError) Error() string {
	return fmt.Sprintf("%s:%d:%d %s", e.File, e.Line, e.Col, e.Msg)
}
