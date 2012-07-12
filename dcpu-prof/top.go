// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/prof"
	"path/filepath"
	"strings"
)

const (
	DefaultTopCount = 5
	DefaultTopSort  = "count"
)

// Display sorted list of profile data for every function call.
func top(p *prof.Profile, count int, sort string) {
	var filename, source string
	var scount, scost string
	var counttotal, costtotal float64

	list := p.DataForFunctions()

	if len(list) == 0 {
		fmt.Println("0 samples.")
		return
	}

	switch strings.ToLower(sort) {
	case "cumulative":
		prof.SamplesByCumulativeCost(list).Sort()

	case "count":
		prof.SamplesByCount(list).Sort()
	}

	for _, v := range list {
		counttotal += float64(v.Data.Count)
		costtotal += float64(v.Data.CumulativeCost())
	}

	if len(list) > count {
		list = list[:count]
	}

	for i := range list {
		filename = p.Files[list[i].Data.File]
		source = getSourceLine(p, list[i].PC)

		filename = filepath.Base(filename)
		filename = fmt.Sprintf("%s:%d", filename, list[i].Data.Line)

		if len(filename) > 20 {
			filename = "..." + filename[3:]
		}

		scount = fmt.Sprintf("%.1f%%",
			float64(list[i].Data.Count)/counttotal)

		scost = fmt.Sprintf("%.1f%%",
			float64(list[i].Data.CumulativeCost())/costtotal)

		fmt.Printf(" %8d %6s %8d %6s %s %s\n",
			list[i].Data.Count,
			scount,
			list[i].Data.CumulativeCost(),
			scost,
			filename,
			source,
		)
	}

	fmt.Println()
}
