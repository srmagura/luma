package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func compareOutputs(t *testing.T, expected string, actual string) {
	expected = strings.Trim(expected, " \n\r\t")
	actual = strings.Trim(actual, " \n\r\t")

	t.Logf("EXPECTED:\n%s\n", expected)
	t.Logf("ACTUAL:\n%s\n", actual)

	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	for i := 0; i < min(len(expectedLines), len(actualLines)); i++ {
		if expectedLines[i] != actualLines[i] {
			t.Fatalf("Difference at line %d\n", i+1)
		}
	}

	if len(expectedLines) != len(actualLines) {
		t.Fatalf("Expected had %d lines while actual had %d lines\n", len(expectedLines), len(actualLines))
	}
}

func compareASTs(t *testing.T, expected Node, actual Node) {
	expectedString := StringifyAST(expected)
	actualString := StringifyAST(actual)
	compareOutputs(t, expectedString, actualString)
}

func parseExpr(src string) (Node, error) {
	src = normalizeSource(src)
	tokens := lex(src)

	for _, token := range tokens {
		if token._type == tokenUnknown {
			return nil, &internalCompilerError{
				message: fmt.Sprintf("Unknown token: %s", token.literal),
				pos:     token.pos,
			}
		}
	}

	parser := newParser(tokens)
	return parser.parseExpr()

}

func testFailedExprParse(t *testing.T, src string, expectedMessage string, expectedLine int) {
	_, err := parseExpr(src)

	internalParserErr, ok := errors.AsType[*internalCompilerError](err)
	if !ok {
		t.Fatalf("Could not cast error to InternalParserError: %s", err.Error())
	}

	line, _ := getLineColFromPosition(src, internalParserErr.pos)
	parserErr := CompilerError{
		Message: internalParserErr.message,
		Line:    line,
	}

	if parserErr.Message != expectedMessage {
		t.Fatalf("Unexpected error message: %s", parserErr.Message)
	}

	if parserErr.Line != expectedLine {
		t.Fatalf("Unexpected line: %d", parserErr.Line)
	}

	// Not bothering to test Col
}

func testSuccessfulExprParse(t *testing.T, src string, expected Node) {
	actual, err := parseExpr(src)
	if err != nil {
		t.Fatal(err.Error())
	}

	compareASTs(t, expected, actual)
}

func TestParseInvalidToken(t *testing.T) {
	src := "1@"
	expectedMessage := "Unknown token: @"
	testFailedExprParse(t, src, expectedMessage, 0)
}

func TestParseIntLiteral(t *testing.T) {
	src := "2"
	expected := IntLiteral{Value: 2}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseIdent(t *testing.T) {
	src := "x"
	expected := IdentNode{Name: "x"}
	testSuccessfulExprParse(t, src, expected)
}

func TestParsePositive(t *testing.T) {
	src := "+1"
	expected := UnaryExpr{
		Operator: OperatorPositive,
		Value:    IntLiteral{Value: 1},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseNegate(t *testing.T) {
	src := "+-1"
	expected := UnaryExpr{
		Operator: OperatorPositive,
		Value: UnaryExpr{
			Operator: OperatorNegate,
			Value:    IntLiteral{Value: 1},
		},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParsePostfixIncrement(t *testing.T) {
	src := "x++"
	expected := UnaryExpr{
		Operator: OperatorPostfixIncrement,
		Value:    IdentNode{Name: "x"},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseAddition(t *testing.T) {
	src := "2 - 3 + 4"
	expected := BinaryExpr{
		Operator: OperatorAdd,
		Left: BinaryExpr{
			Operator: OperatorSubtract,
			Left:     IntLiteral{Value: 2},
			Right:    IntLiteral{Value: 3},
		},
		Right: IntLiteral{Value: 4},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseMultiplication(t *testing.T) {
	src := "5 / 2 * 3 ~/ 7"
	expected := BinaryExpr{
		Operator: OperatorDivideInteger,
		Left: BinaryExpr{
			Operator: OperatorMultiply,
			Left: BinaryExpr{
				Operator: OperatorDivide,
				Left:     IntLiteral{Value: 5},
				Right:    IntLiteral{Value: 2},
			},
			Right: IntLiteral{Value: 3},
		},
		Right: IntLiteral{Value: 7},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseCall1(t *testing.T) {
	src := "print()"
	expected := CallExpr{
		Func: IdentNode{Name: "print"},
		Args: []Node{},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseCall2(t *testing.T) {
	src := "print(2 * 3)"
	expected := CallExpr{
		Func: IdentNode{Name: "print"},
		Args: []Node{
			BinaryExpr{
				Operator: OperatorMultiply,
				Left:     IntLiteral{Value: 2},
				Right:    IntLiteral{Value: 3},
			},
		},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseCall3(t *testing.T) {
	src := "print(2 * 3, 4, 5)"
	expected := CallExpr{
		Func: IdentNode{Name: "print"},
		Args: []Node{
			BinaryExpr{
				Operator: OperatorMultiply,
				Left:     IntLiteral{Value: 2},
				Right:    IntLiteral{Value: 3},
			},
			IntLiteral{Value: 4},
			IntLiteral{Value: 5},
		},
	}
	testSuccessfulExprParse(t, src, expected)
}

func TestParseCall4(t *testing.T) {
	src := "print(2 4)"
	expectedMessage := "Arguments were not separated by a comma in a call expression"
	testFailedExprParse(t, src, expectedMessage, 0)
}

func TestParseComparison(t *testing.T) {
	src := "0 * 1 > 2 + 3"
	expected := BinaryExpr{
		Operator: OperatorGreaterThan,
		Left: BinaryExpr{
			Operator: OperatorMultiply,
			Left:     IntLiteral{Value: 0},
			Right:    IntLiteral{Value: 1},
		},
		Right: BinaryExpr{
			Operator: OperatorAdd,
			Left:     IntLiteral{Value: 2},
			Right:    IntLiteral{Value: 3},
		},
	}
	testSuccessfulExprParse(t, src, expected)
}
