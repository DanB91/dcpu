// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
)

// readSource takes the input files and parses their contents into
// the given AST.
func readSource(ast *AST, c *Config) (err error) {
	var fd *os.File

	for i := range c.Input {
		if fd, err = os.Open(c.Input[i]); err != nil {
			return
		}

		err = ast.Parse(fd, c.Input[i])
		fd.Close()

		if err != nil {
			return
		}
	}

	return resolveIncludes(ast, c)
}

// resolveIncludes finds references to undefined labels.
// It then tries to find the code for these labels in the supplied
// include paths. Files should be defined as '<labelname>.dasm'.
func resolveIncludes(ast *AST, c *Config) (err error) {
	var labels []*Label
	var refs []*Name

	findLabels(ast.Root.Children, &labels)
	findRefs(ast.Root.Children, &refs)

	println(len(labels), len(refs))
	return
}

// findLabels recursively finds Label nodes.
func findLabels(n []Node, l *[]*Label) {
	for i := range n {
		switch tt := n[i].(type) {
		case *Expression:
			findLabels(tt.Children, l)

		case *Block:
			findLabels(tt.Children, l)

		case *Instruction:
			findLabels(tt.Children, l)

		case *Label:
			*l = append(*l, tt)
		}
	}
}

// findRefs recursively finds Label references.
func findRefs(n []Node, l *[]*Name) {
	for i := range n {
		switch tt := n[i].(type) {
		case *Expression:
			findRefs(tt.Children, l)

		case *Block:
			findRefs(tt.Children, l)

		case *Instruction:
			findRefs(tt.Children, l)

		case *Name:
			if isRegister(tt.Data) || isInstruction(tt.Data) {
				continue
			}

			*l = append(*l, tt)
		}
	}
}
