// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"io"
)

var newline = []byte{'\n'}

type ByteWriter struct {
	io.Writer
	c int
}

func (w *ByteWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	for n = range p {
		if w.c%12 == 0 {
			w.Writer.Write(newline)
			w.c = 0
		}

		fmt.Fprintf(w.Writer, "0x%02x,", p[n])
		w.c++
	}

	n++

	return
}
