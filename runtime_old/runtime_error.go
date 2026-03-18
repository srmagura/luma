package runtime

import "fmt"

type internalRuntimeError struct {
	message string
	pos     int
}

type RuntimeError struct {
	Message string
	Line    int
	Col     int
}

func (e *internalRuntimeError) Error() string {
	return fmt.Sprintf("%d: %s", e.pos, e.message)
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Col, e.Message)
}
