// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cpu

// An IntFunc is called by a device when it wants to issue
// an interrupt to the CPU.
type IntFunc func(Word)

// DeviceBuilder serves as a common constructor type for
// devices. All devices should implement this.
//
// The channel it receives is to be used as a CPU interrupt.
// Whenever the device is ready to be used, it sends a mask value
// over this channel. The CPU will then determine if the interrupt
// should be handled or not based on this mask.
type DeviceBuilder func(IntFunc) Device

// Device represents an arbitrary hardware module that
// can be plugged into the DCPU system.
//
// Code can communicate with it through use of interupts.
type Device interface {
	// 32 bit code, identifying the hardware manufacturer.
	Manufacturer() uint32

	// 32 bit id, identifying the type of device.
	Id() uint32

	// 16 bit revision number for device.
	Revision() uint16

	// The device's interrupt handler.
	//
	// It gets passed the CPU's storage. This means that during
	// execution of the handler, the device has full control over
	// CPU registers and memory.
	Handler(*Storage)
}
