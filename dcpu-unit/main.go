// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	AppName    = "dcpu-unit"
	AppVersion = "0.3.4"
)

func main() {
	var err error

	cfg := parseArgs()
	tests := collectTests(cfg)

	for {
		select {
		case file := <-tests:
			if len(file) == 0 {
				return
			}

			t := NewTest(file, cfg.Include)

			err = t.Run(cfg)

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}
		}
	}
}

// collectTests traverses the input directory and finds all
// unit test files.
func collectTests(cfg *Config) <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		filepath.Walk(cfg.Input, func(file string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			_, name := filepath.Split(file)
			ok, err := filepath.Match("*_test.dasm", name)

			if !ok || err != nil {
				return err
			}

			c <- file
			return nil
		})
	}()

	return c
}

// process commandline arguments.
func parseArgs() *Config {
	var version, help bool
	var include string

	c := NewConfig()

	flag.Int64Var(&c.Clock, "c", c.Clock, "Clock speed in nanoseconds at which to run the tests.")
	flag.BoolVar(&help, "h", false, "Display this help.")
	flag.StringVar(&include, "i", "", "Colon-separated list of additional include paths.")
	flag.BoolVar(&c.Trace, "t", false, "Print trace output for each instruction as it is executed.")
	flag.BoolVar(&version, "v", false, "Display version information.")
	flag.BoolVar(&c.Verbose, "V", false, "Print additional debug output.")
	flag.Parse()

	if version {
		fmt.Fprintf(os.Stdout, "%s %s (Go runtime %s)\nCopyright (c) 2012, Jim Teeuwen.\n",
			AppName, AppVersion, runtime.Version())
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

	c.Input = path.Clean(flag.Arg(0))

	// Ensure we have an existing directory.
	if stat, err := os.Lstat(c.Input); err != nil {
		fmt.Fprintf(os.Stderr, "Input path: %v\n", err)
		os.Exit(1)
	} else if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Input path %q is not a directory.\n", c.Input)
		os.Exit(1)
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
