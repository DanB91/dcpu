// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import (
	"bytes"
	"github.com/jteeuwen/dcpu/cpu"
	"testing"
)

// Ensure Read/Write roundtrip is the identity function.
func TestIdentity(t *testing.T) {
	var w bytes.Buffer

	a := New([]string{"a.dasm"}, 5)
	a.Update(0, cpu.SET, 0, 0x1f, 0, 1, 1)
	a.Update(2, cpu.ADD, 0, 0x22, 0, 2, 1)
	a.Update(3, cpu.MUL, 0, 0x22, 0, 3, 1)
	a.Update(4, cpu.SET, 1, 0, 0, 4, 1)

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

	if len(a.Usage) != len(b.Usage) {
		t.Fatalf("len(a.Usage) != len(b.Usage)")
	}

	for pc, va := range a.Usage {
		vb := b.Usage[pc]

		if va == nil || vb == nil {
			if va == nil && vb != nil {
				t.Fatalf("va is nil. vb is not.")
			}

			if va != nil && vb == nil {
				t.Fatalf("va is not nil. vb is.")
			}

			continue
		}

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

		if va.Opcode != vb.Opcode {
			t.Fatalf("va.Opcode != vb.Opcode")
		}

		if va.A != vb.A {
			t.Fatalf("va.A != vb.A")
		}

		if va.Opcode != vb.Opcode {
			t.Fatalf("va.B != vb.B")
		}

		if va.Cost() != vb.Cost() {
			t.Fatalf("va.Cost() != vb.Cost()")
		}
	}
}
