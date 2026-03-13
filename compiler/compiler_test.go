package compiler

import (
	"strings"
	"testing"

	"github.com/srmagura/luma/shared"
)

func diffASTs(t *testing.T, expected shared.Node, actual shared.Node) {
	expectedString := shared.StringifyAST(expected)
	actualString := shared.StringifyAST(actual)

	t.Logf("EXPECTED:\n%s\n", expectedString)
	t.Logf("ACTUAL:\n%s\n", actualString)

	expectedLines := strings.Split(expectedString, "\n")
	actualLines := strings.Split(actualString, "\n")

	for i := 0; i < min(len(expectedLines), len(actualLines)); i++ {
		if expectedLines[i] != actualLines[i] {
			t.Fatalf("Difference at line %d\n", i)
		}
	}

	if len(expectedLines) != len(actualLines) {
		t.Fatalf("Expected had %d lines while actual had %d lines\n", len(expectedLines), len(actualLines))
	}
}

func TestIntLiteral(t *testing.T) {
	actual, ok := Compile("2")
	if !ok {
		t.Error("Compilation failed")
	}

	expected := IntLiteral{Value: 2}
	diffASTs(t, expected, actual)
}
