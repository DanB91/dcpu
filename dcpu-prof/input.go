// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// pollInput polls for commandline input.
// Commands are sent over the returned channel.
func pollInput() <-chan []string {
	c := make(chan []string)

	go func() {
		defer close(c)

		r := bufio.NewReader(os.Stdin)

		fmt.Printf("%s\n", Version())
		fmt.Printf("Press 'ctrl-C' to exit or 'help' for help.\n")

		for {
			line, _, err := r.ReadLine()
			if err != nil {
				return
			}

			c <- strings.Fields(string(line))
		}

	}()

	return c
}
