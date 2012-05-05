// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import "github.com/jteeuwen/dcpu/cpu"

type opcode struct {
	code cpu.Word
	ext  bool // Extended opcode?
}

var opcodes = map[string]opcode{
	"set":  {0x01, false},
	"add":  {0x02, false},
	"sub":  {0x03, false},
	"mul":  {0x04, false},
	"mli":  {0x05, false},
	"div":  {0x06, false},
	"dvi":  {0x07, false},
	"mod":  {0x08, false},
	"and":  {0x09, false},
	"bor":  {0x0a, false},
	"xor":  {0x0b, false},
	"shr":  {0x0c, false},
	"asr":  {0x0d, false},
	"shl":  {0x0e, false},
	"ifb":  {0x10, false},
	"ifc":  {0x11, false},
	"ife":  {0x12, false},
	"ifn":  {0x13, false},
	"ifg":  {0x14, false},
	"ifa":  {0x15, false},
	"ifl":  {0x16, false},
	"ifu":  {0x17, false},
	"adx":  {0x1a, false},
	"sbx":  {0x1b, false},
	"sti":  {0x1e, false},
	"std":  {0x1f, false},
	"jsr":  {0x01, true},
	"int":  {0x08, true},
	"iag":  {0x09, true},
	"ias":  {0x0a, true},
	"rfi":  {0x0b, true},
	"iaq":  {0x0c, true},
	"hwn":  {0x10, true},
	"hwq":  {0x11, true},
	"hwi ": {0x12, true},
	"test": {0x1e, true},
	"exit": {0x1f, true},
}
