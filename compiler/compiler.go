package compiler

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/srmagura/luma/shared"
)

func Compile(src string) (shared.Node, error) {
	src = normalizeSource(src)
	tokens := lex(src)

	for _, token := range tokens {
		if token._type == tokenUnknown {
			line, col := getLineColFromPosition(src, token.pos)

			return nil, &CompilerError{
				Message: fmt.Sprintf("Unknown token: %s", token.literal),
				Line:    line,
				Col:     col,
			}
		}
	}

	ast, err := parse(tokens)

	if err != nil {
		internalParserErr, ok := errors.AsType[*internalCompilerError](err)
		if !ok {
			log.Fatalln("Could not cast error to internalCompilerError")
		}

		line, col := getLineColFromPosition(src, internalParserErr.pos)

		return nil, &CompilerError{
			Message: internalParserErr.message,
			Line:    line,
			Col:     col,
		}
	}

	return ast, nil
}

func compileCore(src string) (shared.Node, error) {
	src = normalizeSource(src)
	tokens := lex(src)

	for _, tok := range tokens {
		if tok._type == tokenUnknown {
			return nil, &internalCompilerError{
				message: fmt.Sprintf("Unknown token: %s", tok.literal),
				pos:     tok.pos,
			}
		}
	}

	return parse(tokens)
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
