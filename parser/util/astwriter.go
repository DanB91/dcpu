// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package util

import (
	"fmt"
	"github.com/jteeuwen/dcpu/parser"
	"io"
)

// AstWriter outputs a human-readable version of the input AST.
type AstWriter struct {
	a      *parser.AST
	w      io.Writer
	Indent string
}

// NewAstWriter creates a new AstWriter for the given output stream and ast.
func NewAstWriter(w io.Writer, a *parser.AST) *AstWriter {
	aw := new(AstWriter)
	aw.a = a
	aw.w = w
	aw.Indent = "  "
	return aw
}

// Write writes the out in human-readable form.
func (aw *AstWriter) Write() {
	fmt.Fprintf(aw.w, "Files:\n")

	for i := range aw.a.Files {
		fmt.Fprintf(aw.w, "  %d - %s\n", i, aw.a.Files[i])
	}

	fmt.Fprintf(aw.w, "Nodes:\n")
	aw.writeNode(aw.a.Root, "")
}

func (aw *AstWriter) writeNode(n parser.Node, pad string) {
	switch tt := n.(type) {
	case parser.NodeCollection:
		aw.writeCollection(tt, tt.Children(), pad)
	case *parser.Comment:
		aw.writeString(tt, tt.Data, pad)
	case *parser.Label:
		aw.writeString(tt, tt.Data, pad)
	case *parser.Name:
		aw.writeString(tt, tt.Data, pad)
	case *parser.Number:
		aw.writeString(tt, tt.Data, pad)
	case *parser.Char:
		aw.writeString(tt, tt.Data, pad)
	case *parser.Operator:
		aw.writeString(tt, tt.Data, pad)
	case *parser.String:
		aw.writeString(tt, tt.Data, pad)
	}
}

func (aw *AstWriter) writeNodeBase(n *parser.NodeBase, pad string) {
	fmt.Fprintf(aw.w, "%s%02d:%04d:%03d", pad, n.File(), n.Line(), n.Col())
}

func (aw *AstWriter) writeCollection(n parser.Node, l []parser.Node, pad string) {
	aw.writeNodeBase(n.Base(), pad)
	fmt.Fprintf(aw.w, " %T {\n", n)

	for _, v := range l {
		aw.writeNode(v, pad+aw.Indent)
	}

	fmt.Fprintf(aw.w, "%s}\n", pad)
}

func (aw *AstWriter) writeString(n parser.Node, data, pad string) {
	aw.writeNodeBase(n.Base(), pad)

	if len(data) > 20 {
		fmt.Fprintf(aw.w, " %T(%.20q...)\n", n, data)
	} else {
		fmt.Fprintf(aw.w, " %T(%q)\n", n, data)
	}
}

func (aw *AstWriter) writeNumber(n *parser.Number, pad string) {
	aw.writeNodeBase(n.Base(), pad)
	fmt.Fprintf(aw.w, " %T(0x%04x)\n", n, n.Data)
}
