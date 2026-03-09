package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func registerGobTypes() {
	gob.Register(IntLiteral{})
}

func EncodeAST(buf *bytes.Buffer, ast *Node) {
	registerGobTypes()

	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(ast)
	if err != nil {
		fmt.Println("Encode error:", err)
	}
}

func DecodeAST(buf *bytes.Buffer, ast **Node) {
	registerGobTypes()

	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(ast)
	if err != nil {
		fmt.Println("Decode error:", err)
	}
}
