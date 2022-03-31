package evaluator

import (
	"fmt"
	"glimmer/object"
)

var builtins = map[string]*object.Builtin{
	"print": {Fn: func(args ...object.Object) object.Object {
		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}
		return NULL
	}},
	"len": {Fn: func(args ...object.Object) object.Object {
		if err := enforceNumArgs(1, args...); err != nil {
			return err
		}
		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		default:
			return newError("argument to `len` not supported, got=%s", args[0].Type())
		}
	}},
	"head": {Fn: func(args ...object.Object) object.Object {
		if err := enforceNumArgs(1, args...); err != nil {
			return err
		}
		if typeErr := enforceArgType("head", args, object.ARRAY_OBJ); typeErr != nil {
			return typeErr
		}
		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}
		return NULL
	}},
	"tail": {Fn: func(args ...object.Object) object.Object {
		if err := enforceNumArgs(1, args...); err != nil {
			return err
		}
		if typeErr := enforceArgType("tail", args, object.ARRAY_OBJ); typeErr != nil {
			return typeErr
		}
		arr := args[0].(*object.Array)
		length := len(arr.Elements)
		if length > 0 {
			return arr.Elements[length-1]
		}
		return NULL
	}},
	"slice": {Fn: func(args ...object.Object) object.Object {
		if err := enforceNumArgs(3, args...); err != nil {
			return err
		}
		if typeErr := enforceArgType("slice", args, object.ARRAY_OBJ, object.INTEGER_OBJ, object.INTEGER_OBJ); typeErr != nil {
			return typeErr
		}
		arr := args[0].(*object.Array)
		start := int(args[1].(*object.Integer).Value)
		end := int(args[2].(*object.Integer).Value)
		length := len(arr.Elements)

		if start > end {
			return newError("invalid slice index %d > %d", start, end)
		}
		if start < 0 || start >= length {
			return newError("start index %d out of range for array of length %d", start, length)
		}
		if end < 0 || end >= length {
			return newError("end index %d out of range for array of length %d", end, length)
		}
		return &object.Array{Elements: arr.Elements[start:end]}
	}},
	"push": {Fn: func(args ...object.Object) object.Object {
		if err := enforceNumArgs(2, args...); err != nil {
			return err
		}
		if typeErr := enforceArgType("push", args, object.ARRAY_OBJ); typeErr != nil {
			return typeErr
		}
		arr := args[0].(*object.Array)
		return &object.Array{Elements: append(arr.Elements, args[1])}
	}},
	"pop": {Fn: func(args ...object.Object) object.Object {
		if err := enforceNumArgs(1, args...); err != nil {
			return err
		}
		if typeErr := enforceArgType("pop", args, object.ARRAY_OBJ); typeErr != nil {
			return typeErr
		}
		arr := args[0].(*object.Array)
		length := len(arr.Elements)
		return &object.Array{Elements: arr.Elements[0 : length-1]}
	}},
	// TODO: MAP, FILTER, REDUCE
}

func enforceNumArgs(numArgs int, args ...object.Object) *object.Error {
	if numArgs != -1 && len(args) != numArgs {
		return newError("wrong number of arguments. got=%d, want=%d", len(args), numArgs)
	}
	return nil
}

func enforceArgType(fnName string, args []object.Object, types ...object.ObjectType) *object.Error {
	for i := range types {
		if args[i].Type() != types[i] {
			return newError("argument %d to `%s` not supported, got=%s", i+1, fnName, args[i].Type())
		}
	}
	return nil
}
