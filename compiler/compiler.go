package compiler

import (
	"fmt"

	"github.com/srmagura/luma/shared"
)

func Compile(src string) (shared.Node, bool) {
	tokens := Lex(src)
	ast, ok := Parse(tokens)

	if !ok {
		fmt.Println("Compilation failed")
		return shared.UnknownNode{}, false
	}

	shared.PrintAST(ast)
	return ast, true
}

/* BinaryExpr{
	Op: OpAdd,
	Left: BinaryExpr{
		Op:    OpSubtract,
		Left:  IntLiteral{Value: 1},
		Right: IntLiteral{Value: 2},
	},
	Right: BinaryExpr{
		Op:    OpSubtract,
		Left:  IntLiteral{Value: 1},
		Right: IntLiteral{Value: 2},
	},
}*/
