package main

import (
	"fmt"
	"unicode"
)

type TokenType byte

const (
	TokenUnknown TokenType = iota
	TokenEOF
	TokenNumber
	TokenPlus
)

func (t TokenType) String() string {
	switch t {
	case TokenUnknown:
		return "Unknown"
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "Number"
	case TokenPlus:
		return "Plus"
	default:
		return "Could not convert TokenType to string"
	}
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q}", t.Type, t.Literal)
}

type Lexer struct {
	src []rune
	pos int
}

func NewLexer(src string) *Lexer {
	return &Lexer{src: []rune(src)}
}

func (l *Lexer) Next() Token {
	l.skipWhitespace()

	if l.pos >= len(l.src) {
		return l.make(TokenEOF, "")
	}

	ch := l.src[l.pos]

	switch {
	case unicode.IsDigit(ch):
		return l.readNumber()
	case ch == '+':
		return l.advance(TokenPlus)
	default:
		return l.advance(TokenUnknown)
	}
}

func Lex(src string) []Token {
	var tokens []Token
	lexer := NewLexer(src)

	for {
		token := lexer.Next()
		tokens = append(tokens, token)

		if token.Type == TokenEOF {
			break
		}
	}

	return tokens
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.src) && unicode.IsSpace(l.src[l.pos]) {
		l.pos++
	}
}

// advance consumes one rune and returns a token for it.
func (l *Lexer) advance(typ TokenType) Token {
	tok := l.make(typ, string(l.src[l.pos]))
	l.pos++
	return tok
}

func (l *Lexer) make(typ TokenType, lit string) Token {
	return Token{Type: typ, Literal: lit, Pos: l.pos}
}

func (l *Lexer) makeRange(typ TokenType, start int) Token {
	return Token{Type: typ, Literal: string(l.src[start:l.pos]), Pos: start}
}

func (l *Lexer) readNumber() Token {
	start := l.pos

	for l.pos < len(l.src) && unicode.IsDigit(l.src[l.pos]) {
		l.pos++
	}

	// Optional decimal part
	if l.pos < len(l.src) && l.src[l.pos] == '.' {
		l.pos++
		for l.pos < len(l.src) && unicode.IsDigit(l.src[l.pos]) {
			l.pos++
		}
	}

	return l.makeRange(TokenNumber, start)
}
