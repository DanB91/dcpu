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
	"github.com/jteeuwen/dcpu/parser/util"
	"github.com/jteeuwen/dcpu/prof"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Test represents a single unit test case.
// It covers one test file which may contain multiple unit tests.
type Test struct {
	dbg       *asm.DebugInfo      // Debug symbols for compiled program.
	cache     map[cpu.Word]string // Cache of source lines for trace output.
	profile   *prof.Profile       // Profiling information.
	includes  []string            // Include paths.
	callstack []string            // callstack for the test program.
	file      string              // Test source file.
}

// NewTest creates a new test cases.
func NewTest(file string, inc []string) *Test {
	t := new(Test)
	t.file = file
	t.includes = inc
	t.cache = make(map[cpu.Word]string)
	return t
}

// Run loads up the test sources, compiles it and performs
// the unit tests defined in it. Additionally, it runs a profiler
// on the program which can optionally be written to an output file
// for examination.
func (t *Test) Run() (err error) {
	fmt.Fprintf(os.Stdout, "[*] %s...\n", t.file)

	ast, err := t.parse()
	if err != nil {
		return
	}

	c, err := t.compile(ast)
	if err != nil {
		return
	}

	err = c.Run(0)

	if err != nil {
		if te, ok := err.(*cpu.TestError); ok {
			err = t.formatTestError(te)
			return
		}
	}

	if *profile {
		file := strings.Replace(t.file, ".dasm", ".prof", 1)

		var fd *os.File
		fd, err = os.Create(file)

		if err != nil {
			return
		}

		err = prof.Write(t.profile, fd)
		fd.Close()
	}

	return
}

// formatTestError constructs a full error message like this:
//
//     [E] string/memchr_test.dasm Assertion failed: A != B
//      Call stack:
//      - memchr_test.dasm:7 | jsr asserteq
//
func (t *Test) formatTestError(e *cpu.TestError) error {
	if int(e.PC) >= len(t.dbg.SourceMapping) {
		return errors.New(fmt.Sprintf("No debug symbols available for address %04x.", e.PC))
	}

	var b bytes.Buffer
	fmt.Fprintf(&b, "[E] %s: %s\n", t.file, e.Msg)
	fmt.Fprintln(&b, "    Call stack:")

	for i := len(t.callstack) - 1; i >= 0; i-- {
		fmt.Fprintf(&b, "    - %s\n", t.callstack[i])
	}

	return errors.New(b.String())
}

// parseInstruction builds a callstack for the executing program.
// This is used for adequate source context when an error occurs.
//
// It also optionally prints trace output for each instruction as it is executed.
// It yields current PC, opcode, operands, all register contents and
// appends the original line of sourcecode
func (t *Test) parseInstruction(pc, op, a, b cpu.Word, s *cpu.Storage, verbose bool) {
	// Update callstack
	if op == cpu.EXT && a == cpu.JSR {
		line := t.getSourceLine(pc)
		t.callstack = append(t.callstack, line)

	} else if op == cpu.SET && a == 0x1c /*PC*/ && b == 0x18 /*POP*/ {
		sz := len(t.callstack)

		if sz == 0 {
			return
		}

		t.callstack = t.callstack[:sz-1]
	}

	if !verbose {
		return
	}

	// Print trace output
	if int(pc) >= len(t.dbg.SourceMapping) {
		fmt.Fprintf(os.Stdout,
			"%04x: %04x %04x %04x | %04x %04x %04x %04x %04x %04x %04x %04x | %04x %04x %04x | <unknown>\n",
			pc, op, a, b, s.A, s.B, s.C, s.X, s.Y, s.Z, s.I, s.J, s.SP, s.EX, s.IA)
		return
	}

	line := t.getSourceLine(pc)

	fmt.Fprintf(os.Stdout,
		"%04x: %04x %04x %04x | %04x %04x %04x %04x %04x %04x %04x %04x | %04x %04x %04x | %s\n",
		pc, op, a, b, s.A, s.B, s.C, s.X, s.Y, s.Z,
		s.I, s.J, s.SP, s.EX, s.IA, line)
}

// parse reads the test source and constructs a complete AST.
// This includes importing externally referenced files.
func (t *Test) parse() (*dp.AST, error) {
	var ast dp.AST

	err := util.ReadSource(&ast, t.file, t.includes)
	if err != nil {
		return nil, err
	}

	return &ast, nil
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

	t.profile = prof.New(bin, t.dbg)

	c = cpu.New()
	copy(c.Store.Mem[:], bin)

	c.ClockSpeed = time.Duration(time.Duration(*clock))
	c.Trace = func(pc, op, a, b cpu.Word, s *cpu.Storage) {
		t.parseInstruction(pc, op, a, b, s, *trace)
	}

	c.InstructionHandler = func(pc cpu.Word, s *cpu.Storage) {
		t.profile.Update(pc, s)
	}

	c.NotifyBranchSkip = func(pc, cost cpu.Word) {
		t.profile.UpdateCost(pc, cost)
	}

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
func (t *Test) getSourceLine(pc cpu.Word) string {
	if line, ok := t.cache[pc]; ok {
		return line
	}

	symbol := t.dbg.SourceMapping[pc]
	file := t.dbg.Files[symbol.File]
	fname := file.Name

	fd, err := os.Open(fname)
	if err != nil {
		t.cache[pc] = ""
		return ""
	}

	defer fd.Close()

	r := bufio.NewReader(fd)

	var count int
	var line []byte

	_, fname = filepath.Split(fname)

	for {
		if line, _, err = r.ReadLine(); err != nil {
			t.cache[pc] = ""
			if err == io.EOF {
				err = nil
			}
			return ""
		}

		if count < symbol.Line-1 {
			count++
			continue
		}

		line = bytes.TrimSpace(line)
		t.cache[pc] = fmt.Sprintf("%s:%d | %s", fname, symbol.Line, line)
		return t.cache[pc]
	}

	return ""
}
