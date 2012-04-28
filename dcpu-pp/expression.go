// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

// An Expression is a collection of AST nodes.
type Expression struct {
	*NodeBase
	Children []Node
}

func NewExpression(file, line, col int) *Expression {
	return &Expression{
		NewNodeBase(file, line, col),
		nil,
	}
}
