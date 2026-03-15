package compiler

import "fmt"

type internalCompilerError struct {
	message string
	pos     int
}

type CompilerError struct {
	Message string
	Line    int
	Col     int
}

func (e *internalCompilerError) Error() string {
	return fmt.Sprintf("%d: %s", e.pos, e.message)
}

func (e *CompilerError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Col, e.Message)
}
