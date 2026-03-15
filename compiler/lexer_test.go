package compiler

import (
	"testing"
)

func testCore(t *testing.T, src string, expected []string) {
	actual := lex(src)

	for i := range expected {
		if actual[i]._type == tokenUnknown {
			t.Fatalf("Unknown token: %s", actual[i].literal)
		}

		if actual[i].literal != expected[i] {
			t.Fatalf("Expected: %s    Actual: %s", expected[i], actual[i].literal)
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
