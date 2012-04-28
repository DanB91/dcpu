// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

func init() {
	Register("scramble", "Obfuscate label names and label references.", NewScrambler)
}

// Scrambler obfuscates Label and Label References in the given AST.
type Scrambler struct{}

func NewScrambler() Processor { return new(Scrambler) }

func (p *Scrambler) Process(ast *AST) (err error) {
	var labels []*Label
	var refs []*Name
	var i, j int

	findLabels(ast.Root.Children, &labels)
	findReferences(ast.Root.Children, &refs)

	for i = range labels {
		old := labels[i].Data
		labels[i].Data = fmt.Sprintf("l%04x", i)

		for j = range refs {
			if refs[j].Data == old {
				refs[j].Data = labels[i].Data
			}
		}
	}

	return
}
