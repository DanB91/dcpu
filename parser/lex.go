// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package parser

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

const Eof = -1

type LexFunc func(*Lexer) LexFunc

// A Lexer turns GASM source into a stream of tokens.
type Lexer struct {
	ch     chan *Token // Channel on which we send out new tokens.
	state  LexFunc     // Current state.
	data   []byte      // Data to be parsed.
	line   [2]int      // Current line and line where token started.
	col    [2]int      // Current column and column where token started.
	start  int         // Start of current token.
	pos    int         // Current position in buffer.
	size   int         // Size of last read rune.
	prlnsz int         // Size of previous line. Needed for accurate line/col tracking when rewinding.
}

// Run processes the given data and yields tokens on the returned channel.
func (l *Lexer) Run(data []byte) <-chan *Token {
	// Cheap hack to simplify Eof handling in lexer.
	// The extra whitespace at the end ensures that the lexer doesn't
	// swallow the last real character in an Eof message. 
	if sz := len(data); sz == 0 {
		data = []byte{'\n'}
	} else if data[sz-1] != '\n' {
		data = append(data, '\n')
	}

	l.data = data
	l.line[0], l.line[1] = 1, 1
	l.col[0], l.col[1] = 1, 1
	l.ch = make(chan *Token)
	l.state = lexText

	go func() {
		defer close(l.ch)

		// Loop for as long as we have a valid state.
		for l.state != nil {
			l.state = l.state(l)
		}
	}()

	return l.ch
}

func (l *Lexer) errorf(f string, argv ...interface{}) {
	var tok Token

	tok.Type = TokErr
	tok.Line = l.line[1]
	tok.Col = l.col[1]
	tok.Data = []byte(fmt.Sprintf(f, argv...))

	l.ch <- &tok
	l.ignore()
}

func (l *Lexer) emit(tt TokenType) {
	var tok Token

	tok.Type = tt
	tok.Line = l.line[1]
	tok.Col = l.col[1]
	tok.Data = l.data[l.start:l.pos]

	if tt == TokIdent || tt == TokLabel {
		tok.Data = bytes.ToLower(tok.Data)
	}

	l.ch <- &tok
	l.ignore()
}

// nextRune retuns the nextRune unicode rune in the input.
func (l *Lexer) nextRune() (r rune) {
	if l.pos >= len(l.data) {
		return Eof
	}

	r, l.size = utf8.DecodeRune(l.data[l.pos:])
	l.pos += l.size

	if r == '\n' {
		l.line[0]++
		l.prlnsz, l.col[0] = l.col[0], 0
	}

	l.col[0]++
	return r
}

// ignore the input so far.
func (l *Lexer) ignore() {
	l.start = l.pos
	l.line[1] = l.line[0]
	l.col[1] = l.col[0]
}

// rewind rewinds to the last rune.
// Can be called only once per nextRune() call.
func (l *Lexer) rewind() {
	l.pos -= l.size
	if l.col[0] > 1 {
		l.col[0]--
	} else {
		l.line[0]--
		l.col[0] = l.prlnsz
	}
}

// skip Skips the nextRune character.
func (l *Lexer) skip() {
	l.nextRune()
	l.ignore()
}

// accept consumes the next rune if it is contained in the supplied string.
func (l *Lexer) accept(valid string) int {
	r := l.nextRune()

	if r == Eof {
		return Eof
	}

	if indexRune(valid, r) == -1 {
		l.rewind()
		return 0
	}

	return 1
}

// acceptRun consumes runes for as long they are contained in the supplied string.
// It returns the number of runes consumed or Eof.
func (l *Lexer) acceptRun(valid string) int {
	var r rune
	var n int

	for {
		if r = l.nextRune(); r == Eof {
			return Eof
		}

		if indexRune(valid, r) == -1 {
			l.rewind()
			break
		}

		n++
	}

	return n
}

// acceptUntil consumes runes for as long they are NOT contained in the
// supplied string.
func (l *Lexer) acceptUntil(valid string) int {
	var r rune
	var n int

	for {
		if r = l.nextRune(); r == Eof {
			return Eof
		}

		if indexRune(valid, r) != -1 {
			l.rewind()
			return n
		}

		n++
	}

	return 0
}

// acceptLiteral consumes runes if they are an exact, rune-for-rune match with
// the supplied string.
func (l *Lexer) acceptLiteral(valid string) int {
	if len(valid) == 0 || l.pos+len(valid) >= len(l.data) {
		return 0
	}

	// This is orders of magnitude faster than using bytes.Index().
	for i := range valid {
		if l.data[l.pos+i] != valid[i] {
			return 0
		}
	}

	// Update line/col info by consuming runes.
	c := utf8.RuneCount(l.data[l.pos : l.pos+len(valid)])
	for i := 0; i < c; i++ {
		l.nextRune()
	}

	return 1
}

// acceptIdent consumes runes until it hits anything that does not
// qualify as a valid identifier.
func (l *Lexer) acceptIdent() int {
	var r rune

	for {
		if r = l.nextRune(); r == Eof {
			return Eof
		}

		if isDigit(r) || isIdent(r) {
			continue
		}

		l.rewind()
		break
	}

	return 1
}

// acceptSpace consumes runes for as long as they are whitespace,
// /except newlines/.
func (l *Lexer) acceptSpace() int {
	var r rune
	var n int

	for {
		if r = l.nextRune(); r == Eof {
			return Eof
		}

		if r != '\n' && unicode.IsSpace(r) {
			n++
			continue
		}

		l.rewind()
		break
	}

	return n
}

func indexRune(s string, r rune) int {
	for i, tr := range s {
		if tr == r {
			return i
		}
	}
	return -1
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdent(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isOperator(r rune) bool {
	switch r {
	case '+', '-', '/', '*', '>', '<', '&', '|', '^', '%', '=':
		return true
	}
	return false
}
