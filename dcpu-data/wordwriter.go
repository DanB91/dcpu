// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"io"
)

// Maximum number of words to write per line.
const WordsPerLine = 9

func WriteWords(w io.Writer, p []byte, le bool) (err error) {
	var word uint16
	var count int

	newline := []byte("\ndat ")
	spacing := []byte(", ")

	for n := 0; n < len(p); n += 2 {
		if count%WordsPerLine == 0 {
			if count > 0 {
				_, err = w.Write(newline)
			} else {
				_, err = w.Write(newline[1:])
			}

			if err != nil {
				return
			}

			count = 0
		}

		if le {
			word = uint16(p[n]) | uint16(p[n+1])<<8
		} else {
			word = uint16(p[n])<<8 | uint16(p[n+1])
		}

		if _, err = fmt.Fprintf(w, "0x%04x", word); err != nil {
			return
		}

		if count < WordsPerLine-1 && n < len(p)-2 {
			if _, err = w.Write(spacing); err != nil {
				return
			}
		}

		count++
	}

	return
}
