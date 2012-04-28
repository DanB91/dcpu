// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

var (
	escA  = [][]byte{{'\\', 'a'}, {'\a'}}
	escF  = [][]byte{{'\\', 'f'}, {'\f'}}
	escN  = [][]byte{{'\\', 'n'}, {'\n'}}
	escR  = [][]byte{{'\\', 'r'}, {'\r'}}
	escT  = [][]byte{{'\\', 't'}, {'\t'}}
	escV  = [][]byte{{'\\', 'v'}, {'\v'}}
	escSQ = [][]byte{{'\\', '\''}, {'\''}}
	escDQ = [][]byte{{'\\', '"'}, {'"'}}
	escBS = [][]byte{{'\\', '\\'}, {'\\'}}
)

// escape validates character escape sequences in the given input string
// and returns the result with all sequences processed.
func escape(in []byte) (output []byte, err error) {
	in = bytes.Replace(in, escA[0], escA[1], -1)
	in = bytes.Replace(in, escF[0], escF[1], -1)
	in = bytes.Replace(in, escN[0], escN[1], -1)
	in = bytes.Replace(in, escR[0], escR[1], -1)
	in = bytes.Replace(in, escT[0], escT[1], -1)
	in = bytes.Replace(in, escV[0], escV[1], -1)
	in = bytes.Replace(in, escSQ[0], escSQ[1], -1)
	in = bytes.Replace(in, escDQ[0], escDQ[1], -1)

	var i, base int
	var repl []byte

	for {
		i = bytes.Index(in[i:], escBS[1])
		if i == -1 || len(in[i:]) == 1 {
			break
		}

		if i+1 < len(in) && in[i+1] == '\\' {
			// Escaped backslash. Turn it into a proper backslash and move on.
			in = bytes.Replace(in, in[i:i+2], escBS[1], -1)
			i++
			continue
		}

		if len(in[i:]) >= 4 {
			var offset int

			switch {
			case in[i+1] >= '0' && in[i+1] <= '7':
				base = 8
				offset = 1
			case in[i+1] == 'x':
				base = 16
				offset = 2
			default:
				goto skip
			}

			if repl, err = convertBaseEsc(in[i+offset:i+4], base); err != nil {
				return
			}

			in = bytes.Replace(in, in[i:i+4], repl, -1)
			continue
		}

	skip:
		if len(in[i:]) >= 6 && in[i+1] == 'u' {
			if repl, err = convertBaseEsc(in[i+2:i+6], 16); err != nil {
				return
			}

			in = bytes.Replace(in, in[i:i+6], repl, -1)
			continue
		}

		if len(in[i:]) >= 10 {
			var base int

			switch in[i+1] {
			case 'U':
				base = 16
			case 'b':
				base = 2
			default:
				goto end
			}

			if repl, err = convertBaseEsc(in[i+2:i+10], base); err != nil {
				return
			}

			in = bytes.Replace(in, in[i:i+10], repl, -1)
			continue
		}

	end:
		return nil, errors.New(fmt.Sprintf("Invalid escape sequence: %.10q", in[i:]))
	}

	return in, nil
}

func convertBaseEsc(in []byte, base int) (out []byte, err error) {
	var n uint64
	if n, err = strconv.ParseUint(string(in), base, 64); err == nil {
		out = []byte(string(n))
	}
	return
}
