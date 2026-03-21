package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatalln("File path must be provided as a command-line argument.")
	}

	srcFilename := args[1]
	srcExt := filepath.Ext(srcFilename)
	outFilename := strings.ReplaceAll(srcFilename, srcExt, ".bin")

	srcBytes, err := os.ReadFile(srcFilename)
	if err != nil {
		log.Fatalln("Failed to read the source file.")
	}

	src := string(srcBytes)

	f, err := os.Create(outFilename)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()

	out := bufio.NewWriter(f)

	err = Compile(src, out)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
