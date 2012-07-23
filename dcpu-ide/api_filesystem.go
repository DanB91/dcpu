// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {
	Register("/api/readfile", "POST", apiReadFile)
	Register("/api/dirlist", "POST", apiDirList)
}

// apiReadFile attempts to read a given file and returns its contents.
func apiReadFile(r *http.Request) ([]byte, int) {
	file, err := safeFilePath(r.FormValue("tLocation"))
	if err != nil {
		return err, http.StatusNotAcceptable
	}

	data, e := ioutil.ReadFile(file)
	if e != nil {
		return Error(ErrFileRead), http.StatusInternalServerError
	}

	return Pack(data), 200
}

// apiDirList returns the contents of a given directory.
func apiDirList(r *http.Request) ([]byte, int) {
	var recursive bool
	var err error

	if v := r.FormValue("bRecursive"); len(v) > 0 {
		recursive, err = strconv.ParseBool(v)
		if err != nil {
			recursive = false
		}
	}

	dir, e := safeDirPath(r.FormValue("tLocation"))
	if e != nil {
		return e, http.StatusNotAcceptable
	}

	var list []struct {
		Path string
		Dir  bool
	}

	err = filepath.Walk(dir, func(f string, info os.FileInfo, e error) (err error) {
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

	if err != nil && err != io.EOF {
		return Error(ErrUnknown), http.StatusInternalServerError
	}

	return Pack(list), 200
}

func safeFilePath(file string) (string, []byte) {
	file, stat, err := safePath(file)
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		return "", Error(ErrNotFile)
	}

	return file, nil
}

func safeDirPath(file string) (string, []byte) {
	if len(file) == 0 {
		file = config.ProjectPath
	}

	file, stat, err := safePath(file)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", Error(ErrNotDirectory)
	}

	return file, nil
}

func safePath(file string) (string, os.FileInfo, []byte) {
	if len(file) == 0 {
		return "", nil, Error(ErrInvalidPath)
	}

	if file == "/" || !path.IsAbs(file) {
		file = filepath.Join(config.ProjectPath, file)
	}

	if !strings.HasPrefix(file, path.Clean(config.ProjectPath)) {
		return "", nil, Error(ErrInvalidPath)
	}

	stat, err := os.Lstat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, Error(ErrPathNotExist)
		} else {
			return "", nil, Error(ErrUnknown)
		}
	}

	return file, stat, nil
}
