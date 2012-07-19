// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
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

	// If xdg-open exists, we should use it instead.
	// It will find the default browser for us.
	file, err := exec.LookPath("xdg-open")
	if err == nil {
		return file
	}

	return defaultpath
}

// Find the location of the configuration file.
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
