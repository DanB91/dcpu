// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"github.com/jteeuwen/dcpu/cpu"
	"io"
	"os"
)

func writeProgram(program []cpu.Word, file string) (err error) {
	var w io.Writer

	if len(file) == 0 {
		w = os.Stdout
	} else {
		fd, err := os.Create(file)
		if err != nil {
			return err
		}

		defer fd.Close()
		w = fd
	}

	var b [2]byte
	le := *littleendian

	for _, word := range program {
		if le {
			b[0] = byte(word & 0xff)
			b[1] = byte((word >> 8) & 0xff)
		} else {
			b[0] = byte((word >> 8) & 0xff)
			b[1] = byte(word & 0xff)
		}

		_, err = w.Write(b[:])
		if err != nil {
			return
		}
	}

	return
}
