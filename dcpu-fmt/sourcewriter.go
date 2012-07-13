// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	dp "github.com/jteeuwen/dcpu/parser"
	"io"
)

var (
	comma   = []byte{','}
	space   = []byte{' '}
	newline = []byte{'\n'}
	lbrack  = []byte{'['}
	rbrack  = []byte{']'}
)

// WriteSource writes the given AST out as assembly source code.
func WriteSource(w io.Writer, a *dp.AST) {
	if a.Root == nil || len(a.Root.Children()) == 0 {
		return
	}

	var indent []byte

	if *tabs {
		indent = []byte{'\t'}
	} else {
		indent = make([]byte, *tabwidth)
		for i := range indent {
			indent[i] = ' '
		}
	}

	followsBranch := false
	followsInstruction := false
	followsComment := false
	nestlevel := 1

	for i, v := range a.Root.Children() {
		switch v.(type) {
		case *dp.Comment:
			if i > 0 && !followsComment {
				w.Write(newline)
			}

		case *dp.Label:
			if i > 0 && !followsBranch && followsInstruction && !followsComment {
				w.Write(newline)
			}

		case *dp.Instruction:
			for i := 0; i < nestlevel; i++ {
				w.Write(indent)
			}
		}

		writeNode(w, v)

		followsBranch = false
		followsComment = false
		followsInstruction = false

		switch tt := v.(type) {
		case *dp.Comment:
			w.Write(newline)
			followsComment = true

		case *dp.Label:
			w.Write(newline)

		case *dp.Instruction:
			w.Write(newline)
			followsInstruction = true

			name := tt.Children()[0].(*dp.Name).Data
			if isBranch(name) {
				nestlevel++

			} else if nestlevel > 1 {
				w.Write(newline)
				nestlevel = 1
				followsBranch = true
			}
		}
	}
}

func isBranch(name string) bool {
	switch name {
	case "ifb", "ifc", "ife", "ifn", "ifg", "ifa", "ifl", "ifu":
		return true
	}
	return false
}

func writeNode(w io.Writer, n dp.Node) {
	switch tt := n.(type) {
	case *dp.Block:
		writeBlock(w, tt)
	case *dp.Expression:
		writeExpression(w, tt)
	case *dp.Instruction:
		writeInstruction(w, tt)
	case *dp.Comment:
		writeComment(w, tt.Data)
	case *dp.Label:
		writeLabel(w, tt.Data)
	case *dp.Name:
		writeLiteral(w, tt.Data)
	case *dp.Number:
		writeLiteral(w, tt.Data)
	case *dp.Char:
		writeLiteral(w, tt.Data)
	case *dp.Operator:
		writeLiteral(w, tt.Data)
	case *dp.String:
		writeString(w, tt.Data)
	}
}

func writeBlock(w io.Writer, n *dp.Block) {
	w.Write(lbrack)

	for _, v := range n.Children() {
		writeNode(w, v)
	}

	w.Write(rbrack)
}

func writeInstruction(w io.Writer, n *dp.Instruction) {
	chld := n.Children()
	for i, v := range chld {
		writeNode(w, v)

		if i < len(chld)-1 {
			if i > 0 {
				w.Write(comma)
			}
			w.Write(space)
		}
	}
}

func writeExpression(w io.Writer, n *dp.Expression) {
	for _, v := range n.Children() {
		switch v.(type) {
		case *dp.Comment:
			w.Write(space)
		}

		writeNode(w, v)
	}
}

func writeLabel(w io.Writer, s string) {
	fmt.Fprintf(w, ":%s", s)
}

func writeString(w io.Writer, s string) {
	fmt.Fprintf(w, "%q", s)
}

func writeLiteral(w io.Writer, s string) {
	fmt.Fprintf(w, "%s", s)
}

func writeComment(w io.Writer, s string) {
	if !*strip {
		fmt.Fprintf(w, ";%s", s)
	}
}
