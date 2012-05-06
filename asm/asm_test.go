// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import (
	"bytes"
	//"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"github.com/jteeuwen/dcpu/parser"
	"testing"
)

type testCase struct {
	src string
	bin []cpu.Word
}

var _exit = encode(cpu.EXT, cpu.EXIT, 0)

var tests = []testCase{
	{
		`set a, 1
		 set b, 0xfffe
		 add b, [0x100]
		 add a, [a]
		 ;moo
		 exit`,
		[]cpu.Word{
			encode(cpu.SET, 0, 0x22),
			encode(cpu.SET, 1, 0x1f),
			0xfffe,
			encode(cpu.ADD, 1, 0x1e),
			0x100,
			encode(cpu.ADD, 0, 8),
			_exit,
		},
	},
	{
		`  set a, [sp+end]
		   xor [0xfffe], a
		 :end
		   exit`,
		[]cpu.Word{
			encode(cpu.SET, 0, 0x1a),
			0x4,
			encode(cpu.XOR, 0x1e, 0),
			0xfffe,
			_exit,
		},
	},
	{
		`jsr 3
		 jsr a
		 jsr [z]
		 exit`,
		[]cpu.Word{
			encode(cpu.EXT, cpu.JSR, 0x24),
			encode(cpu.EXT, cpu.JSR, 0),
			encode(cpu.EXT, cpu.JSR, 0xd),
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

	//fmt.Printf("%04x\n", bin)
	//fmt.Printf("%04x\n", tc.bin)

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
