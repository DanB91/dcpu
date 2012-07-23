// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {
	Register("/api/dirlist", "POST", apiDirList)
}

// apiDirList returns the contents of a given directory.
func apiDirList(r *http.Request) ([]byte, int) {
	var recursive bool
	var err error

	root := r.FormValue("tLocation")

	if v := r.FormValue("bRecursive"); len(v) > 0 {
		recursive, err = strconv.ParseBool(v)
		if err != nil {
			recursive = false
		}
	}

	if len(root) == 0 {
		root = config.ProjectPath
	}

	if !path.IsAbs(root) {
		root = filepath.Join(config.ProjectPath, root)
	}

	if !strings.HasPrefix(root, config.ProjectPath) {
		root = config.ProjectPath
	}

	stat, err := os.Lstat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return Error(ErrPathNotExist), http.StatusNotAcceptable
		} else {
			return Error(ErrUnknown), http.StatusInternalServerError
		}
	}

	if !stat.IsDir() {
		return Error(ErrNotDirectory), http.StatusNotAcceptable
	}

	var list []struct {
		Path string
		Dir  bool
	}

	err = filepath.Walk(root, func(f string, info os.FileInfo, e error) (err error) {
		if e != nil {
			return e
		}

		list = append(list, struct {
			Path string
			Dir  bool
		}{
			f, info.IsDir(),
		})

		if info.IsDir() && !recursive {
			err = io.EOF
		}

		return
	})

	return Pack(list), 200
}
