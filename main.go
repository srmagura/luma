package main

import (
	"log"
	"os"

	"github.com/srmagura/luma/compiler"
	"github.com/srmagura/luma/runtime"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatalln("File path must be provided as a command-line argument.")
	}

	srcBytes, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatalln("Failed to read the source file.")
	}

	src := string(srcBytes)

	ast, err := compiler.Compile(src)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = runtime.Execute(ast, os.Stdout)
	if err != nil {
		os.Exit(1)
	}
}
