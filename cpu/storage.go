// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cpu

// All memory and register operations, all instructions are
// done in words.
type Word uint16

// Signed returns signed version of the given word.
func Signed(w Word) SWord {
	sign := int32(w >> 15)
	val := int32(w & 0x7fff)
	return SWord(val + (sign * -0xffff))
}

// Signed word. used in some signed arithmetic instructions.
type SWord int16

// Number of words in memory.
const MemSize = 0x10000

// Storage defines storage elements for a CPU.
type Storage struct {
	Mem                    [MemSize]Word // System memory.
	A, B, C, X, Y, Z, I, J Word          // General Purpose registers.
	PC, SP, EX, IA         Word          // Special registers.
}

// Clear clears the storage.
func (s *Storage) Clear() {
	s.A = 0
	s.B = 0
	s.C = 0
	s.X = 0
	s.Y = 0
	s.Z = 0
	s.I = 0
	s.J = 0
	s.PC = 0
	s.SP = 0xffff
	s.EX = 0
	s.IA = 0
}

// readString reads a string from the given memory address.
func (s *Storage) readString(addr Word) string {
	if int(addr) >= len(s.Mem) {
		return ""
	}

	runes := make([]rune, 0, 128)

	for ; int(addr) < len(s.Mem); addr++ {
		ch := s.Mem[addr] & 0xff

		if ch == 0 {
			break
		}

		runes = append(runes, rune(ch))
	}

	return string(runes)
}
