package compiler

import (
	"fmt"
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
	ast, err := p.parseModule()

	if err != nil {
		return nil, err
	}
	if ast == nil {
		return p.error("Parsing failed")
	}

	return ast, err
}

func (p *Parser) error(message string) (shared.Node, error) {
	return nil, &internalCompilerError{
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

func (p *Parser) consumeExpected(expected TokenType) (Token, error) {
	if p.pos >= len(p.tokens) {
		if expected != TokenUnknown {
			return Token{}, &internalCompilerError{
				message: fmt.Sprintf("Expected token %s but reached end of input", expected),
				pos:     p.pos,
			}
		}

		return Token{}, nil
	}

	token := p.tokens[p.pos]

	if expected != TokenUnknown && token.Type != expected {
		return Token{}, &internalCompilerError{
			message: fmt.Sprintf("Expected token %s but got %s", expected, token.Literal),
			pos:     p.pos,
		}
	}

	p.pos++
	return token, nil
}

func (p *Parser) consume() (Token, error) {
	return p.consumeExpected(TokenUnknown)
}

func (p *Parser) parseModule() (shared.Node, error) {
	var children []shared.Node

	for {
		// TODO handle multiple statements
		child, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if child == nil {
			break
		}
	}

	return shared.ModuleNode{Children: children}, nil
}

func (p *Parser) parseStatement() (shared.Node, error) {
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consumeExpected(TokenSemi)
	if err != nil {
		return p.error("Statements must be terminated by a semicolon")
	}

	return n, nil
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

		opToken, err := p.consume()
		if err != nil {
			return nil, err
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

		opToken, err := p.consume()
		if err != nil {
			return nil, err
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

		token, err = p.consumeExpected(TokenLParen)
		if err != nil {
			return nil, err
		}

		args := []shared.Node{}

		// TODO support more than one arg
		for {
			token, ok = p.peek()
			if !ok {
				return p.error("Reached end while parsing call expression")
			}
			if token.Type == TokenRParen {
				break
			}

			arg, err := p.parseExpr()
			if arg == nil || err != nil {
				return nil, err
			}

			args = append(args, arg)

			token, ok = p.peek()
			if !ok {
				return p.error("Reached end while parsing call expression")
			}
			if token.Type != TokenRParen {
				if token.Type == TokenComma {
					p.consumeExpected(TokenComma)
				} else {
					return p.error("Arguments were not separated by a comma in a call expression")
				}
			}
		}

		token, err = p.consumeExpected(TokenRParen)
		if err != nil {
			return nil, err
		}

		left = shared.CallExpr{
			Func: v,
			Args: args,
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

	token, err := p.consumeExpected(TokenIdent)
	if err != nil {
		return nil, err
	}

	return shared.IdentNode{Name: token.Literal, Pos: token.Pos}, nil
}

func (p *Parser) parseNumber() (shared.Node, error) {
	token, ok := p.peek()
	if !ok || token.Type != TokenNumber {
		return nil, nil
	}

	token, err := p.consumeExpected(TokenNumber)
	if err != nil {
		return nil, err
	}

	n, err := strconv.Atoi(token.Literal)
	if err != nil {
		return p.error("Failed to parse int")
	}

	return shared.IntLiteral{Value: n, Pos: token.Pos}, nil
}
