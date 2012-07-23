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
	files []File

	// Options for Go code generation from arbitrary input files.
	indir   = flag.String("i", "", "Path to input directory.")
	outdir  = flag.String("o", "", "Path to output directory.")
	prefix  = flag.String("p", "", "Prefix generated files with this.")
	tocfile = flag.String("t", "", "File name for index of generated code.")
	dev     = flag.Bool("d", false, "Output code for dev mode.")

	// Options for code generation from templates and data.
	cgdata = flag.String("cgdata", "", "Input data (JSON) file that should be translated to source code.")
	cgin   = flag.String("cgin", "", "Input file with code template.")
	cgout  = flag.String("cgout", "", "Output file for generated code.")
)

func main() {
	var err error

	flag.Parse()

	if len(*cgdata) > 0 && len(*cgin) > 0 && len(*cgout) > 0 {
		if err = generate(*cgdata, *cgin, *cgout); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %v\n", err)
			os.Exit(1)
		}
	}

	if len(*indir) > 0 && len(*outdir) > 0 {
		if err = parseFiles(); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %v\n", err)
			os.Exit(1)
		}

		if err = createTOC(); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %v\n", err)
			os.Exit(1)
		}
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
