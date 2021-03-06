package typechecker

import (
	"glimmer/lexer"
	"glimmer/parser"
	"glimmer/types"
	"testing"
)

func CheckParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestTypeofErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{"x", "Static TypeError at [1,2]: identifier not found: x"},
		{"x = 2 + []bool", "Static TypeError at [1,7]: infix operator for 'int + array[bool]' not found"},
		{"[1, 2, 3, 4.2]", "Static TypeError at [1,1]: array must have matching types"},
		{`{"a": 1, "b": fn(x: int) -> int { x }}`, "Static TypeError at [1,1]: dict must have matching value types"},
		{"fn(x: none, y: int) -> none { }", "Static TypeError at [1,3]: param can not be none type"},
		{"ife true { 1 } else { 2.2 }", "Static TypeError at [1,4]: ife branches must match types"},
		{"if true { x = 2 + fn(x: int) -> int { x + 1} }", "Static TypeError at [1,17]: infix operator for 'int + fn(int) -> int' not found"},
		{"arr = [1,2,3,4]; arr[3.2];", "Static TypeError at [1,21]: index of array must be int"},
		{`dic = {"a": 1}; dic[3.2];`, "Static TypeError at [1,20]: index of dict must be string"},
		{"fn(a: int, b: int) -> int { ife true { false } else { false } }", "Static TypeError at [1,38]: return type mismatching function type"},
		{"fn() -> int { 1 }(true)", "Static TypeError at [1,18]: invalid number of arguments in call"},
		{"fn(x: int) -> int { x } (false)", "Static TypeError at [1,25]: param type mismatch for param 1 in call"},
		{"-[1,2,3,4]", "Static TypeError at [1,1]: input to prefix op '-' must be numeric"},
		{"![1,2,3,4]", "Static TypeError at [1,1]: input to prefix op '!' must be numeric"},
		{"[]int + []int", "Static TypeError at [1,7]: infix operator for 'array[int] + array[int]' not found"},
		{"len(1, 2)", "Static TypeError at [1,4]: Incorrect num of arguments to len, got=2"},
		{"len(1)", "Static TypeError at [1,4]: Argument to len must be array or string, got=int"},
		{"head(1, 2)", "Static TypeError at [1,5]: Incorrect num of arguments to head, got=2"},
		{"head(1)", "Static TypeError at [1,5]: Argument to head must be array, got=int"},
		{"tail(1, 2)", "Static TypeError at [1,5]: Incorrect num of arguments to tail, got=2"},
		{"tail(1)", "Static TypeError at [1,5]: Argument to tail must be array, got=int"},
		{"range(1, 2, 3, 4)", "Static TypeError at [1,6]: Incorrect num of arguments to range, got=4"},
		{"range(3.3)", "Static TypeError at [1,6]: Argument 1 to range must be int, got=float"},
		{"x = [1,2,3,4,5]; slice(x)", "Static TypeError at [1,23]: Incorrect num of arguments to slice, got=1"},
		{"x = [1,2,3,4,5]; slice(1, 2, 3)", "Static TypeError at [1,23]: Argument 1 to slice must be array, got=int"},
		{"x = [1,2,3,4,5]; slice(x, true, 3)", "Static TypeError at [1,23]: Argument 2 to slice must be int, got=bool"},
		{"x = [1,2,3,4,5]; slice(x, 2, true)", "Static TypeError at [1,23]: Argument 3 to slice must be int, got=bool"},
		{"push()", "Static TypeError at [1,5]: Incorrect num of arguments to push, got=0"},
		{"push(1, 2)", "Static TypeError at [1,5]: Argument 1 to push must be array, got=int"},
		{"push([1,2,3], true)", "Static TypeError at [1,5]: Argument 2 to push must be match Argument 1's held type: int, got=bool"},
		{"pop(1, 2)", "Static TypeError at [1,4]: Incorrect num of arguments to pop, got=2"},
		{"pop(1)", "Static TypeError at [1,4]: Argument 1 to pop must be array, got=int"},
		{"a = fn() -> int { return 3.3; 1 } ()", "Static TypeError at [1,17]: return type mismatching function type"},
		{"a = fn() -> int { return 1; 3.3 } ()", "Static TypeError at [1,17]: return type mismatching function type"},
		{"for i, v, k in [1] {}", "Static TypeError at [1,4]: For statements must have at most 2 loop variables"},
		{"while 1 > []int {}", "Static TypeError at [1,9]: infix operator for 'int > array[int]' not found"},
		{"if 1 > []int {}", "Static TypeError at [1,6]: infix operator for 'int > array[int]' not found"},
		{"ife 1 > []int {}", "Static TypeError at [1,7]: infix operator for 'int > array[int]' not found"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != types.ERROR {
			t.Errorf("pType is not ERROR, got=%s", pType.Type())
		}

		if pType.String() != tt.expectedMsg {
			t.Errorf("error string does not match. want=%s, got=%s", tt.expectedMsg, pType.String())
		}
	}
}

func TestTypeofBasicStatements(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   types.GlimmerType
		expectedString string
	}{
		{"1", "INTEGER", "int"},
		{"2.2", "FLOAT", "float"},
		{"true", "BOOLEAN", "bool"},
		{`"hello"`, "STRING", "string"},
		{"x = 5; x", "INTEGER", "int"},
		{"for i in [1,2,3,4,5] { break }", "NONE", "none"},
		{"x = 5", "NONE", "none"},
		{"return 5;", "INTEGER", "int"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != tt.expectedType {
			t.Errorf("pType is not %s, got=%s", tt.expectedString, pType.Type())
		}

		if pType.String() != tt.expectedString {
			t.Errorf("type string does not match. want=%s, got=%s", tt.expectedString, pType.String())
		}
	}
}

func TestTypeofForStatements(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   types.GlimmerType
		expectedString string
	}{
		{"for i, val in [1,2,3,4] { 1 }", "NONE", "none"},
		{"arr = [1,2,3,4]; for val in arr { 1 }", "NONE", "none"},
		{`for key, val in {"a": 1, "b": 2} { 1 }`, "NONE", "none"},
		{`dct = {"a": 1, "b": 2}; for key in dct { 1 }`, "NONE", "none"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != tt.expectedType {
			t.Errorf("pType is not %s, got=%s", tt.expectedString, pType.Type())
		}

		if pType.String() != tt.expectedString {
			t.Errorf("type string does not match. want=%s, got=%s", tt.expectedString, pType.String())
		}
	}
}

func TestTypeofWhileStatements(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   types.GlimmerType
		expectedString string
	}{
		{"while 1 > 2 { fn() -> none {} () }", "NONE", "none"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != tt.expectedType {
			t.Errorf("pType is not %s, got=%s", tt.expectedString, pType.Type())
		}

		if pType.String() != tt.expectedString {
			t.Errorf("type string does not match. want=%s, got=%s", tt.expectedString, pType.String())
		}
	}
}

func TestTypeofFunctionLiteral(t *testing.T) {
	input := "fn(x: int, y: bool) -> array[int] { [1,2] }"
	expected := "fn(int, bool) -> array[int]"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	ctx := types.NewContext()

	pType := Typeof(program, ctx)

	if pType.Type() != types.FUNCTION {
		t.Fatalf("pType is not types.FunctionType, got=%s", pType.Type())
	}

	if pType.String() != expected {
		t.Fatalf("type string does not match. want=%s, got=%s", expected, pType.String())
	}
}

func TestTypeofArrayLiteral(t *testing.T) {
	input := "[1, 2, 3, 4, 5]"
	expected := "array[int]"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	ctx := types.NewContext()

	pType := Typeof(program, ctx)

	if pType.Type() != types.ARRAY {
		t.Fatalf("pType is not types.ArrayType, got=%s", pType.Type())
	}

	if pType.String() != expected {
		t.Fatalf("type string does not match. want=%s, got=%s", expected, pType.String())
	}
}

func TestTypeofDictLiteral(t *testing.T) {
	input := `{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}`
	expected := "dict[int]"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	ctx := types.NewContext()

	pType := Typeof(program, ctx)

	if pType.Type() != types.DICT {
		t.Fatalf("pType is not types.DictType, got=%s", pType.Type())
	}

	if pType.String() != expected {
		t.Fatalf("type string does not match. want=%s, got=%s", expected, pType.String())
	}
}

func TestTypeofIfExpression(t *testing.T) {
	input := `ife true { 1 } else ife true { 1 } else ife true { 1 } else { 1 }`
	expected := "int"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	ctx := types.NewContext()

	pType := Typeof(program, ctx)

	if pType.Type() != types.INTEGER {
		t.Fatalf("pType is not types.IntegerType, got=%s", pType.Type())
	}

	if pType.String() != expected {
		t.Fatalf("type string does not match. want=%s, got=%s", expected, pType.String())
	}
}

func TestTypeofIfStatement(t *testing.T) {
	input := `if true { 1 } else if true { x = "string" } else if true { fn(x: int) -> int { x + 1 } } else { []int }`
	expected := "none"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	ctx := types.NewContext()

	pType := Typeof(program, ctx)

	if pType.Type() != types.NONE {
		t.Fatalf("pType is not types.NONE, got=%s", pType.Type())
	}

	if pType.String() != expected {
		t.Fatalf("type string does not match. want=%s, got=%s", expected, pType.String())
	}
}

func TestTypeofIndexExpression(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   types.GlimmerType
		expectedString string
	}{
		{"[1,2,3,4][1]", "INTEGER", "int"},
		{`{"a": 1, "b": 2}["b"]`, "INTEGER", "int"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != tt.expectedType {
			t.Errorf("pType is not %s, got=%s", tt.expectedString, pType.Type())
		}

		if pType.String() != tt.expectedString {
			t.Errorf("type string does not match. want=%s, got=%s", tt.expectedString, pType.String())
		}
	}
}

func TestTypeofCallExpression(t *testing.T) {
	input := `myFunc = fn(x: array[int], y: int) -> int { x[y] }; myFunc([1,2], 0);`
	expected := "int"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	ctx := types.NewContext()

	pType := Typeof(program, ctx)

	if pType.Type() != types.INTEGER {
		t.Fatalf("pType is not types.IntegerType, got=%s", pType.Type())
	}

	if pType.String() != expected {
		t.Fatalf("type string does not match. want=%s, got=%s", expected, pType.String())
	}
}

func TestTypeofPrefixExpression(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   types.GlimmerType
		expectedString string
	}{
		{"!true", "BOOLEAN", "bool"},
		{"!2", "BOOLEAN", "bool"},
		{"!2.2", "BOOLEAN", "bool"},
		{"-true", "INTEGER", "int"},
		{"-1", "INTEGER", "int"},
		{"-1.2", "FLOAT", "float"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != tt.expectedType {
			t.Errorf("pType is not %s, got=%s", tt.expectedString, pType.Type())
		}

		if pType.String() != tt.expectedString {
			t.Errorf("type string does not match. want=%s, got=%s", tt.expectedString, pType.String())
		}
	}
}

func TestTypeofInfixExpression(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   types.GlimmerType
		expectedString string
	}{
		{"1 + 1 + 1", "INTEGER", "int"},
		{"1 + 1.2 / false", "FLOAT", "float"},
		{"1 + true * 4.3", "FLOAT", "float"},
		{"1 + true - 2", "INTEGER", "int"},
		{"4.2 + 2.3 * true", "FLOAT", "float"},
		{`"hello" + "world"`, "STRING", "string"},
		{`"hello" - "world"`, "STRING", "string"},
		{`"hello" * "world"`, "STRING", "string"},
		{`"hello" / "world"`, "STRING", "string"},
		{`"hello" == "world"`, "BOOLEAN", "bool"},
		{`"hello" != "world"`, "BOOLEAN", "bool"},
		{`"hello" * 3`, "STRING", "string"},
		{"(1 < 1) >= 3 && 3.2", "BOOLEAN", "bool"},
		{"(1 < 1.2) != true == false", "BOOLEAN", "bool"},
		{"(1 < true) || 5 <= 4 > 3", "BOOLEAN", "bool"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != tt.expectedType {
			t.Errorf("pType is not %s, got=%s", tt.expectedString, pType.Type())
		}

		if pType.String() != tt.expectedString {
			t.Errorf("type string does not match. want=%s, got=%s", tt.expectedString, pType.String())
		}
	}
}

func TestTypeofBuiltin(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   types.GlimmerType
		expectedString string
	}{
		{"print(1)", "NONE", "none"},
		{"len([1,2,3])", "INTEGER", "int"},
		{"head([1,2,3])", "INTEGER", "int"},
		{"tail([fn(x: int) -> int { 1 }, fn(y: int) -> int { 1 }])", "FUNCTION", "fn(int) -> int"},
		{"x = [1,2,3,4,5]; slice(x, 2, 3)", "ARRAY", "array[int]"},
		{"push([1,2,3,4], 5)", "ARRAY", "array[int]"},
		{"pop([ [1,2], [3,4] ])", "ARRAY", "array[int]"},
		{"range(5)", "ARRAY", "array[int]"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		ctx := types.NewContext()

		pType := Typeof(program, ctx)

		if pType.Type() != tt.expectedType {
			t.Errorf("pType is not %s, got=%s", tt.expectedString, pType.Type())
		}

		if pType.String() != tt.expectedString {
			t.Errorf("type string does not match. want=%s, got=%s", tt.expectedString, pType.String())
		}
	}
}
