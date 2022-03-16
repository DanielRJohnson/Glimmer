package typechecker

import (
	"fmt"
	"glimmer/ast"
	"glimmer/types"
)

// flyweights
var (
	INT_T    = &types.IntegerType{}
	FLOAT_T  = &types.FloatType{}
	BOOL_T   = &types.BooleanType{}
	STRING_T = &types.StringType{}
	NONE_T   = &types.NoneType{}
)

func Typeof(node ast.Node, ctx *types.Context) types.TypeNode {
	switch node := node.(type) {
	case *ast.Program:
		return typeofProgram(node, ctx)

	case *ast.ReturnStatement:
		return Typeof(node.ReturnValue, ctx)

	case *ast.AssignStatement:
		var valType types.TypeNode
		if fun, ok := node.Value.(*ast.FunctionLiteral); ok {
			valType = typeofFunctionLiteral(fun, ctx, &node.Name.Value) // handle recursion case
		} else {
			valType = Typeof(node.Value, ctx)
		}

		ctx.Set(node.Name.Value, valType)
		return NONE_T

	case *ast.ExpressionStatement:
		return Typeof(node.Expression, ctx)

	case *ast.BlockStatement:
		return typeofBlockStatement(node, ctx)

	case *ast.IfExpression:
		return typeofIfExpression(node, ctx)

	case *ast.ForExpression:
		return Typeof(node.Body, ctx)

	case *ast.BreakStatement:
		return NONE_T

	case *ast.ContinueStatement:
		return NONE_T

	case *ast.PrefixExpression:
		return typeofPrefixExpression(node, ctx)

	case *ast.InfixExpression:
		return typeofInfixExpression(node, ctx)

	case *ast.CallExpression:
		return typeofCallExpression(node, ctx)

	case *ast.IndexExpression:
		return typeofIndexExpression(node, ctx)

	case *ast.FunctionLiteral:
		return typeofFunctionLiteral(node, ctx, nil)

	case *ast.Identifier:
		typ, ok := ctx.Get(node.Value)
		if !ok {
			return &types.ErrorType{Msg: fmt.Sprintf("identifier not found: %s", node.Value),
				Line: node.Token.Line, Col: node.Token.Col}
		}
		return typ

	case *ast.ArrayLiteral:
		return typeofArrayLiteral(node, ctx)

	case *ast.DictLiteral:
		return typeofDictLiteral(node, ctx)

	case *ast.StringLiteral:
		return STRING_T

	case *ast.IntegerLiteral:
		return INT_T

	case *ast.FloatLiteral:
		return FLOAT_T

	case *ast.Boolean:
		return BOOL_T
	}

	return nil
}

func typeofProgram(program *ast.Program, ctx *types.Context) types.TypeNode {
	var result types.TypeNode

	// TODO: match returns of entire program
	for _, stmt := range program.Statements {
		result = Typeof(stmt, ctx)
		// early exit condition(s)
		switch result := result.(type) {
		case *types.ErrorType:
			return result
		}
	}

	return result
}
