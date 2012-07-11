// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bufio"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/prof"
	"os"
	"strings"
)

var sourceCache map[cpu.Word]string

func init() {
	sourceCache = make(map[cpu.Word]string)
}

func getSourceLine(prof *prof.Profile, pc cpu.Word) string {
	if line, ok := sourceCache[pc]; ok {
		return line
	}

	fileno := prof.Usage[pc].File
	lineno := prof.Usage[pc].Line
	file := prof.Files[fileno]

	fd, err := os.Open(file)
	if err != nil {
		sourceCache[pc] = ""
		return ""
	}

	defer fd.Close()

	r := bufio.NewReader(fd)

	var count int
	var data []byte

	for {
		if data, _, err = r.ReadLine(); err != nil {
			sourceCache[pc] = ""
			return ""
		}

		if count < lineno-1 {
			count++
			continue
		}

		sourceCache[pc] = strings.TrimSpace(string(data))
		break
	}

	return sourceCache[pc]
}
