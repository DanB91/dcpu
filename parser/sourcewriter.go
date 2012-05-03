// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"fmt"
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

// WriteSource writes the given AST out as assembly source code.
func WriteSource(w io.Writer, a *AST) {
	for _, v := range a.Root.Children() {
		writeSourceNode(w, v)
	}
}

func writeSourceNode(w io.Writer, n Node) {
	switch tt := n.(type) {
	case *Block:
		writeSourceBlock(w, tt)
	case *Comment:
		writeSourceComment(w, tt)
	case *Expression:
		writeSourceExpression(w, tt)
	case *Instruction:
		writeSourceInstruction(w, tt)
	case *Label:
		writeSourceLabel(w, tt)
	case *Name:
		writeSourceLiteral(w, tt.Data)
	case *Number:
		writeSourceNumber(w, tt.Data)
	case *Operator:
		writeSourceLiteral(w, tt.Data)
	case *String:
		writeSourceString(w, tt)
	}
}

func writeSourceBlock(w io.Writer, n *Block) {
	w.Write(lbrack)

	for _, v := range n.Children() {
		writeSourceNode(w, v)
	}

	w.Write(rbrack)
}

func writeSourceInstruction(w io.Writer, n *Instruction) {
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

func writeSourceExpression(w io.Writer, n *Expression) {
	for _, v := range n.Children() {
		writeSourceNode(w, v)
	}
}

func writeSourceComment(w io.Writer, n *Comment) {
	fmt.Fprintf(w, ";%s\n", n.Data)
}

func writeSourceLabel(w io.Writer, n *Label) {
	fmt.Fprintf(w, ":%s\n", n.Data)
}

func writeSourceString(w io.Writer, n *String) {
	fmt.Fprintf(w, "%q", n.Data)
}

func writeSourceLiteral(w io.Writer, s string) {
	fmt.Fprintf(w, "%s", s)
}

func writeSourceNumber(w io.Writer, n Word) {
	fmt.Fprintf(w, "0x%x", n)
}
