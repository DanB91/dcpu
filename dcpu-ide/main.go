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
	nobrowser := parseArgs()
	tracker = NewStateTracker(config.Timeout)

	go launchServer(config.Address)

	if !nobrowser {
		go launchBrowser(config.Address)
	}

	return tracker.Poll()
}

func shutdown() {
	log.Printf("Idle for %d second(s). Shutting down.", config.Timeout)

	if len(cfgpath) > 0 {
		log.Printf("Saving %s", cfgpath)
		config.Save(cfgpath)
	}
}

func parseArgs() bool {
	cfgpath = getConfigPath(AppName)
	config = NewConfig()

	flag.StringVar(&cfgpath, "c", cfgpath, "Path to configuration file.")

	addr := flag.String("a", "", "The HTTP service address on which to run the server.")
	projdir := flag.String("p", "", "Path to directory where DASM projects are stored.")
	timeout := flag.Uint("t", 0, "Shut the server down when it has been idle for t seconds.")
	version := flag.Bool("v", false, "Display version information.")
	nobrowser := flag.Bool("n", false, "Do not launch a browser.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if len(cfgpath) > 0 {
		log.Printf("Loading %s", cfgpath)
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

	if len(*projdir) > 0 {
		config.ProjectPath = *projdir
	}

	return *nobrowser
}
