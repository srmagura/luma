package compiler

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

func Parse(tokens []Token) (Node, error) {
	p := NewParser(tokens)
	ast, err := p.parseNumber()
	return ast, err
}

func (p *Parser) error(message string) (Node, error) {
	return nil, &InternalParserError{
		message: message,
		pos:     p.pos,
	}
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

func (p *Parser) parseNumber() (Node, error) {
	token, ok := p.peek()
	if !ok || token.Type != TokenNumber {
		return nil, nil
	}

	token, ok = p.consume(TokenNumber)
	if !ok {
		return nil, nil
	}

	n, err := strconv.Atoi(token.Literal)
	if err != nil {
		return p.error("Failed to parse int")
	}

	return IntLiteral{Value: n, Pos: token.Pos}, nil
}
