// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cpu

import (
	"io"
	"time"
)

// The OverloadError occurs when we exceed the capacity
// of the interrupt queue. The CPU should catch fire here and
// do all manner of strange things to your memory/registers.
//
// For the sake of testing, we simply return an error.
type OverloadError struct{}

func (e OverloadError) Error() string {
	return "Interrupt overload: System on fire!"
}

// Maximum interrupt queue size.
const MaxIntQueue = 0xff

// This signature represents a Debug trace handler.
type TraceFunc func(pc, op, a, b Word, store *Storage)

// A CPU can run a single program.
type CPU struct {
	Store           *Storage      // Memory and registers
	devices         []Device      // List of hardware devices.
	intQueue        chan Word     // Interrupt queue.
	Trace           TraceFunc     // Allow tracing of instructions as they are executed. For debug only.
	ClockSpeed      time.Duration // Speed of CPU clock.
	size            Word          // Size of last instruction (in words).
	queueInterrupts bool          // Use interrupt queueing or not.
}

// New creates and initializes a new CPU instance.
func New() *CPU {
	c := new(CPU)
	c.ClockSpeed = 1000 // 100khz
	c.Store = new(Storage)
	c.Reset()
	return c
}

// Devices returns the current list of registered devices.
func (c *CPU) Devices() []Device { return c.devices }

// CountDevices returns the number of currently registered devices.
func (c *CPU) CountDevices() int { return len(c.devices) }

// Register adds a new device. If capacity has been reached,
// this is silently ignored. We can have a maximum of MaxUint16 number
// of devices at any given time.
func (c *CPU) RegisterDevice(db DeviceBuilder) {
	if len(c.devices) < 1<<16-1 {
		c.devices = append(c.devices, db(func(w Word) { c.interrupt(w) }))
	}
}

// ClearDevices removes all registered devices.
func (c *CPU) ClearDevices() { c.devices = nil }

// Clears CPU state.
func (c *CPU) Reset() {
	c.Store.Clear()
	c.intQueue = make(chan Word, MaxIntQueue)
	c.queueInterrupts = false
}

// interrupt either queues or triggers an interrupt with the given message.
// These can come from hardware or software.
func (c *CPU) interrupt(msg Word) {
	if c.queueInterrupts {
		if len(c.intQueue) >= MaxIntQueue {
			panic(OverloadError{})
		}

		c.intQueue <- msg
		return
	}

	if c.Store.IA != 0 {
		c.triggerInterrupt(msg)
	}
}

// triggerInterrupt triggers an interrupt handler with the given message.
//
// This should always be balanced by an RFI instruction at the end
// of each interrupt handler.
func (c *CPU) triggerInterrupt(msg Word) {
	s := c.Store
	c.queueInterrupts = true
	s.Mem[s.SP], s.SP = s.PC, s.SP-1 // PUSH PC
	s.Mem[s.SP], s.SP = s.A, s.SP-1  // PUSH A
	s.PC = s.IA
	s.A = msg
}

// Run runs code, starting at the given entrypoint.
// This repeatedly calls cpu.Step for as long as the program is valid.
//
// In order to step through code for debugging, call cpu.Step
// manually.
func (c *CPU) Run(entrypoint Word) (err error) {
	s := c.Store
	s.PC = entrypoint

	for {
		select {
		case <-time.After(c.ClockSpeed):
			err = c.Step()

			if err != nil {
				if err == io.EOF {
					err = nil // No need to propagate this one.
				}
				return
			}
		}
	}

	return
}

// nextInstruction decodes the next instruction in the
// supplied word pointers.
func (c *CPU) nextInstruction() (op, a, b Word) {
	s := c.Store
	w := s.Mem[s.PC]
	op, a, b = w&0x1f, (w>>5)&0x1f, (w>>10)&0x3f
	c.size = wordCount(op, a, b)
	s.PC++
	return
}

// skip skips a single instruction
func (c *CPU) skip() Word {
	s := c.Store
	w := s.Mem[s.PC]
	op, a, b := w&0x1f, (w>>5)&0x1f, (w>>10)&0x3f
	s.PC += wordCount(op, a, b)
	return op
}

// When a branching instruction fails its check, we need to skip
// all subsequent branching instructions + the next non-branching instruction.
// This behaviour allows us to use nested IF[X] blocks.
//
// For example, the following should skip straight to `exit 0`
// when the first `ifg` fails.
//
//  [pc]
//   |-> set a, 0
//   |-> ifg a, 1
//   |    ifg a, 2
//   |     ife a, 0
//   |      set a, 4
//   |-> exit 0
//
func (c *CPU) skipBranch() {
	for {
		if op := c.skip(); op < IFB || op > IFU {
			break
		}
	}
}

// Step performs a single cycle.
func (c *CPU) Step() (err error) {
	var va, vb *Word

	// Handle any queued interrupts.
	// The way this is handled, is not entirely as defined in
	// the spec. We can only handle one interrupt per clock cycle.
	// 
	// However, this implementation allows one queued interrupt,
	// as well as at least one non-ququed interrupt to be triggered
	// in a single cycle. Notably when we are about to execute
	// a non-cached HWI or INT.
	if sz := len(c.intQueue); !c.queueInterrupts && sz > 0 {
		c.triggerInterrupt(<-c.intQueue)
	}

	s := c.Store
	op, a, b := c.nextInstruction()

	// Resolve operands.
	if op != EXT {
		va = c.decodeOperand(a, true)
	}

	vb = c.decodeOperand(b, false)

	// Trace output for debugging?
	if c.Trace != nil {
		c.Trace(s.PC-c.size, op, a, b, s)
	}

	switch op {
	case SET:
		*va = *vb

	case ADD:
		u32 := uint32(*va) + uint32(*vb)
		*va = Word(u32 & 0xffff)
		s.EX = Word(u32 >> 16)

	case SUB:
		u32 := uint32(*va) - uint32(*vb)
		*va = Word(u32 & 0xffff)
		s.EX = Word(u32 >> 16)

	case MUL:
		u32 := uint32(*va) * uint32(*vb)
		*va = Word(u32 & 0xffff)
		s.EX = Word(u32 >> 16)

	case MLI:
		i32 := Signed(*va) * Signed(*vb)
		*va = Word(i32 & 0x7fff)
		s.EX = Word(i32 >> 16)

	case DIV:
		if *vb == 0 {
			*va, s.EX = 0, 0
			break
		}

		*va /= *vb
		s.EX = Word(((uint32(*va) << 16) / uint32(*vb)))

	case DVI:
		if *vb == 0 {
			*va, s.EX = 0, 0
			break
		}

		sa, sb := Signed(*va), Signed(*vb)
		*va = Word(sa / sb)
		s.EX = Word((sa << 16) / sb)

	case MOD:
		if *vb == 0 {
			*va = 0
		} else {
			*va = Word(Signed(*va) % Signed(*vb))
		}

	case AND:
		*va &= *vb

	case BOR:
		*va |= *vb

	case XOR:
		*va ^= *vb

	case SHR:
		*va >>= *vb
		s.EX = Word(((uint32(*va) << 16) >> uint32(*vb)))

	case ASR:
		sa := Signed(*va)
		*va = Word(sa >> *vb)
		s.EX = Word((int32(sa) << 16) >> uint32(*vb))

	case SHL:
		u32 := uint32(*va) << uint32(*vb)
		*va = Word(u32 & 0xffff)
		s.EX = Word(u32 >> 16)

	case IFB:
		if (*va & *vb) == 0 {
			c.skipBranch()
		}

	case IFC:
		if (*va & *vb) != 0 {
			c.skipBranch()
		}

	case IFE:
		if *va != *vb {
			c.skipBranch()
		}

	case IFN:
		if *va == *vb {
			c.skipBranch()
		}

	case IFG:
		if *va <= *vb {
			c.skipBranch()
		}

	case IFA:
		if Signed(*va) <= Signed(*vb) {
			c.skipBranch()
		}

	case IFL:
		if *va >= *vb {
			c.skipBranch()
		}

	case IFU:
		if Signed(*va) >= Signed(*vb) {
			c.skipBranch()
		}

	case ADX:
		u32 := uint32(*va) + uint32(*vb) + uint32(s.EX)
		*va = Word(u32 & 0xffff)
		s.EX = Word(u32 >> 16)

	case SBX:
		u32 := uint32(*va) - uint32(*vb) + uint32(s.EX)
		*va = Word(u32 & 0xffff)
		s.EX = Word(u32 >> 16)

	case STI:
		*va = *vb
		s.I++
		s.J++

	case STD:
		*va = *vb
		s.I--
		s.J--

	case EXT:
		switch a {
		case JSR:
			s.Mem[s.SP] = s.PC
			s.SP--
			s.PC = *vb

		case INT:
			c.interrupt(*vb)

		case IAG:
			*vb = s.IA

		case IAS:
			s.IA = *vb

		case RFI:
			c.queueInterrupts = false
			s.SP++
			s.A = s.Mem[s.SP]
			s.SP++
			s.PC = s.Mem[s.SP]

		case IAQ:
			c.queueInterrupts = *vb != 0

		case HWN:
			*vb = Word(len(c.devices))

		case HWQ:
			if *vb >= Word(len(c.devices)) {
				s.A, s.B, s.C, s.X, s.Y = 0, 0, 0, 0, 0
				break
			}

			dev := c.devices[*vb]

			w := dev.Id()
			s.A = Word(w & 0xffff)
			s.B = Word((w >> 16) & 0xffff)

			s.C = Word(dev.Revision())

			w = dev.Manufacturer()
			s.X = Word(w & 0xffff)
			s.Y = Word((w >> 16) & 0xffff)

		case HWI:
			if *vb < Word(len(c.devices)) {
				c.devices[*vb].Handler(s)
			}

		case TEST:

		case EXIT:
			return io.EOF
		}
	}

	return
}

// decodeOperand interprets the given instruction operand and returns a pointer
// to the appropriate storage bit along with its address. 
//
// isTarget deterines if this operand is  the write target.
// This is necessary to properly decode the PUSH/POP operands (0x18).
func (c *CPU) decodeOperand(w Word, isTarget bool) *Word {
	// literal value 0xffff-0x1e (-1..30)
	if w >= 0x20 && w <= 0x3f {
		w -= 0x21
		return &w
	}

	var a Word
	s := c.Store

	switch w {
	// register (A, B, C, X, Y, Z, I or J)
	case 0x0:
		return &s.A
	case 0x1:
		return &s.B
	case 0x2:
		return &s.C
	case 0x3:
		return &s.X
	case 0x4:
		return &s.Y
	case 0x5:
		return &s.Z
	case 0x6:
		return &s.I
	case 0x7:
		return &s.J

	// [register]
	case 0x8:
		return &s.Mem[s.A]
	case 0x9:
		return &s.Mem[s.B]
	case 0xa:
		return &s.Mem[s.C]
	case 0xb:
		return &s.Mem[s.X]
	case 0xc:
		return &s.Mem[s.Y]
	case 0xd:
		return &s.Mem[s.Z]
	case 0xe:
		return &s.Mem[s.I]
	case 0xf:
		return &s.Mem[s.J]

	// [next word + register]
	case 0x10:
		a, s.PC = s.Mem[s.PC]+s.A, s.PC+1
		return &s.Mem[a]
	case 0x11:
		a, s.PC = s.Mem[s.PC]+s.B, s.PC+1
		return &s.Mem[a]
	case 0x12:
		a, s.PC = s.Mem[s.PC]+s.C, s.PC+1
		return &s.Mem[a]
	case 0x13:
		a, s.PC = s.Mem[s.PC]+s.X, s.PC+1
		return &s.Mem[a]
	case 0x14:
		a, s.PC = s.Mem[s.PC]+s.Y, s.PC+1
		return &s.Mem[a]
	case 0x15:
		a, s.PC = s.Mem[s.PC]+s.Z, s.PC+1
		return &s.Mem[a]
	case 0x16:
		a, s.PC = s.Mem[s.PC]+s.I, s.PC+1
		return &s.Mem[a]
	case 0x17:
		a, s.PC = s.Mem[s.PC]+s.J, s.PC+1
		return &s.Mem[a]

	// isTarget ? (PUSH / [--SP]) : (POP / [SP++])
	case 0x18:
		if isTarget {
			s.SP--
			return &s.Mem[s.SP+1]
		}

		s.SP++
		return &s.Mem[s.SP]

	// [SP] / PEEK
	case 0x19:
		return &s.Mem[s.SP]

	// [SP + next word] / PICK n
	case 0x1a:
		a, s.PC = s.Mem[s.PC], s.PC+1
		return &s.Mem[a+s.SP]

	case 0x1b:
		return &s.SP

	case 0x1c:
		return &s.PC

	case 0x1d:
		return &s.EX

	// [next word]
	case 0x1e:
		a, s.PC = s.Mem[s.PC], s.PC+1
		return &s.Mem[a]

	// Next word (literal)
	case 0x1f:
		s.PC++
		return &s.Mem[s.PC-1]
	}

	return &w
}

// wordCount counts the number of words occupied by the next instruction.
// This is needed to correctly skip the appropriate amount of words
// in the IF(x) instructions.
func wordCount(op, a, b Word) (count Word) {
	count = 1

	if op != EXT && (a == 0x1e || a == 0x1f || (a >= 0x10 && a <= 0x17)) {
		count++
	}

	if b == 0x1e || b == 0x1f || (b >= 0x10 && b <= 0x17) {
		count++
	}

	return
}
