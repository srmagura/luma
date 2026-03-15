package compiler

import (
	"testing"
)

func testSuccessfulCompilation(t *testing.T, src string, expected Node) {
	actual, err := Compile(src)
	if err != nil {
		t.Fatalf("Compilation failed\n%s", err.Error())
	}

	compareASTs(t, expected, actual)
}

func TestCallStatement(t *testing.T) {
	src := "print(11);"
	expected := ModuleNode{
		Children: []Node{
			CallExpr{
				Func: IdentNode{Name: "print"},
				Args: []Node{IntLiteral{Value: 11}},
			},
		},
	}
	testSuccessfulCompilation(t, src, expected)
}
