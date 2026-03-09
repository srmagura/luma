package main

import (
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func Parse(tokens []Token) (Node, bool) {
	p := NewParser(tokens)
	return p.parseNumber()
}

func (p *Parser) peek() (Token, bool) {
	if p.pos >= len(p.tokens) {
		return Token{}, false
	}

	return p.tokens[p.pos], true
}

func (p *Parser) consume(expected TokenType) (Token, bool) {
	if p.pos >= len(p.tokens) {
		return Token{}, false
	}

	token := p.tokens[p.pos]

	if expected != TokenUnknown && token.Type != expected {
		return Token{}, false
	}

	p.pos++
	return token, true
}

// Handle + and -
// func (p *Parser) parseAdditiveExpr() (Node, bool) {
// 	left, ok := p.parseFactor()
// }

// Handle literals and parenthesized expressions
// func (p *Parser) parseFactor() (Node, bool) {
// 	token, ok := p.peek()
// 	if !ok {
// 		return Node{}, false
// 	}

// 	if token.Type == TokenNumber {
// 		return p.parseNumber()
// 	}

// 	return Node{}, false
// }

func (p *Parser) parseNumber() (Node, bool) {
	token, ok := p.consume(TokenNumber)
	if !ok {
		return UnknownNode{}, false
	}

	n, err := strconv.Atoi(token.Literal)
	if err != nil {
		return UnknownNode{}, false
	}

	return IntLiteral{Value: n, Pos: token.Pos}, true
}
