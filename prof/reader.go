// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jteeuwen/dcpu/cpu"
	"io"
)

// Read reads the binary version of a profile into a Profile structure.
//
// See the documentation on prof.Write for details on the file format.
func Read(r io.Reader) (p *Profile, err error) {
	var size uint32
	p = new(Profile)

	// [1]
	if err = binary.Read(r, be, &size); err != nil {
		return
	}

	p.Files = make([]string, size)

	// [2]
	for i := range p.Files {
		var size uint16
		if err = binary.Read(r, be, &size); err != nil {
			return
		}

		d := make([]byte, size)
		if _, err = r.Read(d); err != nil {
			return
		}

		p.Files[i] = string(d)
	}

	// [3]
	if err = binary.Read(r, be, &size); err != nil {
		return
	}

	p.Data = make([]*ProfileData, size)

	// [4]
	if err = binary.Read(r, be, &size); err != nil {
		return
	}

	// [5]
	var d [33]byte
	for i := uint32(0); i < size; i++ {
		if _, err = r.Read(d[:]); err != nil {
			return
		}

		pc := (uint16(d[0]) << 8) | uint16(d[1])
		if int(pc) >= len(p.Data) {
			err = errors.New(fmt.Sprintf("Invalid program counter value: 0x%04x", pc))
			return
		}

		pd := new(ProfileData)
		p.Data[pc] = pd

		pd.Opcode = cpu.Word(d[2])
		pd.A = cpu.Word(d[3])
		pd.B = cpu.Word(d[4])
		pd.File = int(d[5])<<24 | int(d[6])<<16 | int(d[7])<<8 | int(d[8])
		pd.Line = int(d[9])<<24 | int(d[10])<<16 | int(d[11])<<8 | int(d[12])
		pd.Col = int(d[13])<<24 | int(d[14])<<16 | int(d[15])<<8 | int(d[16])
		pd.Count = uint64(d[17])<<56 | uint64(d[18])<<48 | uint64(d[19])<<40 |
			uint64(d[20])<<32 | uint64(d[21])<<24 | uint64(d[22])<<16 |
			uint64(d[23])<<8 | uint64(d[24])
		pd.Penalty = uint64(d[25])<<56 | uint64(d[26])<<48 | uint64(d[27])<<40 |
			uint64(d[28])<<32 | uint64(d[29])<<24 | uint64(d[30])<<16 |
			uint64(d[31])<<8 | uint64(d[32])
	}

	return
}
