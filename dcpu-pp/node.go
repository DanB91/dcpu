// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

// Represents a single AST node.
type Node interface {
	File() int
	Line() int
	Col() int

	// Used for debug purposes
	Dump(pad string) string
}

// Base-type for all nodes. This takes care of some common
// aspects they all need in order to qualify as a Node interface.
type NodeBase struct {
	file int
	line int
	col  int // Location of this node in source code.
}

func NewNodeBase(file, line, col int) *NodeBase {
	return &NodeBase{file, line, col}
}

func (n *NodeBase) File() int { return n.file }
func (n *NodeBase) Line() int { return n.line }
func (n *NodeBase) Col() int  { return n.col }

func (n *NodeBase) Dump(pad string) string {
	return fmt.Sprintf("%s%02d:%04d:%03d", pad, n.file, n.line, n.col)
}
