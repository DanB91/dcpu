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

type CodeData struct {
	Errors []struct {
		Name   string
		String string
		Code   int
	}
	Stylesheets []string
	Scripts     []string
}

var (
	files []File

	// Misc flag.
	dev  = flag.Bool("d", false, "Output for debug mode.")
	data = flag.String("data", "", "Input data (JSON) for merging and template parsing.")

	// Options for file coversion to Go code.
	convin  = flag.String("convin", "", "Path to conversion input directory.")
	convout = flag.String("convout", "", "Path to conversion output directory.")
	convpre = flag.String("convpre", "", "Prefix generated conversion files with this.")
	convtoc = flag.String("convtoc", "", "File name for index of converted code.")

	// Options for code generation from templates and data.
	cgin  = flag.String("cgin", "", "Input file with code template.")
	cgout = flag.String("cgout", "", "Output file for generated code.")

	// Merge files.
	mergetype = flag.String("mergetype", "", "Merge type: js or css.")
	mergeout  = flag.String("mergeout", "", "Output file for file merge.")
	mergepre  = flag.String("mergepre", "", "Path prefix for each merge target.")
)

func main() {
	var err error

	flag.Parse()

	if len(*data) > 0 && len(*mergetype) > 0 && len(*mergeout) > 0 {
		if err = merge(*data, *mergetype, *mergepre, *mergeout); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %v\n", err)
			os.Exit(1)
		}
	}

	if len(*data) > 0 && len(*cgin) > 0 && len(*cgout) > 0 {
		if err = generate(*data, *cgin, *cgout); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %v\n", err)
			os.Exit(1)
		}
	}

	if len(*convin) > 0 && len(*convout) > 0 {
		err = filepath.Walk(*convin, func(fin string, info os.FileInfo, e error) (err error) {
			if e != nil || info.IsDir() {
				return e
			}

			return translate(fin)
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "[e] %v\n", err)
			os.Exit(1)
		}

		if err = createTOC(*convtoc, files); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %v\n", err)
			os.Exit(1)
		}
	}
}
