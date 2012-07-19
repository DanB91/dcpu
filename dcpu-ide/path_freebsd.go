// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

// Find path to browser.
func getBrowserPath(defaultpath string) string {
	// The $BROWSER environment variable takes precedence.
	file := os.Getenv("BROWSER")
	if len(file) > 0 {
		return file
	}

	return defaultpath
}

// Find suitable location for the configuration file.
// For unix systems this is usually in $HOME.
func getConfigPath() string {
	var file string
	var err error
	var usr *user.User

	if file = os.Getenv("HOME"); len(file) > 0 {
		goto ret
	}

	if usr, err = user.Current(); err == nil {
		file = path.Join("/home", usr.Username)
		goto ret
	}

	return ""

ret:
	return path.Join(file, fmt.Sprintf(".%s", AppName))
}
