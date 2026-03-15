package compiler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/srmagura/luma/shared"
)

func compileExpr(src string) (Node, error) {
	src = NormalizeSource(src)
	tokens := Lex(src)

	for _, token := range tokens {
		if token.Type == TokenUnknown {
			return nil, &internalCompilerError{
				message: fmt.Sprintf("Unknown token: %s", token.Literal),
				pos:     token.Pos,
			}
		}
	}

	parser := NewParser(tokens)
	return parser.parseExpr()

}

func testFailedExprCompilation(t *testing.T, src string, expectedMessage string, expectedLine int) {
	_, err := compileExpr(src)

	internalParserErr, ok := errors.AsType[*internalCompilerError](err)
	if !ok {
		t.Fatalf("Could not cast error to InternalParserError: %s", err.Error())
	}

	line, col := GetLineColFromPosition(src, internalParserErr.pos)
	parserErr := CompilerError{
		Message: internalParserErr.message,
		Line:    line,
		Col:     col,
	}

	if parserErr.Message != expectedMessage {
		t.Fatalf("Unexpected error message: %s", parserErr.Message)
	}

	if parserErr.Line != expectedLine {
		t.Fatalf("Unexpected line: %d", parserErr.Line)
	}

	// Not bothering to test Col
}

// func testSuccessfulExprCompilation(t *testing.T, src string, expected shared.Node) {
// 	actual, err := Compile(src)
// 	if err != nil {
// 		t.Fatalf("Compilation failed\n%s", err.Error())
// 	}

// 	compareASTs(t, expected, actual)
// }

func testSuccessfulExprCompilation(t *testing.T, src string, expected shared.Node) {
	actual, err := compileExpr(src)
	if err != nil {
		t.Fatal(err.Error())
	}

	compareASTs(t, expected, actual)
}

func TestInvalidToken(t *testing.T) {
	src := "1@"
	expectedMessage := "Unknown token: @"
	testFailedExprCompilation(t, src, expectedMessage, 0)
}

func TestIntLiteral(t *testing.T) {
	src := "2"
	expected := IntLiteral{Value: 2}
	testSuccessfulExprCompilation(t, src, expected)
}

func TestAddition(t *testing.T) {
	src := "2 - 3 + 4"
	expected := BinaryExpr{
		Op: shared.OpAdd,
		Left: BinaryExpr{
			Op:    shared.OpSubtract,
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
		Op: shared.OpDivideInteger,
		Left: BinaryExpr{
			Op: shared.OpMultiply,
			Left: BinaryExpr{
				Op:    shared.OpDivide,
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
				Op:    shared.OpMultiply,
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
				Op:    shared.OpMultiply,
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
