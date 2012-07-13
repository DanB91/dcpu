// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "github.com/jteeuwen/dcpu/parser"

// StripComments removes all code comment nodes from the AST.
func StripComments(a *parser.AST) {
	a.Root.SetChildren(stripList(a.Root.Children()))
}

func stripList(n []parser.Node) []parser.Node {
	var list []parser.Node
	var ok bool

	for _, v := range n {
		if _, ok = v.(*parser.Comment); ok {
			continue
		}

		switch tt := v.(type) {
		case parser.NodeCollection:
			tt.SetChildren(stripList(tt.Children()))
		}

		list = append(list, v)
	}

	return list
}
