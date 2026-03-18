package compiler

import (
	"fmt"
	"strings"
)

type Op byte

const (
	OpPositive Op = iota
	OpNegate
	OpPostfixIncrement
	OpPostfixDecrement
	OpAdd
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
	case OpPositive:
		return "+"
	case OpNegate:
		return "-"
	case OpPostfixIncrement:
		return "x++"
	case OpPostfixDecrement:
		return "x--"
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

type UnaryExpr struct {
	Op    Op
	Value Node
	Pos   int
}

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

type DeclarationStatement struct {
	Name  string
	Value Node // optional
	Pos   int
}

type AssignmentStatement struct {
	Name  string
	Value Node
	Pos   int
}

type ForBlock struct {
	Statement1 Node // optional
	Expr2      Node // optional
	Expr3      Node // optional
	Children   []Node
	Pos        int
}

type ModuleNode struct {
	Children []Node
	Pos      int
}

// --- Implement the sealed interface ---

func (IntLiteral) nodeTag()           {}
func (IdentNode) nodeTag()            {}
func (UnaryExpr) nodeTag()            {}
func (BinaryExpr) nodeTag()           {}
func (CallExpr) nodeTag()             {}
func (DeclarationStatement) nodeTag() {}
func (AssignmentStatement) nodeTag()  {}
func (ForBlock) nodeTag()             {}
func (ModuleNode) nodeTag()           {}

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

	case UnaryExpr:
		fmt.Fprintf(sb, "%s%sUnaryExpr(%s)\n", prefix, connector, v.Op)
		stringifyNode(sb, v.Value, childPrefix, false, true)

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

	case DeclarationStatement:
		fmt.Fprintf(sb, "%s%sDeclarationStatement(%s)\n", prefix, connector, v.Name)

		if v.Value != nil {
			stringifyNode(sb, v.Value, childPrefix, false, true)
		}

	case AssignmentStatement:
		fmt.Fprintf(sb, "%s%sAssignmentStatement(%s)\n", prefix, connector, v.Name)
		stringifyNode(sb, v.Value, childPrefix, false, true)

	case ForBlock:
		fmt.Fprintf(sb, "%s%sFor\n", prefix, connector)
		stringifyNode(sb, v.Statement1, childPrefix, false, false)
		stringifyNode(sb, v.Expr2, childPrefix, false, false)
		stringifyNode(sb, v.Expr3, childPrefix, false, len(v.Children) == 0)

		for i, child := range v.Children {
			stringifyNode(sb, child, childPrefix, false, i == len(v.Children)-1)
		}

	case ModuleNode:
		fmt.Fprintf(sb, "%s%sModule\n", prefix, connector)

		for i, child := range v.Children {
			stringifyNode(sb, child, childPrefix, false, i == len(v.Children)-1)
		}

	default:
		fmt.Fprintf(sb, "stringifyNode: unknown node\n")
	}
}
