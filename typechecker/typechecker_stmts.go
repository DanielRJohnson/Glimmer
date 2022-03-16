package typechecker

import (
	"glimmer/ast"
	"glimmer/types"
)

func typeofBlockStatement(node *ast.BlockStatement, ctx *types.Context) types.TypeNode {
	// get type of last statement and all returns
	// error if they dont match, return matched
	if len(node.Statements) == 0 {
		return NONE_T
	}

	retTypes := []types.TypeNode{}
	for _, stmt := range node.Statements {
		stmtType := Typeof(stmt, ctx)

		if stmtType.Type() == types.ERROR {
			return stmtType
		}

		if _, ok := stmt.(*ast.ReturnStatement); ok {
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
