// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/prof"
	"strings"
)

const (
	DefaultTopCount = 10
	DefaultTopSort  = "cost"
)

// Display sorted list of profile data for every function call.
func top(p *prof.Profile, filemode bool, count uint, sort string) {
	var counttotal, costtotal float64
	var blocks prof.BlockList

	if filemode {
		blocks = p.ListFiles()
	} else {
		blocks = p.ListFunctions()
	}

	if len(blocks) == 0 {
		fmt.Println("[*] 0 samples.")
		if !filemode {
			fmt.Println("[*] This most likely means that there are no function")
			fmt.Println("    definitions in the source code. Try using -file mode.")
		}
		return
	}

	for i := range blocks {
		if len(blocks[i].Label) == 0 {
			blocks[i].Label = GetLabel(p, blocks[i].Addr)
		}

		a, b := blocks[i].Cost()
		counttotal += float64(a)
		costtotal += float64(b)
	}

	switch strings.ToLower(sort) {
	case "count":
		prof.BlockListByCount(blocks).Sort()
	case "cost":
		prof.BlockListByCost(blocks).Sort()
	}

	if uint(len(blocks)) > count {
		blocks = blocks[:count]
	}

	fmt.Printf("[*] %.0f sample(s), %.0f cycle(s)\n", counttotal, costtotal)

	for i := range blocks {
		count, cost := blocks[i].Cost()
		scount := fmt.Sprintf("%.2f%%", float64(count)/(counttotal*0.01))
		scost := fmt.Sprintf("%.2f%%", float64(cost)/(costtotal*0.01))

		fmt.Printf(" %8d %7s %8d %7s %s\n",
			count, scount, cost, scost, blocks[i].Label)
	}

	fmt.Println()
}
