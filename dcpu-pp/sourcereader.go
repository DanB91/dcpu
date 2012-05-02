// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	dp "github.com/jteeuwen/dcpu/parser"
	"io"
	"os"
	"path/filepath"
)

// parseInput takes the input files and parses their contents into
// the given AST.
func parseInput(ast *dp.AST, c *Config) (err error) {
	if err = readSource(ast, c.Input); err != nil {
		return
	}

	return resolveIncludes(ast, c)
}

// readSource reads the given file and parses its contents
// into the given AST.
func readSource(ast *dp.AST, file string) error {
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
func resolveIncludes(ast *dp.AST, c *Config) (err error) {
	var labels []*dp.Label
	var refs []*dp.Name

	findLabels(ast.Root.Children(), &labels)
	findReferences(ast.Root.Children(), &refs)

	refs = findUndefinedRefs(refs, labels)
	refs = stripDuplicateNames(refs)

	if len(refs) == 0 {
		// No undefined references. We're done here.
		return
	}

	if len(c.Include) == 0 {
		// We have unresolved references, but no places to look
		// for their implementation. This constitutes a booboo.
		return dp.NewParseError(ast.Files[refs[0].File()], refs[0].Line(), refs[0].Col(),
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
func loadInclude(ast *dp.AST, c *Config, r *dp.Name) (err error) {
	var file string

	name := r.Data + ".dasm"
	walker := func(f string, info os.FileInfo, err error) error {
		if filepath.Base(f) == name {
			file = f
			return io.EOF // Signal walker to stop.
		}

		return nil
	}

	for i := range c.Include {
		filepath.Walk(c.Include[i], walker)

		if len(file) > 0 {
			break
		}
	}

	if len(file) == 0 {
		return dp.NewParseError(ast.Files[r.File()], r.Line(), r.Col(),
			"Undefined reference: %q", r.Data)
	}

	if err = readSource(ast, file); err != nil {
		return
	}

	// Test if the code actually contains the label we are looking for.
	if !includeHasLabel(ast, file, r.Data) {
		return dp.NewParseError(ast.Files[r.File()], r.Line(), r.Col(),
			"Undefined reference: %q. Include file was found, but "+
				"it did not define the desired label.", r.Data)
	}

	// This new file may hold its own include requirements.
	return resolveIncludes(ast, c)
}

// findUndefinedRefs compares both given lists of labels and
// label references. Any reference that is not present in the
// label list, is considered unresolved and added to the 
// returned list.
func findUndefinedRefs(refs []*dp.Name, labels []*dp.Label) []*dp.Name {
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
func includeHasLabel(ast *dp.AST, file string, target string) bool {
	index := fileIndex(ast, file)
	if index == -1 {
		return false
	}

	return hasLabel(ast.Root.Children(), index, target)
}

// hasLabel recursively finds a specific label definition.
// Returns true if it was found. False otherwise.
func hasLabel(n []dp.Node, file int, target string) bool {
	for i := range n {
		switch tt := n[i].(type) {
		case dp.NodeCollection:
			if hasLabel(tt.Children(), file, target) {
				return true
			}

		case *dp.Label:
			if tt.File() == file && tt.Data == target {
				return true
			}

			// @_@ These are not the labels you are looking for. @_@
		}
	}

	return false
}

// fileIndex returns the given file's index as it is stored in the AST.
func fileIndex(ast *dp.AST, file string) int {
	for i := range ast.Files {
		if ast.Files[i] == file {
			return i
		}
	}
	return -1
}
