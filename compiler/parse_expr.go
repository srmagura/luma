package compiler

import (
	"strconv"

	"github.com/srmagura/luma/shared"
)

func (p *parser) parseExpr() (shared.Node, error) {
	return p.parseComparisonExpr()
}

func (p *parser) parseComparisonExpr() (shared.Node, error) {
	left, err := p.parseAdditiveExpr()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek()
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

		var op shared.Op
		switch opTok._type {
		case tokenLAngle:
			op = shared.OpLessThan
		case tokenRAngle:
			op = shared.OpGreaterThan
		case tokenLAngleEq:
			op = shared.OpLessThanEq
		case tokenRAngleEq:
			op = shared.OpGreaterThanEq
		}

		right, err := p.parseAdditiveExpr()
		if err != nil {
			return nil, err
		}

		left = shared.BinaryExpr{Op: op, Left: left, Right: right, Pos: opTok.pos}
	}

	return left, nil
}

// Handle + and -
func (p *parser) parseAdditiveExpr() (shared.Node, error) {
	left, err := p.parseMultiplicativeExpr()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek()
		if !ok || (peeked._type != tokenPlus && peeked._type != tokenMinus) {
			break
		}

		opTok, err := p.consume()
		if err != nil {
			return nil, err
		}

		var op shared.Op
		switch opTok._type {
		case tokenPlus:
			op = shared.OpAdd
		case tokenMinus:
			op = shared.OpSubtract
		}

		right, err := p.parseMultiplicativeExpr()
		if err != nil {
			return nil, err
		}

		left = shared.BinaryExpr{Op: op, Left: left, Right: right, Pos: opTok.pos}
	}

	return left, nil
}

// Handle *, /, and ~/
func (p *parser) parseMultiplicativeExpr() (shared.Node, error) {
	left, err := p.parseNegateExpr()
	if err != nil {
		return nil, err
	}

	for {
		peeked, ok := p.peek()
		if !ok || (peeked._type != tokenStar && peeked._type != tokenFSlash && peeked._type != tokenTildeFSlash) {
			break
		}

		opTok, err := p.consume()
		if err != nil {
			return nil, err
		}

		var op shared.Op
		switch opTok._type {
		case tokenStar:
			op = shared.OpMultiply
		case tokenFSlash:
			op = shared.OpDivide
		case tokenTildeFSlash:
			op = shared.OpDivideInteger
		}

		right, err := p.parseNegateExpr()
		if err != nil {
			return nil, err
		}

		left = shared.BinaryExpr{Op: op, Left: left, Right: right, Pos: opTok.pos}
	}

	return left, nil
}

// Handle +x and -x
func (p *parser) parseNegateExpr() (shared.Node, error) {
	peeked, ok := p.peek()
	if !ok || (peeked._type != tokenPlus && peeked._type != tokenMinus) {
		return p.parsePostfixExpr()
	}

	opTok, err := p.consume()
	if err != nil {
		return nil, err
	}

	var op shared.Op
	switch opTok._type {
	case tokenPlus:
		op = shared.OpPositive
	case tokenMinus:
		op = shared.OpNegate
	}

	value, err := p.parseNegateExpr()
	if err != nil {
		return nil, err
	}

	return shared.UnaryExpr{Op: op, Value: value, Pos: opTok.pos}, nil
}

// Handle x++ and x--
func (p *parser) parsePostfixExpr() (shared.Node, error) {
	value, err := p.parseCallExpr()
	if err != nil {
		return nil, err
	}

	peeked, ok := p.peek()
	if !ok || (peeked._type != tokenPlusPlus && peeked._type != tokenMinusMinus) {
		return value, nil
	}

	opTok, err := p.consume()
	if err != nil {
		return nil, err
	}

	var op shared.Op
	switch opTok._type {
	case tokenPlusPlus:
		op = shared.OpPostfixIncrement
	case tokenMinusMinus:
		op = shared.OpPostfixDecrement
	}

	return shared.UnaryExpr{Op: op, Value: value, Pos: opTok.pos}, nil
}

func (p *parser) parseCallExpr() (shared.Node, error) {
	left, err := p.parseLeaf()
	if err != nil {
		return nil, err
	}

	// Does not handle parenthesized functions (yet?)
	switch v := left.(type) {
	case shared.IdentNode:
		tok, ok := p.peek()
		if !ok || tok._type != tokenLParen {
			// This is just an ident, not a call expr
			return left, nil
		}

		tok, err = p.consumeExpected(tokenLParen)
		if err != nil {
			return nil, err
		}

		args := []shared.Node{}

		for {
			tok, ok = p.peek()
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

			tok, ok = p.peek()
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

		left = shared.CallExpr{
			Func: v,
			Args: args,
			Pos:  v.Pos,
		}
	}

	return left, nil
}

func (p *parser) parseLeaf() (shared.Node, error) {
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

func (p *parser) parseIdent() (shared.Node, error) {
	tok, ok := p.peek()
	if !ok || tok._type != tokenIdent {
		return nil, nil
	}

	tok, err := p.consume()
	if err != nil {
		return nil, err
	}

	return shared.IdentNode{Name: tok.literal, Pos: tok.pos}, nil
}

func (p *parser) parseNumber() (shared.Node, error) {
	tok, ok := p.peek()
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

	return shared.IntLiteral{Value: n, Pos: tok.pos}, nil
}
