package compiler

import (
	"strings"
	"testing"

	"github.com/srmagura/luma/shared"
)

type (
	Node       = shared.Node
	IntLiteral = shared.IntLiteral
	IdentNode  = shared.IdentNode
	BinaryExpr = shared.BinaryExpr
	CallExpr   = shared.CallExpr
	ModuleNode = shared.ModuleNode
)

func compareASTs(t *testing.T, expected Node, actual Node) {
	expectedString := shared.StringifyAST(expected)
	actualString := shared.StringifyAST(actual)

	t.Logf("EXPECTED:\n%s\n", expectedString)
	t.Logf("ACTUAL:\n%s\n", actualString)

	expectedLines := strings.Split(expectedString, "\n")
	actualLines := strings.Split(actualString, "\n")

	for i := 0; i < min(len(expectedLines), len(actualLines)); i++ {
		if expectedLines[i] != actualLines[i] {
			t.Fatalf("Difference at line %d\n", i+1)
		}
	}

	if len(expectedLines) != len(actualLines) {
		t.Fatalf("Expected had %d lines while actual had %d lines\n", len(expectedLines), len(actualLines))
	}
}
