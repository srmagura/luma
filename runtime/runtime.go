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
	case shared.ForBlock:
		return env.evalForBlock(v)
	case shared.DeclarationStatement:
		return env.evalDeclarationStatement(v)
	case shared.AssignmentStatement:
		return env.evalAssignmentStatement(v)
	case shared.CallExpr:
		return env.evalCallExpr(v)
	case shared.UnaryExpr:
		return env.evalUnaryExpr(v)
	case shared.BinaryExpr:
		return env.evalBinaryExpr(v)
	case shared.IdentNode:
		return env.evalIdent(v)
	case shared.IntLiteral:
		return v.Value, nil
	default:
		return 0, &internalRuntimeError{
			message: "evalNode: Unexpected node type",
			pos:     0, // OK to set position to 0 since this error indicates a bug in luma
		}
	}
}

func (env *env) evalModule(n shared.ModuleNode) (int, error) {
	return env.evalManyBlocks(n.Children)
}

func (env *env) evalManyBlocks(blocks []shared.Node) (int, error) {
	for _, block := range blocks {
		_, err := env.evalNode(block)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil // TODO change to returning void
}

func (env *env) evalForBlock(n shared.ForBlock) (int, error) {
	_, err := env.evalNode(n.Statement1)
	if err != nil {
		return 0, err
	}

	for {
		condition, err := env.evalNode(n.Expr2)
		if err != nil {
			return 0, err
		}
		if condition == 0 {
			break
		}

		env.evalManyBlocks(n.Children)

		_, err = env.evalNode(n.Expr3)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func (env *env) evalDeclarationStatement(n shared.DeclarationStatement) (int, error) {
	value, err := env.evalNode(n.Value)
	if err != nil {
		return 0, err
	}

	env.variables[n.Name] = value

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
