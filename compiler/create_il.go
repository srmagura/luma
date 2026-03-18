package main

import (
	"io"
)

type ilCreator struct {
	out io.Writer
}

func createIL(ast Node, out io.Writer) error {
	cr := &ilCreator{out}
	return cr.evalNode(ast)
}

func (cr *ilCreator) evalNode(n Node) error {
	switch v := n.(type) {
	case ModuleNode:
		return cr.evalModule(v)
	// case shared.ForBlock:
	// 	return cr.evalForBlock(v)
	// case shared.DeclarationStatement:
	// 	return cr.evalDeclarationStatement(v)
	// case shared.AssignmentStatement:
	// 	return cr.evalAssignmentStatement(v)
	//case CallExpr:
	// 	return cr.evalCallExpr(v)
	// case shared.UnaryExpr:
	// 	return cr.evalUnaryExpr(v)
	// case shared.BinaryExpr:
	// 	return cr.evalBinaryExpr(v)
	// case shared.IdentNode:
	// 	return cr.evalIdent(v)
	case IntLiteral:
		return cr.evalIntLiteral(v)
	default:
		return &internalCompilerError{
			message: "evalNode: Unexpected node type",
			pos:     0, // OK to set position to 0 since this error indicates a bug in luma
		}
	}
}

func (cr *ilCreator) evalModule(n ModuleNode) error {
	// TODO
	return cr.evalNode(n.Children[0])
}

func (cr *ilCreator) evalIntLiteral(n IntLiteral) error {
	return nil
}
