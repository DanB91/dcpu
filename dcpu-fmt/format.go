// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"github.com/jteeuwen/dcpu/parser"
	"io"
	"os"
)

const (
	Stdin  = "<stdin>"
	Stdout = "<stdout>"
)

var (
	_stdin  = &NopCloseReader{os.Stdin}
	_stdout = &NopCloseWriter{os.Stdout}
)

func Format(in, out string) (err error) {
	var fout io.WriteCloser
	var fin io.ReadCloser
	var ast parser.AST

	// Open input stream.
	if in == Stdin {
		fin = _stdin
	} else {
		fin, err = os.Open(in)
		if err != nil {
			return
		}
	}

	// Parse source into AST.
	err = ast.Parse(fin, in)
	fin.Close()

	if err != nil {
		return
	}

	if *strip {
		StripComments(&ast)
	}

	// Open output stream.
	if out == Stdout {
		fout = _stdout
	} else {
		fout, err = os.Create(out)
		if err != nil {
			return
		}

		defer fout.Close()
	}

	// Write source.
	sw := parser.NewSourceWriter(fout, &ast)
	sw.Tabs = *tabs
	sw.TabWidth = *tabwidth
	sw.Comments = !*strip
	sw.Write()
	return
}
