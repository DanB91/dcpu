// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

// verify performs some sanity checks on the ast nodes.
func verify(a *AST, list []Node) (err error) {
	for i := range list {
		switch tt := list[i].(type) {
		case *Instruction:
			err = verifyInstruction(a, tt)
		}

		if err != nil {
			return
		}
	}

	return
}

func verifyInstruction(a *AST, n *Instruction) (err error) {
	list := n.children
	name := list[0].(*Name)

	switch name.Data {
	case "equ":
		return verifyConstant(a, n)
	case "def":
		return verifyDef(a, n)
	}

	var expr *Expression
	var ok bool

	if len(n.children) > 1 {
		expr, ok = n.children[1].(*Expression)
		if !ok || len(expr.children) == 0 {
			goto fail
		}
	}

	if len(n.children) > 2 {
		expr, ok = n.children[2].(*Expression)
		if !ok || len(expr.children) == 0 {
			goto fail
		}
	}

	return

fail:
	return NewParseError(a.Files[n.File()], n.Line(), n.Col(),
		"Invalid instruction %q. Expected: <name> [<expression> [, <expression>]]", name.Data)
}

func verifyConstant(a *AST, n *Instruction) (err error) {
	var expr *Expression
	var ok bool

	if len(n.children) != 3 {
		goto fail
	}

	expr, ok = n.children[1].(*Expression)
	if !ok || len(expr.children) != 1 {
		goto fail
	}

	if _, ok = expr.children[0].(*Name); !ok {
		goto fail
	}

	expr, ok = n.children[2].(*Expression)
	if !ok || len(expr.children) == 0 {
		goto fail
	}

	return

fail:
	return NewParseError(a.Files[n.File()], n.Line(), n.Col(),
		"Invalid constant definition. Expected: equ <name>, <expression>")
}

func verifyDef(a *AST, n *Instruction) (err error) {
	var expr *Expression
	var ok bool

	if len(n.children) != 2 {
		goto fail
	}

	expr, ok = n.children[1].(*Expression)
	if !ok || len(expr.children) != 1 {
		goto fail
	}

	if _, ok = expr.children[0].(*Name); !ok {
		goto fail
	}

	return

fail:
	return NewParseError(a.Files[n.File()], n.Line(), n.Col(),
		"Invalid function definition. Expected: def <name>")
}
