// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

func NewLexer() ILexer {
	return &lexer{start: 0, pos: 0, state: nil, funcs: 0}
}

type ILexer interface {
	Parse(input string, parser IParser)
}

type LexState func(l *lexer) LexState

type lexer struct {
	input  string
	parser IParser
	start  int
	pos    int
	state  LexState
	funcs  uint
}

func (l *lexer) Parse(input string, parser IParser) {
	l.input = input
	l.parser = parser
	l.state = LexEntryState

	for {
		l.state = l.state(l)
		if l.state == nil {
			break
		}
	}
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) inc() {
	l.pos++
}

func (l *lexer) dec() {
	l.pos--
}

func (l *lexer) in() {
	l.funcs++
}

func (l *lexer) out() {
	l.funcs--
}

func (l *lexer) isfunc() bool {
	return l.funcs > 0
}

func (l *lexer) emit(t TokenType) {
	l.parser.Interpret(t, l.input[l.start:l.pos])
	l.start = l.pos
}

func (l *lexer) skipSpaces() {
	for ; l.pos < len(l.input); l.pos++ {
		c := l.input[l.pos]
		if c != SPACE {
			break
		}
	}

	l.start = l.pos
}

// --- StateFunc ---

func LexEntryState(l *lexer) LexState {
	if l.isfunc() {
		l.skipSpaces()
	}

	if l.pos >= len(l.input) {
		return nil
	}

	switch l.input[l.pos] {
	case PREFIX:
		l.inc()
		return LexPrefixState
	default:
		return LexValueState
	}
}

func LexValueState(l *lexer) LexState {
	for ; l.pos < len(l.input); l.pos++ {
		c := l.input[l.pos]
		switch c {
		case PREFIX:
			l.emit(TOKEN_VALUE)
			l.inc()
			return LexPrefixState
		case RIGHT_BRACKET:
			if l.isfunc() {
				if l.pos != l.start {
					l.emit(TOKEN_VALUE)
				}
				l.emit(TOKEN_END_METHOD)
				l.out()
				l.inc()
				l.ignore()
				return LexEntryState
			}
		case COMMA:
			if l.isfunc() {
				if l.pos != l.start {
					l.emit(TOKEN_VALUE)
				}
				l.inc()
				l.ignore()
				l.emit(TOKEN_END_PARAM)
				return LexEntryState
			}
		}
	}

	l.emit(TOKEN_VALUE)
	return nil
}

func LexPrefixState(l *lexer) LexState {
	for ; l.pos < len(l.input); l.pos++ {
		c := l.input[l.pos]
		if !NameExp.Match([]byte{c}) {
			if l.pos == l.start {
				return LexValueState
			} else {
				switch c {
				case LEFT_BRACKET:
					l.start++
					l.emit(TOKEN_BEGIN_METHOD)
					l.in()
					l.inc()
					l.ignore()
					return LexEntryState
				case RIGHT_BRACKET:
					l.start++
					l.emit(TOKEN_VARIABLE)
					if l.isfunc() {
						l.emit(TOKEN_END_METHOD)
						l.out()
						l.inc()
					}
					return LexEntryState
				default:
					l.start++
					l.emit(TOKEN_VARIABLE)
					return LexValueState
				}
			}
		}
	}

	l.emit(TOKEN_VARIABLE)
	return nil
}
