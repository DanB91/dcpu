// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import (
	"bytes"
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/parser"
	"testing"
)

type testCase struct {
	src string
	bin []cpu.Word
}

var _exit = encode(cpu.EXT, cpu.EXIT, 0x21)

var tests = []testCase{
	{
		`set a, 1
		 set b, 0xfffe
		 add a, b
		 exit`,
		[]cpu.Word{
			_exit,
		},
	},
}

func Test(t *testing.T) {
	for i := range tests {
		doTest(t, i, &tests[i])
	}
}

func doTest(t *testing.T, index int, tc *testCase) {
	var ast parser.AST
	var bin []cpu.Word

	buf := bytes.NewBufferString(tc.src)
	err := ast.Parse(buf, "")

	if err != nil {
		t.Fatal(err)
	}

	bin, err = Assemble(&ast)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%04x\n", bin)
	fmt.Printf("%04x\n", tc.bin)

	if len(bin) != len(tc.bin) {
		t.Fatalf("test %d: Size mismatch. Expect %d, got %d",
			index, len(tc.bin), len(bin))
	}

	for i := range bin {
		if bin[i] != tc.bin[i] {
			t.Fatalf("test %d: code mismatch at %d. Expect %04x, got %04x",
				index, i, tc.bin[i], bin[i])
		}
	}
}
