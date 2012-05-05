// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cpu

import (
	"fmt"
	"testing"
)

type TestDevice struct {
	f IntFunc
}

func NewTestDevice(f IntFunc) Device {
	return &TestDevice{f}
}

func (d *TestDevice) Manufacturer() uint32 { return 0x12345678 }
func (d *TestDevice) Id() uint32           { return 0x87654321 }
func (d *TestDevice) Revision() uint16     { return 0xabcd }
func (d *TestDevice) Handler(s *Storage) {
	switch s.A {
	case 0x1: // SetMem
		s.Mem[s.B] = 0xbeef
	}
}

var _exit = enc(EXT, EXIT, 0x20) // EXIT 0

func enc(a, b, c Word) Word { return a | (b << 5) | (c << 10) }

func trace(pc, op, a, b Word, s *Storage) {
	fmt.Printf("%04x: %04x %04x %04x | %04x %04x %04x %04x %04x %04x %04x %04x | %04x %04x %04x\n",
		pc, op, a, b, s.A, s.B, s.C, s.X, s.Y, s.Z, s.I, s.J, s.SP, s.EX, s.IA)
}

func doTest(t *testing.T, c *CPU, result, overflow Word) {
	s := c.Store
	s.A, s.EX = 0, 0

	if err := c.Run(0); err != nil {
		t.Fatal(err)
	}

	if s.A != result {
		t.Fatalf("Want result 0x%04x, got 0x%04x", result, s.A)
	}

	if s.EX != overflow {
		t.Fatalf("Want overflow 0x%04x, got 0x%04x", overflow, s.EX)
	}
}

func TestSet(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0x30
	s.Mem[1] = 0x30
	s.Mem[2] = _exit
	doTest(t, c, 0x30, 0)
}

func TestAdd(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0xffff
	s.Mem[1] = 0xffff
	s.Mem[2] = enc(ADD, 0, 0x22) // ADD A, 0x1
	s.Mem[3] = _exit
	doTest(t, c, 0, 1)
}

func TestSub(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0x0, 0x21) // SET A, 0
	s.Mem[1] = enc(SUB, 0x0, 0x22) // SUB A, 1
	s.Mem[2] = _exit
	doTest(t, c, 0xffff, 0xffff)
}

func TestMul(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0x18
	s.Mem[1] = 0x18
	s.Mem[2] = enc(MUL, 0, 0x23) // MUL A, 0x2
	s.Mem[3] = _exit
	doTest(t, c, 0x18*0x2, 0)
}

func TestMli(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0x18
	s.Mem[1] = 0x18
	s.Mem[2] = enc(MLI, 0, 0x23) // MLI A, 0x2
	s.Mem[3] = _exit
	doTest(t, c, 0x18*0x2, 0)
}

func TestDiv(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0x60
	s.Mem[1] = 0x60
	s.Mem[2] = enc(DIV, 0, 0x23) // DIV A, 0x2
	s.Mem[3] = _exit
	doTest(t, c, 0x60/0x2, 0)
}

func TestMod(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x30) // SET A, 0xf
	s.Mem[1] = enc(MOD, 0, 0x23) // MOD A, 0x2
	s.Mem[2] = _exit
	doTest(t, c, 0xf%0x2, 0)
}

func TestAnd(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x24) // SET A, 0x3
	s.Mem[1] = enc(AND, 0, 0x22) // AND A, 0x1
	s.Mem[2] = _exit
	doTest(t, c, 0x3&0x1, 0)
}

func TestBor(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x22) // SET A, 0x1
	s.Mem[1] = enc(BOR, 0, 0x23) // OR A, 0x2
	s.Mem[2] = _exit
	doTest(t, c, 0x1|0x2, 0)
}

func TestXor(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[1] = enc(XOR, 0, 0x27) // XOR A, 0x6
	s.Mem[2] = _exit
	doTest(t, c, 0x4^0x6, 0)
}

func TestShr(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x31) // SET A, 0x10
	s.Mem[1] = enc(SHR, 0, 0x22) // SHR A, 0x1
	s.Mem[2] = _exit
	doTest(t, c, 0x10>>0x1, 0)
}

func TestAsr(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x31) // SET A, 0x10
	s.Mem[1] = enc(ASR, 0, 0x22) // ASR A, 0x1
	s.Mem[2] = _exit
	doTest(t, c, 0x10>>0x1, 0)
}

func TestShl(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x31) // SET A, 0x10
	s.Mem[1] = enc(SHL, 0, 0x22) // SHL A, 0x1
	s.Mem[2] = _exit
	doTest(t, c, 0x10<<0x1, 0)
}

func TestIfb(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[1] = enc(IFB, 0, 0x25) // IFB A, 0x4
	s.Mem[2] = enc(SET, 0, 0x21) // SET A, 0x0
	s.Mem[3] = _exit
	doTest(t, c, 0, 0)
}

func TestIfc(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[1] = enc(IFC, 0, 0x25) // IFC A, 0x4
	s.Mem[2] = enc(SET, 0, 0x21) // SET A, 0x0
	s.Mem[3] = _exit
	doTest(t, c, 0x4, 0)
}

func TestIfe(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[1] = enc(IFE, 0, 0x25) // IFE A, 0x4
	s.Mem[2] = enc(SET, 0, 0x21) // SET A, 0x0
	s.Mem[3] = _exit
	doTest(t, c, 0, 0)
}

func TestIfn(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[1] = enc(IFN, 0, 0x25) // IFN A, 0x4
	s.Mem[2] = enc(SET, 0, 0x21) // SET A, 0x0
	s.Mem[3] = _exit
	doTest(t, c, 0x4, 0)
}

func TestIfg(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[1] = enc(IFG, 0, 0x25) // IFG A, 0x4
	s.Mem[2] = enc(SET, 0, 0x21) // SET A, 0x0
	s.Mem[3] = _exit
	doTest(t, c, 0x4, 0)
}

func TestIfl(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[1] = enc(IFL, 0, 0x25) // IFL A, 0x4
	s.Mem[2] = enc(SET, 0, 0x21) // SET A, 0x0
	s.Mem[3] = _exit
	doTest(t, c, 0x4, 0)
}

func TestNestedIf(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x21) // SET A, 0x0
	s.Mem[1] = enc(IFG, 0, 0x22) // IFG A, 0x1
	s.Mem[2] = enc(IFG, 0, 0x23) // IFG A, 0x2
	s.Mem[3] = enc(IFE, 0, 0x21) // IFE A, 0x0
	s.Mem[4] = enc(SET, 0, 0x25) // SET A, 0x4
	s.Mem[5] = _exit
	doTest(t, c, 0, 0)
}

func TestADX(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0xffff
	s.Mem[1] = 0xffff
	s.Mem[2] = enc(ADD, 0, 0x22) // ADD A, 0x1
	s.Mem[3] = enc(SET, 0, 0x1f) // SET A, 0xffff
	s.Mem[4] = 0xffff
	s.Mem[5] = enc(ADX, 0, 0x22) // ADX A, 0x1
	s.Mem[6] = _exit
	doTest(t, c, 1, 1)
}

func TestSBX(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0xffff
	s.Mem[1] = 0xffff
	s.Mem[2] = enc(ADD, 0, 0x22) // ADD A, 0x1
	s.Mem[3] = enc(SBX, 0, 0x23) // SBX A, 0x2
	s.Mem[4] = _exit
	doTest(t, c, 0xffff, 0xffff)
}

func TestSTI(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 1, 0x22) // SET B, 0x1
	s.Mem[1] = enc(STI, 0, 1)    // STI A, B
	s.Mem[2] = enc(ADD, 0, 6)    // ADD A, I
	s.Mem[3] = enc(ADD, 0, 7)    // ADD A, J
	s.Mem[4] = _exit
	doTest(t, c, 3, 0)
}

func TestSTD(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 1, 0x22) // SET B, 0x1
	s.Mem[1] = enc(STD, 0, 1)    // STD A, B
	s.Mem[2] = enc(ADD, 0, 6)    // ADD A, I
	s.Mem[3] = enc(ADD, 0, 7)    // ADD A, J
	s.Mem[4] = _exit
	doTest(t, c, 0xffff, 0)
}

func TestJsr(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x1f) // SET A, 0xffff
	s.Mem[1] = 0xffff
	s.Mem[2] = enc(EXT, JSR, 0x25) // JSR my_sub
	s.Mem[3] = _exit
	s.Mem[4] = enc(ADD, 0, 0x22)    // :my_sub ADD A, 0x1
	s.Mem[5] = enc(SET, 0x1c, 0x18) // SET PC, POP
	doTest(t, c, 0, 1)
}

func TestIntRfi(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(EXT, IAS, 0x26) // IAS my_handler
	s.Mem[1] = enc(EXT, INT, 0x1f) // INT 0xbeef
	s.Mem[2] = 0xbeef
	s.Mem[3] = enc(SET, 0, 1) // SET A, B
	s.Mem[4] = _exit

	// :my_handler
	s.Mem[5] = enc(SET, 1, 0)    // SET B, A
	s.Mem[6] = enc(ADD, 1, 0x22) // ADD B, 1
	s.Mem[7] = enc(EXT, RFI, 0)  // RFI A
	doTest(t, c, 0xbef0, 0)
}

func TestIasIag(t *testing.T) {
	c := NewCPU()
	s := c.Store
	s.Mem[0] = enc(SET, 0, 0x22) // SET A, 1
	s.Mem[1] = enc(EXT, IAS, 0)  // IAS A
	s.Mem[2] = enc(SET, 0, 0x23) // SET A, 2
	s.Mem[3] = enc(EXT, IAG, 0)  // IAG A
	s.Mem[4] = _exit
	doTest(t, c, 1, 0)
}

func TestHwn(t *testing.T) {
	c := NewCPU()
	s := c.Store
	c.RegisterDevice(NewTestDevice)
	s.Mem[0] = enc(EXT, HWN, 0) // HWN A
	s.Mem[1] = _exit
	doTest(t, c, 1, 0)
}

func TestHwq(t *testing.T) {
	c := NewCPU()
	s := c.Store
	c.RegisterDevice(NewTestDevice)
	s.Mem[0] = enc(SET, 0, 0x21) // SET A, 0
	s.Mem[1] = enc(EXT, HWQ, 0)  // HWQ A
	s.Mem[2] = _exit
	doTest(t, c, 0x4321, 0)
}

func TestHwi(t *testing.T) {
	c := NewCPU()
	s := c.Store
	c.RegisterDevice(NewTestDevice)
	s.Mem[0] = enc(SET, 0, 0x22) // SET A, 0x1
	s.Mem[1] = enc(SET, 1, 0x1f) // SET B, 0x100
	s.Mem[2] = 0x100
	s.Mem[3] = enc(EXT, HWI, 0x21) // HWI 0  ; 0 = TestDevice
	s.Mem[4] = enc(SET, 0, 0x1e)   // SET A, [0x100]
	s.Mem[5] = 0x100
	s.Mem[6] = _exit
	doTest(t, c, 0xbeef, 0)
}
