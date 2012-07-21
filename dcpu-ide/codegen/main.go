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

	// Options for javascript & Go generation from templates and data.
	data  = flag.String("data", "", "Input data (JSON) file that should be translated to javascript and/or Go.")
	jsin  = flag.String("jsin", "", "Input javascript template.")
	jsout = flag.String("jsout", "", "Output generated javascript file.")
	goin  = flag.String("goin", "", "Input Go template.")
	goout = flag.String("goout", "", "Output generated Go file.")
)

func main() {
	var err error

	flag.Parse()

	if len(*data) > 0 {
		if len(*jsin) > 0 && len(*jsout) > 0 {
			if err = generate(*data, *jsin, *jsout); err != nil {
				fmt.Fprintf(os.Stderr, "[e] %s: %v\n", *jsin, err)
				os.Exit(1)
			}
		}

		if len(*goin) > 0 && len(*goout) > 0 {
			if err = generate(*data, *goin, *goout); err != nil {
				fmt.Fprintf(os.Stderr, "[e] %s: %v\n", *goin, err)
				os.Exit(1)
			}
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
