package compiler

import (
	"testing"
)

func testCore(t *testing.T, src string, expected []string) {
	actual := Lex(src)

	for i := range expected {
		if actual[i].Type == tokenUnknown {
			t.Fatalf("Unknown token: %s", actual[i].Literal)
		}

		if actual[i].Literal != expected[i] {
			t.Fatalf("Expected: %s    Actual: %s", expected[i], actual[i].Literal)
		}
	}

	if len(actual) != len(expected) {
		t.Fatalf("Expected length: %d   Actual length: %d", len(expected), len(actual))
	}
}

func TestLexOperators(t *testing.T) {
	src := "+-*/~/"
	expected := []string{"+", "-", "*", "/", "~/"}
	testCore(t, src, expected)
}

func TestLexIdent(t *testing.T) {
	src := "1_test2+abC a_b"
	expected := []string{"1", "_test2", "+", "abC", "a_b"}
	testCore(t, src, expected)
}

func TestLexDelimiters(t *testing.T) {
	src := ";,()"
	expected := []string{";", ",", "(", ")"}
	testCore(t, src, expected)
}
