// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bufio"
	"github.com/jteeuwen/dcpu/cpu"
	"os"
)

var sourceCache map[cpu.Word]string

func init() {
	sourceCache = make(map[cpu.Word]string)
}

// GetSourceLines returns a range of lines from the given file.
func GetSourceLines(file string, start, end int) []string {
	fd, err := os.Open(file)
	if err != nil {
		return nil
	}

	defer fd.Close()

	r := bufio.NewReader(fd)

	var lines []string
	var data []byte

	count := 1 // line number start at 1, not 0.

	for {
		data, _, err = r.ReadLine()

		if count >= start && count <= end {
			lines = append(lines, string(data))
		}

		if err != nil || count > end {
			break
		}

		count++
	}

	return lines
}
