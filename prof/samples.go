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

// List of samples, sortable by Opcode.
type SamplesByOpcode SampleList

func (s SamplesByOpcode) Len() int           { return len(s) }
func (s SamplesByOpcode) Less(i, j int) bool { return s[i].Data.Opcode >= s[j].Data.Opcode }
func (s SamplesByOpcode) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SamplesByOpcode) Sort()              { sort.Sort(s) }

// List of samples, sortable by cost.
type SamplesByCost SampleList

func (s SamplesByCost) Len() int           { return len(s) }
func (s SamplesByCost) Less(i, j int) bool { return s[i].Data.Cost() >= s[j].Data.Cost() }
func (s SamplesByCost) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SamplesByCost) Sort()              { sort.Sort(s) }

// List of samples, sortable by cumulative cost.
type SamplesByCumulativeCost SampleList

func (s SamplesByCumulativeCost) Len() int { return len(s) }
func (s SamplesByCumulativeCost) Less(i, j int) bool {
	return s[i].Data.CumulativeCost() >= s[j].Data.CumulativeCost()
}
func (s SamplesByCumulativeCost) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SamplesByCumulativeCost) Sort()         { sort.Sort(s) }
