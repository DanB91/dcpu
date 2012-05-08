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
	"set":   {cpu.SET, 2, false},
	"add":   {cpu.ADD, 2, false},
	"sub":   {cpu.SUB, 2, false},
	"mul":   {cpu.MUL, 2, false},
	"mli":   {cpu.MLI, 2, false},
	"div":   {cpu.DIV, 2, false},
	"dvi":   {cpu.DVI, 2, false},
	"mod":   {cpu.MOD, 2, false},
	"mdi":   {cpu.MDI, 2, false},
	"and":   {cpu.AND, 2, false},
	"bor":   {cpu.BOR, 2, false},
	"xor":   {cpu.XOR, 2, false},
	"shr":   {cpu.SHR, 2, false},
	"asr":   {cpu.ASR, 2, false},
	"shl":   {cpu.SHL, 2, false},
	"ifb":   {cpu.IFB, 2, false},
	"ifc":   {cpu.IFC, 2, false},
	"ife":   {cpu.IFE, 2, false},
	"ifn":   {cpu.IFN, 2, false},
	"ifg":   {cpu.IFG, 2, false},
	"ifa":   {cpu.IFA, 2, false},
	"ifl":   {cpu.IFL, 2, false},
	"ifu":   {cpu.IFU, 2, false},
	"adx":   {cpu.ADX, 2, false},
	"sbx":   {cpu.SBX, 2, false},
	"sti":   {cpu.STI, 2, false},
	"std":   {cpu.STD, 2, false},
	"jsr":   {cpu.JSR, 1, true},
	"int":   {cpu.INT, 1, true},
	"iag":   {cpu.IAG, 1, true},
	"ias":   {cpu.IAS, 1, true},
	"rfi":   {cpu.RFI, 1, true},
	"iaq":   {cpu.IAQ, 1, true},
	"hwn":   {cpu.HWN, 1, true},
	"hwq":   {cpu.HWQ, 1, true},
	"hwi":   {cpu.HWI, 1, true},
	"panic": {cpu.PANIC, 1, true},
	"exit":  {cpu.EXIT, 0, true},
	"dat":   {0, 0, false}, // Pseudo-instruction
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
	"pop":  0x18,
	"push": 0x18,
	"peek": 0x19,
	"sp":   0x1b,
	"pc":   0x1c,
	"ex":   0x1d,
}
