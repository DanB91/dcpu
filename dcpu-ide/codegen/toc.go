// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
)

func createTOC() (err error) {
	fd, err := os.Create(*tocfile)
	if err != nil {
		return
	}

	defer fd.Close()

	fmt.Fprintln(fd, `
// This file was automatically generated.
// Any changes to it will not be preserved.

package main

const DefaultPage = "/index.html"

type File struct {
	Data        func() []byte
	ContentType string
}

var static map[string]*File

func init() {
	static = make(map[string]*File)`)

	for i := range files {
		fmt.Fprintf(fd, `
	static[%q] = &File {
		Data: %s,
		ContentType: %q,
	}`, files[i].Path, files[i].Var, files[i].Type)
	}

	fmt.Fprintln(fd, `
}`)
	return
}
