package typechecker

import (
	"fmt"
	"glimmer/ast"
	"glimmer/types"
)

func typeofIfExpression(node *ast.IfExpression, ctx *types.Context) types.TypeNode {
	// get types of all branches
	// error if they do not match
	// return the matched
	branchTypes := []types.TypeNode{}

	for _, stmt := range node.Condition {
		condType := Typeof(stmt, ctx)
		if condType.Type() == types.ERROR {
			return condType
		}
	}
	trueType := Typeof(node.TrueBranch, ctx)
	branchTypes = append(branchTypes, trueType)
	if trueType.Type() == types.ERROR {
		return trueType
	}
	for i, branch := range node.ElifBranches {
		for _, stmt := range node.ElifConditions[i] {
			condType := Typeof(stmt, ctx)
			if condType.Type() == types.ERROR {
				return condType
			}
		}
		elifType := Typeof(branch, ctx)
		branchTypes = append(branchTypes, elifType)
		if elifType.Type() == types.ERROR {
			return elifType
		}
	}
	if node.FalseBranch != nil {
		falseType := Typeof(node.FalseBranch, ctx)
		branchTypes = append(branchTypes, falseType)
		if falseType.Type() == types.ERROR {
			return falseType
		}
	} else {
		branchTypes = append(branchTypes, NONE_T) // nonexistant else is type none
	}

	for _, typ := range branchTypes {
		if typ.String() != branchTypes[0].String() {
			return &types.ErrorType{Msg: "ife branches must match types", Line: node.Token.Line, Col: node.Token.Col}
		}
	}

	return trueType
}

func typeofIndexExpression(node *ast.IndexExpression, ctx *types.Context) types.TypeNode {
	// error if not array or index is not int
	// return inner type of array
	contType := Typeof(node.Left, ctx)

	if contType.Type() != types.ARRAY && contType.Type() != types.DICT {
		return &types.ErrorType{Msg: "indexed type must be array or dict", Line: node.Token.Line, Col: node.Token.Col}
	}

	indexType := Typeof(node.Index, ctx)

	switch typ := contType.(type) {
	case *types.ArrayType:
		if indexType.Type() != types.INTEGER {
			return &types.ErrorType{Msg: "index of array must be int", Line: node.Token.Line, Col: node.Token.Col}
		} else {
			return typ.HeldType
		}
	case *types.DictType:
		if indexType.Type() != types.STRING {
			return &types.ErrorType{Msg: "index of dict must be string", Line: node.Token.Line, Col: node.Token.Col}
		} else {
			return typ.HeldType
		}
	}
	return nil // should never happen, to please the compiler
}

func typeofCallExpression(node *ast.CallExpression, ctx *types.Context) types.TypeNode {
	// return builtinType if function is builtin, else
	// error if not function or params dont match
	// return the ret type
	if fnIdent, ok := node.Function.(*ast.Identifier); ok {
		if _, ok := builtinExists[fnIdent.Value]; ok {
			return typeofBuiltin(node, ctx)
		}
	}

	funTypeNode := Typeof(node.Function, ctx)

	if funTypeNode.Type() == types.ERROR {
		return funTypeNode
	}

	funType, ok := funTypeNode.(*types.FunctionType)
	if !ok || funType.Type() != types.FUNCTION {
		return &types.ErrorType{Msg: "called object must be function", Line: node.Token.Line, Col: node.Token.Col}
	}

	if len(funType.ParamTypes) != len(node.Arguments) {
		return &types.ErrorType{Msg: "invalid number of arguments in call", Line: node.Token.Line, Col: node.Token.Col}
	}

	for idx, pt := range funType.ParamTypes {
		argType := Typeof(node.Arguments[idx], ctx)
		if argType.Type() == types.ERROR {
			return argType
		}
		if argType.String() != pt.String() {
			return &types.ErrorType{Msg: fmt.Sprintf("param type mismatch for param %d in call", idx+1),
				Line: node.Token.Line, Col: node.Token.Col}
		}
	}

	return funType.ReturnType
}

func typeofPrefixExpression(node *ast.PrefixExpression, ctx *types.Context) types.TypeNode {
	// look at operator and types of operands and return the correct type
	switch node.Operator {
	case "!":
		inputType := Typeof(node.Right, ctx)
		if !typeIsNumeric(inputType) {
			return &types.ErrorType{Msg: "input to prefix op '!' must be numeric", Line: node.Token.Line, Col: node.Token.Col}
		}
		return BOOL_T
	case "-":
		inputType := Typeof(node.Right, ctx)
		if !typeIsNumeric(inputType) {
			return &types.ErrorType{Msg: "input to prefix op '-' must be numeric", Line: node.Token.Line, Col: node.Token.Col}
		}
		switch inputType.Type() {
		case "BOOLEAN":
			return INT_T
		case "INTEGER":
			return INT_T
		case "FLOAT":
			return FLOAT_T
		default:
			return nil // never happens
		}
	default:
		return &types.ErrorType{Msg: fmt.Sprintf("prefix operator for %s not found", node.Operator),
			Line: node.Token.Line, Col: node.Token.Col}
	}
}

func typeofInfixExpression(node *ast.InfixExpression, ctx *types.Context) types.TypeNode {
	// look at operator and types of operands and return the correct type
	leftType := Typeof(node.Left, ctx)
	rightType := Typeof(node.Right, ctx)

	switch node.Operator {
	case "+": // defined over numeric types and (string, string)
		if leftType.Type() == types.STRING && rightType.Type() == types.STRING {
			return STRING_T
		} else {
			return typeofNumericOp(node, leftType, rightType, highestPromotion(leftType, rightType))
		}
	case "-": // defined over numeric types and (string, string)
		if leftType.Type() == types.STRING && rightType.Type() == types.STRING {
			return STRING_T
		} else {
			return typeofNumericOp(node, leftType, rightType, highestPromotion(leftType, rightType))
		}
	case "*": // defined over numeric types and (string, string) and (string, int)
		if leftType.Type() == types.STRING && rightType.Type() == types.STRING {
			return STRING_T
		} else if leftType.Type() == types.STRING && rightType.Type() == types.INTEGER {
			return STRING_T
		} else {
			return typeofNumericOp(node, leftType, rightType, highestPromotion(leftType, rightType))
		}
	case "/": // defined over numeric types and (string, string)
		if leftType.Type() == types.STRING && rightType.Type() == types.STRING {
			return STRING_T
		} else {
			return typeofNumericOp(node, leftType, rightType, highestPromotion(leftType, rightType))
		}
	case "<": // defined over numeric types
		return typeofNumericOp(node, leftType, rightType, BOOL_T)
	case ">": // defined over numeric types
		return typeofNumericOp(node, leftType, rightType, BOOL_T)
	case "<=": // defined over numeric types
		return typeofNumericOp(node, leftType, rightType, BOOL_T)
	case ">=": // defined over numeric types
		return typeofNumericOp(node, leftType, rightType, BOOL_T)
	case "==": // defined over numeric types and (string, string)
		if leftType.Type() == types.STRING && rightType.Type() == types.STRING {
			return BOOL_T
		} else {
			return typeofNumericOp(node, leftType, rightType, BOOL_T)
		}
	case "!=": // defined over numeric types and (string, string)
		if leftType.Type() == types.STRING && rightType.Type() == types.STRING {
			return BOOL_T
		} else {
			return typeofNumericOp(node, leftType, rightType, BOOL_T)
		}
	case "&&": // defined over numeric types
		return typeofNumericOp(node, leftType, rightType, BOOL_T)
	case "||": // defined over numeric types
		return typeofNumericOp(node, leftType, rightType, BOOL_T)
	default:
		return &types.ErrorType{Msg: fmt.Sprintf("infix operator for '%s %s %s' not found", leftType.String(),
			node.Operator, rightType.String()), Line: node.Token.Line, Col: node.Token.Col}
	}

}

func typeofNumericOp(node *ast.InfixExpression, left, right types.TypeNode, retType types.TypeNode) types.TypeNode {
	if typeIsNumeric(left) && typeIsNumeric(right) {
		return retType
	} else {
		return &types.ErrorType{Msg: fmt.Sprintf("infix operator for '%s %s %s' not found", left.String(),
			node.Operator, right.String()), Line: node.Token.Line, Col: node.Token.Col}
	}
}

func typeIsNumeric(typ types.TypeNode) bool {
	return typ.Type() == types.BOOLEAN || typ.Type() == types.INTEGER || typ.Type() == types.FLOAT
}

func highestPromotion(typ1, typ2 types.TypeNode) types.TypeNode {
	if typ1.Type() == types.FLOAT || typ2.Type() == types.FLOAT {
		return FLOAT_T
	} else {
		return INT_T
	}
}
