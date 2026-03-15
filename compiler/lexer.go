package compiler

import (
	"fmt"
	"unicode"
)

type tokenType byte

const (
	tokenUnknown tokenType = iota
	tokenEOF
	tokenNumber
	tokenPlus
	tokenMinus
	tokenStar
	tokenFSlash
	tokenTildeFSlash
	tokenLParen
	tokenRParen
	tokenIdent
	tokenComma
	tokenSemi
)

func (t tokenType) String() string {
	switch t {
	case tokenUnknown:
		return "Unknown"
	case tokenEOF:
		return "EOF"
	case tokenNumber:
		return "Number"
	case tokenPlus:
		return "Plus"
	case tokenMinus:
		return "Minus"
	case tokenStar:
		return "Star"
	case tokenFSlash:
		return "FSlash"
	case tokenTildeFSlash:
		return "TildeFSlash"
	case tokenLParen:
		return "LParen"
	case tokenRParen:
		return "RParen"
	case tokenIdent:
		return "Ident"
	case tokenComma:
		return "Comma"
	case tokenSemi:
		return "Semi"
	default:
		return "Could not convert tokenType to string"
	}
}

type token struct {
	_type   tokenType
	literal string
	pos     int
}

func (t token) String() string {
	return fmt.Sprintf("token{%s, %q}", t._type, t.literal)
}

type lexer struct {
	src []rune
	pos int
}

func newLexer(src string) *lexer {
	return &lexer{src: []rune(src)}
}

func (l *lexer) next() token {
	l.skipWhitespace()

	if l.pos >= len(l.src) {
		return l.make(tokenEOF, "")
	}

	ch := l.src[l.pos]

	switch {
	case unicode.IsLetter(ch) || ch == '_':
		return l.readIdent()
	case unicode.IsDigit(ch):
		return l.readNumber()
	case ch == '+':
		return l.advance(tokenPlus)
	case ch == '-':
		return l.advance(tokenMinus)
	case ch == '*':
		return l.advance(tokenStar)
	case ch == '/':
		return l.advance(tokenFSlash)
	case ch == '~':
		return l.readTildeFSlash()
	case ch == '(':
		return l.advance(tokenLParen)
	case ch == ')':
		return l.advance(tokenRParen)
	case ch == ',':
		return l.advance(tokenComma)
	case ch == ';':
		return l.advance(tokenSemi)
	default:
		return l.advance(tokenUnknown)
	}
}

func lex(src string) []token {
	var tokens []token
	lexer := newLexer(src)

	for {
		token := lexer.next()
		if token._type == tokenEOF {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func (l *lexer) skipWhitespace() {
	for l.pos < len(l.src) && unicode.IsSpace(l.src[l.pos]) {
		l.pos++
	}
}

// advance consumes one rune and returns a token for it.
func (l *lexer) advance(typ tokenType) token {
	tok := l.make(typ, string(l.src[l.pos]))
	l.pos++
	return tok
}

func (l *lexer) make(typ tokenType, lit string) token {
	return token{_type: typ, literal: lit, pos: l.pos}
}

func (l *lexer) makeRange(typ tokenType, start int) token {
	return token{_type: typ, literal: string(l.src[start:l.pos]), pos: start}
}

func (l *lexer) readIdent() token {
	start := l.pos

	// First character must be letter or underscore
	if unicode.IsLetter(l.src[l.pos]) || l.src[l.pos] == '_' {
		l.pos++
	} else {
		return token{}
	}

	// Additional characters can be letter, underscore, or number
	for l.pos < len(l.src) && (unicode.IsLetter(l.src[l.pos]) || l.src[l.pos] == '_' || unicode.IsDigit(l.src[l.pos])) {
		l.pos++
	}

	return l.makeRange(tokenIdent, start)
}

func (l *lexer) readNumber() token {
	start := l.pos

	for l.pos < len(l.src) && unicode.IsDigit(l.src[l.pos]) {
		l.pos++
	}

	if l.pos == start {
		return token{}
	}

	// Optional decimal part
	if l.pos < len(l.src) && l.src[l.pos] == '.' {
		l.pos++
		for l.pos < len(l.src) && unicode.IsDigit(l.src[l.pos]) {
			l.pos++
		}
	}

	return l.makeRange(tokenNumber, start)
}

func (l *lexer) readTildeFSlash() token {
	start := l.pos

	if l.pos+1 < len(l.src) && l.src[l.pos] == '~' && l.src[l.pos+1] == '/' {
		l.pos += 2
		return l.makeRange(tokenTildeFSlash, start)
	}

	return token{}
}
