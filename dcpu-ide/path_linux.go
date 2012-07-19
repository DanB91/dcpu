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

	// If xdg-open exists, we should use it instead.
	// It will find the default browser for us.
	file, err := exec.LookPath("xdg-open")
	if err == nil {
		return file
	}

	return defaultpath
}
