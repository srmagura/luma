package main

import (
	"log"

	"github.com/srmagura/luma/compiler"
	"github.com/srmagura/luma/runtime"
)

func main() {
	src := "3 + 4"

	ast, err := compiler.Compile(src)
	if err != nil {
		log.Fatalln(err.Error())
	}

	runtime.Execute(ast)
}
