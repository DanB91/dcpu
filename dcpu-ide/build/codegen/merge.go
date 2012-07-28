// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// merge merges all file contents from the input directory and
// writes it to a single output file.
func merge(data, mtype, prefix, out string) (err error) {
	bytes, err := ioutil.ReadFile(data)
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %v", err)
	}

	var cd CodeData

	err = json.Unmarshal(bytes, &cd)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	fout, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	defer fout.Close()

	var list []string
	switch mtype {
	case "js":
		list = cd.Scripts
	case "css":
		list = cd.Stylesheets
	}

	var file string
	for i := range list {
		file = filepath.Join(prefix, list[i])

		err = mergeFile(fout, file)
		if err != nil {
			return fmt.Errorf("mergeFile: %v", err)
		}

		err = os.Remove(file)
		if err != nil {
			return fmt.Errorf("os.Remove: %v", err)
		}
	}

	return
}

func mergeFile(fout io.Writer, in string) (err error) {
	fin, err := os.Open(in)
	if err != nil {
		return
	}

	defer fin.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		return
	}

	_, err = fout.Write([]byte{'\n'})
	return
}
