// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package util

import (
	"fmt"
	"github.com/jteeuwen/dcpu/parser"
	"io"
)

var (
	comma   = []byte{','}
	space   = []byte{' '}
	newline = []byte{'\n'}
	lbrack  = []byte{'['}
	rbrack  = []byte{']'}
)

// SourceWriter allows us to write an AST out as source code
// with configurable style.
type SourceWriter struct {
	w             io.Writer
	a             *parser.AST
	indentTokens  []byte
	nestLevel     int
	TabWidth      uint
	Tabs          bool
	Comments      bool
	Indent        bool
	inInstruction bool
	skipLine      bool
	hasComment    bool
}

// NewSourceWriter creates a new source writer for the given ast
// and output stream.
func NewSourceWriter(w io.Writer, a *parser.AST) *SourceWriter {
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

	if sw.Indent {
		if sw.Tabs {
			sw.indentTokens = []byte{'\t'}
		} else {
			sw.indentTokens = make([]byte, sw.TabWidth)
			for i := range sw.indentTokens {
				sw.indentTokens[i] = ' '
			}
		}
	}

	sw.nestLevel = 1
	sw.inInstruction = false
	sw.skipLine = false
	sw.hasComment = false
	sw.writeList(sw.a.Root.Children())
}

func (sw *SourceWriter) writeList(list []parser.Node) {
	for i := range list {
		sw.writeNode(i, list[i])
	}
}

func (sw *SourceWriter) writeNode(i int, n parser.Node) {
	switch tt := n.(type) {
	case *parser.Block:
		sw.writeBlock(i, tt)
	case *parser.Expression:
		sw.writeExpression(i, tt)
	case *parser.Instruction:
		sw.writeInstruction(i, tt)
	case *parser.Function:
		sw.writeFunction(i, tt)
	case *parser.Comment:
		sw.writeComment(i, tt.Data)
	case *parser.Label:
		sw.writeLabel(i, tt.Data)
	case *parser.Name:
		sw.writeLiteral(i, tt.Data)
	case *parser.Number:
		sw.writeLiteral(i, tt.Data)
	case *parser.Char:
		sw.writeChar(i, tt.Data)
	case *parser.Operator:
		sw.writeLiteral(i, tt.Data)
	case *parser.String:
		sw.writeString(i, tt.Data)
	}
}

func (sw *SourceWriter) writeBlock(i int, n *parser.Block) {
	sw.w.Write(lbrack)

	for i, v := range n.Children() {
		sw.writeNode(i, v)
	}

	sw.w.Write(rbrack)
}

func (sw *SourceWriter) writeInstruction(i int, n *parser.Instruction) {
	sw.inInstruction = true

	chld := n.Children()
	name := chld[0].(*parser.Name)
	doSkip := name.Data == "equ"

	if sw.skipLine {
		if !doSkip {
			sw.w.Write(newline)
		}
		sw.skipLine = false
	}

	for i := 0; i < sw.nestLevel; i++ {
		sw.w.Write(sw.indentTokens)
	}

	if parser.IsBranch(name.Data) {
		sw.nestLevel++
	} else {
		if sw.nestLevel > 1 {
			doSkip = true
		}

		sw.nestLevel = 1
	}

	for i, v := range chld {
		sw.writeNode(i, v)

		if i < len(chld)-1 {
			if i > 0 {
				sw.w.Write(comma)
			}
			sw.w.Write(space)
		}
	}

	sw.w.Write(newline)
	sw.inInstruction = false
	sw.skipLine = doSkip
}

func (sw *SourceWriter) writeFunction(i int, n *parser.Function) {
	var name string

	if sw.skipLine {
		sw.w.Write(newline)
		sw.skipLine = false
	}

	chld := n.Children()

	switch tt := chld[0].(type) {
	case *parser.Name:
		name = tt.Data
	case *parser.Label:
		name = tt.Data
	}

	fmt.Fprintf(sw.w, "def %s\n", name)

	sw.writeList(chld[1:])

	sw.w.Write([]byte("end\n"))
	sw.skipLine = true
}

func (sw *SourceWriter) writeExpression(i int, n *parser.Expression) {
	for i, v := range n.Children() {
		switch v.(type) {
		case *parser.Comment:
			sw.w.Write(space)
		}

		sw.writeNode(i, v)
	}
}

func (sw *SourceWriter) writeLabel(i int, s string) {
	if !sw.hasComment {
		sw.w.Write(newline)
		sw.skipLine = false
	} else {
		sw.hasComment = false
	}

	fmt.Fprintf(sw.w, ":%s\n", s)
}

func (sw *SourceWriter) writeString(i int, s string) {
	fmt.Fprintf(sw.w, "%q", s)
}

func (sw *SourceWriter) writeChar(i int, s string) {
	fmt.Fprintf(sw.w, "'%s'", s)
}

func (sw *SourceWriter) writeLiteral(i int, s string) {
	fmt.Fprintf(sw.w, "%s", s)
}

func (sw *SourceWriter) writeComment(i int, s string) {
	if !sw.Comments {
		return
	}

	if !sw.inInstruction && i > 0 && sw.skipLine {
		sw.w.Write(newline)
		sw.skipLine = false
	}

	fmt.Fprintf(sw.w, ";%s", s)

	if !sw.inInstruction {
		sw.w.Write(newline)
		sw.hasComment = true
	}

	sw.hasComment = !sw.inInstruction
}
