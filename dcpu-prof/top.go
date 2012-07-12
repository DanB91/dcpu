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
	DefaultTopSort  = "cost"
)

// Display sorted list of profile data for every function call.
func top(p *prof.Profile, count int, sort string) {
	var counttotal, costtotal float64

	funcs := p.Functions()

	if len(funcs) == 0 {
		fmt.Println("0 samples.")
		return
	}

	for i := range funcs {
		if len(funcs[i].Label) == 0 {
			funcs[i].Label = GetLabel(p, funcs[i].Addr)
		}

		a, b := funcs[i].Cost()
		counttotal += float64(a)
		costtotal += float64(b)
	}

	switch strings.ToLower(sort) {
	case "count":
		prof.FuncListByCount(funcs).Sort()
	case "cost":
		prof.FuncListByCost(funcs).Sort()
	}

	if len(funcs) > count {
		funcs = funcs[:count]
	}

	fmt.Printf("%.0f sample(s), %.0f cycle(s)\n", counttotal, costtotal)

	for i := range funcs {
		count, cost := funcs[i].Cost()
		scount := fmt.Sprintf("%.2f%%", float64(count)/(counttotal*0.01))
		scost := fmt.Sprintf("%.2f%%", float64(cost)/(costtotal*0.01))

		fmt.Printf(" %8d %7s %8d %7s %s\n",
			count, scount, cost, scost, funcs[i].Label)
	}

	fmt.Println()
}
