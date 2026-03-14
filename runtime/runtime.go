package runtime

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/srmagura/luma/shared"
)

type env struct {
	out io.Writer
}

func Execute(ast shared.Node, out io.Writer) {
	env := &env{out: out}
	env.evalExpr(ast)
}

func (env *env) evalExpr(n shared.Node) int {
	switch v := n.(type) {
	case shared.IntLiteral:
		return v.Value
	case shared.BinaryExpr:
		return env.evalBinaryExpr(v)
	case shared.CallExpr:
		return env.evalCallExpr(v)
	default:
		log.Fatal("eval: Unexpected node type: ", n)
		return 0
	}
}

func (env *env) evalBinaryExpr(n shared.BinaryExpr) int {
	left := env.evalExpr(n.Left)
	right := env.evalExpr(n.Right)

	switch n.Op {
	case shared.OpAdd:
		return left + right
	case shared.OpSubtract:
		return left - right
	case shared.OpMultiply:
		return left * right
		// TODO implement
		//case shared.OpDivide:
		//	return left / right
	case shared.OpDivideInteger:
		return left / right
	default:
		// TODO error handling
		log.Fatal("Unexpected binary operator: ", n.Op)
		return 0
	}
}

func (env *env) evalCallExpr(n shared.CallExpr) int {
	args := make([]int, len(n.Args))

	for i, v := range n.Args {
		args[i] = env.evalExpr(v)
	}

	switch v := n.Func.(type) {
	case shared.IdentNode:
		if v.Name == "print" {
			return env.evalPrint(args)
		}

		// TODO allow other functions
		log.Fatal("print is the only valid function currently")
		return 0
	default:
		// TODO error handling
		log.Fatal("Only identifiers are valid as a function")
		return 0
	}
}

func (env *env) evalPrint(args []int) int {
	var sb strings.Builder

	for i, arg := range args {
		sb.WriteString(strconv.Itoa(arg))
		if i != len(args)-1 {
			sb.WriteString(" ")
		}
	}

	fmt.Fprintln(env.out, sb.String())

	// TODO print should not return anything
	return 0
}
