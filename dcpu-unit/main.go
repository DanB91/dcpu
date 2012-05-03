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
	"sync"
)

const (
	AppName    = "dcpu-unit"
	AppVersion = "0.1.0"
)

func main() {
	var wg sync.WaitGroup

	cfg := parseArgs()
	log := NewLog(os.Stdout)
	tests := collectTests(cfg)

	defer log.Close()
	defer wg.Wait()

	for {
		select {
		case file := <-tests:
			if len(file) == 0 {
				return
			}

			wg.Add(1)
			go runTest(file, cfg.Include, &wg, log)
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
			if info.IsDir() || path.Ext(file) != ".test" {
				return nil
			}

			// We need a matching compare file.
			stat, err := os.Stat(file[:len(file)-5] + ".cmp")
			if err != nil || stat.IsDir() {
				return nil
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

	flag.BoolVar(&help, "h", false, "Display this help.")
	flag.StringVar(&include, "i", "", "Colon-separated list of additional include paths.")
	flag.BoolVar(&version, "v", false, "Display version information.")
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
