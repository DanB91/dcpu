// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/cpu"
	dp "github.com/jteeuwen/dcpu/parser"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Test represents a single unit test case.
// It covers one test file which may contain multiple unit tests.
type Test struct {
	dbg       *asm.DebugInfo    // Debug symbols for compiled program.
	includes  []string          // Include paths.
	callstack []*asm.SourceInfo // callstack for the test program.
	file      string            // Test source file.
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
func (t *Test) Run(cfg *Config) (err error) {
	if cfg.Verbose {
		fmt.Fprintf(os.Stdout, "[*] %s...\n", t.file)
	}

	var ast *dp.AST
	if ast, err = t.parse(); err != nil {
		return
	}

	var c *cpu.CPU
	if c, err = t.compile(ast); err != nil {
		return
	}

	c.Trace = func(pc, op, a, b cpu.Word, s *cpu.Storage) {
		t.trace(pc, op, a, b, s, cfg.Trace)
	}

	c.ClockSpeed = time.Duration(cfg.Clock)
	if err = c.Run(0); err != nil {
		if te, ok := err.(*cpu.TestError); ok {
			return t.formatTestError(te)
		}
	}

	return
}

// formatTestError constructs a full error message like this:
//
//     [E] string/memchr_test.dasm Assertion failed: A != B
//      Call stack:
//      - memchr_test.dasm:7 | jsr asserteq
//      - memchr_test.dasm:4 | jsr memchr
//
func (t *Test) formatTestError(e *cpu.TestError) error {
	var b bytes.Buffer

	if int(e.PC) >= len(t.dbg.Data) {
		return errors.New(fmt.Sprintf("No debug symbols available for address %04x.", e.PC))
	}

	s := t.dbg.Data[e.PC]
	file := t.dbg.Files[s.File]
	_, file = filepath.Split(file)

	fmt.Fprintf(&b, "[E] %s: %s\n", t.file, e.Msg)
	fmt.Fprintln(&b, "    Call stack:")

	for i := len(t.callstack) - 1; i >= 0; i-- {
		s = t.callstack[i]
		file = t.dbg.Files[s.File]

		fmt.Fprintf(&b, "    - %s\n", t.getSourceLine(file, s.Line))
	}

	return errors.New(b.String())
}

// trace builds a callstack for the executing program.
// This is used for adequate source context when an error occurs.
//
// It also optionally prints trace output for each instruction as it is executed.
// It yields current PC, opcode, operands, all register contents and
// appends the original line of sourcecode
func (t *Test) trace(pc, op, a, b cpu.Word, s *cpu.Storage, verbose bool) {
	if op == cpu.EXT && a == cpu.JSR {
		t.callstack = append(t.callstack, t.dbg.Data[pc])
	}

	if op == cpu.SET && a == 0x1c /*PC*/ && b == 0x18 /*POP*/ {
		sz := len(t.callstack)

		if sz == 0 {
			return
		}

		t.callstack = t.callstack[:sz-1]
	}

	if verbose {
		if int(pc) >= len(t.dbg.Data) {
			fmt.Fprintf(os.Stdout,
				"%04x: %04x %04x %04x | %04x %04x %04x %04x %04x %04x %04x %04x | %04x %04x %04x | <unknown>\n",
				pc, op, a, b, s.A, s.B, s.C, s.X, s.Y, s.Z, s.I, s.J, s.SP, s.EX, s.IA)
			return
		}

		symbol := t.dbg.Data[pc]
		file := t.dbg.Files[symbol.File]
		line := t.getSourceLine(file, symbol.Line)

		fmt.Fprintf(os.Stdout,
			"%04x: %04x %04x %04x | %04x %04x %04x %04x %04x %04x %04x %04x | %04x %04x %04x | %s\n",
			pc, op, a, b, s.A, s.B, s.C, s.X, s.Y, s.Z,
			s.I, s.J, s.SP, s.EX, s.IA, line)
	}
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

	if !hasExit(bin) {
		return nil, errors.New(fmt.Sprintf(
			"%s: Program has no unconditional EXIT. This means the test will run indefinitely.", t.file))
	}

	c = cpu.New()
	copy(c.Store.Mem[:], bin)
	return
}

// hasExit determines if the given program has at least one
// unconditional EXIT instruction.
func hasExit(bin []cpu.Word) bool {
	var op1, op2, a, b cpu.Word
	var size int

	for i := 0; i < len(bin); i++ {
		op1, a, b = cpu.Decode(bin[i])

		// Not an EXIT. Skip with next word.
		if !(op1 == cpu.EXT && a == cpu.EXIT) {
			continue
		}

		// If the previous instruction is not a branching instruction,
		// we have found a non-conditional EXIT.
		op2, _, _ = cpu.Decode(bin[i-size])

		if op2 < cpu.IFB || op2 > cpu.IFU {
			return true
		}

		// Determine size of instruction.
		// Some of them occupy multiple words.
		size = 1

		if op1 != cpu.EXT && (a == 0x1e || a == 0x1f || (a >= 0x10 && a <= 0x17)) {
			size++
		}

		if b == 0x1e || b == 0x1f || (b >= 0x10 && b <= 0x17) {
			size++
		}
	}

	return false
}

// getSourceLine fetches the line of sourcecode from the
// file defined by the given PC value. This data is stored in the
// debug symbol table.
func (t *Test) getSourceLine(file string, lineno int) string {
	fd, err := os.Open(file)
	if err != nil {
		return ""
	}

	defer fd.Close()

	r := bufio.NewReader(fd)

	var count int
	var line []byte

	_, file = filepath.Split(file)

	for {
		if line, _, err = r.ReadLine(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return ""
		}

		if count < lineno-1 {
			count++
			continue
		}

		line = bytes.TrimSpace(line)
		return fmt.Sprintf("%s:%d | %s", file, lineno, line)
	}

	return ""
}
