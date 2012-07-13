// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"fmt"
	"io"
)

var (
	writeComments bool
	comma         = []byte{','}
	space         = []byte{' '}
	newline       = []byte{'\n'}
	lbrack        = []byte{'['}
	rbrack        = []byte{']'}
)

// WriteSource writes the given AST out as assembly source code.
func WriteSource(w io.Writer, a *AST, tabs bool, tabwidth uint, comments bool) {
	if a.Root == nil || len(a.Root.Children()) == 0 {
		return
	}

	writeComments = comments

	var indent []byte

	if tabs {
		indent = []byte{'\t'}
	} else {
		indent = make([]byte, tabwidth)
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
		case *Comment:
			if i > 0 && !followsComment {
				w.Write(newline)
			}

		case *Label:
			if i > 0 && !followsBranch && followsInstruction && !followsComment {
				w.Write(newline)
			}

		case *Instruction:
			for i := 0; i < nestlevel; i++ {
				w.Write(indent)
			}
		}

		writeNode(w, v)

		followsBranch = false
		followsComment = false
		followsInstruction = false

		switch tt := v.(type) {
		case *Comment:
			w.Write(newline)
			followsComment = true

		case *Label:
			w.Write(newline)

		case *Instruction:
			w.Write(newline)
			followsInstruction = true

			name := tt.Children()[0].(*Name).Data
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

func writeNode(w io.Writer, n Node) {
	switch tt := n.(type) {
	case *Block:
		writeBlock(w, tt)
	case *Expression:
		writeExpression(w, tt)
	case *Instruction:
		writeInstruction(w, tt)
	case *Comment:
		writeComment(w, tt.Data)
	case *Label:
		writeLabel(w, tt.Data)
	case *Name:
		writeLiteral(w, tt.Data)
	case *Number:
		writeLiteral(w, tt.Data)
	case *Char:
		writeLiteral(w, tt.Data)
	case *Operator:
		writeLiteral(w, tt.Data)
	case *String:
		writeString(w, tt.Data)
	}
}

func writeBlock(w io.Writer, n *Block) {
	w.Write(lbrack)

	for _, v := range n.Children() {
		writeNode(w, v)
	}

	w.Write(rbrack)
}

func writeInstruction(w io.Writer, n *Instruction) {
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

func writeExpression(w io.Writer, n *Expression) {
	for _, v := range n.Children() {
		switch v.(type) {
		case *Comment:
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
	if writeComments {
		fmt.Fprintf(w, ";%s", s)
	}
}
