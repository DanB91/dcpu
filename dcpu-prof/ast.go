// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"github.com/jteeuwen/dcpu/parser"
	"os"
)

var astCache map[string]*parser.AST

func init() {
	astCache = make(map[string]*parser.AST)
}

// GetAST returns the AST for the given file.
func GetAST(file string) (*parser.AST, error) {
	if a, ok := astCache[file]; ok {
		return a, nil
	}

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	astCache[file] = new(parser.AST)
	err = astCache[file].Parse(fd, file)
	return astCache[file], err
}
