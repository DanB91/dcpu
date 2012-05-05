// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cpu

// List of known basic opcodes.
const (
	EXT = 0x00
	SET = 0x01 // SET b, a | sets b to a
	ADD = 0x02 // ADD b, a | sets b to b+a, sets EX to 0x0001 if there's an overflow,  0x0 otherwise
	SUB = 0x03 // SUB b, a | sets b to b-a, sets EX to 0xffff if there's an underflow, 0x0 otherwise
	MUL = 0x04 // MUL b, a | sets b to b*a, sets EX to ((b*a)>>16)&0xffff (treats b, a as unsigned)
	MLI = 0x05 // MLI b, a | like MUL, but treat b, a as signed
	DIV = 0x06 // DIV b, a | sets b to b/a, sets EX to ((b<<16)/a)&0xffff. if a==0, sets b and EX to 0 instead. (treats b, a as unsigned)
	DVI = 0x07 // DVI b, a | like DIV, but treat b, a as signed
	MOD = 0x08 // MOD b, a | sets b to b%a. if a==0, sets b to 0 instead.
	AND = 0x09 // AND b, a | sets b to b&a
	BOR = 0x0a // BOR b, a | sets b to b|a
	XOR = 0x0b // XOR b, a | sets b to b^a
	SHR = 0x0c // SHR b, a | sets b to b>>>a, sets EX to ((b<<16)>>a)&0xffff  (logical shift)
	ASR = 0x0d // ASR b, a | sets b to b>>a, sets EX to ((b<<16)>>>a)&0xffff  (arithmetic shift) (treats b as signed)
	SHL = 0x0e // SHL b, a | sets b to b<<a, sets EX to ((b<<a)>>16)&0xffff
	IFB = 0x10 // IFB b, a | performs next instruction only if (b&a)!=0
	IFC = 0x11 // IFC b, a | performs next instruction only if (b&a)==0
	IFE = 0x12 // IFE b, a | performs next instruction only if b==a 
	IFN = 0x13 // IFN b, a | performs next instruction only if b!=a 
	IFG = 0x14 // IFG b, a | performs next instruction only if b>a 
	IFA = 0x15 // IFA b, a | performs next instruction only if b>a (signed)
	IFL = 0x16 // IFL b, a | performs next instruction only if b<a 
	IFU = 0x17 // IFU b, a | performs next instruction only if b<a (signed)
	ADX = 0x1a // ADX b, a | sets b to b+a+EX, sets EX to 0x0001 if there is an overflow, 0x0 otherwise
	SBX = 0x1b // SBX b, a | sets b to b-a+EX, sets EX to 0xFFFF if there is an underflow, 0x0 otherwise
	STI = 0x1e // STI b, a | sets b to a, then increases I and J by 1
	STD = 0x1f // STD b, a | sets b to a, then decreases I and J by 1
)

// Extended instructions.
const (
	JSR = 0x01 // JSR a - pushes the address of the next instruction to the stack, then sets PC to a.
	INT = 0x08 // INT a | triggers a software interrupt with message a
	IAG = 0x09 // IAG a | sets a to IA 
	IAS = 0x0a // IAS a | sets IA to a
	RFI = 0x0b // RFI a | disables interrupt queueing, pops A from the stack, then  pops PC from the stack
	IAQ = 0x0c /*
		IAQ a | if a is nonzero, interrupts will be added to the queue
		instead of triggered. if a is zero, interrupts will be
		triggered as normal again.
	*/
	HWN = 0x10 // HWN a | sets a to number of connected hardware devices
	HWQ = 0x11 /*
		HWQ a sets A, B, C, X, Y registers to information about hardware a.
		 - a+b is a 32 bit word identifying the hardware.
		 - type c is the hardware revision
		 - x+y is a 32 bit word identifying the manufacturer.
	*/
	HWI  = 0x12 // HWI a | sends an interrupt to hardware a
	TEST = 0x1e // (Non-standard) TEST instruction. Used to perform unit tests.
	EXIT = 0x1f // (Non-standard) Exit instruction. Stops the world.
)
