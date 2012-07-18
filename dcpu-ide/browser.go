// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"log"
	"os/exec"
	"strings"
)

const DefaultBrowser = "chromium"

// launchBrowser attempts to open a webbrowser in a platform compatible manner.
// It will attempt various standard modes of opening webpages for various
// platforms and fall back to environment variables like $BROWSER.
func launchBrowser(url string) {
	var args []string

	if len(url) == 0 {
		url = "localhost"
	}

	if url[0] == ':' {
		// We have only a port number.
		url = "localhost" + url
	}

	url = "http://" + url

	app := getBrowserPath()
	if strings.Index(app, "chromium") != -1 {
		// Chromium can run in 'app mode'.
		// This gets rid of url bars, buttons, etc.
		// This is more suitable to our needs.
		args = []string{"--app=" + url}
	} else {
		args = []string{url}
	}

	log.Printf("Launching %s...", app)

	err := exec.Command(app, args...).Run()
	if err != nil {
		log.Printf("Failed to launch a suitable browser.")
		log.Printf("Either open the url %q manually, or set the $BROWSER "+
			"environment variable and run this program again.",
			url)
	}
}
