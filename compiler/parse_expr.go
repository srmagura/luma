package compiler

import (
	"strconv"
)

func (p *parser) parseExpr() (Node, error) {
	return p.parseComparisonExpr()
}

func (p *parser) parseComparisonExpr() (Node, error) {
	left, err := p.parseAdditiveExpr()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek(0)
		if !ok || (peeked._type != tokenLAngle &&
			peeked._type != tokenRAngle &&
			peeked._type != tokenLAngleEq &&
			peeked._type != tokenRAngleEq) {
			break
		}

		opTok, err := p.consume()
		if err != nil {
			return nil, err
		}

		var op Op
		switch opTok._type {
		case tokenLAngle:
			op = OpLessThan
		case tokenRAngle:
			op = OpGreaterThan
		case tokenLAngleEq:
			op = OpLessThanEq
		case tokenRAngleEq:
			op = OpGreaterThanEq
		}

		right, err := p.parseAdditiveExpr()
		if err != nil {
			return nil, err
		}

		left = BinaryExpr{Op: op, Left: left, Right: right, Pos: opTok.pos}
	}

	return left, nil
}

// Handle + and -
func (p *parser) parseAdditiveExpr() (Node, error) {
	left, err := p.parseMultiplicativeExpr()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek(0)
		if !ok || (peeked._type != tokenPlus && peeked._type != tokenMinus) {
			break
		}

		opTok, err := p.consume()
		if err != nil {
			return nil, err
		}

		var op Op
		switch opTok._type {
		case tokenPlus:
			op = OpAdd
		case tokenMinus:
			op = OpSubtract
		}

		right, err := p.parseMultiplicativeExpr()
		if err != nil {
			return nil, err
		}

		left = BinaryExpr{Op: op, Left: left, Right: right, Pos: opTok.pos}
	}

	return left, nil
}

// Handle *, /, and ~/
func (p *parser) parseMultiplicativeExpr() (Node, error) {
	left, err := p.parseNegateExpr()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek(0)
		if !ok || (peeked._type != tokenStar && peeked._type != tokenFSlash && peeked._type != tokenTildeFSlash) {
			break
		}

		opTok, err := p.consume()
		if err != nil {
			return nil, err
		}

		var op Op
		switch opTok._type {
		case tokenStar:
			op = OpMultiply
		case tokenFSlash:
			op = OpDivide
		case tokenTildeFSlash:
			op = OpDivideInteger
		}

		right, err := p.parseNegateExpr()
		if err != nil {
			return nil, err
		}

		left = BinaryExpr{Op: op, Left: left, Right: right, Pos: opTok.pos}
	}

	return left, nil
}

// Handle +x and -x
func (p *parser) parseNegateExpr() (Node, error) {
	peeked, ok := p.peek(0)
	if !ok || (peeked._type != tokenPlus && peeked._type != tokenMinus) {
		return p.parsePostfixExpr()
	}

	opTok, err := p.consume()
	if err != nil {
		return nil, err
	}

	var op Op
	switch opTok._type {
	case tokenPlus:
		op = OpPositive
	case tokenMinus:
		op = OpNegate
	}

	value, err := p.parseNegateExpr()
	if err != nil {
		return nil, err
	}

	return UnaryExpr{Op: op, Value: value, Pos: opTok.pos}, nil
}

// Handle x++ and x--
func (p *parser) parsePostfixExpr() (Node, error) {
	value, err := p.parseCallExpr()
	if err != nil {
		return nil, err
	}

	peeked, ok := p.peek(0)
	if !ok || (peeked._type != tokenPlusPlus && peeked._type != tokenMinusMinus) {
		return value, nil
	}

	opTok, err := p.consume()
	if err != nil {
		return nil, err
	}

	var op Op
	switch opTok._type {
	case tokenPlusPlus:
		op = OpPostfixIncrement
	case tokenMinusMinus:
		op = OpPostfixDecrement
	}

	return UnaryExpr{Op: op, Value: value, Pos: opTok.pos}, nil
}

func (p *parser) parseCallExpr() (Node, error) {
	left, err := p.parseLeaf()
	if err != nil {
		return nil, err
	}

	// Does not handle parenthesized functions (yet?)
	switch v := left.(type) {
	case IdentNode:
		tok, ok := p.peek(0)
		if !ok || tok._type != tokenLParen {
			// This is just an ident, not a call expr
			return left, nil
		}

		tok, err = p.consumeExpected(tokenLParen)
		if err != nil {
			return nil, err
		}

		args := []Node{}

		for {
			tok, ok = p.peek(0)
			if !ok {
				return p.error("Reached end while parsing call expression")
			}
			if tok._type == tokenRParen {
				break
			}

			arg, err := p.parseExpr()
			if arg == nil || err != nil {
				return nil, err
			}

			args = append(args, arg)

			tok, ok = p.peek(0)
			if !ok {
				return p.error("Reached end while parsing call expression")
			}
			if tok._type != tokenRParen {
				if tok._type == tokenComma {
					p.consumeExpected(tokenComma)
				} else {
					return p.error("Arguments were not separated by a comma in a call expression")
				}
			}
		}

		tok, err = p.consumeExpected(tokenRParen)
		if err != nil {
			return nil, err
		}

		left = CallExpr{
			Func: v,
			Args: args,
			Pos:  v.Pos,
		}
	}

	return left, nil
}

func (p *parser) parseLeaf() (Node, error) {
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

func (p *parser) parseIdent() (Node, error) {
	tok, ok := p.peek(0)
	if !ok || tok._type != tokenIdent {
		return nil, nil
	}

	tok, err := p.consume()
	if err != nil {
		return nil, err
	}

	return IdentNode{Name: tok.literal, Pos: tok.pos}, nil
}

func (p *parser) parseNumber() (Node, error) {
	tok, ok := p.peek(0)
	if !ok || tok._type != tokenNumber {
		return nil, nil
	}

	tok, err := p.consumeExpected(tokenNumber)
	if err != nil {
		return nil, err
	}

	n, err := strconv.Atoi(tok.literal)
	if err != nil {
		return p.error("Failed to parse int")
	}

	return IntLiteral{Value: n, Pos: tok.pos}, nil
}
