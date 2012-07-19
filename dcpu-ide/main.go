// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This is a client/server based development environment
// for DCPU assembly projects.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	cfgpath string
	config  *Config
	tracker *StateTracker
)

func main() {
	parseArgs()

	tracker = NewStateTracker(config.Timeout)

	quit := tracker.Poll()

	go launchServer(config.Address)
	go launchBrowser(config.Address)

	<-quit

	shutdown()
}

func shutdown() {
	log.Printf("Shutting down.")
	config.Save(cfgpath)
}

func parseArgs() {
	cfgpath = getConfigPath()
	config = NewConfig()

	flag.StringVar(&cfgpath, "c", cfgpath, "Path to configuration file.")
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

	config.Load(cfgpath)
}
