// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
	"os"
)

// Config holds configuration and state data for the preprocessor.
type Config struct {
	Include []string       // List of paths where we look to resolve source file references.
	Input   []string       // Names of input source files.
	Output  io.WriteCloser // Output source writer. Defaults to stdout.
	DumpAST bool           // Whether to dump the final AST, or assembly source.
}

// NewConfig creates a new, standard configuration instance.
func NewConfig() *Config {
	c := new(Config)
	c.Output = os.Stdout
	c.DumpAST = false
	return c
}
