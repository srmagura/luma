package compiler

import (
	"testing"

	"github.com/srmagura/luma/shared"
)

func testSuccessfulCompilation(t *testing.T, src string, expected Node) {
	actual, err := Compile(src)
	if err != nil {
		t.Fatalf("Compilation failed\n%s", err.Error())
	}

	compareASTs(t, expected, actual)
}

func TestEmptyModule(t *testing.T) {
	src := ""
	expected := ModuleNode{
		Children: []Node{},
	}
	testSuccessfulCompilation(t, src, expected)
}

func TestExprStatements(t *testing.T) {
	src := "print(11);\n1*2;"
	expected := ModuleNode{
		Children: []Node{
			CallExpr{
				Func: IdentNode{Name: "print"},
				Args: []Node{IntLiteral{Value: 11}},
			},
			BinaryExpr{
				Op:    shared.OpMultiply,
				Left:  IntLiteral{Value: 1},
				Right: IntLiteral{Value: 2},
			},
		},
	}
	testSuccessfulCompilation(t, src, expected)
}
