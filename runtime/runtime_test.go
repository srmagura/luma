package runtime

import (
	"strings"
	"testing"

	"github.com/srmagura/luma/compiler"
)

func TestWriteToStringBuilder(t *testing.T) {
	src := "print(2 - 3 + 4)"
	expected := "3\n"

	ast, err := compiler.Compile(src)
	if err != nil {
		t.Fatal(err.Error())
	}

	var sb strings.Builder
	err = Execute(ast, &sb)
	if err != nil {
		t.Fatal(err.Error())
	}

	result := sb.String()

	if result != expected {
		t.Fatalf("Result was: %s", result)
	}
}
