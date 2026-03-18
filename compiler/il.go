package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Op byte

const (
	//  ldc.i4 <int32> - Push 4-byte int to the stack
	OpLdcI4 = 0x20
)

func PrintIL(code []byte) {
	fmt.Print(StringifyIL(code))
}

func StringifyIL(code []byte) string {
	var sb strings.Builder

	for i := 0; i < len(code); {
		op := code[i]
		i++

		switch op {
		case OpLdcI4:
			value := int32(binary.BigEndian.Uint32(code[i : i+4]))
			i += 4

			fmt.Fprintf(&sb, "ldc.i4 %d", value)
		default:
			fmt.Fprintf(&sb, "Unknown opcode: %#X", op)
			break
		}
	}

	return sb.String()
}
