// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type KeyValuePair struct {
	Key   string
	Value interface{}
}

type CodeData struct {
	Const []KeyValuePair // Constant definitions.
	Vars  []KeyValuePair // Variable definitions.
}

// generate reads data entries from the data file.
// It passes the data through the Go template defined in infile,
// and writes the result to outfile.
func generate(data, infile, outfile string) (err error) {
	bytes, err := ioutil.ReadFile(data)
	if err != nil {
		return
	}

	var cd CodeData

	err = json.Unmarshal(bytes, &cd)
	if err != nil {
		return
	}

	t, err := template.New("page").ParseFiles(infile)
	if err != nil {
		return
	}

	fd, err := os.Create(outfile)
	if err != nil {
		return
	}

	defer fd.Close()
	return t.ExecuteTemplate(fd, filepath.Base(infile), cd)
}
