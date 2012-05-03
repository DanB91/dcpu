// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	dp "github.com/jteeuwen/dcpu/parser"
	"sync"
)

// Test represents a single unit test case.
// It covers one test file which may contain multiple unit tests.
type Test struct {
	*Log              // Output log for errors and status messages.
	includes []string // Include paths.
	file     string   // Test source file.
}

// NewTest creates a new test for the given inputs.
func NewTest(file string, inc []string, log *Log) *Test {
	t := new(Test)
	t.Log = log
	t.includes = inc
	t.file = file
	return t
}

// Run loads up the test sources, compiles it and performs
// the unit tests defined in it.
func (t *Test) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	t.Printf("[*] %s...", t.file)

	ast, err := t.readAST()
	if err != nil {
		t.Errorf("[e] %s", err)
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
