package compiler

import "fmt"

type InternalParserError struct {
	message string
	pos     int
}

type ParserError struct {
	Message string
	Line    int
	Col     int
}

func (e *InternalParserError) Error() string {
	return fmt.Sprintf("%d: %s", e.pos, e.message)
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Col, e.Message)
}
