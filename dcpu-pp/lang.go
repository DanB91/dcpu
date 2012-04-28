// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

// Instructions and registers according to spec v1.7:
// http://pastebin.com/raw.php?i=Q4JvQvnM

var (
	instructions = []string{
		"set", "add", "sub", "mul", "mli", "div", "dvi", "mod", "mdi", "and",
		"bor", "xor", "shr", "asr", "shl", "ifb", "ifc", "ife", "ifn", "ifg",
		"ifa", "ifl", "ifu", "adx", "sbx", "sti", "std", "jsr", "int", "iag",
		"ias", "rfi", "iaq", "hwn", "hwq", "hwi",
	}

	regs = []string{
		"a", "b", "c", "x", "y", "z", "i", "j", "ia",
		"ex", "peek", "push", "pop", "pc", "sp",
	}
)

// isRegister returns true if the given value qualifies as
// a valid register name.
func isRegister(v string) bool {
	for i := range regs {
		if regs[i] == v {
			return true
		}
	}
	return false
}

// isInstruction returns true if the given value qualifies as
// a valid instruction name.
func isInstruction(v string) bool {
	for i := range instructions {
		if instructions[i] == v {
			return true
		}
	}
	return false
}
