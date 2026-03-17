package runtime

import (
	"fmt"

	"github.com/srmagura/luma/shared"
)

func (env *env) evalCallExpr(n shared.CallExpr) (int, error) {
	args := make([]int, len(n.Args))

	for i, v := range n.Args {
		arg, err := env.evalNode(v)
		if err != nil {
			return 0, err
		}

		args[i] = arg
	}

	switch v := n.Func.(type) {
	case shared.IdentNode:
		if v.Name == "print" {
			return env.evalPrint(args)
		}

		// TODO allow other functions
		return 0, &internalRuntimeError{
			message: "print is the only valid function currently",
			pos:     v.Pos,
		}
	default:
		return 0, &internalRuntimeError{
			message: "Only identifiers are valid as a function",
			pos:     n.Pos,
		}
	}
}

func (env *env) evalUnaryExpr(n shared.UnaryExpr) (int, error) {
	switch n.Op {
	case shared.OpPositive:
		value, err := env.evalNode(n.Value)
		if err != nil {
			return 0, err
		}

		return value, nil
	case shared.OpNegate:
		value, err := env.evalNode(n.Value)
		if err != nil {
			return 0, err
		}

		return -value, nil
	case shared.OpPostfixIncrement:
		return env.evalPostfix(n)
	case shared.OpPostfixDecrement:
		return env.evalPostfix(n)
	default:
		return 0, &internalRuntimeError{
			message: fmt.Sprintf("Unexpected unary operator: %s", n.Op),
			pos:     n.Pos,
		}
	}
}

func (env *env) evalPostfix(n shared.UnaryExpr) (int, error) {
	switch v := n.Value.(type) {
	case shared.IdentNode:
		value := env.variables[v.Name]

		switch n.Op {
		case shared.OpPostfixIncrement:
			env.variables[v.Name]++
		case shared.OpPostfixDecrement:
			env.variables[v.Name]--
		}

		return value, nil
	default:
		return 0, &internalRuntimeError{
			message: "Postfix operators can only be applied to identifiers",
			pos:     n.Pos,
		}
	}
}

func (env *env) evalBinaryExpr(n shared.BinaryExpr) (int, error) {
	left, err := env.evalNode(n.Left)
	if err != nil {
		return 0, err
	}

	right, err := env.evalNode(n.Right)
	if err != nil {
		return 0, err
	}

	switch n.Op {
	case shared.OpAdd:
		return left + right, nil
	case shared.OpSubtract:
		return left - right, nil
	case shared.OpMultiply:
		return left * right, nil
		// TODO implement
		//case shared.OpDivide:
		//	return left / right
	case shared.OpDivideInteger:
		return left / right, nil
		// TODO implement booleans
	case shared.OpLessThan:
		if left < right {
			return 1, nil
		} else {
			return 0, nil
		}
	case shared.OpGreaterThan:
		if left > right {
			return 1, nil
		} else {
			return 0, nil
		}
	case shared.OpLessThanEq:
		if left <= right {
			return 1, nil
		} else {
			return 0, nil
		}
	case shared.OpGreaterThanEq:
		if left >= right {
			return 1, nil
		} else {
			return 0, nil
		}
	default:
		return 0, &internalRuntimeError{
			message: fmt.Sprintf("Unexpected binary operator: %s", n.Op),
			pos:     n.Pos,
		}
	}
}

func (env *env) evalIdent(n shared.IdentNode) (int, error) {
	value, ok := env.variables[n.Name]
	if !ok {
		return 0, &internalRuntimeError{
			message: fmt.Sprintf("Undefined variable: %s", n.Name),
			pos:     n.Pos,
		}
	}

	return value, nil
}
