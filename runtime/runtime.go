package runtime

import (
	"fmt"
	"io"
	"log"

	"github.com/srmagura/luma/shared"
)

func Execute(ast shared.Node, out io.Writer) {
	i := evalExpr(ast)
	fmt.Fprintln(out, i)
}

func evalExpr(n shared.Node) int {
	switch v := n.(type) {
	case shared.IntLiteral:
		return v.Value
	case shared.BinaryExpr:
		return evalBinaryExpr(v)
	default:
		log.Fatal("eval: Unexpected node type: ", n)
		return 0
	}
}

func evalBinaryExpr(n shared.BinaryExpr) int {
	left := evalExpr(n.Left)
	right := evalExpr(n.Right)

	switch n.Op {
	case shared.OpAdd:
		return left + right
	case shared.OpSubtract:
		return left - right
	default:
		log.Fatal("Unexpected binary operator: ", n.Op)
		return 0
	}
}
