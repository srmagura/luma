package main

import (
	"bufio"
	"bytes"
	"testing"
)

func testSuccessfulCompilation(t *testing.T, src string, expected string) {
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	err := Compile(src, out)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = out.Flush()
	if err != nil {
		t.Fatal(err.Error())
	}

	actual := StringifyBytecode(buf.Bytes())
	compareOutputs(t, expected, actual)
}

func TestHelloWorld(t *testing.T) {
	src := "print(1);"
	expected := `
ldc.i4 1
print
	`
	testSuccessfulCompilation(t, src, expected)
}
