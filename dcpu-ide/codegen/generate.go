// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// generate reads data entries from the data file.
// It passes the data through the Go template defined in infile,
// and writes the result to outfile.
func generate(data, infile, outfile string) (err error) {
	bytes, err := ioutil.ReadFile(data)
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %v", err)
	}

	var cd CodeData

	err = json.Unmarshal(bytes, &cd)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if !*dev {
		// In release mode, we'll be having all style and script
		// files merged into one.
		cd.Scripts = []string{"/app.js"}
		cd.Stylesheets = []string{"/app.css"}
	}

	t, err := template.New("page").ParseFiles(infile)
	if err != nil {
		return fmt.Errorf("template.ParseFiles: %v", err)
	}

	fd, err := os.Create(outfile)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	defer fd.Close()
	err = t.ExecuteTemplate(fd, filepath.Base(infile), cd)
	if err != nil {
		err = fmt.Errorf("template.Execute: %v", err)
	}

	return
}
