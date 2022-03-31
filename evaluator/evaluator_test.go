package evaluator

import (
	"fmt"
	"glimmer/lexer"
	"glimmer/object"
	"glimmer/parser"
	"glimmer/typechecker"
	"glimmer/types"
	"testing"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	ctx := types.NewContext()
	typechecker.Typeof(program, ctx)
	// programType := typechecker.Typeof(program, ctx)
	// if programType.Type() == types.ERROR {
	// 	return &object.Error{Message: programType.(*types.ErrorType).Msg}
	// }

	return Eval(program, env)
}

func testLiteralObject(t *testing.T, obj object.Object, expected interface{}) bool {
	if expected == nil {
		return testNullObject(t, obj)
	}
	switch expected := expected.(type) {
	case string:
		return testStringObject(t, obj, expected)
	case int:
		return testIntegerObject(t, obj, int64(expected))
	case int64:
		return testIntegerObject(t, obj, expected)
	case float64:
		return testFloatObject(t, obj, expected)
	case bool:
		return testBooleanObject(t, obj, expected)
	default:
		t.Fatalf("Literal type %T not covered in testLiteralObject", expected)
		return false
	}
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}

	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x: int) -> int { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	if fn.Body.String() != "{ (x + 2) }" {
		t.Fatalf("body is not %q. got=%q", "{ (x + 2) }", fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"ID = fn(x: int) -> int {x}; ID(5);", 5},
		{"ID = fn(x: int) -> int { return x }; ID(5);", 5},
		{"double = fn(x: int) -> int { x * 2 }; double(5);", 10},
		{"add = fn(x: int, y: int) -> int { x + y }; add(5, 5);", 10},
		{"add = fn(x: int, y: int) -> int { x + y }; add(5 + 5, add(5, 5));", 20},
		{"fn(x: int) -> int { x }(5)", 5},
	}

	for _, tt := range tests {
		testLiteralObject(t, testEval(tt.input), tt.expected)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got=INTEGER"},
		{`len(1, 2)`, "wrong number of arguments. got=2, want=1"},
		{"head([1,2,3,4])", 1},
		{"head([])", NULL},
		{"head(1)", "argument 1 to `head` not supported, got=INTEGER"},
		{"tail([1,2,3,4])", 4},
		{"tail([])", NULL},
		{"tail(1)", "argument 1 to `tail` not supported, got=INTEGER"},
		{"slice([1,2,3,4], 6, 3)", fmt.Sprintf("invalid slice index %d > %d", 6, 3)},
		{"slice([1,2,3,4], -1, 5)", fmt.Sprintf("start index %d out of range for array of length %d", -1, 4)},
		{"slice([1,2,3,4], 1, 5)", fmt.Sprintf("end index %d out of range for array of length %d", 5, 4)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestClosures(t *testing.T) {
	input := `
	newAdder = fn(x: int) -> fn(int) -> int {
		fn(y: int) -> int { x + y };
	};

	addTwo = newAdder(2);
	addTwo(2);`

	testLiteralObject(t, testEval(input), 4)
}

func TestStaticScoping(t *testing.T) {
	input := `
	n = 5;
	addN = fn(x: int) -> int { x + n };
	n = 6;
	addN(5);`

	testLiteralObject(t, testEval(input), 10)
}

func TestRecursion(t *testing.T) {
	input := `
	fib = fn(fibnum: int) -> int {
		ife fibnum == 0 {
			0
		} else ife fibnum == 1 {
			1
		} else {
			return fib(fibnum - 1) + fib(fibnum - 2)
		}
	}
	fib(5)`

	testLiteralObject(t, testEval(input), 5)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.", 5},
		{"10.3", 10.3},
		{"-5.", -5},
		{"-10.3", -10.3},
		{"-10.2 + 4.65565 - 101.3 * 0.25 / 2.56", -15.436928125},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalNumberPromotion(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5 + 2.3", 7.3},
		{"5 - 2.3", 2.7},
		{"5 * 2.3", 11.5},
		{"5 / 2.5", 2.},
		{"true + 1", 2},
		{"true + true", 2},
		{"1 - false", 1},
		{"1.2 * true", 1.2},
		{"1.2 / true", 1.2},
		{"true > false", true},
		{"2 && 3", true},
		{"false || 17", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testLiteralObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},

		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1.", false},
		{"1 > 1.", false},

		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1 <= 1.", true},
		{"1 >= 1.", true},

		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2.", false},
		{"1 != 1.", false},

		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2.) == true", false},
		{"(1 > 2.) == false", true},

		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},

		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestNotOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!0", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!0", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElifElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"ife (true) { 10 }", 10},
		{"ife (false) { 10 }", nil},
		{"ife (1) { 10 }", 10},
		{"ife (0) { 10 }", nil},
		{"ife (1 < 2) { 10 }", 10},
		{"ife (1 > 2) { 10 }", nil},
		{"ife (1 > 2) { 10 } else { 20 }", 20},
		{"ife (1 < 2) { 10 } else { 20 }", 10},

		{"ife (1 < 2) { 10 } else ife (1 < 2) { 20 } else { 30 }", 10},
		{"ife (1 > 2) { 10 } else ife (1 < 2) { 20 } else { 30 }", 20},
		{"ife (1 > 2) { 10 } else ife (1 > 2) { 20 } else { 30 }", 30},

		{"ife (1 > 2) { 10 } else ife (1 < 2) { 20 }", 20},
		{"ife (1 > 2) { 10 } else ife (1 > 2) { 20 }", nil},

		{"ife (1 > 2) { 10 } else ife (1 > 2) { 20 } else ife (1 < 2) { 30 }", 30},
		{"ife (1 > 2) { 10 } else ife (1 > 2) { 20 } else ife (1 > 2) { 30 }", nil},

		{"x = 0; ife x += 1; x == 1 { 30 }", 30},
		{"x = 0; ife x += 1; x != 1 { 30 } else ife x += 1; x != 1 { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestForExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"for x = 0, x < 10, x = x + 1 { }; x", 10},
		{"x = 0; for x < 10, x = x + 1 { }; x", 10},
		{"x = 0; for x < 10 { x = x + 1 }", 10},
		{"for x {}", "identifier not found: x"},
		{"for x = 1, x < 10, x += 1 { break; }; x", 1},
		{"for i = 0; x = 0, i < 10, i += 1 { x += 1; continue; x += 1 }; x", 10},
		{"x = 0; for { x += 1; ife x >= 10 { break } }; x", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5;", 10},
		{"ife 1 { ife 1 { return 10 } return 1; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestAssignExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"x = 5; x = 6; x", 6},
		{"x = 5; x += 6; x", 11},
		{"x = z", "identifier not found: z"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foobar", "identifier not found: foobar"},
		{`"hello" - 4`, "unknown operator: STRING - INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q. got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	testLiteralObject(t, testEval(input), "Hello World!")
}

func TestStringOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello" + " " + "World" + "!"`, "Hello World!"},
		{`"hello" - "el"`, "hlo"},
		{`"hellohello" - "hello"`, "hello"},
		{`"hello" - "io"`, "hello"},
		{`"hello" - ""`, "hello"},
		{`"hi" * "world"`, "hwhohrhlhdiwioirilid"},
		{`"hellohellohello" / "el"`, "hlohlohlo"},
		{`"hellohellohello" / "hello"`, ""},
		{`"hello" / ""`, "hello"},
		{`"hello" * 4`, "hellohellohellohello"},
	}

	for _, tt := range tests {
		testLiteralObject(t, testEval(tt.input), tt.expected)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testLiteralObject(t, result.Elements[0], 1)
	testLiteralObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1,2,3][0]", 1},
		{"i = 2; [1,2,3][i]", 3},
		{"[1,2,3][1 + 1]", 3},
		{"myArray = [1,2,3]; myArray[1]", 2},
		{"[1,2,3][3]", "Index 3 out of range for array of length 3"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testLiteralObject(t, evaluated, int64(integer))
		} else {
			errObj, _ := evaluated.(*object.Error)
			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q. got=%q", tt.expected, errObj.Message)
			}
		}
	}
}

func TestDictLiterals(t *testing.T) {
	input := `two = "two"; {"one": 1, two: 2 + 2}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Dict)
	if !ok {
		t.Fatalf("Eval didn't return Dict. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[string]int64{
		"one": 1,
		"two": 4,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Dict has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expKey, expVal := range expected {
		val, ok := result.Pairs[expKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, val, expVal)
	}
}

func TestDictIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`key = "foo"; {"foo": 5}[key]`, 5},
		{`{"foo": 5}["bar"]`, "key `bar` not found in dict"},
		{`{"foo": 5}[5]`, "index operator not supported: DICT[INTEGER]"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}
