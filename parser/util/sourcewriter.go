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
	w            io.Writer
	a            *parser.AST
	indentTokens []byte
	nestLevel    int
	TabWidth     uint
	Tabs         bool
	Comments     bool
	Indent       bool
	inInstr      bool
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
	sw.inInstr = false
	sw.writeList(sw.a.Root.Children())
}

func (sw *SourceWriter) indent() {
	for i := 0; i < sw.nestLevel; i++ {
		sw.w.Write(sw.indentTokens)
	}
}

func (sw *SourceWriter) writeList(list []parser.Node) {
	preceedingComment := false

	for i := range list {
		switch tt := list[i].(type) {
		case *parser.Block:
			sw.writeBlock(tt)

		case *parser.Expression:
			if i > 0 && i < len(list) {
				sw.w.Write(space)
			}

			sw.writeExpression(tt)

			if i > 0 && i < len(list)-1 {
				sw.w.Write(comma)
			}

		case *parser.Instruction:
			sw.inInstr = true
			sw.indent()
			sw.writeInstruction(tt)
			sw.w.Write(newline)
			sw.inInstr = false

		case *parser.Function:
			sw.writeFunction(tt)
			sw.w.Write(newline)

		case *parser.Comment:
			if i > 0 && !sw.inInstr && !preceedingComment {
				sw.w.Write(newline)
			}

			if sw.inInstr && i < len(list) {
				sw.w.Write(space)
			}

			sw.writeComment(tt.Data)

			if !sw.inInstr {
				sw.w.Write(newline)
			}

		case *parser.Label:
			if i > 0 && !preceedingComment {
				sw.w.Write(newline)
			}

			sw.writeLabel(tt.Data)
			sw.w.Write(newline)

		case *parser.Name:
			sw.writeLiteral(tt.Data)

		case *parser.Number:
			sw.writeLiteral(tt.Data)

		case *parser.Char:
			sw.writeChar(tt.Data)

		case *parser.Operator:
			sw.writeLiteral(tt.Data)

		case *parser.String:
			sw.writeString(tt.Data)
		}

		switch list[i].(type) {
		case *parser.Comment:
			preceedingComment = true
		case *parser.Label:
			preceedingComment = true
		default:
			preceedingComment = false
		}
	}
}

func (sw *SourceWriter) writeBlock(n *parser.Block) {
	sw.w.Write(lbrack)
	sw.writeList(n.Children())
	sw.w.Write(rbrack)
}

func (sw *SourceWriter) writeInstruction(n *parser.Instruction) {
	chld := n.Children()
	name := chld[0].(*parser.Name)

	if parser.IsBranch(name.Data) {
		sw.nestLevel++
	} else {
		sw.nestLevel = 1
	}

	sw.writeList(chld)
}

func (sw *SourceWriter) writeFunction(n *parser.Function) {
	var name string
	chld := n.Children()

	switch tt := chld[0].(type) {
	case *parser.Name:
		name = tt.Data
	case *parser.Label:
		name = tt.Data
	}

	fmt.Fprintf(sw.w, "def %s\n", name)
	sw.writeList(chld[1:])
	sw.w.Write([]byte("end"))
}

func (sw *SourceWriter) writeExpression(n *parser.Expression) {
	sw.writeList(n.Children())
}

func (sw *SourceWriter) writeLabel(s string) {
	fmt.Fprintf(sw.w, ":%s", s)
}

func (sw *SourceWriter) writeString(s string) {
	fmt.Fprintf(sw.w, "%q", s)
}

func (sw *SourceWriter) writeChar(s string) {
	fmt.Fprintf(sw.w, "'%s'", s)
}

func (sw *SourceWriter) writeLiteral(s string) {
	fmt.Fprintf(sw.w, "%s", s)
}

func (sw *SourceWriter) writeComment(s string) {
	if sw.Comments {
		fmt.Fprintf(sw.w, ";%s", s)
	}
}
