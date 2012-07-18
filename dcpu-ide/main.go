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

var (
	config  *Config
	tracker *StateTracker
)

func main() {
	parseArgs()

	tracker = NewStateTracker(config.Timeout)
	go tracker.Poll()

	err := Run(config.Address)

	if err != nil {
		fmt.Fprintf(os.Stderr, "RunServer: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() {
	config = NewConfig()

	flag.UintVar(&config.Timeout, "t", config.Timeout,
		"Shut the server down when it has been idle for t seconds.")

	flag.StringVar(&config.Address, "a", config.Address,
		"The HTTP service address on which to run the server.")

	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}
}
