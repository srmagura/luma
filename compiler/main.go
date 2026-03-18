package main

import (
	"log"
	"os"
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

	ast, err := Compile(src)
	if err != nil {
		log.Fatalln(err.Error())
	}

	PrintAST(ast)
}
