package runtime

import (
	"testing"

	"github.com/srmagura/luma/compiler"
)

func testEvalExpr(t *testing.T, src string, expected int) {
	ast, err := compiler.Compile(src)
	if err != nil {
		t.Fatal(err.Error())
	}

	result := evalExpr(ast)
	if result != expected {
		t.Fatalf("Result was: %d", result)
	}
}

func TestAddition(t *testing.T) {
	testEvalExpr(t, "2 - 3 + 4", 3)
}
