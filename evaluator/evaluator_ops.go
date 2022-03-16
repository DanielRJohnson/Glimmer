package evaluator

import (
	"glimmer/object"
)

func evalNotOperator(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Boolean:
		return boolToBoolObj(!right.Value)
	case *object.Null:
		return TRUE
	case *object.Integer:
		return boolToBoolObj(!intToBool(right.Value))
	case *object.Float:
		return boolToBoolObj(!floatToBool(right.Value))
	default:
		return newError("unknown operator: !%s", right.Type())
	}
}

func evalNegOperator(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.Float:
		return &object.Float{Value: -right.Value}
	case *object.Boolean:
		return &object.Integer{Value: -(boolToInt(right.Value))}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}
