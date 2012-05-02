// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	dp "github.com/jteeuwen/dcpu/parser"
	"io"
)

var (
	comma    = []byte{','}
	space    = []byte{' '}
	dblspace = []byte{' ', ' '}
	newline  = []byte{'\n'}
	lbrack   = []byte{'['}
	rbrack   = []byte{']'}
)

// writeSource writes the given AST out as assembly source code.
func writeSource(w io.Writer, a *dp.AST) {
	for _, v := range a.Root.Children() {
		writeSourceNode(w, v)
	}
}

func writeSourceNode(w io.Writer, n dp.Node) {
	switch tt := n.(type) {
	case *dp.Block:
		writeSourceBlock(w, tt)
	case *dp.Comment:
		writeSourceComment(w, tt)
	case *dp.Expression:
		writeSourceExpression(w, tt)
	case *dp.Instruction:
		writeSourceInstruction(w, tt)
	case *dp.Label:
		writeSourceLabel(w, tt)
	case *dp.Name:
		writeSourceLiteral(w, tt.Data)
	case *dp.Number:
		writeSourceNumber(w, tt.Data)
	case *dp.Operator:
		writeSourceLiteral(w, tt.Data)
	case *dp.String:
		writeSourceString(w, tt)
	}
}

func writeSourceBlock(w io.Writer, n *dp.Block) {
	w.Write(lbrack)

	for _, v := range n.Children() {
		writeSourceNode(w, v)
	}

	w.Write(rbrack)
}

func writeSourceInstruction(w io.Writer, n *dp.Instruction) {
	w.Write(dblspace)

	chld := n.Children()
	for i, v := range chld {
		writeSourceNode(w, v)

		if i < len(chld)-1 {
			if i > 0 {
				w.Write(comma)
			}
			w.Write(space)
		}
	}

	w.Write(newline)
}

func writeSourceExpression(w io.Writer, n *dp.Expression) {
	for _, v := range n.Children() {
		writeSourceNode(w, v)
	}
}

func writeSourceComment(w io.Writer, n *dp.Comment) {
	fmt.Fprintf(w, ";%s\n", n.Data)
}

func writeSourceLabel(w io.Writer, n *dp.Label) {
	fmt.Fprintf(w, ":%s\n", n.Data)
}

func writeSourceString(w io.Writer, n *dp.String) {
	fmt.Fprintf(w, "%q", n.Data)
}

func writeSourceLiteral(w io.Writer, s string) {
	fmt.Fprintf(w, "%s", s)
}

func writeSourceNumber(w io.Writer, n dp.Word) {
	fmt.Fprintf(w, "0x%x", n)
}
