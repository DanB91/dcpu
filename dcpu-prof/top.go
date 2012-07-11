// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/prof"
	"path/filepath"
	"strings"
)

func top(prof *prof.Profile, count int, sort string) {
	var list SampleList

	for i, v := range prof.Usage {
		if v != nil {
			list = append(list, Sample{PC: cpu.Word(i), Data: v})
		}
	}

	switch strings.ToLower(sort) {
	case "cost":
		SamplesByCost(list).Sort()

	case "cumulative":
		SamplesByCumulativeCost(list).Sort()

	case "count":
		SamplesByCount(list).Sort()
	}

	if len(list) > count {
		list = list[:count]
	}

	var filename, source string

	fmt.Println("      COUNT | COST |  CUM. COST |                 FILE | SOURCE")
	fmt.Println("============================================================================")
	for i := range list {
		filename = prof.Files[list[i].Data.File]
		source = getSourceLine(prof, list[i].PC)

		filename = filepath.Base(filename)
		filename = fmt.Sprintf("%s:%d", filename, list[i].Data.Line)

		if len(filename) > 20 {
			filename = "..." + filename[3:]
		}

		fmt.Printf(" %10d | %4d | %10d | %20s | %s\n",
			list[i].Data.Count,
			list[i].Data.Cost(),
			list[i].Data.CumulativeCost(),
			filename,
			source,
		)
	}
}
