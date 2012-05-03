// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"sync"
)

type Test struct {
}

// runTest loads up the test sources, compiles it and performs
// the unit tests defined in it.
func runTest(file string, wg *sync.WaitGroup, log *Log) {
	defer wg.Done()
	log.Write(1, "%s...", file)
}
