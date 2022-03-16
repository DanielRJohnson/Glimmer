package typechecker

import (
	"fmt"
	"glimmer/ast"
	"glimmer/types"
)

func typeofBuiltin(node *ast.CallExpression, ctx *types.Context) types.TypeNode {
	switch node.Function.(*ast.Identifier).Value {
	case "print":
		return &types.NoneType{}
	case "len":
		if len(node.Arguments) != 1 {
			return &types.ErrorType{Msg: fmt.Sprintf("Incorrect num of arguments to len, got=%d", len(node.Arguments)),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		argType := Typeof(node.Arguments[0], ctx)
		if argType.Type() != types.ARRAY && argType.Type() != types.STRING {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument to len must be array or string, got=%s", argType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		return &types.IntegerType{}
	case "head":
		if len(node.Arguments) != 1 {
			return &types.ErrorType{Msg: fmt.Sprintf("Incorrect num of arguments to head, got=%d", len(node.Arguments)),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		argType := Typeof(node.Arguments[0], ctx)
		if argType.Type() != types.ARRAY {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument to head must be array, got=%s", argType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		return argType.(*types.ArrayType).HeldType
	case "tail":
		if len(node.Arguments) != 1 {
			return &types.ErrorType{Msg: fmt.Sprintf("Incorrect num of arguments to tail, got=%d", len(node.Arguments)),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		argType := Typeof(node.Arguments[0], ctx)
		if argType.Type() != types.ARRAY {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument to tail must be array, got=%s", argType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		return argType.(*types.ArrayType).HeldType
	case "slice":
		if len(node.Arguments) != 3 {
			return &types.ErrorType{Msg: fmt.Sprintf("Incorrect num of arguments to slice, got=%d", len(node.Arguments)),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		arrType := Typeof(node.Arguments[0], ctx)
		if arrType.Type() != types.ARRAY {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument 1 to slice must be array, got=%s", arrType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		beginType := Typeof(node.Arguments[1], ctx)
		if beginType.Type() != types.INTEGER {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument 2 to slice must be int, got=%s", beginType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		endType := Typeof(node.Arguments[2], ctx)
		if endType.Type() != types.INTEGER {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument 3 to slice must be int, got=%s", endType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		return arrType
	case "push":
		if len(node.Arguments) != 2 {
			return &types.ErrorType{Msg: fmt.Sprintf("Incorrect num of arguments to push, got=%d", len(node.Arguments)),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		arrType := Typeof(node.Arguments[0], ctx)
		if arrType.Type() != types.ARRAY {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument 1 to push must be array, got=%s", arrType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		pushedType := Typeof(node.Arguments[1], ctx)
		held := arrType.(*types.ArrayType).HeldType
		if pushedType.String() != held.String() {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument 2 to push must be match Argument 1's held type: %s, got=%s",
				held.String(), pushedType.String()), Line: node.Token.Line, Col: node.Token.Col}
		}
		return arrType
	case "pop":
		if len(node.Arguments) != 1 {
			return &types.ErrorType{Msg: fmt.Sprintf("Incorrect num of arguments to pop, got=%d", len(node.Arguments)),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		arrType := Typeof(node.Arguments[0], ctx)
		if arrType.Type() != types.ARRAY {
			return &types.ErrorType{Msg: fmt.Sprintf("Argument 1 to pop must be array, got=%s", arrType.String()),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		return arrType.(*types.ArrayType).HeldType
	}
	return nil
}

// map acting as set
var builtinExists = map[string]bool{
	"print": true,
	"len":   true,
	"head":  true,
	"tail":  true,
	"slice": true,
	"push":  true,
	"pop":   true,
}
