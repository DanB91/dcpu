// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"fmt"
	"io"
)

var (
	comma   = []byte{','}
	space   = []byte{' '}
	newline = []byte{'\n'}
	lbrack  = []byte{'['}
	rbrack  = []byte{']'}
)

func isBranch(name string) bool {
	switch name {
	case "ifb", "ifc", "ife", "ifn", "ifg", "ifa", "ifl", "ifu":
		return true
	}
	return false
}

// SourceWriter allows us to write an AST out as source code
// with configurable style.
type SourceWriter struct {
	w        io.Writer
	a        *AST
	TabWidth uint
	Tabs     bool
	Comments bool
	Indent   bool
}

// NewSourceWriter creates a new source writer for the given ast
// and output stream.
func NewSourceWriter(w io.Writer, a *AST) *SourceWriter {
	s := new(SourceWriter)
	s.w = w
	s.a = a
	s.Indent = true
	s.Tabs = false
	s.TabWidth = 3
	s.Comments = true
	return s
}

// WriteSource writes the given AST out as assembly source code.
func (sw *SourceWriter) Write() {
	if sw.a.Root == nil || len(sw.a.Root.Children()) == 0 {
		return
	}

	var indent []byte

	if sw.Indent {
		if sw.Tabs {
			indent = []byte{'\t'}
		} else {
			indent = make([]byte, sw.TabWidth)
			for i := range indent {
				indent[i] = ' '
			}
		}
	}

	followsBranch := false
	followsInstruction := false
	followsComment := false
	nestlevel := 1

	for i, v := range sw.a.Root.Children() {
		switch v.(type) {
		case *Comment:
			if i > 0 && !followsBranch && followsInstruction && !followsComment {
				sw.w.Write(newline)
			}

		case *Label:
			if i > 0 && !followsBranch && followsInstruction && !followsComment {
				sw.w.Write(newline)
			}

		case *Instruction:
			for i := 0; i < nestlevel; i++ {
				sw.w.Write(indent)
			}
		}

		sw.writeNode(v)

		followsBranch = false
		followsComment = false
		followsInstruction = false

		switch tt := v.(type) {
		case *Comment:
			sw.w.Write(newline)
			followsComment = true

		case *Label:
			sw.w.Write(newline)

		case *Instruction:
			sw.w.Write(newline)
			followsInstruction = true

			name := tt.Children()[0].(*Name).Data
			if isBranch(name) {
				nestlevel++

			} else if nestlevel > 1 {
				sw.w.Write(newline)
				nestlevel = 1
				followsBranch = true
			}
		}
	}
}

func (sw *SourceWriter) writeNode(n Node) {
	switch tt := n.(type) {
	case *Block:
		sw.writeBlock(tt)
	case *Expression:
		sw.writeExpression(tt)
	case *Instruction:
		sw.writeInstruction(tt)
	case *Comment:
		sw.writeComment(tt.Data)
	case *Label:
		sw.writeLabel(tt.Data)
	case *Name:
		sw.writeLiteral(tt.Data)
	case *Number:
		sw.writeLiteral(tt.Data)
	case *Char:
		sw.writeLiteral(tt.Data)
	case *Operator:
		sw.writeLiteral(tt.Data)
	case *String:
		sw.writeString(tt.Data)
	}
}

func (sw *SourceWriter) writeBlock(n *Block) {
	sw.w.Write(lbrack)

	for _, v := range n.Children() {
		sw.writeNode(v)
	}

	sw.w.Write(rbrack)
}

func (sw *SourceWriter) writeInstruction(n *Instruction) {
	chld := n.Children()
	for i, v := range chld {
		sw.writeNode(v)

		if i < len(chld)-1 {
			if i > 0 {
				sw.w.Write(comma)
			}
			sw.w.Write(space)
		}
	}
}

func (sw *SourceWriter) writeExpression(n *Expression) {
	for _, v := range n.Children() {
		switch v.(type) {
		case *Comment:
			sw.w.Write(space)
		}

		sw.writeNode(v)
	}
}

func (sw *SourceWriter) writeLabel(s string) {
	fmt.Fprintf(sw.w, ":%s", s)
}

func (sw *SourceWriter) writeString(s string) {
	fmt.Fprintf(sw.w, "%q", s)
}

func (sw *SourceWriter) writeLiteral(s string) {
	fmt.Fprintf(sw.w, "%s", s)
}

func (sw *SourceWriter) writeComment(s string) {
	if sw.Comments {
		fmt.Fprintf(sw.w, ";%s", s)
	}
}
