// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"bytes"
	"github.com/jteeuwen/dcpu/cpu"
	"io"
	"strconv"
	"unicode/utf8"
)

// An Abstract Syntax Tree.
type AST struct {
	Files []string // List of file names from which this tree was built.
	Root  *Block   // Root node.
}

// Parse takes the given input stream and merges its AST nodes with
// the current AST instance.
//
// The filename is optional and used for error messages and debugging.
func (a *AST) Parse(r io.Reader, filename string) (err error) {
	var lex Lexer
	var buf bytes.Buffer

	if len(filename) == 0 {
		filename = a.tempName()
	}

	if a.hasFile(filename) {
		return // Already parsed. Ignore it.
	}

	a.Files = append(a.Files, filename)

	if a.Root == nil {
		a.Root = NewBlock(0, 0, 0)
	}

	io.Copy(&buf, r)

	if err = a.readDocument(lex.Run(buf.Bytes()), &a.Root.children); err != nil {
		return
	}

	constants := make(map[string]Node)

	if err = a.findConstants(constants); err != nil {
		return
	}

	a.replaceConstants(a.Root.children, constants)
	return
}

// replaceConstants finds al constant references and replaces them
// with constant values.
func (a *AST) replaceConstants(nodes []Node, c map[string]Node) {
	for i := range nodes {
		instr, ok := a.Root.children[i].(*Instruction)

		if !ok {
			continue
		}

		argv := instr.children[1:]
		for j := range argv {
			expr := argv[j].(*Expression)

			for k := range expr.children {
				switch tt := expr.children[k].(type) {
				case *Name:
					if node, ok := c[tt.Data]; ok {
						expr.children[k] = node
					}
				}
			}
		}
	}

	return
}

// findConstants finds all constant (EQU) definitions.
// Adds them to the supplied map and removes them from the AST.
func (a *AST) findConstants(list map[string]Node) (err error) {
	for i := 0; i < len(a.Root.children); i++ {
		instr, ok := a.Root.children[i].(*Instruction)
		if !ok {
			continue
		}

		name := instr.children[0].(*Name)

		if name.Data != "equ" {
			continue
		}

		if len(instr.children) != 3 {
			return NewParseError(
				a.Files[instr.File()], instr.Line(), instr.Col(),
				"Invalid argument count for EQU. Want 2, got %d.",
				len(instr.children)-1,
			)
		}

		exp := instr.children[1].(*Expression)
		key, ok := exp.children[0].(*Name)

		if !ok {
			return NewParseError(
				a.Files[key.File()], key.Line(), key.Col(),
				"Invalid Node %T. Expected *Name", key,
			)
		}

		if _, ok = list[key.Data]; ok {
			return NewParseError(
				a.Files[key.File()], key.Line(), key.Col(),
				"Duplicate constant %q", key.Data,
			)
		}

		val := instr.children[2].(*Expression).children[0]

		switch val.(type) {
		case *Name:
		case *Number:
		case *String:
		default:
			return NewParseError(
				a.Files[val.File()], val.Line(), val.Col(),
				"Invalid Node %T. Expected *Name", val,
			)
		}

		list[key.Data] = val

		copy(a.Root.children[i:], a.Root.children[i+1:])
		a.Root.children = a.Root.children[:len(a.Root.children)-1]
		i--
	}

	return
}

// readDocument reads tokens from the given channel and turns them into AST nodes.
// It deals with top-level language constructs.
func (a *AST) readDocument(c <-chan *Token, n *[]Node) (err error) {
	file := len(a.Files) - 1

	for {
		select {
		case tok := <-c:
			if tok == nil {
				return
			}

			if tok.Type == TokErr {
				return a.errorf(tok, "%s", tok.Data)
			}

			switch tok.Type {
			case TokComment:
				*n = append(*n, NewComment(file, tok.Line, tok.Col, string(tok.Data)))

			case TokLabel:
				*n = append(*n, NewLabel(file, tok.Line, tok.Col, string(tok.Data)))

			case TokIdent:
				err = a.readInstruction(c, n, tok)

			default:
				return a.errorf(tok, "Unexpected token %s. Want Comment, Label or Ident", tok)
			}

			if err != nil {
				return
			}
		}
	}

	return
}

// readInstruction reads tokens from the given channel and turns them into AST nodes.
// It deals with individual instructions.
func (a *AST) readInstruction(c <-chan *Token, n *[]Node, tok *Token) (err error) {
	file := len(a.Files) - 1
	instr := NewInstruction(len(a.Files)-1, tok.Line, tok.Col)
	expr := NewExpression(file, 0, 0)

	instr.children = append(instr.children,
		NewName(file, tok.Line, tok.Col, string(tok.Data)),
	)

	defer func() { *n = append(*n, instr) }()

	for {
		select {
		case tok := <-c:
			if tok == nil {
				return
			}

			if tok.Type == TokErr {
				return a.errorf(tok, "%s", tok.Data)
			}

			switch tok.Type {
			case TokEndLine:
				if len(expr.children) > 0 {
					// Correct expression source location
					expr.line = expr.children[0].Line()
					expr.col = expr.children[0].Col()
					instr.children = append(instr.children, expr)
				}
				return

			case TokComma:
				if len(expr.children) == 0 {
					return a.errorf(tok, "Expected expression")
				}

				// Correct expression source location
				expr.line = expr.children[0].Line()
				expr.col = expr.children[0].Col()
				instr.children = append(instr.children, expr)
				expr = NewExpression(file, tok.Line, tok.Col)

			case TokComment:
				expr.children = append(expr.children,
					NewComment(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokIdent:
				expr.children = append(expr.children,
					NewName(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokOperator:
				expr.children = append(expr.children,
					NewOperator(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokNumber:
				num := NewNumber(file, tok.Line, tok.Col, 0)
				num.Parse(tok.Data)
				expr.children = append(expr.children, num)

			case TokString:
				data, err := escape(tok.Data)
				if err != nil {
					return err
				}

				expr.children = append(expr.children,
					NewString(file, tok.Line, tok.Col, string(data)),
				)

			case TokRawString:
				expr.children = append(expr.children,
					NewString(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokChar:
				data, err := escape(tok.Data)
				if err != nil {
					return err
				}

				r, size := utf8.DecodeRune(data)

				// If the encoding is invalid, utf8.DecodeRune yields (RuneError, 1).
				// This constitutes an impossible result for correct UTF-8.
				if r == utf8.RuneError && size == 1 {
					return a.errorf(tok, "Invalid utf8 character literal: %s", tok)
				}

				expr.children = append(expr.children,
					NewNumber(file, tok.Line, tok.Col, cpu.Word(r)),
				)

			case TokBlockStart:
				if err = a.readBlock(c, &expr.children, tok); err != nil {
					return
				}

			case TokExprStart:
				if err = a.readExpr(c, &expr.children, tok); err != nil {
					return
				}

			default:
				return a.errorf(tok,
					"Unexpected token %s. Want Comment, Ident, Number, Block or Expression", tok)
			}
		}
	}

	return
}

// readBlock reads tokens from the given channel and turns them into AST nodes.
// It deals with block statements.
func (a *AST) readBlock(c <-chan *Token, n *[]Node, tok *Token) (err error) {
	file := len(a.Files) - 1
	block := NewBlock(file, tok.Line, tok.Col)

	defer func() { *n = append(*n, block) }()

	for {
		select {
		case tok := <-c:
			if tok == nil {
				return
			}

			if tok.Type == TokErr {
				return a.errorf(tok, "%s", tok.Data)
			}

			switch tok.Type {
			case TokBlockEnd:
				if len(block.children) == 0 {
					return a.errorf(tok, "Expected expression")
				}
				return

			case TokIdent:
				block.children = append(block.children,
					NewName(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokNumber:
				num := NewNumber(file, tok.Line, tok.Col, 0)
				num.Parse(tok.Data)
				block.children = append(block.children, num)

			case TokOperator:
				block.children = append(block.children,
					NewOperator(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokExprStart:
				if err = a.readExpr(c, &block.children, tok); err != nil {
					return
				}

			default:
				return a.errorf(tok,
					"Unexpected token %s. Want Ident, Number or Expression", tok)
			}
		}
	}
	return
}

// readExpr reads tokens from the given channel and turns them into AST nodes.
// It deals with expression statements.
func (a *AST) readExpr(c <-chan *Token, n *[]Node, tok *Token) (err error) {
	file := len(a.Files) - 1
	expr := NewExpression(file, tok.Line, tok.Col)

	defer func() { *n = append(*n, expr) }()

	for {
		select {
		case tok := <-c:
			if tok == nil {
				return
			}

			if tok.Type == TokErr {
				return a.errorf(tok, "%s", tok.Data)
			}

			switch tok.Type {
			case TokExprEnd, TokComma:
				return

			case TokNumber:
				num := NewNumber(file, tok.Line, tok.Col, 0)
				num.Parse(tok.Data)
				expr.children = append(expr.children, num)

			case TokOperator:
				expr.children = append(expr.children,
					NewOperator(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokIdent:
				expr.children = append(expr.children,
					NewName(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokExprStart:
				if err = a.readExpr(c, &expr.children, tok); err != nil {
					return
				}

			default:
				return a.errorf(tok,
					"Unexpected token %s. Want Number, Expression or Operator", tok)
			}
		}
	}
	return
}

// hasFile returns true if the given filename is already listed.
func (a *AST) hasFile(s string) bool {
	for i := range a.Files {
		if a.Files[i] == s {
			return true
		}
	}
	return false
}

// errorf formats an error messsage from the given token context and returns it.
func (a *AST) errorf(t *Token, f string, argv ...interface{}) error {
	var file string

	if len(a.Files) > 0 {
		file = a.Files[len(a.Files)-1]
	}

	return NewParseError(file, t.Line, t.Col, f, argv...)
}

// tempName creates a unique file name.
func (a *AST) tempName() string {
	var s string
	var n int

	for {
		s = strconv.Itoa(n) + ".tmp"

		if !a.hasFile(s) {
			break
		}

		n++
	}

	return s
}
