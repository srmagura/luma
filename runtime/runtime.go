package runtime

import (
	"fmt"
	"io"

	"github.com/srmagura/luma/shared"
)

type env struct {
	out       io.Writer
	variables map[string]int
}

func Execute(ast shared.Node, out io.Writer) error {
	env := &env{
		out:       out,
		variables: make(map[string]int),
	}

	_, err := env.evalNode(ast)
	if err != nil {
		fmt.Fprintln(env.out, err.Error())
	}

	return err
}

func (env *env) evalNode(n shared.Node) (int, error) {
	switch v := n.(type) {
	case shared.ModuleNode:
		return env.evalModule(v)
	case shared.AssignmentStatement:
		return env.evalAssignmentStatement(v)
	case shared.CallExpr:
		return env.evalCallExpr(v)
	case shared.BinaryExpr:
		return env.evalBinaryExpr(v)
	case shared.IdentNode:
		return env.evalIdent(v)
	case shared.IntLiteral:
		return v.Value, nil
	default:
		return 0, &internalRuntimeError{
			message: fmt.Sprintf("evalNode: Unexpected node type: %s", n),
			pos:     0, // OK to set position to 0 since this error indicates a bug in luma
		}
	}
}

func (env *env) evalModule(n shared.ModuleNode) (int, error) {
	for _, v := range n.Children {
		_, err := env.evalNode(v)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil // TODO change to returning void
}

func (env *env) evalAssignmentStatement(n shared.AssignmentStatement) (int, error) {
	value, err := env.evalNode(n.Value)
	if err != nil {
		return 0, err
	}

	env.variables[n.Name] = value

	return 0, nil // TODO change to returning void
}

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
