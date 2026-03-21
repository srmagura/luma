package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type bytecodeCreator struct {
	out io.Writer
}

func createBytecode(ast Node, out io.Writer) error {
	cr := &bytecodeCreator{out}
	return cr.evalNode(ast)
}

func (cr *bytecodeCreator) evalNode(n Node) error {
	switch v := n.(type) {
	case ModuleNode:
		return cr.evalModule(v)
	// case shared.ForBlock:
	// 	return cr.evalForBlock(v)
	// case shared.DeclarationStatement:
	// 	return cr.evalDeclarationStatement(v)
	// case shared.AssignmentStatement:
	// 	return cr.evalAssignmentStatement(v)
	case CallExpr:
		return cr.evalCallExpr(v)
	// case shared.UnaryExpr:
	// 	return cr.evalUnaryExpr(v)
	case BinaryExpr:
		return cr.evalBinaryExpr(v)
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

func (cr *bytecodeCreator) evalModule(n ModuleNode) error {
	// TODO
	return cr.evalNode(n.Children[0])
}

func (cr *bytecodeCreator) evalBinaryExpr(n BinaryExpr) error {
	err := cr.evalNode(n.Left)
	if err != nil {
		return err
	}

	err = cr.evalNode(n.Right)
	if err != nil {
		return err
	}

	switch n.Operator {
	case OperatorAdd:
		cr.out.Write([]byte{OpAdd})
	// case shared.OpSubtract:
	// 	return left - right, nil
	// case shared.OpMultiply:
	// 	return left * right, nil
	// 	// TODO implement
	// 	//case shared.OpDivide:
	// 	//	return left / right
	// case shared.OpDivideInteger:
	// 	return left / right, nil
	// 	// TODO implement booleans
	// case shared.OpLessThan:
	// 	if left < right {
	// 		return 1, nil
	// 	} else {
	// 		return 0, nil
	// 	}
	// case shared.OpGreaterThan:
	// 	if left > right {
	// 		return 1, nil
	// 	} else {
	// 		return 0, nil
	// 	}
	// case shared.OpLessThanEq:
	// 	if left <= right {
	// 		return 1, nil
	// 	} else {
	// 		return 0, nil
	// 	}
	// case shared.OpGreaterThanEq:
	// 	if left >= right {
	// 		return 1, nil
	// 	} else {
	// 		return 0, nil
	// 	}
	default:
		return &internalCompilerError{
			message: fmt.Sprintf("Unexpected binary operator: %s", n.Operator),
			pos:     n.Pos,
		}
	}

	return nil
}

func (cr *bytecodeCreator) evalCallExpr(n CallExpr) error {
	for _, v := range n.Args {
		err := cr.evalNode(v)
		if err != nil {
			return err
		}
	}

	switch v := n.Func.(type) {
	case IdentNode:
		if v.Name == "print" {
			cr.out.Write([]byte{OpPrint})
			return nil
		}

		return &internalCompilerError{
			message: "print is the only valid function currently",
			pos:     v.Pos,
		}
	default:
		return &internalCompilerError{
			message: "Only identifiers are valid as a function",
			pos:     n.Pos,
		}
	}
}

func (cr *bytecodeCreator) evalIntLiteral(n IntLiteral) error {
	code := []byte{OpLdcI4}
	code = binary.BigEndian.AppendUint32(code, uint32(n.Value))
	_, err := cr.out.Write(code)
	return err
}
