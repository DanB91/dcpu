// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

// Proprocessor mode of operation
type ParseMode uint8

// Known modes
const (
	ModeAssemble ParseMode = iota // Default mode - BUild output source.
	ModeDumpAST                   // Dump AST for all source code.
)

// Config holds configuration and state data for the preprocessor.
type Config struct {
	Include []string  // List of paths where we look to resolve source file references.
	Input   []string  // Names of input source files.
	Output  string    // Name of output source file. Defaults to nil (stdout).
	Mode    ParseMode // Selected mode of operation.
}

// NewConfig creates a new, standard configuration instance.
func NewConfig() *Config {
	c := new(Config)
	c.Mode = ModeAssemble
	return c
}
