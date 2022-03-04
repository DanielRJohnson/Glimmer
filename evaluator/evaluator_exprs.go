package evaluator

import (
	"glimmer/ast"
	"glimmer/object"
	"strings"
)

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.TrueBranch, env)
	} else if branch, ok := trueElifBranch(ie, env); ok {
		return Eval(branch, env)
	} else if ie.FalseBranch != nil {
		return Eval(ie.FalseBranch, env)
	} else {
		return NULL
	}
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	var returnVal object.Object = NULL

	pre := evalStatements(fe.ForPrecondition, env)
	if isError(pre) {
		return pre
	}

	for {
		cond := evalStatements(fe.ForCondition, env)
		if isError(cond) {
			return cond
		}
		if len(fe.ForCondition) > 0 && !isTruthy(cond) {
			return returnVal
		}
		returnVal = evalLoopStatements(fe.Body.Statements, env)
		if isError(returnVal) || returnVal.Type() == object.BREAK_OBJ {
			return returnVal
		}

		post := evalStatements(fe.ForPostcondition, env)
		if isError(post) {
			return post
		}
	}
}

// if a condition is true, trueElifBranch returns (the first true branch, true), else (nil, false)
func trueElifBranch(ie *ast.IfExpression, env *object.Environment) (*ast.BlockStatement, bool) {
	for index, cond := range ie.ElifConditions {
		evaledCond := Eval(cond, env)
		if isTruthy(evaledCond) {
			return ie.ElifBranches[index], true
		}
	}
	return nil, false
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.DICT_OBJ && index.Type() == object.STRING_OBJ:
		return evalDictIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return newError("Index %d out of range for array of length %d", idx, len(arrayObj.Elements))
	}
	return arrayObj.Elements[idx]
}

func evalDictIndexExpression(dict, index object.Object) object.Object {
	dictObj := dict.(*object.Dict)

	key, ok := index.(*object.String)
	if !ok {
		return newError("key is not of type string. got=%s", key.Type())
	}

	val, ok := dictObj.Pairs[key.Value]
	if !ok {
		return newError("key `%s` not found in dict", key.Value)
	}

	return val
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNotOperator(right)
	case "-":
		return evalNegOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// left and right are both integers
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	// left and right are both booleans
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)

	// one of left and right are boolean, the other integer
	case isIntegerOrBoolType(left) && isIntegerOrBoolType(right):
		left = promoteToInt(left)
		right = promoteToInt(right)
		return evalIntegerInfixExpression(operator, left, right)

	// both left and right are floats
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)

	// one of left and right is float, other is int or bool
	case isNumericType(left) && isNumericType(right):
		left = promoteToFloat(left)
		right = promoteToFloat(right)
		return evalFloatInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalStringIntInfixExpression(operator, left, right)

	// TODO: Pipe
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("divide by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return boolToBoolObj(leftVal < rightVal)
	case ">":
		return boolToBoolObj(leftVal > rightVal)
	case "<=":
		return boolToBoolObj(leftVal <= rightVal)
	case ">=":
		return boolToBoolObj(leftVal >= rightVal)
	case "==":
		return boolToBoolObj(leftVal == rightVal)
	case "!=":
		return boolToBoolObj(leftVal != rightVal)
	case "&&":
		leftBool := promoteToBoolean(left)
		rightBool := promoteToBoolean(right)
		return boolToBoolObj(leftBool.Value && rightBool.Value)
	case "||":
		leftBool := promoteToBoolean(left)
		rightBool := promoteToBoolean(right)
		return boolToBoolObj(leftBool.Value || rightBool.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return boolToBoolObj(leftVal < rightVal)
	case ">":
		return boolToBoolObj(leftVal > rightVal)
	case "<=":
		return boolToBoolObj(leftVal <= rightVal)
	case ">=":
		return boolToBoolObj(leftVal >= rightVal)
	case "==":
		return boolToBoolObj(leftVal == rightVal)
	case "!=":
		return boolToBoolObj(leftVal != rightVal)
	case "&&":
		leftBool := promoteToBoolean(left)
		rightBool := promoteToBoolean(right)
		return boolToBoolObj(leftBool.Value && rightBool.Value)
	case "||":
		leftBool := promoteToBoolean(left)
		rightBool := promoteToBoolean(right)
		return boolToBoolObj(leftBool.Value || rightBool.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "+":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return &object.Integer{Value: leftInt.Value + rightInt.Value}
	case "-":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return &object.Integer{Value: leftInt.Value - rightInt.Value}
	case "*":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return &object.Integer{Value: leftInt.Value * rightInt.Value}
	case "/":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return &object.Integer{Value: leftInt.Value / rightInt.Value}
	case "<":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return boolToBoolObj(leftInt.Value < rightInt.Value)
	case ">":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return boolToBoolObj(leftInt.Value > rightInt.Value)
	case "<=":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return boolToBoolObj(leftInt.Value <= rightInt.Value)
	case ">=":
		leftInt := promoteToInt(left)
		rightInt := promoteToInt(right)
		return boolToBoolObj(leftInt.Value >= rightInt.Value)
	case "==":
		return boolToBoolObj(leftVal == rightVal)
	case "!=":
		return boolToBoolObj(leftVal != rightVal)
	case "&&":
		return boolToBoolObj(leftVal && rightVal)
	case "||":
		return boolToBoolObj(leftVal || rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "-":
		return &object.String{Value: strings.Replace(leftVal, rightVal, "", 1)}
	case "*":
		crossProduct := ""
		for _, leftChar := range leftVal {
			for _, rightChar := range rightVal {
				crossProduct += string(leftChar) + string(rightChar)
			}
		}
		return &object.String{Value: crossProduct}
	case "/":
		return &object.String{Value: strings.Replace(leftVal, rightVal, "", -1)}
	case "==":
		return boolToBoolObj(leftVal == rightVal)
	case "!=":
		return boolToBoolObj(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringIntInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "*":
		return &object.String{Value: strings.Repeat(leftVal, int(rightVal))}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
