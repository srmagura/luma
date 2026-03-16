package compiler

import (
	"testing"
)

func testSuccessfulLex(t *testing.T, src string, expected []string, expectedTypes []tokenType) {
	actual := lex(src)

	for i := range expected {
		if actual[i]._type == tokenUnknown {
			t.Fatalf("Unknown token: %s", actual[i].literal)
		}

		if actual[i].literal != expected[i] || actual[i]._type != expectedTypes[i] {
			t.Fatalf("Expected: %s (%s)    Actual: %s (%s)", expected[i], expectedTypes[i], actual[i].literal, actual[i]._type)
		}
	}

	if len(actual) != len(expected) {
		t.Fatalf("Expected length: %d   Actual length: %d", len(expected), len(actual))
	}
}

func TestLexUnaryOperators(t *testing.T) {
	src := "+ - ++ --"
	expected := []string{"+", "-", "++", "--"}
	expectedTypes := []tokenType{tokenPlus, tokenMinus, tokenPlusPlus, tokenMinusMinus}
	testSuccessfulLex(t, src, expected, expectedTypes)
}

func TestLexBinaryOperators(t *testing.T) {
	src := "+-*/~/="
	expected := []string{"+", "-", "*", "/", "~/", "="}
	expectedTypes := []tokenType{tokenPlus, tokenMinus, tokenStar, tokenFSlash, tokenTildeFSlash, tokenEqual}
	testSuccessfulLex(t, src, expected, expectedTypes)
}

func TestLexComparisonOperators(t *testing.T) {
	src := "<<= >>="
	expected := []string{"<", "<=", ">", ">="}
	expectedTypes := []tokenType{tokenLAngle, tokenLAngleEq, tokenRAngle, tokenRAngleEq}
	testSuccessfulLex(t, src, expected, expectedTypes)
}

func TestLexIdent(t *testing.T) {
	src := "1_test2+abC a_b"
	expected := []string{"1", "_test2", "+", "abC", "a_b"}
	expectedTypes := []tokenType{tokenNumber, tokenIdent, tokenPlus, tokenIdent, tokenIdent}
	testSuccessfulLex(t, src, expected, expectedTypes)
}

func TestKeyword(t *testing.T) {
	src := "var for"
	expected := []string{"var", "for"}
	expectedTypes := []tokenType{tokenVar, tokenFor}
	testSuccessfulLex(t, src, expected, expectedTypes)
}

func TestLexDelimiters(t *testing.T) {
	src := ";,()"
	expected := []string{";", ",", "(", ")"}
	expectedTypes := []tokenType{tokenSemi, tokenComma, tokenLParen, tokenRParen}
	testSuccessfulLex(t, src, expected, expectedTypes)
}
