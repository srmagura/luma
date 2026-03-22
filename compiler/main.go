package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	args := os.Args[1:]

	var srcFilename string

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			srcFilename = arg
			break
		}
	}

	if len(srcFilename) == 0 {
		log.Fatalln("File path must be provided as a command-line argument.")
	}

	shouldOutputAssembly := slices.Contains(args, "--dasm")

	srcExt := filepath.Ext(srcFilename)
	assemblyFilename := strings.ReplaceAll(srcFilename, srcExt, ".asm")
	outFilename := strings.ReplaceAll(srcFilename, srcExt, ".out")

	srcBytes, err := os.ReadFile(srcFilename)
	if err != nil {
		log.Fatalln("Failed to read the source file.")
	}

	src := string(srcBytes)

	outFile, err := os.Create(outFilename)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer outFile.Close()

	out := bufio.NewWriter(outFile)

	err = Compile(src, out)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = out.Flush()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if shouldOutputAssembly {
		bytecode, err := os.ReadFile(outFilename)
		if err != nil {
			log.Fatalln(err.Error())
		}

		assembly := StringifyBytecode(bytecode)

		assemblyFile, err := os.Create(assemblyFilename)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer assemblyFile.Close()

		assemblyFile.WriteString(assembly)
	}
}
