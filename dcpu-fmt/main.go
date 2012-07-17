// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This tool formats DCPU source files according to some predefined styling rules.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var (
	isdir     bool
	infile    string
	outfile   string
	strip     = flag.Bool("s", false, "Strip comments from the input source.")
	tabs      = flag.Bool("tabs", false, "Indent using tabs instead of space.")
	tabwidth  = flag.Uint("tabwidth", 3, "Width of a tab.")
	writefile = flag.Bool("w", false, "Write result to source file instead of stdout.")
)

func main() {
	var err error
	parseArgs()

	if isdir {
		err = filepath.Walk(infile, func(f string, info os.FileInfo, e error) (err error) {
			if e != nil {
				return e
			}

			if info.IsDir() || path.Ext(f) != ".dasm" {
				return nil
			}

			if outfile == Stdout {
				return Format(f, Stdout)
			}

			return Format(f, f)
		})
	} else {
		err = Format(infile, outfile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func parseArgs() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <file>\n", os.Args[0])
		fmt.Printf("   or: cat <file> | %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	outfile = Stdout
	if flag.NArg() > 0 {
		infile = filepath.Clean(flag.Arg(0))

		stat, err := os.Lstat(infile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Input file: %v\n", err)
			os.Exit(1)
		}

		isdir = stat.IsDir()

		if *writefile {
			outfile = infile
		}
	} else {
		infile = Stdin
	}

}
