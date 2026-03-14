package runtime

import (
	"strings"
	"testing"

	"github.com/srmagura/luma/compiler"
)

func testExecute(t *testing.T, src string, expected string) {
	ast, err := compiler.Compile(src)
	if err != nil {
		t.Fatal(err.Error())
	}

	var sb strings.Builder
	Execute(ast, &sb)
	result := sb.String()

	if result != expected {
		t.Fatalf("Result was: %s", result)
	}
}

func TestAddition(t *testing.T) {
	testExecute(t, "2 - 3 + 4", "3")
}
