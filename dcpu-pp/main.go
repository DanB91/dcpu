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
	AppName    = "dcpu-pp"
	AppVersion = "0.1"
)

func main() {
	cfg := parseArgs()
	fmt.Printf("%+v\n", cfg)
}

// process commandline arguments.
func parseArgs() *Config {
	var version, help bool
	var include string

	c := NewConfig()

	flag.BoolVar(&help, "h", false, "Display this help.")
	flag.BoolVar(&version, "v", false, "Display version information.")
	flag.StringVar(&include, "i", "", "Colon-separated list of additional include paths.")
	flag.Parse()

	if version {
		fmt.Fprintf(os.Stdout, "%s %s (Go runtime %s)\nCopyright (c) 2012, Jim Teeuwen.\n",
			AppName, AppVersion, runtime.Version())
		os.Exit(0)
	}

	if help {
		flag.Usage()
		os.Exit(0)
	}

	// Collect source files
	root, err := os.Getwd()
	if err == nil {
		err = filepath.Walk(root, func(file string, info os.FileInfo, err error) error {
			if info.IsDir() || path.Ext(file) != ".dasm" {
				return nil
			}

			c.Input = append(c.Input, file)
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Collect source files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(c.Input) == 0 {
		fmt.Fprintf(os.Stderr, "No source files.\n")
		os.Exit(1)
	}

	// Parse include paths.
	if len(include) > 0 {
		c.Include = strings.Split(include, ":")

		for i := range c.Include {
			c.Include[i] = path.Clean(c.Include[i])
			v := c.Include[i]

			stat, err := os.Lstat(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Test import paths: %v\n", err)
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
