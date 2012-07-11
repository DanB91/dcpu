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
	_ = prof

	input := pollInput()

	for {
		select {
		case str := <-input:
			if len(str) == 0 {
				return
			}

			println(">", str)
		}
	}
}

func parseArgs() *prof.Profile {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <file>\n", os.Args[0])
		fmt.Printf("   or: cat <file> | %s [options]\n\n", os.Args[0])
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

	var input io.Reader

	// See if have an input file. If not, read data from stdin.
	if flag.NArg() == 0 {
		input = os.Stdin
	} else {
		fd, err := os.Open(filepath.Clean(flag.Arg(0)))

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		defer fd.Close()
		input = fd
	}

	p, err := prof.Read(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read: %v\n", err)
		os.Exit(1)
	}

	return p
}
