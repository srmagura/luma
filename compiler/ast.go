package main

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

type BinaryOp byte

const (
	Add BinaryOp = iota
	Subtract
)

type BinaryExpr struct {
	Op    BinaryOp
	Left  Node
	Right Node
}

// --- Implement the sealed interface ---

func (UnknownNode) nodeTag() {}
func (IntLiteral) nodeTag()  {}
func (BinaryExpr) nodeTag()  {}
