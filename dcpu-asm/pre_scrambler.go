// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/parser"
	"github.com/jteeuwen/dcpu/parser/util"
)

func init() {
	RegisterPreProcessor("scramble",
		"Obfuscate label names and label references.", NewScrambler, false)
}

// Scrambler obfuscates labels and label references in the given AST.
type Scrambler struct{}

func NewScrambler() PreProcessor { return new(Scrambler) }

func (p *Scrambler) Process(ast *parser.AST) (err error) {
	var labels []*parser.Label
	var refs []*parser.Name
	var i, j int

	util.FindLabels(ast.Root.Children(), &labels)
	util.FindReferences(ast.Root.Children(), &refs)

	for i = range labels {
		old := labels[i].Data
		labels[i].Data = fmt.Sprintf("l%x", i)

		for j = range refs {
			if refs[j].Data == old {
				refs[j].Data = labels[i].Data
			}
		}
	}

	return
}
