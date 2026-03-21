package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

func Compile(src string, out io.Writer) error {
	src = normalizeSource(src)
	tokens := lex(src)

	for _, token := range tokens {
		if token._type == tokenUnknown {
			line, col := getLineColFromPosition(src, token.pos)

			return &CompilerError{
				Message: fmt.Sprintf("Unknown token: %s", token.literal),
				Line:    line,
				Col:     col,
			}
		}
	}

	ast, err := parse(tokens)

	if err != nil {
		internalCompilerErr, ok := errors.AsType[*internalCompilerError](err)
		if !ok {
			log.Fatalln("Could not cast error to internalCompilerError")
		}

		line, col := getLineColFromPosition(src, internalCompilerErr.pos)

		return &CompilerError{
			Message: internalCompilerErr.message,
			Line:    line,
			Col:     col,
		}
	}

	err = createBytecode(ast, out)
	if err != nil {
		internalCompilerErr, ok := errors.AsType[*internalCompilerError](err)
		if !ok {
			log.Fatalln("Could not cast error to internalCompilerError")
		}

		line, col := getLineColFromPosition(src, internalCompilerErr.pos)

		return &CompilerError{
			Message: internalCompilerErr.message,
			Line:    line,
			Col:     col,
		}
	}

	return nil
}

func normalizeSource(src string) string {
	return strings.ReplaceAll(src, "\r\n", "\n")
}

func getLineColFromPosition(src string, pos int) (int, int) {
	line := 0
	col := 0

	for i := 0; i < len(src); i++ {
		if i == pos {
			return line, col
		}

		if src[i] == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}

	return 0, 0
}
