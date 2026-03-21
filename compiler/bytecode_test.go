package main

import (
	"testing"
)

func TestStringifyBytecode(t *testing.T) {
	actual := StringifyBytecode([]byte{OpLdcI4, 0x00, 0x00, 0x04, 0xD2})
	expected := "ldc.i4 1234\n"
	if actual != expected {
		t.Fatalf("Expected: %s    Actual: %s", expected, actual)
	}
}
