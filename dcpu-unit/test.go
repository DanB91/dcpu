// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	dp "github.com/jteeuwen/dcpu/parser"
	"sync"
)

// A single unit is a complete program for a single test unit.
// Represented as AST nodes.
type Unit []dp.Node

// Test represents a single unit test case.
// It covers one test file which may contain multiple unit tests.
type Test struct {
	includes []string // Include paths.
	file     string   // Test source file.
}

// Run loads up the test sources, compiles it and performs
// the unit tests defined in it.
func runTest(file string, inc []string, wg *sync.WaitGroup, status chan<- error) {
	defer wg.Done()

	var t Test
	t.file = file
	t.includes = inc

	ast, err := t.readAST()
	if err != nil {
		status <- err
		return
	}

	_ = ast
}

// readAST reads the test source and constructs a complete AST.
// This includes importing externally referenced files.
func (t *Test) readAST() (*dp.AST, error) {
	var ast dp.AST

	err := readSource(&ast, t.file)
	if err != nil {
		return nil, err
	}

	return &ast, resolveIncludes(&ast, t.includes)
}
