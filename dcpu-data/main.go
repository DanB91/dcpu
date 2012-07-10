// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This tool generates DCPU code from any input file.
// Used to embed data files in a DCPU program.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	input        io.ReadCloser
	output       io.WriteCloser
	label        = flag.String("l", "", "Name of the label that should be written before the actual data.")
	littleendian = flag.Bool("LE", false, "Create output in Little Endian format. Defaults to Big Endian.")
	pad          = flag.Bool("p", false, "Pad data with a NULL word.")
)

func main() {
	parseArgs()

	defer input.Close()
	defer output.Close()

	if len(*label) > 0 {
		fmt.Fprintf(output, ":%s\n", *label)
	}

	// Read input data.
	var b bytes.Buffer
	io.Copy(&b, input)

	// Ensure we have even number of bytes.
	if b.Len()%2 != 0 {
		b.Write([]byte{0})
	}

	// Apply padding if necessary.
	if *pad {
		b.Write([]byte{0, 0})
	}

	// Write assembly source.
	err := WriteWords(output, b.Bytes(), *littleendian)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Finish with a newline.
	_, err = output.Write([]byte{'\n'})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func parseArgs() {
	outfile := flag.String("o", "", "path to output file. Defaults to stdout.")
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

	// See if have an input file. If not, read data from stdin.
	if flag.NArg() == 0 {
		input = os.Stdin
	} else {
		fd, err := os.Open(filepath.Clean(flag.Arg(0)))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Input path: %v\n", err)
			os.Exit(1)
		}

		input = fd
	}

	if len(*outfile) == 0 {
		output = os.Stdout
	} else {
		var err error
		output, err = os.Create(*outfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Output path: %v\n", err)
			input.Close()
			os.Exit(1)
		}
	}
}
