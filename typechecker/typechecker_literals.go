package typechecker

import (
	"glimmer/ast"
	"glimmer/types"
)

func typeofFunctionLiteral(node *ast.FunctionLiteral, ctx *types.Context, bindName *string) types.TypeNode {
	// create function type
	// error if param is none
	// error if body does not result in return type
	fun := &types.FunctionType{}

	for _, pt := range node.ParamTypes {
		if pt == NONE_T {
			return &types.ErrorType{Msg: "param can not be none type", Line: node.Token.Line, Col: node.Token.Col}
		}
		fun.ParamTypes = append(fun.ParamTypes, pt)
	}

	fun.ReturnType = node.ReturnType

	fun.FnCtx = types.NewEnclosedContext(ctx.DeepCopy())
	for idx, param := range node.Parameters {
		fun.FnCtx.Set(param.Value, node.ParamTypes[idx])
	}
	if bindName != nil {
		fun.FnCtx.Set(*bindName, fun)
	}

	bodyType := Typeof(node.Body, fun.FnCtx)

	if bodyType.String() != node.ReturnType.String() {
		return &types.ErrorType{Msg: "function body type does not match return type", Line: node.Token.Line, Col: node.Token.Col}
	}

	return fun
}

func typeofArrayLiteral(node *ast.ArrayLiteral, ctx *types.Context) types.TypeNode {
	// create array type
	// error if type mismatch
	arr := &types.ArrayType{}

	if len(node.Elements) == 0 {
		arr.HeldType = node.ExplicitType // i.e. []int
		return arr
	}

	arr.HeldType = Typeof(node.Elements[0], ctx)
	for _, item := range node.Elements {
		if Typeof(item, ctx).String() != arr.HeldType.String() {
			return &types.ErrorType{Msg: "array must have matching types", Line: node.Token.Line, Col: node.Token.Col}
		}
	}

	return arr
}

func typeofDictLiteral(node *ast.DictLiteral, ctx *types.Context) types.TypeNode {
	// create dict type
	// error if type mismatch
	dict := &types.DictType{}

	if len(node.Pairs) == 0 {
		return dict
	}

	firstIter := true
	for _, value := range node.Pairs {
		if firstIter {
			dict.HeldType = Typeof(value, ctx)
			firstIter = false
			continue
		}
		if Typeof(value, ctx).String() != dict.HeldType.String() {
			return &types.ErrorType{Msg: "dict must have matching value types", Line: node.Token.Line, Col: node.Token.Col}
		}
	}

	return dict
}
