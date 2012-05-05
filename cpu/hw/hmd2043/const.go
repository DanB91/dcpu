// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package hmd2043

// Known interrupt messages.
const (
	QueryMediaPresent = iota
	QueryMediaParams
	QueryDeviceFlags
	UpdateDeviceFlags
	QueryInterruptType
	SetInterruptId
	ReadSectors       = 0x10
	WriteSectors      = 0x11
	QueryMediaQuality = 0xffff
)

// Known error codes
const (
	ErrorNone = iota
	ErrorNoMedia
	ErrorInvalidSector
	ErrorPending
)

// Device flags
const (
	NonBlocking = 1 << iota
	MediaStatusInterrupt
)

// Interrupt types
const (
	TypeNone = iota
	TypeMediaStatus
	TypeReadComplete
	TypeWriteComplete
)

// Media identifiers
const (
	AuthenticHITMedia = 0x7ff
	OtherMedia        = 0xffff
)
