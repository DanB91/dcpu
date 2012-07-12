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
//    [1] 32-bit unsigned integer:
//        Number of source file strings.
//    
//    [2] N number of utf-8 strings:
//        The source file names. Where N is the amount described in [1].
//        - Each string starts with an unsigned 16 bit value indicating the
//          length of the string in bytes.
//    
//        - Each string is written out as raw bytes.
//    
//    [3] 32-bit unsigned integer:
//        Total number of ProfileData entries. One for instruction in
//        the program.
//    
//    [4] 32-bit unsigned integer:
//        Number of ProfileData entries actually being used. Meaning for
//        each PC value we encountered.
//    
//    [5] N number of ProfileData entries.
//        Where N is the amount described in [4].
//    
//        - 16-bit unsigned integer:
//          The PC value to which this entry applies.
//    
//        - 8-bit unsigned integer:
//          The opcode to which this entry applies.
//    
//        - 8-bit unsigned integer:
//          Operand A for this instruction (op a, b).
//    
//        - 8-bit unsigned integer:
//          Operand B for this instruction (op a, b).
//    
//        - 16-bit unsigned integer:
//          Value of operand A (next word).
//    
//        - 16-bit unsigned integer:
//          Value of operand B (next word).
//    
//        - 32-bit unsigned int:
//          The file index for the original source code.
//          - This is an index into the list of files in section [2].
//    
//        - 32-bit unsigned int:
//          The line number for the original source code.
//    
//        - 32-bit unsigned int:
//          The column number for the original source code.
//    
//        - 64-bit unsigned int:
//          Number of times we executed this instruction.
//    
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
		size := uint16(len(p.Files[i]))

		err = binary.Write(w, be, size)
		if err != nil {
			return
		}

		_, err = w.Write([]byte(p.Files[i]))
		if err != nil {
			return
		}
	}

	// [3]
	size = uint32(len(p.Data))
	if err = binary.Write(w, be, size); err != nil {
		return
	}

	// [4]
	size = uint32(p.CountUses())
	if err = binary.Write(w, be, size); err != nil {
		return
	}

	// [5]
	var d [37]byte
	for pc, v := range p.Data {
		if v == nil {
			continue
		}

		d[0] = byte(pc >> 8)
		d[1] = byte(pc)

		d[2] = byte(v.Opcode)
		d[3] = byte(v.A)
		d[4] = byte(v.B)

		d[5] = byte(v.AValue >> 8)
		d[6] = byte(v.AValue)

		d[7] = byte(v.BValue >> 8)
		d[8] = byte(v.BValue)

		d[9] = byte(v.File >> 24)
		d[10] = byte(v.File >> 16)
		d[11] = byte(v.File >> 8)
		d[12] = byte(v.File)

		d[13] = byte(v.Line >> 24)
		d[14] = byte(v.Line >> 16)
		d[15] = byte(v.Line >> 8)
		d[16] = byte(v.Line)

		d[17] = byte(v.Col >> 24)
		d[18] = byte(v.Col >> 16)
		d[19] = byte(v.Col >> 8)
		d[20] = byte(v.Col)

		d[21] = byte(v.Count >> 56)
		d[22] = byte(v.Count >> 48)
		d[23] = byte(v.Count >> 40)
		d[24] = byte(v.Count >> 32)
		d[25] = byte(v.Count >> 24)
		d[26] = byte(v.Count >> 16)
		d[27] = byte(v.Count >> 8)
		d[28] = byte(v.Count)

		d[29] = byte(v.Penalty >> 56)
		d[30] = byte(v.Penalty >> 48)
		d[31] = byte(v.Penalty >> 40)
		d[32] = byte(v.Penalty >> 32)
		d[33] = byte(v.Penalty >> 24)
		d[34] = byte(v.Penalty >> 16)
		d[35] = byte(v.Penalty >> 8)
		d[36] = byte(v.Penalty)

		_, err = w.Write(d[:])
		if err != nil {
			return
		}
	}

	return
}
