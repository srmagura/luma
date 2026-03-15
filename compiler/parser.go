package compiler

import (
	"fmt"

	"github.com/srmagura/luma/shared"
)

type parser struct {
	tokens []token
	pos    int
}

func newParser(tokens []token) *parser {
	return &parser{tokens: tokens, pos: 0}
}

func parse(tokens []token) (shared.Node, error) {
	p := newParser(tokens)
	ast, err := p.parseModule()

	if err != nil {
		return nil, err
	}
	if ast == nil {
		return p.error("Parsing failed")
	}

	return ast, err
}

func (p *parser) error(message string) (shared.Node, error) {
	return nil, &internalCompilerError{
		message: message,
		pos:     p.pos,
	}
}

func (p *parser) peek() (token, bool) {
	if p.pos >= len(p.tokens) {
		return token{}, false
	}

	return p.tokens[p.pos], true
}

func (p *parser) consumeExpected(expected tokenType) (token, error) {
	if p.pos >= len(p.tokens) {
		if expected != tokenUnknown {
			return token{}, &internalCompilerError{
				message: fmt.Sprintf("Expected token %s but reached end of input", expected),
				pos:     p.pos,
			}
		}

		return token{}, nil
	}

	tok := p.tokens[p.pos]

	if expected != tokenUnknown && tok._type != expected {
		return token{}, &internalCompilerError{
			message: fmt.Sprintf("Expected token %s but got %s", expected, tok.literal),
			pos:     p.pos,
		}
	}

	p.pos++
	return tok, nil
}

func (p *parser) consume() (token, error) {
	return p.consumeExpected(tokenUnknown)
}
