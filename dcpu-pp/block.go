// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"fmt"
)

// An Block is a collection of AST nodes.
type Block struct {
	*NodeBase
	Children []Node
}

func NewBlock(file, line, col int) *Block {
	return &Block{
		NewNodeBase(file, line, col),
		nil,
	}
}

func (n *Block) Dump(pad string) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s %T {\n", n.NodeBase.Dump(pad), n)

	for _, v := range n.Children {
		b.WriteString(v.Dump(pad + "  "))
	}

	fmt.Fprintf(&b, "%s}\n", pad)
	return b.String()
}
