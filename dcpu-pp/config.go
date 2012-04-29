// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
	"os"
	"path"
)

// Config holds configuration and state data for the preprocessor.
type Config struct {
	Include []string       // List of paths where we look to resolve source file references.
	Output  io.WriteCloser // Output source writer. Defaults to stdout.
	Input   string         // Name of input source file.
	DumpAST bool           // Whether to dump the final AST, or assembly source.
}

// NewConfig creates a new, standard configuration instance.
func NewConfig() *Config {
	c := new(Config)
	c.Output = os.Stdout
	c.DumpAST = false

	wd, err := os.Getwd()
	if err == nil {
		c.Include = append(c.Include, path.Clean(wd))
	}

	return c
}
