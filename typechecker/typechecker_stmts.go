package typechecker

import (
	"glimmer/ast"
	"glimmer/types"
)

func typeofIfStatement(node *ast.IfStatement, ctx *types.Context) types.TypeNode {
	// get types of all branches
	// error if they contain error
	// return none

	trueType := Typeof(node.TrueBranch, ctx)
	if trueType.Type() == types.ERROR {
		return trueType
	}
	for _, branch := range node.ElifBranches {
		elifType := Typeof(branch, ctx)
		if elifType.Type() == types.ERROR {
			return elifType
		}
	}
	if node.FalseBranch != nil {
		falseType := Typeof(node.FalseBranch, ctx)
		if falseType.Type() == types.ERROR {
			return falseType
		}
	}

	return NONE_T
}

func typeofForStatement(node *ast.ForStatement, ctx *types.Context) types.TypeNode {
	// loopvar is guarenteed to be ID through parser
	// collection must be arr or dict
	// eval body statements and error if they error
	collType := Typeof(node.Collection, ctx)
	if collType.Type() != types.ARRAY && collType.Type() != types.DICT {
		return &types.ErrorType{Msg: "For statements must iterate over a collection",
			Line: node.Token.Line, Col: node.Token.Col}
	}

	if len(node.LoopVars) > 2 {
		return &types.ErrorType{Msg: "For statements must have at most 2 loop variables",
			Line: node.Token.Line, Col: node.Token.Col}
	}

	if collType.Type() == types.ARRAY {
		if len(node.LoopVars) == 1 {
			ctx.Set(node.LoopVars[0].Value, collType.(*types.ArrayType).HeldType)
		} else { // 2
			ctx.Set(node.LoopVars[0].Value, INT_T)
			ctx.Set(node.LoopVars[1].Value, collType.(*types.ArrayType).HeldType)
		}
	} else if collType.Type() == types.DICT {
		ctx.Set(node.LoopVars[0].Value, STRING_T)
		if len(node.LoopVars) > 1 { // len==2
			ctx.Set(node.LoopVars[1].Value, collType.(*types.DictType).HeldType)
		}
	}

	if bt := Typeof(node.Body, ctx); bt.Type() == types.ERROR {
		return bt
	}

	return NONE_T
}

func typeofBlockStatement(node *ast.BlockStatement, ctx *types.Context) types.TypeNode {
	// get type of last statement and all returns
	// error if they dont match, return matched
	if len(node.Statements) == 0 {
		return NONE_T
	}

	retTypes := []types.TypeNode{}
	for i, stmt := range node.Statements {
		stmtType := Typeof(stmt, ctx)

		if stmtType.Type() == types.ERROR {
			return stmtType
		}

		if _, ok := stmt.(*ast.ReturnStatement); ok || (i == len(node.Statements)-1) {
			if ctx.FnType != nil && (stmtType.Type() != (*ctx.FnType).Type()) {
				return &types.ErrorType{Msg: "return type mismatching function type",
					Line: node.Token.Line, Col: node.Token.Col}
			}
			retTypes = append(retTypes, stmtType)
		}
	}
	retTypes = append(retTypes, Typeof(node.Statements[len(node.Statements)-1], ctx))

	for _, ret := range retTypes {
		if ret.String() != retTypes[0].String() {
			return &types.ErrorType{Msg: "block does not have unified return types",
				Line: node.Token.Line, Col: node.Token.Col}
		}
	}

	return retTypes[0]
}
