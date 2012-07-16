// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"os"
	"testing"
)

func TestDef(t *testing.T) {
	var ast AST

	file := "../testdata/def.dasm"
	fd, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}

	defer fd.Close()

	err = ast.Parse(fd, file)
	if err != nil {
		t.Fatal(err)
	}
}
