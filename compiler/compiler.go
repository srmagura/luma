package compiler

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/srmagura/luma/shared"
)

func Compile(src string) (shared.Node, error) {
	src = NormalizeSource(src)
	tokens := lex(src)

	for _, token := range tokens {
		if token.Type == TokenUnknown {
			line, col := GetLineColFromPosition(src, token.Pos)

			return nil, &ParserError{
				Message: fmt.Sprintf("Unknown token: %s", token.Literal),
				Line:    line,
				Col:     col,
			}
		}
	}

	ast, err := Parse(tokens)

	if err != nil {
		internalParserErr, ok := errors.AsType[*InternalParserError](err)
		if !ok {
			log.Fatalln("Could not cast error to ParserError")
		}

		line, col := GetLineColFromPosition(src, internalParserErr.pos)

		return nil, &ParserError{
			Message: internalParserErr.message,
			Line:    line,
			Col:     col,
		}
	}

	return ast, nil
}

func NormalizeSource(src string) string {
	return strings.ReplaceAll(src, "\r\n", "\n")
}

func GetLineColFromPosition(src string, pos int) (int, int) {
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
