// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// DCPU AST parser package.
package parser

import (
	"bytes"
	"io"
	"path"
	"path/filepath"
	"strconv"
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

	if !path.IsAbs(filename) {
		filename, err = filepath.Abs(filename)
		if err != nil {
			return
		}
	}

	if a.hasFile(filename) {
		return // Already parsed. Ignore it.
	}

	a.Files = append(a.Files, filename)

	if a.Root == nil {
		a.Root = NewBlock(0, 0, 0)
	}

	io.Copy(&buf, r)

	err = a.readDocument(lex.Run(buf.Bytes()), &a.Root.children)
	if err != nil {
		return
	}

	err = verify(a, a.Root.children)
	if err != nil {
		return
	}

	a.Root.children, err = a.parseFunctions(a.Root.children)
	return
}

// Functions returns all function definitions.
func (a *AST) Functions() []*Function {
	var list []*Function
	var f *Function
	var ok bool

	for _, v := range a.Root.children {
		if f, ok = v.(*Function); ok {
			list = append(list, f)
		}
	}

	return list
}

func indexOfEnd(in []Node) int {
	var instr *Instruction
	var name *Name
	var ok bool

	for i := range in {
		instr, ok = in[i].(*Instruction)
		if !ok {
			continue
		}

		name = instr.children[0].(*Name)

		switch name.Data {
		case "def":
			return -1
		case "end":
			return i
		}
	}
	return -1
}

// parseFunctions finds function definitions and turns them into proper
// Function nodes.
func (a *AST) parseFunctions(in []Node) (out []Node, err error) {
	var instr *Instruction
	var expr *Expression
	var f *Function
	var name *Name
	var ok bool
	var s, e int

	for s = 0; s < len(in); s++ {
		instr, ok = in[s].(*Instruction)
		if !ok {
			continue
		}

		name = instr.children[0].(*Name)
		if name.Data != "def" {
			continue
		}

		e = indexOfEnd(in[s+1:])
		if e == -1 {
			return nil, NewParseError(
				a.Files[instr.File()], instr.Line(), instr.Col(),
				"Unmatched 'def'.",
			)
		}

		e += s + 1

		expr = instr.children[1].(*Expression)
		name = expr.children[0].(*Name)

		f = NewFunction(name.File(), name.Line(), name.Col())
		f.children = append(f.children, name)
		f.children = append(f.children, in[s+1:e]...)
		in[s] = f

		copy(in[s+1:], in[e+1:])
		in = in[:len(in)-e+s]
	}

	return in, nil
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
	instr := NewInstruction(file, tok.Line, tok.Col)
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
				expr.children = append(expr.children,
					NewNumber(file, tok.Line, tok.Col, string(tok.Data)),
				)

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

				expr.children = append(expr.children,
					NewChar(file, tok.Line, tok.Col, string(data)),
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
				block.children = append(block.children,
					NewNumber(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokChar:
				data, err := escape(tok.Data)
				if err != nil {
					return err
				}

				block.children = append(block.children,
					NewChar(file, tok.Line, tok.Col, string(data)),
				)

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
				expr.children = append(expr.children,
					NewNumber(file, tok.Line, tok.Col, string(tok.Data)),
				)

			case TokChar:
				data, err := escape(tok.Data)
				if err != nil {
					return err
				}

				expr.children = append(expr.children,
					NewChar(file, tok.Line, tok.Col, string(data)),
				)

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
