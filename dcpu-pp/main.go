// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
	var buf bytes.Buffer
	var err error

	cfg := parseArgs()

	if err = collectFiles(&buf, cfg.Input); err != nil {
		fmt.Fprintf(os.Stderr, "Collect source: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", buf.String())
}

// collectFiles takes the input files and loads their contents
// into the supplied buffer.
func collectFiles(w io.Writer, files []string) (err error) {
	var fd *os.File

	for i := range files {
		if fd, err = os.Open(files[i]); err != nil {
			return
		}

		io.Copy(w, fd)
		fd.Close()

		// Ensure file content is delimited by a newline.
		w.Write([]byte{'\n'})
	}

	return
}

// process commandline arguments.
func parseArgs() *Config {
	var version, help bool
	var include string

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

	c := NewConfig()

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
			v := path.Clean(c.Include[i])
			c.Include[i] = v

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
