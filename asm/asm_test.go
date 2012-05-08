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

func doTest(t *testing.T, src string, sbin ...cpu.Word) {
	var ast parser.AST
	var dbin []cpu.Word

	buf := bytes.NewBufferString(src)
	err := ast.Parse(buf, "")

	if err != nil {
		t.Fatal(err)
	}

	dbin, _, err = Assemble(&ast)
	if err != nil {
		t.Fatal(err)
	}

	if len(dbin) != len(sbin) {
		fmt.Printf("%04x\n", dbin)
		fmt.Printf("%04x\n", sbin)
		t.Fatalf("Size mismatch. Expect %d, got %d", len(sbin), len(dbin))
	}

	for i := range sbin {
		if dbin[i] != sbin[i] {
			fmt.Printf("%04x\n", dbin)
			fmt.Printf("%04x\n", sbin)
			t.Fatalf("Code mismatch at %d. Expect %04x, got %04x", i, sbin[i], dbin[i])
		}
	}
}

func TestSet(t *testing.T) {
	doTest(t,
		`set a, 0x30`,
		cpu.Encode(cpu.SET, 0, 0x1f),
		0x30,
	)
}

func TestAdd(t *testing.T) {
	doTest(t,
		`set a, 0xffff
		 add a, 1`,
		cpu.Encode(cpu.SET, 0, 0x20),
		cpu.Encode(cpu.ADD, 0, 0x22),
	)
}

func TestSub(t *testing.T) {
	doTest(t,
		`set a, 0
		 sub a, 1`,
		cpu.Encode(cpu.SET, 0x0, 0x21),
		cpu.Encode(cpu.SUB, 0x0, 0x22),
	)
}

func TestAsr(t *testing.T) {
	doTest(t,
		`SET A, 10
		 ASR A, 1`,
		cpu.Encode(cpu.SET, 0, 0x2b),
		cpu.Encode(cpu.ASR, 0, 0x22),
	)
}

func TestNestedIf(t *testing.T) {
	doTest(t,
		`SET A, 0
		 IFG A, 1
		   IFG A, 2
		     IFE A, 0
               SET A, 4
         SET [100], A
		`,
		cpu.Encode(cpu.SET, 0, 0x21),
		cpu.Encode(cpu.IFG, 0, 0x22),
		cpu.Encode(cpu.IFG, 0, 0x23),
		cpu.Encode(cpu.IFE, 0, 0x21),
		cpu.Encode(cpu.SET, 0, 0x25),
		cpu.Encode(cpu.SET, 0x1e, 0),
		100,
	)
}

func TestJsr(t *testing.T) {
	doTest(t,
		`  set a, 0xffff
		   jsr my_sub
		   exit
		 :my_sub
		   add a, 1
		   set pc, pop
		`,
		cpu.Encode(cpu.SET, 0, 0x20),
		cpu.Encode(cpu.EXT, cpu.JSR, 0x1f),
		0x4,
		cpu.Encode(cpu.EXT, cpu.EXIT, 0),
		cpu.Encode(cpu.ADD, 0, 0x22),
		cpu.Encode(cpu.SET, 0x1c, 0x18),
	)
}

func TestIntRfi(t *testing.T) {
	doTest(t,
		` ias my_handler
		  int 0xbeef
		  set a, b
		  exit
		:my_handler
		  set b, a
		  add b, 1
		  rfi a
		`,
		cpu.Encode(cpu.EXT, cpu.IAS, 0x1f),
		0x6,
		cpu.Encode(cpu.EXT, cpu.INT, 0x1f),
		0xbeef,
		cpu.Encode(cpu.SET, 0, 1),
		cpu.Encode(cpu.EXT, cpu.EXIT, 0),
		cpu.Encode(cpu.SET, 1, 0),
		cpu.Encode(cpu.ADD, 1, 0x22),
		cpu.Encode(cpu.EXT, cpu.RFI, 0),
	)
}

func TestHwi(t *testing.T) {
	doTest(t,
		`  set a,1
		   set b, 0x100
		   hwi 0
		   set a, [0x100]
		   exit
		`,
		cpu.Encode(cpu.SET, 0, 0x22),
		cpu.Encode(cpu.SET, 1, 0x1f),
		0x100,
		cpu.Encode(cpu.EXT, cpu.HWI, 0x21),
		cpu.Encode(cpu.SET, 0, 0x1e),
		0x100,
		cpu.Encode(cpu.EXT, cpu.EXIT, 0),
	)
}

func TestDat(t *testing.T) {
	doTest(t,
		`:end
		    set pc, end
		 :dat
		    dat 0x170, "Hello, universe", 0
		`,
		cpu.Encode(cpu.SET, 0x1c, 0x21),
		0x170, 'H', 'e', 'l', 'l', 'o', ',', ' ',
		'u', 'n', 'i', 'v', 'e', 'r', 's', 'e', 0,
	)
}
