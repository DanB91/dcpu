// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
	"os/user"
	"path"
	"runtime"
)

// Find the user's home directory.
func getHomeDir() string {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("USERPROFILE")

		if len(dir) == 0 {
			hd := os.Getenv("HOMEDRIVE")
			hp := os.Getenv("HOMEPATH")

			if len(hd) > 0 && len(hp) > 0 {
				dir = path.Join(hd, hp)
			}
		}

		if len(dir) == 0 {
			// This is only relevant if we are running under cygwin.
			dir = os.Getenv("HOME")
		}

	case "freebsd", "linux", "darwin":
		dir = os.Getenv("HOME")

		if len(dir) == 0 {
			usr, err := user.Current()
			if err == nil {
				dir = path.Join("/home", usr.Username)
			}
		}
	}

	return dir
}

// Find suitable location for the configuration file.
func getConfigPath(file string) string {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LOCALAPPDATA")

		if len(dir) == 0 {
			dir = os.Getenv("APPDATA")
		}

		if len(dir) == 0 {
			hd := os.Getenv("HOMEDRIVE")
			hp := os.Getenv("HOMEPATH")

			if len(hd) > 0 && len(hp) > 0 {
				dir = path.Join(hd, hp)
			}
		}

		if len(dir) == 0 {
			// This is only relevant if we are running under cygwin.
			dir = os.Getenv("HOME")
		}

	case "freebsd", "linux", "darwin":
		dir = os.Getenv("HOME")

		if len(dir) == 0 {
			usr, err := user.Current()
			if err == nil {
				dir = path.Join("/home", usr.Username)
			}
		}
	}

	if len(dir) == 0 {
		return ""
	}

	switch runtime.GOOS {
	case "freebsd", "linux", "darwin":
		file = "." + file
	case "windows":
		file += ".cfg"
	}

	return path.Join(dir, file)
}
