// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"io"
)

const PadString = "  "

// writeAst writes the given AST out in human-readable form.
func writeAst(w io.Writer, a *AST) {
	fmt.Fprintf(w, "Files:\n")

	for i := range a.Files {
		fmt.Fprintf(w, "  %d - %s\n", i, a.Files[i])
	}

	fmt.Fprintf(w, "Nodes:\n")
	writeAstNode(w, a.Root, "")
}

func writeAstNode(w io.Writer, n Node, pad string) {
	switch tt := n.(type) {
	case *Block:
		writeAstCollection(w, tt, tt.Children, pad)
	case *Comment:
		writeAstString(w, tt, tt.Data, pad)
	case *Expression:
		writeAstCollection(w, tt, tt.Children, pad)
	case *Instruction:
		writeAstCollection(w, tt, tt.Children, pad)
	case *Label:
		writeAstString(w, tt, tt.Data, pad)
	case *Name:
		writeAstString(w, tt, tt.Data, pad)
	case *Number:
		writeAstNumber(w, tt, pad)
	case *Operator:
		writeAstString(w, tt, tt.Data, pad)
	case *String:
		writeAstString(w, tt, tt.Data, pad)
	}
}

func writeAstNodeBase(w io.Writer, n *NodeBase, pad string) {
	fmt.Fprintf(w, "%s%02d:%04d:%03d", pad, n.file, n.line, n.col)
}

func writeAstCollection(w io.Writer, n Node, l []Node, pad string) {
	writeAstNodeBase(w, n.Base(), pad)
	fmt.Fprintf(w, " %T {\n", n)

	for _, v := range l {
		writeAstNode(w, v, pad+PadString)
	}

	fmt.Fprintf(w, "%s}\n", pad)
}

func writeAstString(w io.Writer, n Node, data, pad string) {
	writeAstNodeBase(w, n.Base(), pad)

	if len(data) > 20 {
		fmt.Fprintf(w, " %T(%.20q...)\n", n, data)
	} else {
		fmt.Fprintf(w, " %T(%q)\n", n, data)
	}
}

func writeAstNumber(w io.Writer, n *Number, pad string) {
	writeAstNodeBase(w, n.Base(), pad)
	fmt.Fprintf(w, " %T(0x%04x)\n", n, n.Data)
}
