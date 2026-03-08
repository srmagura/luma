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

	fmt.Println(ast)
}
