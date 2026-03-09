package main

import (
	"fmt"
)

type Op byte

const (
	OpAdd Op = iota
	OpSubtract
)

func (op Op) String() string {
	switch op {
	case OpAdd:
		return "+"
	case OpSubtract:
		return "+"
	default:
		return "UnknownOp"
	}
}

type Node interface {
	nodeTag()
}

// --- Leaf nodes ---

type UnknownNode struct{}

type IntLiteral struct {
	Value int
	Pos   int
}

// --- Interior nodes ---

type BinaryExpr struct {
	Op    Op
	Left  Node
	Right Node
}

// --- Implement the sealed interface ---

func (UnknownNode) nodeTag() {}
func (IntLiteral) nodeTag()  {}
func (BinaryExpr) nodeTag()  {}

// --- Pretty printer: indented tree view ---

func PrintAST(n Node) {
	printNode(n, "", true, true)
}

func printNode(n Node, prefix string, isRoot bool, isLast bool) {
	connector := ""
	childPrefix := ""

	if !isRoot {
		connector = "├── "
		childPrefix = prefix + "│   "

		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}
	}

	switch v := n.(type) {
	case IntLiteral:
		fmt.Printf("%s%sIntLiteral(%d)\n", prefix, connector, v.Value)

	case BinaryExpr:
		fmt.Printf("%s%sBinaryExpr(%q)\n", prefix, connector, v.Op)
		printNode(v.Left, childPrefix, false, false)
		printNode(v.Right, childPrefix, false, true)

	default:
		fmt.Printf("UnknownNode\n")
	}
}
