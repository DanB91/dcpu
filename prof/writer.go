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
	size = uint32(len(p.Usage))
	if err = binary.Write(w, be, size); err != nil {
		return
	}

	// [4]
	size = uint32(p.CountUses())
	if err = binary.Write(w, be, size); err != nil {
		return
	}

	// [5]
	var d [25]byte
	for pc, v := range p.Usage {
		if v == nil {
			continue
		}

		d[0], d[1] = byte(pc>>8), byte(pc)
		d[2], d[3], d[4] = byte(v.Opcode), byte(v.A), byte(v.B)
		d[5], d[6], d[7], d[8] = byte(v.File>>24), byte(v.File>>16), byte(v.File>>8), byte(v.File)
		d[9], d[10], d[11], d[12] = byte(v.Line>>24), byte(v.Line>>16), byte(v.Line>>8), byte(v.Line)
		d[13], d[14], d[15], d[16] = byte(v.Col>>24), byte(v.Col>>16), byte(v.Col>>8), byte(v.Col)
		d[17], d[18], d[19], d[20], d[21], d[22], d[23], d[24] =
			byte(v.Count>>56), byte(v.Count>>48), byte(v.Count>>40), byte(v.Count>>32),
			byte(v.Count>>24), byte(v.Count>>16), byte(v.Count>>8), byte(v.Count)

		_, err = w.Write(d[:])
		if err != nil {
			return
		}
	}

	return
}
