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

func evalLoopStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ ||
				rt == object.BREAK_OBJ || rt == object.CONT_OBJ {
				return result
			}
		}
	}

	return result
}
