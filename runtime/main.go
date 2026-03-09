package runtime

import (
	"fmt"
	"log"

	"github.com/srmagura/luma/shared"
)

type Node = shared.Node
type IntLiteral = shared.IntLiteral

func Execute(ast Node) {
	i := eval(ast)
	fmt.Println(i)
}

func eval(n Node) int {
	switch v := n.(type) {
	case IntLiteral:
		return v.Value
	default:
		log.Fatal("Unexpected node:", v)
		return 0
	}
}
