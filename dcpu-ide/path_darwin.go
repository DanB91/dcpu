// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
	"os/exec"
)

// Find path to browser.
func getBrowserPath(defaultpath string) string {
	// The $BROWSER environment variable takes precedence.
	file := os.Getenv("BROWSER")
	if len(file) > 0 {
		return file
	}

	file, err := exec.LookPath("open")
	if err == nil {
		return file
	}

	return defaultpath
}
