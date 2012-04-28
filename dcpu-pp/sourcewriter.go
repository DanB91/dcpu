// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
)

// writeSource writes the given AST out as assembly source to the
// supplied writer.
func writeSource(w io.WriteCloser, ast *AST) (err error) {
	defer w.Close()
	return
}
