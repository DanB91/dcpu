// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

const (
	runeComment    = ';'
	runeLabel      = ':'
	runeBlockStart = '['
	runeBlockEnd   = ']'
	runeExprStart  = '('
	runeExprEnd    = ')'
	runeComma      = ','
	runeNewline    = '\n'
	runeString     = '"'
	runeRawString  = '`'
	runeChar       = '\''
)

// This processes the top-level document.
// Filters expected bits and takes it from there.
// This is the initial lexer state.
func lexText(l *Lexer) LexFunc {
	if l.acceptSpace() < 0 {
		return nil
	}

	l.ignore()

	r := l.nextRune()

	switch {
	case r == Eof:
		return nil

	case r == runeNewline:
		l.ignore()
		return lexText

	case r == runeComment:
		l.ignore()
		l.acceptUntil("\r\n")
		l.emit(TokComment)
		return lexText

	case r == runeLabel:
		l.ignore()
		return lexLabel

	case isIdent(r):
		l.rewind()
		return lexInstruction
	}

	l.errorf("Expected comment, label or identifier")
	return nil
}

// lexInstruction processes an instruction definition.
func lexLabel(l *Lexer) LexFunc {
	r := l.acceptIdent()

	switch {
	case r == Eof:
		return nil

	case r > 0:
		l.emit(TokLabel)
		return lexText
	}

	l.errorf("Invalid label definition. Missing name.")
	return nil
}

// lexInstruction processes instructions and arguments.
func lexInstruction(l *Lexer) LexFunc {
	// Clear any remaining space up until a newline.
	if l.acceptSpace() < 0 {
		return nil
	}

	l.ignore()

	// Figure out what the next rune represents.
	r := l.nextRune()

	switch {
	case r == Eof:
		return nil

	case r == runeNewline:
		l.emit(TokEndLine)
		return lexText

	case r == runeComment:
		l.ignore()
		l.acceptUntil("\r\n")
		l.emit(TokComment)
		return lexInstruction

	case r == runeComma:
		l.emit(TokComma)
		return lexInstruction

	case r == runeBlockStart:
		l.emit(TokBlockStart)
		return lexInstruction

	case r == runeBlockEnd:
		l.emit(TokBlockEnd)
		return lexInstruction

	case r == runeExprStart:
		l.emit(TokExprStart)
		return lexInstruction

	case r == runeExprEnd:
		l.emit(TokExprEnd)
		return lexInstruction

	case r == runeString:
		l.ignore()
		return lexString

	case r == runeRawString:
		l.ignore()
		return lexRawString

	case r == runeChar:
		l.ignore()
		return lexChar

	case isOperator(r):
		l.rewind()
		return lexOperator

	case isDigit(r):
		l.rewind()
		return lexNumber

	case isIdent(r):
		l.rewind()
		l.acceptIdent()
		l.emit(TokIdent)
		return lexInstruction
	}

	l.errorf("Expected identifier, number or expression. Got %c", r)
	return nil
}

// lexString processes string literals
func lexString(l *Lexer) LexFunc {
	ret := l.acceptUntil("\"")

	if ret > 0 {
		l.emit(TokString)
		l.skip()
		return lexInstruction
	}

	l.errorf("Newline in string")
	return nil
}

// lexRawString processes raw string literals.
// Same as a regular string, but escape sequences are not resolved.
func lexRawString(l *Lexer) LexFunc {
	ret := l.acceptUntil("`")

	if ret > 0 {
		l.emit(TokRawString)
		l.skip()
		return lexInstruction
	}

	l.errorf("EOF in string")
	return nil
}

// lexChar processes character literals
func lexChar(l *Lexer) LexFunc {
	r := l.acceptUntil("'")

	switch {
	case r == Eof:
		return nil

	case r > 1:
		// Account for simple escape sequence.
		if r == 2 && l.data[l.start+1] == '\\' {
			l.errorf("Missing ' in character literal")
			return nil
		}

		fallthrough

	case r == 1:
		l.emit(TokChar)
		l.skip()
		return lexInstruction
	}

	l.errorf("Newline in string")
	return nil
}

// lexOperator processes arithmatic operators
func lexOperator(l *Lexer) LexFunc {
	var r int

	if r = l.acceptLiteral(">>"); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	if r = l.acceptLiteral("<<"); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	if r = l.acceptLiteral("<="); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	if r = l.acceptLiteral(">="); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	if r = l.acceptLiteral("=="); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	if r = l.acceptLiteral("||"); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	if r = l.acceptLiteral("&&"); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	if r = l.accept("+-*/%^&|><="); r == Eof {
		return nil
	} else if r > 0 {
		l.emit(TokOperator)
		return lexInstruction
	}

	l.errorf("Expected operator")
	return nil
}

// lexNumber tests if the given input might qualify as a number. This is not
// a guarantee, but tests for a reasonable likeness.
//
// This finds integers of the following formats:
//
//    0b0101101011011001 (binary)
//    0777               (octal)
//    1234               (decimal)
//    0xff12             (hexadecimal)
func lexNumber(l *Lexer) LexFunc {
	const Digits = "0123456789abcdefABCDEF"

	var ret int

	base := 10
	if ret = l.accept("0"); ret > 0 {
		if l.accept("xX") > 0 {
			base = len(Digits)

		} else if l.accept("bB") > 0 {
			base = 2

		} else {
			base = 8
			l.rewind()
		}
	}

	if ret = l.acceptRun(Digits[:base]); ret <= 0 {
		if ret == 0 {
			l.errorf("Invalid number")
		}
		return nil
	}

	l.emit(TokNumber)
	return lexInstruction
}
