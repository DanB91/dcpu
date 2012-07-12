// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import (
	"encoding/binary"
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

	p.Data = make([]ProfileData, size)

	// [4]
	var d [30]byte
	for i := range p.Data {
		if _, err = r.Read(d[:]); err != nil {
			return
		}

		var pd ProfileData

		pd.Data = cpu.Word(d[0])<<8 | cpu.Word(d[1])
		pd.File = int(d[2])<<24 | int(d[3])<<16 | int(d[4])<<8 | int(d[5])
		pd.Line = int(d[6])<<24 | int(d[7])<<16 | int(d[8])<<8 | int(d[9])
		pd.Col = int(d[10])<<24 | int(d[11])<<16 | int(d[12])<<8 | int(d[13])

		pd.Count = uint64(d[14])<<56 | uint64(d[15])<<48 | uint64(d[16])<<40 |
			uint64(d[17])<<32 | uint64(d[18])<<24 | uint64(d[19])<<16 |
			uint64(d[20])<<8 | uint64(d[21])

		pd.Penalty = uint64(d[22])<<56 | uint64(d[23])<<48 | uint64(d[24])<<40 |
			uint64(d[25])<<32 | uint64(d[26])<<24 | uint64(d[27])<<16 |
			uint64(d[28])<<8 | uint64(d[29])

		p.Data[i] = pd
	}

	return
}
