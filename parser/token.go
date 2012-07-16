// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

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
	TokDef
	TokEndDef
)

func (t TokenType) String() string {
	switch t {
	case TokEof:
		return "EOF"
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
	case TokDef:
		return "Def"
	case TokEndDef:
		return "EndDef"
	}
	panic("unreachable")
}

type Token struct {
	Data []byte
	Line int
	Col  int
	Type TokenType
}

func (s *Token) String() string {
	return s.Type.String()
}
