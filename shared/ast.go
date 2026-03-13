package shared

import (
	"fmt"
	"strings"
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

func (IntLiteral) nodeTag() {}
func (BinaryExpr) nodeTag() {}

// --- Pretty printer: indented tree view ---

func PrintAST(n Node) {
	fmt.Print(StringifyAST(n))
}

func StringifyAST(n Node) string {
	var sb strings.Builder
	stringifyNode(&sb, n, "", true, true)

	return sb.String()
}

func stringifyNode(sb *strings.Builder, n Node, prefix string, isRoot bool, isLast bool) {
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
		fmt.Fprintf(sb, "%s%sIntLiteral(%d)\n", prefix, connector, v.Value)

	case BinaryExpr:
		fmt.Fprintf(sb, "%s%sBinaryExpr(%q)\n", prefix, connector, v.Op)
		stringifyNode(sb, v.Left, childPrefix, false, false)
		stringifyNode(sb, v.Right, childPrefix, false, true)

	default:
		fmt.Fprintf(sb, "UnknownNode\n")
	}
}
