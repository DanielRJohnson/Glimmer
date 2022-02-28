package evaluator

import (
	"glimmer/ast"
	"glimmer/object"
)

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalDictLiteral(node *ast.DictLiteral, env *object.Environment) object.Object {
	pairs := make(map[string]object.Object)

	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env)
		keyStr, ok := key.(*object.String)
		if !ok {
			return newError("key is not of type string. got=%s", key.Type())
		}
		if isError(key) {
			return key
		}

		val := Eval(valNode, env)
		if isError(val) {
			return val
		}

		pairs[keyStr.Value] = val
	}
	return &object.Dict{Pairs: pairs}
}
