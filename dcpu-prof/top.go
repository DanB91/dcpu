// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/prof"
	"strings"
)

const (
	DefaultTopCount = 5
	DefaultTopSort  = "count"
)

// Display sorted list of profile data for every function call.
func top(p *prof.Profile, count int, sort string) {
	var source string
	var scount, scost string
	var counttotal, costtotal float64

	list := p.FunctionCosts()

	if len(list) == 0 {
		fmt.Println("0 samples.")
		return
	}

	switch strings.ToLower(sort) {
	case "cost":
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

	fmt.Printf("%.0f sample(s), %.0f cycle(s)\n", counttotal, costtotal)

	for i := range list {
		source = getLabel(p, list[i].PC)

		scount = fmt.Sprintf("%.2f%%",
			float64(list[i].Data.Count)/(counttotal*0.01))

		scost = fmt.Sprintf("%.2f%%",
			float64(list[i].Data.CumulativeCost())/(costtotal*0.01))

		fmt.Printf(" %8d %7s %8d %7s %s\n",
			list[i].Data.Count,
			scount,
			list[i].Data.CumulativeCost(),
			scost,
			source,
		)
	}

	fmt.Println()
}
