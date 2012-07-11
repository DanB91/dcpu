// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This tool queries and analyses profiling data from a given input file.
package main

import (
	"bufio"
	"bytes"
	"os"
)

// pollInput polls for commandline input.
// Commands are sent over the returned channel.
func pollInput() <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		r := bufio.NewReader(os.Stdin)

		for {
			line, _, err := r.ReadLine()
			if err != nil {
				return
			}

			c <- string(bytes.TrimSpace(line))
		}

	}()

	return c
}
