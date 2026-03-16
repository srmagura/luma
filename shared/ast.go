package shared

import (
	"fmt"
	"strings"
)

type Op byte

const (
	OpAdd Op = iota
	OpSubtract
	OpMultiply
	OpDivide
	OpDivideInteger
	OpLessThan
	OpGreaterThan
	OpLessThanEq
	OpGreaterThanEq
)

func (op Op) String() string {
	switch op {
	case OpAdd:
		return "+"
	case OpSubtract:
		return "-"
	case OpMultiply:
		return "*"
	case OpDivide:
		return "/"
	case OpDivideInteger:
		return "~/"
	case OpLessThan:
		return "<"
	case OpGreaterThan:
		return ">"
	case OpLessThanEq:
		return "<="
	case OpGreaterThanEq:
		return ">="
	default:
		return "UnknownOp"
	}
}

type Node interface {
	nodeTag()
	//GetPos() int
}

// --- Leaf nodes ---

type IntLiteral struct {
	Value int
	Pos   int
}

type IdentNode struct {
	Name string
	Pos  int
}

// --- Interior nodes ---

type BinaryExpr struct {
	Op    Op
	Left  Node
	Right Node
	Pos   int
}

type CallExpr struct {
	Func Node
	Args []Node
	Pos  int
}

type AssignmentStatement struct {
	// TODO declared type
	Name  string
	Value Node
	Pos   int
}

type ModuleNode struct {
	Children []Node
	Pos      int
}

// --- Implement the sealed interface ---

func (IntLiteral) nodeTag() {}

// func (n IntLiteral) GetPos() int { return n.Pos }
func (IdentNode) nodeTag() {}

// func (n IdentNode) GetPos() int  { return n.Pos }
func (BinaryExpr) nodeTag() {}

// func (n BinaryExpr) GetPos() int { return n.Pos }
func (CallExpr) nodeTag()            {}
func (AssignmentStatement) nodeTag() {}

// func (n CallExpr) GetPos() int   { return n.Pos }
func (ModuleNode) nodeTag() {}

//func (n ModuleNode) GetPos() int { return n.Pos }

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

	case IdentNode:
		fmt.Fprintf(sb, "%s%sIdent(%s)\n", prefix, connector, v.Name)

	case BinaryExpr:
		fmt.Fprintf(sb, "%s%sBinaryExpr(%s)\n", prefix, connector, v.Op)
		stringifyNode(sb, v.Left, childPrefix, false, false)
		stringifyNode(sb, v.Right, childPrefix, false, true)

	case CallExpr:
		fmt.Fprintf(sb, "%s%sCallExpr\n", prefix, connector)
		stringifyNode(sb, v.Func, childPrefix, false, len(v.Args) == 0)

		for i, arg := range v.Args {
			stringifyNode(sb, arg, childPrefix, false, i == len(v.Args)-1)
		}

	case AssignmentStatement:
		fmt.Fprintf(sb, "%s%sAssignmentStatement(%s)\n", prefix, connector, v.Name)
		stringifyNode(sb, v.Value, childPrefix, false, true)

	case ModuleNode:
		fmt.Fprintf(sb, "%s%sModule\n", prefix, connector)

		for i, child := range v.Children {
			stringifyNode(sb, child, childPrefix, false, i == len(v.Children)-1)
		}

	default:
		fmt.Fprintf(sb, "stringifyNode: unknown node\n")
	}
}
