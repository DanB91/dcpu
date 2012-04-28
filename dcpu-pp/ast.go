// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"io"
	"strconv"
	"unicode/utf8"
)

// An Abstract Syntax Tree.
type AST struct {
	Files []string // List of files names from which this tree was built.
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

	return a.readDocument(lex.Run(buf.Bytes()), &a.Root.Children)
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

	instr.Children = append(instr.Children,
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
				if len(expr.Children) == 0 {
					return a.errorf(tok, "Expected expression")
				}

				// Correct expression source location
				expr.line = expr.Children[0].Line()
				expr.col = expr.Children[0].Col()
				instr.Children = append(instr.Children, expr)
				return

			case TokComma:
				if len(expr.Children) == 0 {
					return a.errorf(tok, "Expected expression")
				}

				// Correct expression source location
				expr.line = expr.Children[0].Line()
				expr.col = expr.Children[0].Col()
				instr.Children = append(instr.Children, expr)
				expr = NewExpression(file, tok.Line, tok.Col)

			case TokComment:
				expr.Children = append(expr.Children,
					NewComment(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokIdent:
				expr.Children = append(expr.Children,
					NewName(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokOperator:
				expr.Children = append(expr.Children,
					NewOperator(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokNumber:
				num := NewNumber(file, tok.Line, tok.Col, 0)
				num.Parse(tok.Data)
				expr.Children = append(expr.Children, num)

			case TokString:
				data, err := escape(tok.Data)
				if err != nil {
					return err
				}

				expr.Children = append(expr.Children,
					NewString(file, tok.Line, tok.Col, string(data)),
				)

			case TokRawString:
				expr.Children = append(expr.Children,
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

				expr.Children = append(expr.Children,
					NewNumber(file, tok.Line, tok.Col, Word(r)),
				)

			case TokBlockStart:
				if err = a.readBlock(c, &expr.Children, tok); err != nil {
					return
				}

			case TokExprStart:
				if err = a.readExpr(c, &expr.Children, tok); err != nil {
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
				if len(block.Children) == 0 {
					return a.errorf(tok, "Expected expression")
				}
				return

			case TokIdent:
				block.Children = append(block.Children,
					NewName(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokNumber:
				num := NewNumber(file, tok.Line, tok.Col, 0)
				num.Parse(tok.Data)
				block.Children = append(block.Children, num)

			case TokOperator:
				block.Children = append(block.Children,
					NewOperator(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokExprStart:
				if err = a.readExpr(c, &block.Children, tok); err != nil {
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
				expr.Children = append(expr.Children, num)

			case TokOperator:
				expr.Children = append(expr.Children,
					NewOperator(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokIdent:
				expr.Children = append(expr.Children,
					NewName(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokExprStart:
				if err = a.readExpr(c, &expr.Children, tok); err != nil {
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
