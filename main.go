package main

import (
	"github.com/srmagura/luma/compiler"
	"github.com/srmagura/luma/runtime"
)

func main() {
	src := "3 + 4"

	ast := compiler.Compile(src)
	runtime.Execute(ast)
}
