// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import "github.com/jteeuwen/dcpu/cpu"

type opcode struct {
	code cpu.Word // Binary opcode number.
	argc int      // Number of arguments.
	ext  bool     // Extended opcode?
}

var opcodes = map[string]opcode{
	"set":  {0x01, 2, false},
	"add":  {0x02, 2, false},
	"sub":  {0x03, 2, false},
	"mul":  {0x04, 2, false},
	"mli":  {0x05, 2, false},
	"div":  {0x06, 2, false},
	"dvi":  {0x07, 2, false},
	"mod":  {0x08, 2, false},
	"and":  {0x09, 2, false},
	"bor":  {0x0a, 2, false},
	"xor":  {0x0b, 2, false},
	"shr":  {0x0c, 2, false},
	"asr":  {0x0d, 2, false},
	"shl":  {0x0e, 2, false},
	"ifb":  {0x10, 2, false},
	"ifc":  {0x11, 2, false},
	"ife":  {0x12, 2, false},
	"ifn":  {0x13, 2, false},
	"ifg":  {0x14, 2, false},
	"ifa":  {0x15, 2, false},
	"ifl":  {0x16, 2, false},
	"ifu":  {0x17, 2, false},
	"adx":  {0x1a, 2, false},
	"sbx":  {0x1b, 2, false},
	"sti":  {0x1e, 2, false},
	"std":  {0x1f, 2, false},
	"jsr":  {0x01, 1, true},
	"int":  {0x08, 1, true},
	"iag":  {0x09, 1, true},
	"ias":  {0x0a, 1, true},
	"rfi":  {0x0b, 1, true},
	"iaq":  {0x0c, 1, true},
	"hwn":  {0x10, 1, true},
	"hwq":  {0x11, 1, true},
	"hwi ": {0x12, 1, true},
	"test": {0x1e, 0, true},
	"exit": {0x1f, 0, true},
}

var registers = map[string]cpu.Word{
	"a":    0x0,
	"b":    0x1,
	"c":    0x2,
	"x":    0x3,
	"y":    0x4,
	"z":    0x5,
	"i":    0x6,
	"j":    0x7,
	"peek": 0x19,
	"pop":  0x18,
	"push": 0x18,
	"sp":   0x1b,
	"pc":   0x1c,
	"ex":   0x1d,
}
