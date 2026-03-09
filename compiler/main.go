package main

import (
	"fmt"
)

func main() {
	src := "3 + 4"

	tokens := Lex(src)
	parser := NewParser(tokens)

	ast, ok := parser.parseNumber()
	if !ok {
		fmt.Println("Parse failed")
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
	PrintAst(ast)
}
