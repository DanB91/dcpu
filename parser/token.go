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
)

type Token struct {
	Data []byte
	Line int
	Col  int
	Type TokenType
}
