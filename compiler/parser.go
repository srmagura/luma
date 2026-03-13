package compiler

import (
	"strconv"

	"github.com/srmagura/luma/shared"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func Parse(tokens []Token) (shared.Node, error) {
	p := NewParser(tokens)
	ast, err := p.parseAdditiveExpr()
	return ast, err
}

func (p *Parser) error(message string) (shared.Node, error) {
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

func (p *Parser) consumeExpected(expected TokenType) (Token, bool) {
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

func (p *Parser) consume() (Token, bool) {
	return p.consumeExpected(TokenUnknown)
}

// Handle + and -
func (p *Parser) parseAdditiveExpr() (shared.Node, error) {
	result, err := p.parseNumber()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek()
		if !ok || (peeked.Type != TokenPlus && peeked.Type != TokenMinus) {
			break
		}

		opToken, ok := p.consume()
		if !ok {
			return nil, nil
		}

		var op shared.Op
		switch opToken.Type {
		case TokenPlus:
			op = shared.OpAdd
		case TokenMinus:
			op = shared.OpSubtract
		}

		right, err := p.parseNumber()
		if err != nil {
			return nil, err
		}

		result = shared.BinaryExpr{Op: op, Left: result, Right: right}
	}

	return result, nil
}

func (p *Parser) parseNumber() (shared.Node, error) {
	token, ok := p.peek()
	if !ok || token.Type != TokenNumber {
		return nil, nil
	}

	token, ok = p.consumeExpected(TokenNumber)
	if !ok {
		return nil, nil
	}

	n, err := strconv.Atoi(token.Literal)
	if err != nil {
		return p.error("Failed to parse int")
	}

	return shared.IntLiteral{Value: n, Pos: token.Pos}, nil
}
