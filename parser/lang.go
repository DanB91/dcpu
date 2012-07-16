// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import "strings"

var (
	instructions = [...]string{
		// Standard instructions.
		"set", "add", "sub", "mul", "mli", "div", "dvi", "mod", "mdi", "and",
		"bor", "xor", "shr", "asr", "shl", "ifb", "ifc", "ife", "ifn", "ifg",
		"ifa", "ifl", "ifu", "adx", "sbx", "sti", "std", "jsr", "int", "iag",
		"ias", "rfi", "iaq", "hwn", "hwq", "hwi",

		// Non-standard and pseudo instructions.
		"dat", "panic", "exit", "equ", "return",
	}

	registers = [...]string{
		"a", "b", "c", "x", "y", "z", "i", "j", "ia",
		"ex", "peek", "push", "pop", "pc", "sp",
	}

	branches = [...]string{
		"ifb", "ifc", "ife", "ifn", "ifg", "ifa", "ifl", "ifu",
	}
)

// IsBranch returns true if the given name qualifies as a branching instruction.
func IsBranch(v string) bool {
	for i := range branches {
		if strings.EqualFold(branches[i], v) {
			return true
		}
	}
	return false
}

// IsRegister returns true if the given value qualifies as a valid register name.
func IsRegister(v string) bool {
	for i := range registers {
		if strings.EqualFold(registers[i], v) {
			return true
		}
	}
	return false
}

// IsInstruction returns true if the given value qualifies as a valid instruction name.
func IsInstruction(v string) bool {
	for i := range instructions {
		if strings.EqualFold(instructions[i], v) {
			return true
		}
	}
	return false
}

// IsSpecialName returns true if the given name is a register or
// instruction name. This is short for:
//
//     IsInstruction(v) || IsRegister(v)
//
func IsSpecialName(v string) bool {
	return IsInstruction(v) || IsRegister(v)
}
