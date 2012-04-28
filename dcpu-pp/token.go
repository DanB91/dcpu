// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "fmt"

type TokenType uint8

/// Known token types. Do not chabge the order of these.
const (
	TokEof TokenType = iota
	TokErr
	TokOperator
	TokLabel
	TokIdent
	TokString
	TokRawString
	TokChar
	TokNumber
	TokComment
	TokBlockStart
	TokBlockEnd
	TokExprStart
	TokExprEnd
	TokEndLine
	TokComma
)

func (t TokenType) String() string {
	switch t {
	case TokEof:
		return "Eof"
	case TokErr:
		return "Error"
	case TokOperator:
		return "Operator"
	case TokLabel:
		return "Label"
	case TokIdent:
		return "Ident"
	case TokString:
		return "String"
	case TokRawString:
		return "RawString"
	case TokChar:
		return "Char"
	case TokNumber:
		return "Number"
	case TokComment:
		return "Comment"
	case TokBlockStart:
		return "BlockStart"
	case TokBlockEnd:
		return "BlockEnd"
	case TokExprStart:
		return "ExprStart"
	case TokExprEnd:
		return "ExprEnd"
	case TokEndLine:
		return "EndLine"
	case TokComma:
		return "Comma"
	}
	return fmt.Sprintf("0x%02x", uint8(t))
}

type Token struct {
	Data []byte
	Line int
	Col  int
	Type TokenType
}

func (t Token) String() string {
	if len(t.Data) > 30 {
		return fmt.Sprintf("%s(%.30q...)", t.Type, t.Data)
	}
	return fmt.Sprintf("%s(%q)", t.Type, t.Data)
}
