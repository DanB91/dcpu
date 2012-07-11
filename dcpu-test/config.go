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
	Profile string   // Name of profile output file.
	Clock   int64    // Clockspeed at which to run the tests.
	Trace   bool     // Print trace data for each instruction as it is executed.
	Verbose bool     // Print additional debug output.
}

// NewConfig creates a new, standard configuration instance.
func NewConfig() *Config {
	c := new(Config)
	c.Clock = 1000 // 100khz.

	wd, err := os.Getwd()
	if err == nil {
		c.Input = path.Clean(wd)
	}

	c.Include = append(c.Include, c.Input)
	return c
}
