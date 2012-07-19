// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"
)

// Find suitable location for the configuration file.
func getConfigPath() string {
	var f string

	switch runtime.GOOS {
	case "windows":
		f = os.Getenv("LOCALAPPDATA")

		if len(f) == 0 {
			f = os.Getenv("APPDATA")
		}

		if len(f) == 0 {
			hd := os.Getenv("HOMEDRIVE")
			hp := os.Getenv("HOMEPATH")

			if len(hd) > 0 && len(hp) > 0 {
				f = path.Join(hd, hp)
			}
		}

		if len(f) == 0 {
			// This is only relevant if we are running under cygwin.
			f = os.Getenv("HOME")
		}

	case "freebsd", "linux", "darwin":
		f = os.Getenv("HOME")

		if len(f) == 0 {
			usr, err := user.Current()
			if err == nil {
				f = path.Join("/home", usr.Username)
			}
		}
	}

	if len(f) == 0 {
		return ""
	}

	return path.Join(f, fmt.Sprintf("%s.cfg", AppName))
}
