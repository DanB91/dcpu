// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	AppName    = "dcpu-pp"
	AppVersion = "0.4"
)

func main() {
	var err error
	var ast AST

	cfg := parseArgs()
	defer cfg.Output.Close()

	if err = parseInput(&ast, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err = Process(&ast); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if cfg.DumpAST {
		writeAst(cfg.Output, &ast)
	} else {
		writeSource(cfg.Output, &ast)
	}
}

// process commandline arguments.
func parseArgs() *Config {
	var version, help bool
	var include, output string
	var err error

	c := NewConfig()

	flag.BoolVar(&c.DumpAST, "a", false, "Dump the source code parse tree to stdout.")
	flag.BoolVar(&help, "h", false, "Display this help.")
	flag.StringVar(&include, "i", "", "Colon-separated list of additional include paths.")
	flag.StringVar(&output, "o", "", "Name of destination file. Defaults to stdout.")
	flag.BoolVar(&version, "v", false, "Display version information.")

	CreateProcessorFlags()

	flag.Parse()

	if version {
		fmt.Fprintf(os.Stdout, "%s %s (Go runtime %s)\nCopyright (c) 2012, Jim Teeuwen.\n",
			AppName, AppVersion, runtime.Version())
		os.Exit(0)
	}

	if help {
		fmt.Fprintf(os.Stdout, "Usage: %s [options] file\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	// See if have an input file.
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No source file.\n")
		os.Exit(1)
	}

	c.Input = path.Clean(flag.Arg(0))

	// A valid output path?
	if len(output) > 0 {
		output = path.Clean(output)
		c.Output, err = os.Create(output)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Output file: %v\n", err)
			os.Exit(1)
		}
	}

	// Parse include paths.
	if len(include) > 0 {
		c.Include = strings.Split(include, ":")

		for i := range c.Include {
			v := path.Clean(c.Include[i])
			c.Include[i] = v

			stat, err := os.Lstat(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to stat %q: %v\n", v, err)
				os.Exit(1)
			}

			if !stat.IsDir() {
				fmt.Fprintf(os.Stderr, "Import path %q is not a directory.\n", v)
				os.Exit(1)
			}
		}
	}

	return c
}
