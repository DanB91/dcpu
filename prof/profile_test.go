// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import (
	"bytes"
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/cpu"
	"testing"
)

// Ensure Read/Write roundtrip is the identity function.
func TestIdentity(t *testing.T) {
	var w bytes.Buffer

	code := []cpu.Word{
		cpu.Encode(cpu.SET, 0, 0x1f),
		0xffff,
		cpu.Encode(cpu.ADD, 0, 0x22),
		cpu.Encode(cpu.MUL, 0, 0x22),
		cpu.Encode(cpu.SET, 1, 0),
	}

	var dbg asm.DebugInfo
	dbg.Files = []asm.FileInfo{{"a.dasm", 0}}
	dbg.Functions = []asm.FuncInfo{{"main", 0, 5}}
	dbg.SourceMapping = make([]asm.SourceInfo, 5)

	a := New(code, &dbg)
	a.Update(0, nil)
	a.Update(2, nil)
	a.Update(3, nil)
	a.Update(4, nil)

	a.UpdateCost(2, 5)

	err := Write(a, &w)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	b, err := Read(&w)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if len(a.Files) != len(b.Files) {
		t.Fatalf("len(a.Files) != len(b.Files)")
	}

	for i := range a.Files {
		if a.Files[i].Name != b.Files[i].Name {
			t.Fatalf("a.Files[i].Name != b.Files[i].Name")
		}

		if a.Files[i].Start != b.Files[i].Start {
			t.Fatalf("a.Files[i].Start != b.Files[i].Start")
		}
	}

	if len(a.Functions) != len(b.Functions) {
		t.Fatalf("len(a.Functions) != len(b.Functions)")
	}

	for i := range a.Functions {
		if a.Functions[i].Start != b.Functions[i].Start {
			t.Fatalf("a.Functions[i].Start != b.Functions[i].Start")
		}

		if a.Functions[i].End != b.Functions[i].End {
			t.Fatalf("a.Functions[i].End != b.Functions[i].End")
		}

		if a.Functions[i].Name != b.Functions[i].Name {
			t.Fatalf("a.Functions[i].Name != b.Functions[i].Name")
		}
	}

	if len(a.Data) != len(b.Data) {
		t.Fatalf("len(a.Data) != len(b.Data)")
	}

	for pc, va := range a.Data {
		vb := b.Data[pc]

		if va.Count != vb.Count {
			t.Fatalf("va.Count != vb.Count")
		}

		if va.File != vb.File {
			t.Fatalf("va.File != vb.File")
		}

		if va.Line != vb.Line {
			t.Fatalf("va.Line != vb.Line")
		}

		if va.Col != vb.Col {
			t.Fatalf("va.Col != vb.Col")
		}

		if va.Data != vb.Data {
			t.Fatalf("va.Data != vb.Data")
		}

		if va.Penalty != vb.Penalty {
			t.Fatalf("va.Penalty != vb.Penalty")
		}

		if va.Cost() != vb.Cost() {
			t.Fatalf("va.Cost() != vb.Cost()")
		}
	}
}
