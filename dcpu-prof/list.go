// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/prof"
	"regexp"
)

const DefaultListFilter = ".+"

func getLineData(l []prof.ProfileData, start, end cpu.Word, line int) prof.Block {
	var b prof.Block

	for pc := start; pc < end; pc += l[pc].Size {
		if l[pc].Line == line {
			b.Data = append(b.Data, l[pc])
		}
	}

	return b
}

// Display detailed instruction view for the given filter.
func list(p *prof.Profile, filemode bool, filter *regexp.Regexp) {
	var blocks prof.BlockList
	var linedata prof.Block
	var count, cost uint64

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
		if !filter.MatchString(blocks[i].Label) {
			continue
		}

		start, end := blocks[i].StartAddr, blocks[i].EndAddr
		totalcount, totalcost := blocks[i].Cost()

		file := p.Files[p.Data[start].File]
		startline := blocks[i].StartLine
		endline := blocks[i].EndLine
		source := GetSourceLines(file.Name, startline, endline)

		fmt.Printf("[*] ===> %s\n", blocks[i].Label)
		fmt.Printf("[*] %d sample(s), %d cycle(s)\n\n", totalcount, totalcost)

		if startline == 0 {
			startline++
		}

		for j := range source {
			linedata = getLineData(p.Data, start, end, startline+j)
			count, cost = linedata.Cost()

			if count == 0 {
				fmt.Printf("                    %03d: %s\n", startline+j, source[j])
			} else {
				fmt.Printf("%8d %8d   %03d: %s\n", count, cost, startline+j, source[j])
			}
		}

		fmt.Println()
	}
}
