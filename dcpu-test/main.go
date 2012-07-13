// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// DCPU unit-testing framework.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	input    string   // Input source directory.
	includes []string // List of paths where we look to resolve source file references.
	clock    = flag.Int64("c", 1000, "Clock speed in nanoseconds at which to run the tests.")
	profile  = flag.Bool("p", false, "Save profiling data for each test as file.dasm => file.prof.")
	trace    = flag.Bool("t", false, "Print trace output for each instruction as it is executed.")
)

func main() {
	var err error

	parseArgs()
	tests := collectTests()

	for {
		select {
		case file := <-tests:
			if len(file) == 0 {
				return
			}

			t := NewTest(file, includes)

			err = t.Run()

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}
		}
	}
}

var _ = io.Copy

// collectTests traverses the input directory and finds all
// unit test files.
func collectTests() <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		stat, _ := os.Lstat(input)
		if !stat.IsDir() {
			_, name := filepath.Split(input)
			ok, err := filepath.Match("*_test.dasm", name)

			if !ok || err != nil {
				return
			}

			c <- input
			return
		}

		filepath.Walk(input, func(file string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			_, name := filepath.Split(file)
			ok, err := filepath.Match("*_test.dasm", name)
			if !ok || err != nil {
				return err
			}

			parts := strings.Split(file, string(filepath.Separator))

			for i := range parts {
				if len(parts[i]) == 0 {
					continue
				}

				if parts[i][0] == '_' {
					return nil
				}
			}

			c <- file
			return nil
		})
	}()

	return c
}

// process commandline arguments.
func parseArgs() {
	var version, help bool
	var include string

	flag.StringVar(&include, "i", "", "Colon-separated list of additional include paths.")
	flag.BoolVar(&help, "h", false, "Display this help.")
	flag.BoolVar(&version, "v", false, "Display version information.")
	flag.Parse()

	if version {
		fmt.Fprintf(os.Stdout, "%s\n", Version())
		os.Exit(0)
	}

	if help {
		fmt.Fprintf(os.Stdout, "Usage: %s [options] path\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	// See if have an input file.
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No source file.\n")
		os.Exit(1)
	}

	// Ensure we have an existing directory.
	input = filepath.Clean(flag.Arg(0))

	if len(input) == 0 {
		if wd, err := os.Getwd(); err == nil {
			input = wd
		}
	}

	if _, err := os.Lstat(input); err != nil {
		fmt.Fprintf(os.Stderr, "Input path: %v\n", err)
		os.Exit(1)
	}

	includes = append(includes, input)

	// Parse include paths.
	if len(include) > 0 {
		includes = strings.Split(include, ":")

		for i := range includes {
			v := filepath.Clean(includes[i])
			includes[i] = v

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

}
