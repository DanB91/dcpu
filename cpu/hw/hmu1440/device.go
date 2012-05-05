// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package hmu1440

import (
	"errors"
	"github.com/jteeuwen/dcpu/cpu"
	mmap "launchpad.net/gommap"
	"os"
)

// Known error codes (defined in HMD2043 spec)
const (
	ErrorNone = iota
	ErrorNoMedia
	ErrorInvalidSector
	ErrorPending
)

// Implements the 1.44 MB 3.5" Harold Media Unit.
//
// This device is backed by a memory mapped file.
type HMU1440 struct {
	fd *os.File
	m  mmap.MMap
}

// New creates and initializes a new device instance.
func New() *HMU1440 {
	d := new(HMU1440)
	return d
}

// Open opens th given file. It serves as the actual memory backend
// for this disk.
//
// Note that the file should be at lesat 1.44MB in size.
func (h *HMU1440) Open(file string) (err error) {
	if h.fd != nil {
		return errors.New("File already open.")
	}

	h.fd, err = os.Open(file)
	if err != nil {
		return
	}

	h.m, err = mmap.Map(h.fd.Fd(), mmap.PROT_READ|mmap.PROT_WRITE, mmap.MAP_PRIVATE)
	if err != nil {
		return
	}

	return
}

// Close closes the open file and unmaps memory.
func (h *HMU1440) Close() (err error) {
	if h.fd == nil {
		return
	}

	err = h.m.UnsafeUnmap()
	h.m = nil

	h.fd.Close()
	h.fd = nil
	return
}

func (h *HMU1440) WordsPerSector() cpu.Word { return 0 }
func (h *HMU1440) SectorCount() cpu.Word    { return 18 }
func (h *HMU1440) WriteLocked() bool        { return false }

// Read reads `len(sectors)` words from our underlying store into the
// `sectors` slice. Reading starts at offset `address`.
func (h *HMU1440) Read(address cpu.Word, sectors []cpu.Word) (err cpu.Word) {
	if h.fd == nil {
		return ErrorNoMedia
	}

	return
}

// Write writes `len(sectors)` words from the `sectors` slice into
// our underlying data store. Writing starts at offset `address`.
func (h *HMU1440) Write(address cpu.Word, sectors []cpu.Word) (err cpu.Word) {
	if h.fd == nil {
		return ErrorNoMedia
	}

	return
}
