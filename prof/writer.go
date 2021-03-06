// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package prof

import (
	"encoding/binary"
	"io"
)

var be = binary.BigEndian

// Write dumps the profile to the given writer in a binary format.
// All values are encoded as Big Endian.
//
// The layout of the file is as follows:
//
//    [N * Source file descriptors]
//    [N * Function descriptors]
//    [N * Instruction descriptors]
//
// More detailed:
//
//    [1] 32-bit unsigned integer:
//        Number of source files.
//    
//    [2] N number of source file definitions:
//        Where N is the amount described in [1].
//        - An unsigned 16 bit value indicating the start address in the
//          binary code for this file.
//        - An unsigned 16 bit value indicating the length of the
//          file name in bytes.
//        - Each file name is written out as raw bytes.
//
//    [3] 32-bit unsigned integer:
//        Number of function definitions.
//    
//    [4] N number of function definitions:
//        Where N is the amount described in [3].
//        - An unsigned 16 bit int: The functions's start address.
//        - An unsigned 16 bit int: The functions's end address.
//        - An unsigned 32 bit int: The functions's start line in source.
//        - An unsigned 32 bit int: The functions's end line in source.
//        - An unsigned 16 bit int: The length of the function name in bytes.
//        - Each function name is written out as raw bytes.
//    
//    [5] 32-bit unsigned integer:
//        Number of ProfileData entries.
//        One for instruction in the program.
//    
//    [6] N number of ProfileData entries.
//        Where N is the amount described in [5].
//        - 16-bit unsigned integer:
//          The encoded instruction to which this entry applies.
//        - 32-bit unsigned int:
//          The file index for the original source code.
//          - This is an index into the list of files in section [2].
//        - 32-bit unsigned int:
//          The line number for the original source code.
//        - 32-bit unsigned int:
//          The column number for the original source code.
//        - 64-bit unsigned int:
//          Number of times we executed this instruction.
//        - 64-bit unsigned int:
//          Cost penalty incurred at runtime.
//    
func Write(p *Profile, w io.Writer) (err error) {
	// [1]
	size := uint32(len(p.Files))
	if err = binary.Write(w, be, size); err != nil {
		return
	}

	// [2]
	for i := range p.Files {
		err = binary.Write(w, be, p.Files[i].StartAddr)
		if err != nil {
			return
		}

		size := uint16(len(p.Files[i].Name))

		err = binary.Write(w, be, size)
		if err != nil {
			return
		}

		_, err = w.Write([]byte(p.Files[i].Name))
		if err != nil {
			return
		}
	}

	// [3]
	size = uint32(len(p.Functions))
	if err = binary.Write(w, be, size); err != nil {
		return
	}

	// [4]
	for i := range p.Functions {
		err = binary.Write(w, be, p.Functions[i].StartAddr)
		if err != nil {
			return
		}

		err = binary.Write(w, be, p.Functions[i].EndAddr)
		if err != nil {
			return
		}

		err = binary.Write(w, be, uint32(p.Functions[i].StartLine))
		if err != nil {
			return
		}

		err = binary.Write(w, be, uint32(p.Functions[i].EndLine))
		if err != nil {
			return
		}

		size := uint16(len(p.Functions[i].Name))

		err = binary.Write(w, be, size)
		if err != nil {
			return
		}

		_, err = w.Write([]byte(p.Functions[i].Name))
		if err != nil {
			return
		}
	}

	// [5]
	size = uint32(len(p.Data))
	if err = binary.Write(w, be, size); err != nil {
		return
	}

	// [6]
	var d [30]byte
	for _, v := range p.Data {
		d[0] = byte(v.Data >> 8)
		d[1] = byte(v.Data)

		d[2] = byte(v.File >> 24)
		d[3] = byte(v.File >> 16)
		d[4] = byte(v.File >> 8)
		d[5] = byte(v.File)

		d[6] = byte(v.Line >> 24)
		d[7] = byte(v.Line >> 16)
		d[8] = byte(v.Line >> 8)
		d[9] = byte(v.Line)

		d[10] = byte(v.Col >> 24)
		d[11] = byte(v.Col >> 16)
		d[12] = byte(v.Col >> 8)
		d[13] = byte(v.Col)

		d[14] = byte(v.Count >> 56)
		d[15] = byte(v.Count >> 48)
		d[16] = byte(v.Count >> 40)
		d[17] = byte(v.Count >> 32)
		d[18] = byte(v.Count >> 24)
		d[19] = byte(v.Count >> 16)
		d[20] = byte(v.Count >> 8)
		d[21] = byte(v.Count)

		d[22] = byte(v.Penalty >> 56)
		d[23] = byte(v.Penalty >> 48)
		d[24] = byte(v.Penalty >> 40)
		d[25] = byte(v.Penalty >> 32)
		d[26] = byte(v.Penalty >> 24)
		d[27] = byte(v.Penalty >> 16)
		d[28] = byte(v.Penalty >> 8)
		d[29] = byte(v.Penalty)

		_, err = w.Write(d[:])
		if err != nil {
			return
		}
	}

	return
}
