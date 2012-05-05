// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package hmu1440

import (
	"errors"
	"github.com/jteeuwen/dcpu/cpu"
	"io/ioutil"
	"reflect"
	"unsafe"
)

const (
	WordCount   = 737280
	SectorCount = 1440
	SectorSize  = WordCount / SectorCount // Sector size in words.
)

var (
	ErrNoMedia       = errors.New("No media present.")
	ErrSizeMismatch  = errors.New("Buffer size should be a multiple of SectorSize words.")
	ErrInvalidSector = errors.New("Access violation: invalid sector.")
)

// Implements the 1.44 MB 3.5" Harold Media Unit.
//
// This device is backed by a file.
type HMU1440 struct {
	data []byte
	file string
}

// New creates and initializes a new device instance.
func New() *HMU1440 {
	return new(HMU1440)
}

// Open opens th given file. It serves as the actual memory backend
// for this disk.
//
// Note that the file should be at lesat 1.44MB in size.
func (h *HMU1440) Open(file string) (err error) {
	h.file = file
	h.data, err = ioutil.ReadFile(file)
	return
}

// Close writes the data to the underlying file and cleans things up.
func (h *HMU1440) Close() (err error) {
	err = ioutil.WriteFile(h.file, h.data, 0600)
	h.data = nil
	h.file = ""
	return
}

func (h *HMU1440) SectorSize() cpu.Word  { return SectorSize }
func (h *HMU1440) SectorCount() cpu.Word { return 18 }
func (h *HMU1440) WriteLocked() bool     { return false }

// Read reads `len(buffer)` words from our underlying store into the
// `buffer` slice. Reading starts at sector `sector`.
//
// The supplied buffer length should be a multiple of `SectorSize`.
func (h *HMU1440) Read(sector cpu.Word, buffer []cpu.Word) error {
	return h.copy(sector, buffer, true)
}

// Write writes `len(buffer)` words from the `buffer` slice into
// our underlying data store. Writing starts at sector `sector`.
//
// The supplied buffer length should be a multiple of `SectorSize`.
func (h *HMU1440) Write(sector cpu.Word, buffer []cpu.Word) error {
	return h.copy(sector, buffer, false)
}

// copy copies a number of sectors to/from the given buffer.
// The read value edtermines in which direction the operation goes.
func (h *HMU1440) copy(sector cpu.Word, buffer []cpu.Word, read bool) (err error) {
	if len(h.data) == 0 {
		return ErrNoMedia
	}

	sz := len(buffer)
	if sz%SectorSize != 0 {
		return ErrSizeMismatch
	}

	sz = len(buffer) / SectorSize
	if sector >= SectorCount || int(sector)+sz >= SectorCount {
		return ErrInvalidSector
	}

	// File stores data in bytes. We have 16 bit words.
	// Unfortunately, this means we can't use copy() as things are now.
	//
	// We are going to use some black magic to turn the buffer into a byte
	// slice so copy() will become an option.
	ptr := *(*[]byte)(unsafe.Pointer(&buffer))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&ptr))

	// Alter the slice header to reflect the new len/capacity.
	// One word equals two bytes, so we need to double both
	// in order to make the copy() call work.
	sh.Len *= 2
	sh.Cap *= 2

	// Find the offset in bytes.
	offset := sector * SectorSize * 2

	if read {
		copy(ptr, h.data[offset:])
	} else {
		copy(h.data[offset:], ptr)
	}

	// Return len/capacity to old values.
	sh.Len /= 2
	sh.Cap /= 21
	return
}
