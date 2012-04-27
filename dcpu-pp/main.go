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

	c := NewConfig()

	flag.BoolVar(&help, "h", false, "Display this help.")
	flag.BoolVar(&version, "v", false, "Display version information.")
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

	return c
}
