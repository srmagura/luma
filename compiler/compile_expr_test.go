package compiler

import (
	"errors"
	"fmt"
	"testing"
)

func compileExpr(src string) (Node, error) {
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

func testFailedExprCompilation(t *testing.T, src string, expectedMessage string, expectedLine int) {
	_, err := compileExpr(src)

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

func testSuccessfulExprCompilation(t *testing.T, src string, expected Node) {
	actual, err := compileExpr(src)
	if err != nil {
		t.Fatal(err.Error())
	}

	compareASTs(t, expected, actual)
}

func TestInvalidtoken(t *testing.T) {
	src := "1@"
	expectedMessage := "Unknown token: @"
	testFailedExprCompilation(t, src, expectedMessage, 0)
}

func TestIntLiteral(t *testing.T) {
	src := "2"
	expected := IntLiteral{Value: 2}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestIdent(t *testing.T) {
	src := "x"
	expected := IdentNode{Name: "x"}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestPositive(t *testing.T) {
	src := "+1"
	expected := UnaryExpr{
		Op:    OpPositive,
		Value: IntLiteral{Value: 1},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestNegate(t *testing.T) {
	src := "+-1"
	expected := UnaryExpr{
		Op: OpPositive,
		Value: UnaryExpr{
			Op:    OpNegate,
			Value: IntLiteral{Value: 1},
		},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestPostfixIncrement(t *testing.T) {
	src := "x++"
	expected := UnaryExpr{
		Op:    OpPostfixIncrement,
		Value: IdentNode{Name: "x"},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestAddition(t *testing.T) {
	src := "2 - 3 + 4"
	expected := BinaryExpr{
		Op: OpAdd,
		Left: BinaryExpr{
			Op:    OpSubtract,
			Left:  IntLiteral{Value: 2},
			Right: IntLiteral{Value: 3},
		},
		Right: IntLiteral{Value: 4},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestMultiplication(t *testing.T) {
	src := "5 / 2 * 3 ~/ 7"
	expected := BinaryExpr{
		Op: OpDivideInteger,
		Left: BinaryExpr{
			Op: OpMultiply,
			Left: BinaryExpr{
				Op:    OpDivide,
				Left:  IntLiteral{Value: 5},
				Right: IntLiteral{Value: 2},
			},
			Right: IntLiteral{Value: 3},
		},
		Right: IntLiteral{Value: 7},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestCall1(t *testing.T) {
	src := "print()"
	expected := CallExpr{
		Func: IdentNode{Name: "print"},
		Args: []Node{},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestCall2(t *testing.T) {
	src := "print(2 * 3)"
	expected := CallExpr{
		Func: IdentNode{Name: "print"},
		Args: []Node{
			BinaryExpr{
				Op:    OpMultiply,
				Left:  IntLiteral{Value: 2},
				Right: IntLiteral{Value: 3},
			},
		},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestCall3(t *testing.T) {
	src := "print(2 * 3, 4, 5)"
	expected := CallExpr{
		Func: IdentNode{Name: "print"},
		Args: []Node{
			BinaryExpr{
				Op:    OpMultiply,
				Left:  IntLiteral{Value: 2},
				Right: IntLiteral{Value: 3},
			},
			IntLiteral{Value: 4},
			IntLiteral{Value: 5},
		},
	}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestCall4(t *testing.T) {
	src := "print(2 4)"
	expectedMessage := "Arguments were not separated by a comma in a call expression"
	testFailedExprCompilation(t, src, expectedMessage, 0)
}

func TestComparison(t *testing.T) {
	src := "0 * 1 > 2 + 3"
	expected := BinaryExpr{
		Op: OpGreaterThan,
		Left: BinaryExpr{
			Op:    OpMultiply,
			Left:  IntLiteral{Value: 0},
			Right: IntLiteral{Value: 1},
		},
		Right: BinaryExpr{
			Op:    OpAdd,
			Left:  IntLiteral{Value: 2},
			Right: IntLiteral{Value: 3},
		},
	}
	testSuccessfulExprCompilation(t, src, expected)
}
