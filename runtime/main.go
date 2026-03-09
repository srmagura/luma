package main

import (
	"bytes"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile("../compiler/a.out")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(data)
	var ast *Node
	DecodeAST(buf, &ast)

	PrintAST(*ast)

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
