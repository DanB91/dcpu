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
	<-startup()
	shutdown()
}

func startup() <-chan struct{} {
	parseArgs()
	tracker = NewStateTracker(config.Timeout)

	go launchServer(config.Address)
	go launchBrowser(config.Address)

	return tracker.Poll()
}

func shutdown() {
	log.Printf("Idle for %d second(s). Shutting down.", config.Timeout)

	if len(cfgpath) > 0 {
		config.Save(cfgpath)
	}
}

func parseArgs() {
	cfgpath = getConfigPath()
	config = NewConfig()

	flag.StringVar(&cfgpath, "c", cfgpath, "Path to configuration file.")

	timeout := flag.Uint("t", 0, "Shut the server down when it has been idle for t seconds.")
	addr := flag.String("a", "", "The HTTP service address on which to run the server.")
	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if len(cfgpath) > 0 {
		config.Load(cfgpath)
	}

	// These commandline flags take precedence over the values
	// defined in the config file. So if they are set, overwrite 
	// config with their value.
	if len(*addr) > 0 {
		config.Address = *addr
	}

	if *timeout > 0 {
		config.Timeout = *timeout
	}
}
