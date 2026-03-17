package compiler

import (
	"github.com/srmagura/luma/shared"
)

func (p *parser) parseModule() (shared.Node, error) {
	children, err := p.parseManyBlock()
	if err != nil {
		return nil, err
	}

	return shared.ModuleNode{Children: children, Pos: 0}, nil
}

func (p *parser) parseManyBlock() ([]shared.Node, error) {
	var blocks []shared.Node

	for {
		block, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		if block == nil {
			break
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (p *parser) parseBlock() (shared.Node, error) {
	n, err := p.parseForBlock()
	if err != nil {
		return nil, err
	}
	if n != nil {
		return n, err
	}

	return p.parseStatement()
}

func (p *parser) parseForBlock() (shared.Node, error) {
	forTok, ok := p.peek(0)
	if !ok || forTok._type != tokenFor {
		return nil, nil
	}

	_, err := p.consume()
	if err != nil {
		return nil, err
	}

	statement1, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	expr2, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consumeExpected(tokenSemi)
	if err != nil {
		return nil, err
	}

	expr3, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consumeExpected(tokenLBrace)
	if err != nil {
		return nil, err
	}

	children, err := p.parseManyBlock()

	_, err = p.consumeExpected(tokenRBrace)
	if err != nil {
		return nil, err
	}

	return shared.ForBlock{
		Statement1: statement1,
		Expr2:      expr2,
		Expr3:      expr3,
		Children:   children,
		Pos:        forTok.pos,
	}, nil
}
