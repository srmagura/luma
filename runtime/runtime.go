package runtime

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/srmagura/luma/shared"
)

type env struct {
	out io.Writer
}

func Execute(ast shared.Node, out io.Writer) error {
	env := &env{out: out}

	_, err := env.evalExpr(ast)
	if err != nil {
		fmt.Fprintln(env.out, err.Error())
	}

	return err
}

func (env *env) evalExpr(n shared.Node) (int, error) {

	switch v := n.(type) {
	case shared.IntLiteral:
		return v.Value, nil
	case shared.BinaryExpr:
		return env.evalBinaryExpr(v)
	case shared.CallExpr:
		return env.evalCallExpr(v)
	default:
		return 0, &internalRuntimeError{
			message: fmt.Sprintf("evalExpr: Unexpected node type: %s", n),
			pos:     0, // TODO
		}
	}
}

func (env *env) evalBinaryExpr(n shared.BinaryExpr) (int, error) {
	left, err := env.evalExpr(n.Left)
	if err != nil {
		return 0, err
	}

	right, err := env.evalExpr(n.Right)
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
			pos:     0, // TODO
		}
	}
}

func (env *env) evalCallExpr(n shared.CallExpr) (int, error) {
	args := make([]int, len(n.Args))

	for i, v := range n.Args {
		arg, err := env.evalExpr(v)
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
			pos:     0, // TODO
		}
	}
}

func (env *env) evalPrint(args []int) (int, error) {
	var sb strings.Builder

	for i, arg := range args {
		sb.WriteString(strconv.Itoa(arg))
		if i != len(args)-1 {
			sb.WriteString(" ")
		}
	}

	fmt.Fprintln(env.out, sb.String())

	// TODO print should not return anything
	return 0, nil
}
