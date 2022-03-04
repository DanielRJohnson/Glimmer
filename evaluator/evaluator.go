package evaluator

import (
	"glimmer/ast"
	"glimmer/object"
)

// flyweights
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	BREAK = &object.Break{}
	CONT  = &object.Continue{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		if fun, ok := val.(*object.Function); ok {
			fun.Env.Set(node.Name.Value, val) // fn name goes in fn's environment, allows recursion
		}
		return val

	case *ast.AssignStatement:
		prevVal, ok := env.Get(node.Name.Value)
		if !ok {
			return newError("identifier not found: %s", node.Name.Value)
		}
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if node.Type != "=" {
			val = evalInfixExpression(string(node.Type[0]), prevVal, val)
		}

		env.Set(node.Name.Value, val)
		if fun, ok := val.(*object.Function); ok {
			fun.Env.Set(node.Name.Value, val) // fn name goes in fn's environment, allows recursion
		}
		return val

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.BreakStatement:
		return BREAK

	case *ast.ContinueStatement:
		return CONT

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env.DeepCopy(), Body: body}
		// deepcopy for static scoping, no copy = dynamic scoping

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.DictLiteral:
		return evalDictLiteral(node, env)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
		return boolToBoolObj(node.Value)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value

		case *object.Error:
			return result
		}
	}

	return result
}
