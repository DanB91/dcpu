// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	Path string // File path
	Var  string // Variable name by which we access the data.
	Type string // Content/mime type
}

var (
	files   []File
	indir   = flag.String("i", "", "Path to input directory.")
	outdir  = flag.String("o", "", "Path to output directory.")
	prefix  = flag.String("p", "", "Prefix generated files with this.")
	tocfile = flag.String("t", "", "File name for index of generated code.")
	dev     = flag.Bool("d", false, "Output code for dev mode.")
)

func main() {
	flag.Parse()

	err := parseFiles()

	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		os.Exit(1)
	}

	err = createTOC()

	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		os.Exit(1)
	}
}

func parseFiles() error {
	return filepath.Walk(*indir, func(fin string, info os.FileInfo, e error) (err error) {
		if e != nil || info.IsDir() {
			return e
		}

		return translate(fin)
	})
}
