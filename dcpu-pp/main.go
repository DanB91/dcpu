// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
	"runtime"
	"flag"
)

const (
	AppName = "dcpu-pp"
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

	return c
}
