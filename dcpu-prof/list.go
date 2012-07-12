// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/prof"
	"regexp"
)

var DefaultListFilter = regexp.MustCompile(".+")

func getLineData(l []prof.ProfileData, start, end cpu.Word, line int) *prof.ProfileData {
	for pc := start; pc <= end; pc += l[pc].Size {
		if l[pc].Line == line && l[pc].Count > 0 {
			return &l[pc]
		}
	}

	return nil
}

// Display detailed instruction view for the given filter.
func list(p *prof.Profile, filter *regexp.Regexp) {
	funcs := p.Functions()

	if len(funcs) == 0 {
		fmt.Println("0 samples.")
		return
	}

	for i := range funcs {
		if len(funcs[i].Label) == 0 {
			funcs[i].Label = GetLabel(p, funcs[i].Addr)
		}

		if !filter.MatchString(funcs[i].Label) {
			continue
		}

		start, end := funcs[i].Range()
		totalcount, totalcost := funcs[i].Cost()

		file := p.Files[p.Data[start].File]
		startline := p.Data[start].Line
		endline := p.Data[end].Line
		source := GetSourceLines(file, startline, endline)

		fmt.Printf("===> %s %d-%d\n", file, startline, endline)
		fmt.Printf("%d sample(s), %d cycle(s)\n\n", totalcount, totalcost)

		for j := range source {
			dp := getLineData(p.Data, start, end, startline+j)

			if dp == nil {
				fmt.Printf("                    %03d: %s\n", startline+j, source[j])
			} else {
				fmt.Printf("%8d %8d   %03d: %s\n",
					dp.Count, dp.CumulativeCost(), startline+j, source[j])
			}
		}

		fmt.Println()
	}
}
