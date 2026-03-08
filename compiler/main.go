package main

import "fmt"

func main() {
	src := "3 + 4"

	lexer := NewLexer(src)
	for {
		token := lexer.Next()
		fmt.Println(token)

		if token.Type == TokenEOF {
			break
		}
	}
}
