// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	dp "github.com/jteeuwen/dcpu/parser"
	"io"
)

const PadString = "  "

// writeAst writes the given AST out in human-readable form.
func writeAst(w io.Writer, a *dp.AST) {
	fmt.Fprintf(w, "Files:\n")

	for i := range a.Files {
		fmt.Fprintf(w, "  %d - %s\n", i, a.Files[i])
	}

	fmt.Fprintf(w, "Nodes:\n")
	writeAstNode(w, a.Root, "")
}

func writeAstNode(w io.Writer, n dp.Node, pad string) {
	switch tt := n.(type) {
	case dp.NodeCollection:
		writeAstCollection(w, tt, tt.Children(), pad)
	case *dp.Comment:
		writeAstString(w, tt, tt.Data, pad)
	case *dp.Label:
		writeAstString(w, tt, tt.Data, pad)
	case *dp.Name:
		writeAstString(w, tt, tt.Data, pad)
	case *dp.Number:
		writeAstNumber(w, tt, pad)
	case *dp.Operator:
		writeAstString(w, tt, tt.Data, pad)
	case *dp.String:
		writeAstString(w, tt, tt.Data, pad)
	}
}

func writeAstNodeBase(w io.Writer, n *dp.NodeBase, pad string) {
	fmt.Fprintf(w, "%s%02d:%04d:%03d", pad, n.File(), n.Line(), n.Col())
}

func writeAstCollection(w io.Writer, n dp.Node, l []dp.Node, pad string) {
	writeAstNodeBase(w, n.Base(), pad)
	fmt.Fprintf(w, " %T {\n", n)

	for _, v := range l {
		writeAstNode(w, v, pad+PadString)
	}

	fmt.Fprintf(w, "%s}\n", pad)
}

func writeAstString(w io.Writer, n dp.Node, data, pad string) {
	writeAstNodeBase(w, n.Base(), pad)

	if len(data) > 20 {
		fmt.Fprintf(w, " %T(%.20q...)\n", n, data)
	} else {
		fmt.Fprintf(w, " %T(%q)\n", n, data)
	}
}

func writeAstNumber(w io.Writer, n *dp.Number, pad string) {
	writeAstNodeBase(w, n.Base(), pad)
	fmt.Fprintf(w, " %T(0x%04x)\n", n, n.Data)
}
