// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
)

// Find full path to browser.
func getBrowserPath() string {
	// The $BROWSER environment variable takes precedence.
	file := os.Getenv("BROWSER")
	if len(file) > 0 {
		return file
	}

	return DefaultBrowser
}
