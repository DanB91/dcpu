// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import (
	"github.com/jteeuwen/dcpu/cpu"
	"sort"
)

type Sample struct {
	Data *ProfileData
	PC   cpu.Word
}

// List of samples.
type SampleList []Sample

func (s SampleList) IndexOfPC(pc cpu.Word) int {
	for i := range s {
		if s[i].PC == pc {
			return i
		}
	}
	return -1
}

// List of samples, sortable by PC.
type SamplesByPC SampleList

func (s SamplesByPC) Len() int           { return len(s) }
func (s SamplesByPC) Less(i, j int) bool { return s[i].PC >= s[j].PC }
func (s SamplesByPC) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SamplesByPC) Sort()              { sort.Sort(s) }

// List of samples, sortable by Count.
type SamplesByCount SampleList

func (s SamplesByCount) Len() int           { return len(s) }
func (s SamplesByCount) Less(i, j int) bool { return s[i].Data.Count >= s[j].Data.Count }
func (s SamplesByCount) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SamplesByCount) Sort()              { sort.Sort(s) }

// List of samples, sortable by cumulative cost.
type SamplesByCumulativeCost SampleList

func (s SamplesByCumulativeCost) Len() int { return len(s) }
func (s SamplesByCumulativeCost) Less(i, j int) bool {
	return s[i].Data.CumulativeCost() >= s[j].Data.CumulativeCost()
}
func (s SamplesByCumulativeCost) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SamplesByCumulativeCost) Sort()         { sort.Sort(s) }
