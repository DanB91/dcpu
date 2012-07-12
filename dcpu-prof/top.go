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
	var filename, source string
	var countpercent, costpercent float64
	var counttotal, costtotal float64
	var list SampleList

	for i, v := range prof.Data {
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

	for _, v := range list {
		counttotal += float64(v.Data.Count)
		costtotal += float64(v.Data.CumulativeCost())
	}

	if len(list) > count {
		list = list[:count]
	}

	fmt.Println("           COUNT | COST |       CUM. COST |                 FILE | SOURCE")
	fmt.Println("============================================================================")
	for i := range list {
		filename = prof.Files[list[i].Data.File]
		source = getSourceLine(prof, list[i].PC)

		filename = filepath.Base(filename)
		filename = fmt.Sprintf("%s:%d", filename, list[i].Data.Line)

		if len(filename) > 20 {
			filename = "..." + filename[3:]
		}

		countpercent = float64(list[i].Data.Count) / counttotal
		costpercent = float64(list[i].Data.CumulativeCost()) / costtotal

		fmt.Printf(" %7d (%.2f%%) | %4d | %7d (%.2f%%) | %20s | %s\n",
			list[i].Data.Count,
			countpercent,
			list[i].Data.Cost(),
			list[i].Data.CumulativeCost(),
			costpercent,
			filename,
			source,
		)
	}

	fmt.Println()
}
