// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
	"path"
)

// parseInput takes the input files and parses their contents into
// the given AST.
func parseInput(ast *AST, c *Config) (err error) {
	for i := range c.Input {
		if err = readSource(ast, c.Input[i]); err != nil {
			return
		}
	}

	return resolveIncludes(ast, c)
}

// readSource reads the given file and parses its contents
// into the given AST.
func readSource(ast *AST, file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}

	defer fd.Close()
	return ast.Parse(fd, file)
}

// resolveIncludes finds references to undefined labels.
// It then tries to find the code for these labels in the supplied
// include paths. Files should be defined as '<labelname>.dasm'.
func resolveIncludes(ast *AST, c *Config) (err error) {
	var labels []*Label
	var refs []*Name

	findLabels(ast.Root.Children, &labels)
	findRefs(ast.Root.Children, &refs)
	refs = findUndefinedRefs(refs, labels)

	if len(refs) == 0 {
		// No undefined references. We're done here.
		return
	}

	if len(c.Include) == 0 {
		// We have unresolved references, but no places to look
		// for their implementation. This constitutes a booboo.
		return NewParseError(ast.Files[refs[0].File()], refs[0].Line(), refs[0].Col(),
			"Undefined reference: %q", refs[0].Data)
	}

	for i := range refs {
		if err = loadInclude(ast, c, refs[i]); err != nil {
			return
		}
	}

	return
}

// loadInclude tries to load the given reference as an include file.
// Parses it into the supplied AST and verifies that it contains what
// we are looking for.
func loadInclude(ast *AST, c *Config, r *Name) (err error) {
	var stat os.FileInfo
	var file string

	for i := range c.Include {
		file = path.Join(c.Include[i], r.Data+".dasm")
		stat, err = os.Lstat(file)

		if err != nil || stat.IsDir() {
			return NewParseError(ast.Files[r.File()], r.Line(), r.Col(),
				"Undefined reference: %q", r.Data)
		}

		if err = readSource(ast, file); err != nil {
			return
		}

		if !includeHasLabel(ast, file, r.Data) {
			return NewParseError(ast.Files[r.File()], r.Line(), r.Col(),
				"Undefined reference: %q. Include file was found, but "+
					"it did not define the desired label.", r.Data)
		}
	}

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

// findUndefinedRefs compares both given lists of labels and
// label references. Any reference that is not present in the
// label list, is considered unresolved and added to the 
// returned list.
func findUndefinedRefs(refs []*Name, labels []*Label) []*Name {
	var i, j int

outer:
	for i = range refs {
		for j = range labels {
			if labels[j].Data == refs[i].Data {
				copy(refs[i:], refs[i+1:])
				refs = refs[:len(refs)-1]
				goto outer
			}
		}
	}

	return refs
}

// includeHasLabel checks if a newly parsed include actually
// contains the label reference we are looking for.
func includeHasLabel(ast *AST, file string, target string) bool {
	index := fileIndex(ast, file)
	if index == -1 {
		return false
	}

	return hasLabel(ast.Root.Children, index, target)
}

// hasLabel recursively finds a specific label definition.
// Returns true if it was found. False otherwise.
func hasLabel(n []Node, file int, target string) bool {
	for i := range n {
		switch tt := n[i].(type) {
		case *Expression:
			if hasLabel(tt.Children, file, target) {
				return true
			}

		case *Block:
			if hasLabel(tt.Children, file, target) {
				return true
			}

		case *Instruction:
			if hasLabel(tt.Children, file, target) {
				return true
			}

		case *Label:
			if tt.File() == file && tt.Data == target {
				return true
			}

			// @_@ These are not the labels you are looking for. @_@
		}
	}

	return false
}

// fileIndex returns the given file's index as it is stored in the AST.
func fileIndex(ast *AST, file string) int {
	for i := range ast.Files {
		if ast.Files[i] == file {
			return i
		}
	}
	return -1
}
