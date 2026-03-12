package main

import (
	"log"

	"github.com/srmagura/luma/compiler"
	"github.com/srmagura/luma/runtime"
)

func main() {
	src := "3 + 4"

	ast, ok := compiler.Compile(src)
	if !ok {
		log.Fatalln("Compilation failed.")
	}

	runtime.Execute(ast)
}
