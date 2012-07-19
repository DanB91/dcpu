// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
)

// Find path to browser.
func getBrowserPath(defaultpath string) string {
	// The $BROWSER environment variable takes precedence.
	file := os.Getenv("%BROWSER%")
	if len(file) > 0 {
		return file
	}

	return defaultpath + ".exe"
}

// Find suitable location for the configuration file.
func getConfigPath() string {
	file := os.Getenv("%LOCALAPPDATA%")
	if len(file) > 0 {
		goto ret
	}

	file = os.Getenv("%APPDATA%")
	if len(file) > 0 {
		goto ret
	}

	hd := os.Getenv("%HOMEDRIVE%")
	hp := os.Getenv("%HOMEPATH%")
	if len(hd) > 0 && len(hp) > 0 {
		file = path.Join(hd, hp)
		goto ret
	}

	return ""

ret:
	return path.Join(file, fmt.Sprintf("%s.cfg", AppName))
}
