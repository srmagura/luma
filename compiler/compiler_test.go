package compiler

import (
	"testing"

	"github.com/srmagura/luma/shared"
)

func areASTsEqual(expected shared.Node, actual shared.Node) bool {
	return false
}

func TestIntLiteral(t *testing.T) {
	actual, ok := Compile("2")
	if !ok {
		t.Error("Compilation failed")
	}

	expected := IntLiteral{Value: 2}
	if !areASTsEqual(expected, actual) {
		t.Log("EXPECTED:\n")

		t.Log("\nACTUAL:\n")
		shared.PrintAST(actual)
		t.Error("ASTs not equal\n\n")
	}
}
