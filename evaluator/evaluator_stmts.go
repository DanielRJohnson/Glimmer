package evaluator

import (
	"glimmer/ast"
	"glimmer/object"
)

func evalIfStatement(is *ast.IfStatement, env *object.Environment) object.Object {
	condition := evalStatements(is.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		tr := Eval(is.TrueBranch, env)
		if isError(tr) || tr.Type() == object.RETURN_VALUE_OBJ {
			return tr
		}
	} else if branch, ok := trueElifBranch_Stmt(is, env); ok {
		elif := Eval(branch, env)
		if isError(elif) || elif.Type() == object.RETURN_VALUE_OBJ {
			return elif
		}
	} else if is.FalseBranch != nil {
		els := Eval(is.FalseBranch, env)
		if isError(els) || els.Type() == object.RETURN_VALUE_OBJ {
			return els
		}
	}
	return NULL
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	evaledCollection := Eval(fs.Collection, env)
	if arr, ok := evaledCollection.(*object.Array); ok {
		return evalForArrayStatement(fs.LoopVars, arr, fs.Body, env)
	} else if dict, ok := evaledCollection.(*object.Dict); ok {
		return evalForDictStatement(fs.LoopVars, dict, fs.Body, env)
	} else {
		return newError("For statement must iterate over collection. got=%T", fs.Collection)
	}
}

func evalForArrayStatement(lvs []*ast.Identifier, arr *object.Array, body *ast.BlockStatement, env *object.Environment) object.Object {
	elemIdx := len(lvs) - 1
	for index, element := range arr.Elements {
		env.Set(lvs[elemIdx].Value, element)

		if len(lvs) > 1 { // len==2
			env.Set(lvs[0].Value, &object.Integer{Value: int64(index)})
		}

		evaledBody := Eval(body, env)
		if isError(evaledBody) || evaledBody.Type() == object.RETURN_VALUE_OBJ {
			return evaledBody
		}
	}
	return NULL
}

func evalForDictStatement(lvs []*ast.Identifier, dict *object.Dict, body *ast.BlockStatement, env *object.Environment) object.Object {
	for key, value := range dict.Pairs {
		env.Set(lvs[0].Value, &object.String{Value: key})

		if len(lvs) > 1 { // len==2
			env.Set(lvs[1].Value, value)
		}

		evaledBody := Eval(body, env)
		if isError(evaledBody) || evaledBody.Type() == object.RETURN_VALUE_OBJ {
			return evaledBody
		}
	}
	return NULL
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	condition := evalStatements(ws.Condition, env)
	if isError(condition) {
		return condition
	}

	for isTruthy(condition) {
		loop := Eval(ws.Body, env)
		if isError(loop) || loop.Type() == object.RETURN_VALUE_OBJ ||
			loop.Type() == object.BREAK_OBJ || loop.Type() == object.CONT_OBJ {
			return loop
		}
		condition = evalStatements(ws.Condition, env)
		if isError(condition) {
			return condition
		}
	}
	return NULL
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}
