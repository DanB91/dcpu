// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	dp "github.com/jteeuwen/dcpu/parser"
	"sync"
)

// runTest loads up the test sources, compiles it and performs
// the unit tests defined in it.
func runTest(file string, includes []string, wg *sync.WaitGroup, log *Log) {
	var err error
	var ast dp.AST

	defer wg.Done()

	log.Write("[*] %s...", file)

	if err = readSource(&ast, file); err != nil {
		log.Write("[e] %s", err)
		return
	}

	if err = resolveIncludes(&ast, includes); err != nil {
		log.Write("[e] %s", err)
		return
	}
}
