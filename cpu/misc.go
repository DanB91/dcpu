// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// DCPU emulator package.
package cpu

// Maximum interrupt queue size.
const MaxIntQueue = 0xff

type TraceFunc func(pc, op, a, b Word, store *Storage)

type InstructionFunc func(pc Word, store *Storage)

type BranchSkipFunc func(pc, cost Word)

// Encode encodes the given opcode and operands into an instruction.
func Encode(a, b, c Word) Word {
	return a | (b << 5) | (c << 10)
}

// Decode decodes the given word into an opcode and operands.
func Decode(w Word) (op, a, b Word) {
	return w & 0x1f, (w >> 5) & 0x1f, (w >> 10) & 0x3f
}
