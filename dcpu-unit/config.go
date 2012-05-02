// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
	"path"
)

// Config holds configuration and state data for the unit tester.
type Config struct {
	Include []string // List of paths where we look to resolve source file references.
	Input   string   // Input source directory.
	Verbose bool     // Verbose test output.
}

// NewConfig creates a new, standard configuration instance.
func NewConfig() *Config {
	c := new(Config)

	wd, err := os.Getwd()
	if err == nil {
		c.Input = path.Clean(wd)
	}

	c.Include = append(c.Include, c.Input)
	return c
}
