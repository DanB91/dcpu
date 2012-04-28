// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

func init() {
	Register("scramble", "Obfuscate label names and label references.", NewScrambler)
}

// Scrambler obfuscates Label and Label References in the given AST.
type Scrambler struct{}

func NewScrambler() Processor { return new(Scrambler) }

func (p *Scrambler) Process(ast *AST) (err error) {
	return
}
