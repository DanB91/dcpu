// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/cpu"
	dp "github.com/jteeuwen/dcpu/parser"
	"io"
	"os"
	"strings"
	"time"
)

const unitString = `Line %d, col %d:
           A    B    C    X    Y    Z    I    J   EX   SP   IA
  Want: %s
  Have: %s`

// Test represents a single unit test case.
// It covers one test file which may contain multiple unit tests.
type Test struct {
	dbg      *asm.DebugInfo // Debug symbols for compiled program.
	includes []string       // Include paths.
	compare  []string       // Lines of compare data.
	file     string         // Test source file.
	count    int            // Unit counter.
}

// NewTest creates a new test cases.
func NewTest(file string, inc []string) *Test {
	t := new(Test)
	t.file = file
	t.includes = inc
	return t
}

// Run loads up the test sources, compiles it and performs
// the unit tests defined in it.
//
// Any errors are returned through the status channel.
func (t *Test) Run(trace cpu.TraceFunc, speed int64, verbose bool) (err error) {
	if verbose {
		fmt.Fprintf(os.Stdout, "> %s...\n", t.file)
	}

	if err = t.loadCompareSet(); err != nil {
		return
	}

	var ast *dp.AST
	if ast, err = t.parse(); err != nil {
		return
	}

	var c *cpu.CPU
	if c, err = t.compile(ast); err != nil {
		return
	}

	c.Trace = trace
	c.ClockSpeed = time.Duration(speed)
	c.Test = func(pc cpu.Word, s *cpu.Storage) error {
		return t.handleTest(pc, s, verbose)
	}

	return c.Run(0)
}

// handleTest is called whenever a TEST instruction fires in a test program.
func (t *Test) handleTest(pc cpu.Word, s *cpu.Storage, verbose bool) (err error) {
	if verbose {
		fmt.Fprintf(os.Stdout, "  - Unit %d...", t.count)
	}

	if t.count >= len(t.compare) {
		return errors.New(fmt.Sprintf("Missing compare data for unit %d.", t.count))
	}

	line := fmt.Sprintf("%04x %04x %04x %04x %04x %04x %04x %04x %04x %04x %04x",
		s.A, s.B, s.C, s.X, s.Y, s.Z, s.I, s.J, s.EX, s.SP, s.IA)

	if line != t.compare[t.count] {
		if verbose {
			fmt.Fprintln(os.Stdout)
		}

		symbol := t.dbg.Data[pc]

		return errors.New(fmt.Sprintf(unitString,
			symbol.Line, symbol.Col, t.compare[t.count], line))
	}

	t.count++

	if verbose {
		fmt.Fprintln(os.Stdout, "  OK")
	}

	return
}

// parse reads the test source and constructs a complete AST.
// This includes importing externally referenced files.
func (t *Test) parse() (*dp.AST, error) {
	var ast dp.AST

	err := readSource(&ast, t.file)
	if err != nil {
		return nil, err
	}

	return &ast, resolveIncludes(&ast, t.includes)
}

// compile compiles the given AST and returns a CPU instance ready to run the code.
func (t *Test) compile(ast *dp.AST) (c *cpu.CPU, err error) {
	var bin []cpu.Word

	if bin, t.dbg, err = asm.Assemble(ast); err != nil {
		return
	}

	c = cpu.New()
	copy(c.Store.Mem[:], bin)
	return
}

// loadCompareSet loads the compare file for this unit test.
func (t *Test) loadCompareSet() (err error) {
	file := strings.Replace(t.file, ".test", ".cmp", 1)

	var fd *os.File
	if fd, err = os.Open(file); err != nil {
		return
	}

	defer fd.Close()

	r := bufio.NewReader(fd)

	var line []byte

	for {
		if line, _, err = r.ReadLine(); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		t.compare = append(t.compare, string(line))
	}

	// First line has column headers. Skip it.
	t.compare = t.compare[1:]
	return
}
