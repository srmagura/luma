package main

import (
	"testing"
)

func TestStringifyIL(t *testing.T) {
	actual := StringifyIL([]byte{OpLdcI4, 0x00, 0x00, 0x04, 0xD2})
	expected := "ldc.i4 1234"
	if actual != expected {
		t.Fatalf("Expected: %s    Actual: %s", expected, actual)
	}
}
