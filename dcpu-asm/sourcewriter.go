// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"github.com/jteeuwen/dcpu/parser"
	"io"
	"os"
)

func writeSource(ast *parser.AST, file string, doast bool) (err error) {
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

	if doast {
		parser.WriteAst(w, ast)
	} else {
		parser.NewSourceWriter(w, ast).Write()
	}

	return
}
