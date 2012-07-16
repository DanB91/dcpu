// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This is a client/server based development environment
// for DCPU assembly projects.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	addr := parseArgs()
	err := Run(addr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "RunServer: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() string {
	addr := ":7070"

	if v := os.Getenv("DCPU_IDE_ADDRESS"); len(v) > 0 {
		addr = v
	}

	flag.StringVar(&addr, "a", addr, "The HTTP service address on which to run the server.")
	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	return addr
}
