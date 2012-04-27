// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
)

// Config holds configuration and state data for the preprocessor.
type Config struct {
	Path   []string // List of paths where we look to resolve source file references.
	Input  []string // Names of input source files.
	Output string   // Name of output source file. Defaults to nil (stdout).
}

// NewConfig creates a new, standard configuration instance.
func NewConfig() *Config {
	c := new(Config)

	wd, err := os.Getwd()
	if err == nil {
		c.Path = append(c.Path, wd)
	}

	return c
}
