package evaluator

import (
	"fmt"
	"glimmer/ast"
	"glimmer/object"
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func boolToBoolObj(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func boolToInt(b bool) int64 {
	if b {
		return 1
	} else {
		return 0
	}
}

func intToBool(i int64) bool {
	return i != 0
}

func floatToBool(f float64) bool {
	return f != 0
}

func promoteToBoolean(obj object.Object) *object.Boolean {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj
	case *object.Integer:
		return &object.Boolean{Value: intToBool(obj.Value)}
	case *object.Float:
		return &object.Boolean{Value: floatToBool(obj.Value)}
	default:
		panic(fmt.Sprintf("Unjust type promotion of %s to bool.", obj.Type()))
	}
}

func promoteToInt(obj object.Object) *object.Integer {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj
	case *object.Boolean:
		return &object.Integer{Value: boolToInt(obj.Value)}
	default:
		panic(fmt.Sprintf("Unjust type promotion of %s to int.", obj.Type()))
	}
}

func promoteToFloat(obj object.Object) *object.Float {
	switch obj := obj.(type) {
	case *object.Float:
		return obj
	case *object.Integer:
		return &object.Float{Value: float64(obj.Value)}
	case *object.Boolean:
		return &object.Float{Value: float64(boolToInt(obj.Value))}
	default:
		panic(fmt.Sprintf("Unjust type promotion of %s to float.", obj.Type()))
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// something is truthy if it is NOT (NULL, FALSE, or 0)
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return obj.Inspect() != "0"
	}
}

func isIntegerOrBoolType(obj object.Object) bool {
	switch obj.(type) {
	case *object.Integer:
		return true
	case *object.Boolean:
		return true
	default:
		return false
	}
}

func isNumericType(obj object.Object) bool {
	switch obj.(type) {
	case *object.Integer:
		return true
	case *object.Float:
		return true
	case *object.Boolean:
		return true
	default:
		return false
	}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError("wrong number of arguments. got=%d, want=%d", len(args), len(fn.Parameters))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}

}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
