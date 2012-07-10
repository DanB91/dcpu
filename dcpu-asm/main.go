// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Comprehensive DCPU assembler.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/parser"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	infile       string
	includes     []string
	dumpast      = flag.Bool("a", false, "Dump pre-processed AST to the output.")
	dumpsrc      = flag.Bool("s", false, "Dump pre-processed source code to the output.")
	outfile      = flag.String("o", "", "Name of the output file.")
	debugfile    = flag.String("d", "", "Name of the file to write debug symbols to.")
	littleendian = flag.Bool("l", false, "Generate Little Endian binary output. Defaults to Big Endian.")
)

func main() {
	parseArgs()

	// Collect all the source code into the given AST.
	// This takes care of resolving includes and identifying
	// unresolved label references.
	var ast parser.AST

	err := parseInput(&ast, infile, includes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Source reader: %v\n", err)
		os.Exit(1)
	}

	// Run pre-processors on the generated AST.
	err = PreProcess(&ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pre processor: %v\n", err)
		os.Exit(1)
	}

	// Dump AST or source code if necessary.
	if *dumpast || *dumpsrc {
		err = writeSource(&ast, *outfile, *dumpast)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Source writer: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Assemble program.
	program, dbg, err := asm.Assemble(&ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Assembler: %v\n", err)
		os.Exit(1)
	}

	// Run post-processors on generated binary code and debug symbols.
	err = PostProcess(program, dbg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Post processor: %v\n", err)
		os.Exit(1)
	}

	// Write debug file.
	if err = writeDebug(dbg, *debugfile); err != nil {
		fmt.Fprintf(os.Stderr, "Debug writer: %v\n", err)
		os.Exit(1)
	}

	// Write binary output.
	if err = writeProgram(program, *outfile); err != nil {
		fmt.Fprintf(os.Stderr, "Binary writer: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() {
	include := flag.String("i", "", "Colon separated list of additional include paths.")
	help := flag.Bool("h", false, "Display this help.")
	version := flag.Bool("v", false, "Display version information.")

	CreatePreProcessorFlags()
	CreatePostProcessorFlags()

	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// See if have an input file. If not, read it from stdin.
	if flag.NArg() == 0 {
		var b bytes.Buffer
		io.Copy(&b, os.Stdin)
		infile = filepath.Clean(strings.TrimSpace(b.String()))
	} else {
		infile = filepath.Clean(flag.Arg(0))
	}

	// Ensure we have an existing file.
	if stat, err := os.Lstat(infile); err != nil {
		fmt.Fprintf(os.Stderr, "Input path: %v\n", err)
		os.Exit(1)
	} else if stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Input path %q is not a file.\n", infile)
		os.Exit(1)
	}

	// A valid output path?
	if stat, err := os.Lstat(*outfile); err != nil {
		if os.IsExist(err) && stat.IsDir() {
			fmt.Fprintf(os.Stderr, "Output file %q exists and is not a file.\n", *outfile)
			os.Exit(1)
		}
	}

	// Parse include paths.
	if len(*include) > 0 {
		includes = strings.Split(*include, ":")

		for i := range includes {
			includes[i] = filepath.Clean(includes[i])

			stat, err := os.Lstat(includes[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to stat %q: %v\n", includes[i], err)
				os.Exit(1)
			}

			if !stat.IsDir() {
				fmt.Fprintf(os.Stderr, "Import path %q is not a directory.\n", includes[i])
				os.Exit(1)
			}
		}
	}

	wd, err := os.Getwd()
	if err == nil {
		includes = append(includes, filepath.Clean(wd))
	}
}

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: %s [options] <file>\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "   or: echo <file> | %s [options]\n\n", os.Args[0])

	fmt.Fprintf(os.Stdout, "[Misc options]\n")
	fmt.Fprintf(os.Stdout, " -o <file> : Path to output. Defaults to stdout.\n")
	fmt.Fprintf(os.Stdout, " -d <file> : Path to debug symbol file.\n")
	fmt.Fprintf(os.Stdout, "        -a : Dump pre-processed AST to the output.\n")
	fmt.Fprintf(os.Stdout, "        -s : Dump pre-processed source code to the output.\n")
	fmt.Fprintf(os.Stdout, "        -l : Generate Little Endian binary output. Defaults to Big Endian.\n")
	fmt.Fprintf(os.Stdout, "        -h : Display this help.\n")
	fmt.Fprintf(os.Stdout, "        -v : Display version information.\n")

	fmt.Fprintf(os.Stdout, "\n  The -a and -s options are mutually exclusive.\n")
	fmt.Fprintf(os.Stdout, "  -d and -l have no effect in combination with -a or -s.\n")

	if len(preprocessors) > 0 {
		fmt.Fprintf(os.Stdout, "\nPre-processors operate on the generated AST.\n")
		fmt.Fprintf(os.Stdout, "Available processors are:\n\n")

		for k, v := range preprocessors {
			fmt.Fprintf(os.Stdout, " -%s : %s\n", k, v.desc)
		}
	}

	if len(postprocessors) > 0 {
		fmt.Fprintf(os.Stdout, "\nPost-processors operate on the generated binary code:\n")
		fmt.Fprintf(os.Stdout, "Available processors are:\n\n")

		for k, v := range postprocessors {
			fmt.Fprintf(os.Stdout, " -%s : %s\n", k, v.desc)
		}
	}
}
