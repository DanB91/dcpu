// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This tool queries and analyses profiling data from a given input file.
package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/dcpu/prof"
	"io"
	"os"
	"path/filepath"
)

func main() {
	prof := parseArgs()

	if prof.CountUses() == 0 {
		fmt.Fprintln(os.Stdout, "Profile has no sample data.\n")
		os.Exit(0)
	}

	input := pollInput()

	for {
		select {
		case cmd := <-input:
			err := Handle(prof, cmd)

			if err != nil {
				if err == io.EOF {
					return
				}

				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		}
	}
}

func parseArgs() *prof.Profile {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <file>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	help := flag.Bool("h", false, "Display this help.")
	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "No input file.")
		os.Exit(1)
	}

	fd, err := os.Open(filepath.Clean(flag.Arg(0)))

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defer fd.Close()

	p, err := prof.Read(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read: %v\n", err)
		os.Exit(1)
	}

	return p
}
