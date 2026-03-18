package main

func (p *parser) parseStatement() (Node, error) {
	n, err := p.parseDeclarationStatement()
	if err != nil {
		return nil, err
	}
	if n != nil {
		return n, err
	}

	n, err = p.parseAssignmentStatement()
	if err != nil {
		return nil, err
	}
	if n != nil {
		return n, err
	}

	return p.parseExprStatement()
}

var semicolonErrorMessage = "Statements must be terminated by a semicolon"

func (p *parser) parseDeclarationStatement() (Node, error) {
	varTok, ok := p.peek(0)
	if !ok || varTok._type != tokenVar {
		return nil, nil
	}

	_, err := p.consume()
	if err != nil {
		return nil, err
	}

	identTok, err := p.consumeExpected(tokenIdent)
	if err != nil {
		return nil, err
	}

	_, err = p.consumeExpected(tokenEqual)
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consumeExpected(tokenSemi)
	if err != nil {
		return p.error(semicolonErrorMessage)
	}

	return DeclarationStatement{
		Name:  identTok.literal,
		Value: expr,
		Pos:   varTok.pos,
	}, nil
}

func (p *parser) parseAssignmentStatement() (Node, error) {
	identTok, ok := p.peek(0)
	if !ok || identTok._type != tokenIdent {
		return nil, nil
	}

	eqTok, ok := p.peek(1)
	if !ok || eqTok._type != tokenEqual {
		return nil, nil
	}

	_, err := p.consume()
	if err != nil {
		return nil, err
	}

	_, err = p.consume()
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consumeExpected(tokenSemi)
	if err != nil {
		return p.error(semicolonErrorMessage)
	}

	return AssignmentStatement{
		Name:  identTok.literal,
		Value: expr,
		Pos:   identTok.pos,
	}, nil
}

func (p *parser) parseExprStatement() (Node, error) {
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}

	_, err = p.consumeExpected(tokenSemi)
	if err != nil {
		return p.error(semicolonErrorMessage)
	}

	return n, nil
}
