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

	dbg := new(asm.DebugInfo)
	dbg.Files = []string{"a.dasm"}
	dbg.SourceMapping = []*asm.SourceInfo{
		new(asm.SourceInfo),
		new(asm.SourceInfo),
		new(asm.SourceInfo),
		new(asm.SourceInfo),
		new(asm.SourceInfo),
	}

	a := New(code, dbg)
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
		if a.Files[i] != b.Files[i] {
			t.Fatalf("a.Files != b.Files")
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
