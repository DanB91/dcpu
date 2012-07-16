// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import "github.com/jteeuwen/dcpu/cpu"

// Represents a single AST node.
type Node interface {
	// Index of the original source file name.
	File() int

	// Line number in the original source file.
	Line() int

	// Column number in the original source file.
	Col() int

	// Yields a pointer to the embedded NodeBase struct.
	Base() *NodeBase

	// Copy creates a deep copy of this node.
	//
	// The file, line and column information is updated for all
	// newly created nodes.
	Copy(file, line, col int) Node
}

// Represents a node that has child nodes.
type NodeCollection interface {
	Node
	Children() []Node
	SetChildren([]Node)
}

type NumericNode interface {
	Node
	Parse() (cpu.Word, error)
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

func (n *NodeBase) File() int       { return n.file }
func (n *NodeBase) Line() int       { return n.line }
func (n *NodeBase) Col() int        { return n.col }
func (n *NodeBase) Base() *NodeBase { return n }
