// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package hmu1440

import (
	"github.com/jteeuwen/dcpu/cpu"
	"testing"
)

func TestWrite(t *testing.T) {
	var fdd HMU1440
	err := fdd.Open("test.fdd")

	if err != nil {
		t.Fatal(err)
	}

	defer fdd.Close()

	var buf [SectorSize]cpu.Word
	for i := range buf {
		buf[i] = cpu.Word(i)
	}

	if err = fdd.Write(2, buf[:]); err != nil {
		t.Errorf("write: %v", err)
		return
	}
}

func TestRead(t *testing.T) {
	var fdd HMU1440
	err := fdd.Open("test.fdd")

	if err != nil {
		t.Fatal(err)
	}

	defer fdd.Close()

	var buf [SectorSize]cpu.Word

	if err = fdd.Read(2, buf[:]); err != nil {
		t.Error(err)
		return
	}

	if buf[0xa] != 0xa {
		t.Errorf("Invalid value. Expected 0x0a, got 0x%02x", buf[0xa])
		return
	}
}
