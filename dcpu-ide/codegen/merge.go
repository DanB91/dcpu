// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
	"os"
	"path/filepath"
)

// merge merges all file contents from the input directory and
// writes it to a single output file.
func merge(in, out string) (err error) {
	fout, err := os.Create(out)
	if err != nil {
		return
	}

	defer fout.Close()

	err = filepath.Walk(in, func(fin string, info os.FileInfo, e error) (err error) {
		if e != nil || info.IsDir() {
			return e
		}

		fd, err := os.Open(fin)
		if err != nil {
			return
		}

		defer fd.Close()

		_, err = io.Copy(fout, fd)
		if err != nil {
			return
		}

		_, err = fout.Write([]byte{'\n'})
		return
	})

	return
}
