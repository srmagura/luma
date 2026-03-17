package compiler

import (
	"github.com/srmagura/luma/shared"
)

func (p *parser) parseModule() (shared.Node, error) {
	var children []shared.Node

	for {
		child, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		if child == nil {
			break
		}

		children = append(children, child)
	}

	return shared.ModuleNode{Children: children, Pos: 0}, nil
}

func (p *parser) parseBlock() (shared.Node, error) {
	return p.parseStatement()
}
