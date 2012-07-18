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

	quit := tracker.Poll()

	go startServer(config.Address)
	go launchBrowser(config.Address)

	<-quit
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
