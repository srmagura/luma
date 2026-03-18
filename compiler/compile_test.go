package main

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

func TestEmptyModule(t *testing.T) {
	src := ""
	expected := ModuleNode{
		Children: []Node{},
	}
	testSuccessfulCompilation(t, src, expected)
}

func TestExprStatement(t *testing.T) {
	src := "print(11);\n1*2;"
	expected := ModuleNode{
		Children: []Node{
			CallExpr{
				Func: IdentNode{Name: "print"},
				Args: []Node{IntLiteral{Value: 11}},
			},
			BinaryExpr{
				Operator: OperatorMultiply,
				Left:     IntLiteral{Value: 1},
				Right:    IntLiteral{Value: 2},
			},
		},
	}
	testSuccessfulCompilation(t, src, expected)
}

func TestDeclarationStatementWithInitialValue(t *testing.T) {
	src := "var x = 7 + 8;"
	expected := ModuleNode{
		Children: []Node{
			DeclarationStatement{
				Name: "x",
				Value: BinaryExpr{
					Operator: OperatorAdd,
					Left:     IntLiteral{Value: 7},
					Right:    IntLiteral{Value: 8},
				},
			},
		},
	}
	testSuccessfulCompilation(t, src, expected)
}

func TestAssignmentStatement(t *testing.T) {
	src := "x = 7 + 8;"
	expected := ModuleNode{
		Children: []Node{
			AssignmentStatement{
				Name: "x",
				Value: BinaryExpr{
					Operator: OperatorAdd,
					Left:     IntLiteral{Value: 7},
					Right:    IntLiteral{Value: 8},
				},
			},
		},
	}
	testSuccessfulCompilation(t, src, expected)
}

func TestForBlock(t *testing.T) {
	src := `
	for var i = 0; i < 10; i++ {
		print(i);
	}
	`
	expected := ModuleNode{
		Children: []Node{
			ForBlock{
				Statement1: DeclarationStatement{
					Name:  "i",
					Value: IntLiteral{Value: 0},
				},
				Expr2: BinaryExpr{
					Operator: OperatorLessThan,
					Left:     IdentNode{Name: "i"},
					Right:    IntLiteral{Value: 10},
				},
				Expr3: UnaryExpr{
					Operator: OperatorPostfixIncrement,
					Value:    IdentNode{Name: "i"},
				},
				Children: []Node{
					CallExpr{
						Func: IdentNode{Name: "print"},
						Args: []Node{
							IdentNode{Name: "i"},
						},
					},
				},
			},
		},
	}
	testSuccessfulCompilation(t, src, expected)
}
