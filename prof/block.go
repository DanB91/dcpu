// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import (
	"github.com/jteeuwen/dcpu/cpu"
	"sort"
)

type Block struct {
	Data  []ProfileData // Profile data for this function's instructions.
	Label string        // Label/name of this function.
	Addr  cpu.Word      // start address of function.
}

// Cost returns the cumulative cycle cost and count for all
// instructions in this function.
func (f *Block) Cost() (count, cost uint64) {
	for pc := 0; pc < len(f.Data); pc++ {
		count += f.Data[pc].Count
		cost += f.Data[pc].CumulativeCost()
	}

	return
}

// Range returns the address range for this function.
func (f *Block) Range() (start, end cpu.Word) {
	return f.Addr, f.Addr + cpu.Word(len(f.Data)-1)
}

type BlockList []Block

// List of functions, sortable by Count.
type BlockListByCount BlockList

func (s BlockListByCount) Len() int { return len(s) }
func (s BlockListByCount) Less(i, j int) bool {
	ca, _ := s[i].Cost()
	cb, _ := s[j].Cost()
	return ca >= cb
}
func (s BlockListByCount) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BlockListByCount) Sort()         { sort.Sort(s) }

// List of functions, sortable by cumulative cost.
type BlockListByCost BlockList

func (s BlockListByCost) Len() int { return len(s) }
func (s BlockListByCost) Less(i, j int) bool {
	_, ca := s[i].Cost()
	_, cb := s[j].Cost()
	return ca >= cb
}
func (s BlockListByCost) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BlockListByCost) Sort()         { sort.Sort(s) }
