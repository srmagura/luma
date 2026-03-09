package compiler

import (
	"log"

	"github.com/srmagura/luma/shared"
)

func Compile(src string) shared.Node {
	tokens := Lex(src)
	ast, ok := Parse(tokens)
	if !ok {
		log.Fatal("Parse failed")
	}

	shared.PrintAST(ast)
	return ast
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
