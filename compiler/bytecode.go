package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Op byte

const (
	//  ldc.i4 <int32> - Load constant 4-byte int to the stack
	OpLdcI4 = 0x20

	// print - Print 4-byte int
	Print = 0x02
)

func PrintBytecode(code []byte) {
	fmt.Print(StringifyBytecode(code))
}

func StringifyBytecode(code []byte) string {
	var sb strings.Builder

outer:
	for i := 0; i < len(code); {
		op := code[i]
		i++

		switch op {
		case OpLdcI4:
			value := int32(binary.BigEndian.Uint32(code[i : i+4]))
			i += 4

			fmt.Fprintf(&sb, "ldc.i4 %d\n", value)
		case Print:
			fmt.Fprint(&sb, "print\n")
		default:
			fmt.Fprintf(&sb, "Unknown opcode: %#X\n", op)
			break outer
		}
	}

	return sb.String()
}
