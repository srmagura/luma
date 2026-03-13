package compiler

import (
	"errors"
	"strings"
	"testing"

	"github.com/srmagura/luma/shared"
)

func testFailedCompilation(t *testing.T, src string, expectedMessage string, expectedLine int) {
	_, err := Compile(src)

	parserErr, ok := errors.AsType[*ParserError](err)
	if !ok {
		t.Fatalf("Could not cast error to ParserError")
	}

	if parserErr.Message != expectedMessage {
		t.Fatalf("Unexpected error message: %s", parserErr.Message)
	}

	if parserErr.Line != expectedLine {
		t.Fatalf("Unexpected line: %d", parserErr.Line)
	}

	// Not bothering to test Col
}

func compareASTs(t *testing.T, expected shared.Node, actual shared.Node) {
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

func testSuccessfulCompilation(t *testing.T, src string, expected shared.Node) {
	actual, err := Compile(src)
	if err != nil {
		t.Fatalf("Compilation failed\n%s", err.Error())
	}

	compareASTs(t, expected, actual)
}

func TestInvalidToken(t *testing.T) {
	src := "1@"
	expectedMessage := "Unknown token: @"
	testFailedCompilation(t, src, expectedMessage, 0)
}

func TestIntLiteral(t *testing.T) {
	src := "2"
	expected := IntLiteral{Value: 2}
	testSuccessfulCompilation(t, src, expected)
}

/* BinaryExpr{
	Op: OpAdd,
	Left: BinaryExpr{
		Op:    OpSubtract,
		Left:  IntLiteral{Value: 1},
		Right: IntLiteral{Value: 2},
	},
	Right: BinaryExpr{
		Op:    OpSubtract,
		Left:  IntLiteral{Value: 1},
		Right: IntLiteral{Value: 2},
	},
}*/
