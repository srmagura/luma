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

	if ast == nil {
		return p.error("Parsing failed")
	}

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

func (p *Parser) parseExpr() (shared.Node, error) {
	return p.parseAdditiveExpr()
}

// Handle + and -
func (p *Parser) parseAdditiveExpr() (shared.Node, error) {
	left, err := p.parseMultiplicativeExpr()
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

		right, err := p.parseMultiplicativeExpr()
		if err != nil {
			return nil, err
		}

		left = shared.BinaryExpr{Op: op, Left: left, Right: right}
	}

	return left, nil
}

// Handle *, /, and ~/
func (p *Parser) parseMultiplicativeExpr() (shared.Node, error) {
	left, err := p.parseCall()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek()
		if !ok || (peeked.Type != TokenStar && peeked.Type != TokenFSlash && peeked.Type != TokenTildeFSlash) {
			break
		}

		opToken, ok := p.consume()
		if !ok {
			return nil, nil
		}

		var op shared.Op
		switch opToken.Type {
		case TokenStar:
			op = shared.OpMultiply
		case TokenFSlash:
			op = shared.OpDivide
		case TokenTildeFSlash:
			op = shared.OpDivideInteger
		}

		right, err := p.parseCall()
		if err != nil {
			return nil, err
		}

		left = shared.BinaryExpr{Op: op, Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parseCall() (shared.Node, error) {
	left, err := p.parseLeaf()
	if err != nil {
		return nil, err
	}

	// Does not handle parenthesized functions (yet?)
	switch v := left.(type) {
	case shared.IdentNode:
		token, ok := p.peek()
		if !ok || token.Type != TokenLParen {
			return nil, nil
		}

		token, ok = p.consumeExpected(TokenLParen)
		if !ok {
			return nil, nil
		}

		token, ok = p.peek()
		if !ok || token.Type != TokenRParen {
			return nil, nil
		}

		token, ok = p.consumeExpected(TokenRParen)
		if !ok {
			return nil, nil
		}

		left = shared.CallExpr{
			Func: v,
			Args: []shared.Node{},
		}
	}

	return left, nil
}

func (p *Parser) parseLeaf() (shared.Node, error) {
	node, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if node != nil {
		return node, nil
	}

	node, err = p.parseNumber()
	if err != nil {
		return nil, err
	}
	if node != nil {
		return node, nil
	}

	return nil, nil
}

func (p *Parser) parseIdent() (shared.Node, error) {
	token, ok := p.peek()
	if !ok || token.Type != TokenIdent {
		return nil, nil
	}

	token, ok = p.consumeExpected(TokenIdent)
	if !ok {
		return nil, nil
	}

	return shared.IdentNode{Name: token.Literal, Pos: token.Pos}, nil
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
