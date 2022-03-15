package parser

import (
	"fmt"
	"glimmer/ast"
	"glimmer/lexer"
	"strconv"
	"testing"
)

/*
* GO IS STUPID AND THE BEST WAY TO HAVE TESTS IS IN A FEW TOP LEVEL FILES.
* THIS IS TO SAY THIS IS A LONG FILE, SO HERE'S THE ORDER OF THINGS:
* 1. HELPERS
* 2. PRECEDENCE
* 3. STATEMENTS
* 4. LITERAL EXPRESSIONS
* 5. PREFIX EXPRESSIONS
* 6. INFIX EXPRESSIONS
* 7. IF EXPRESSIONS
* 8. FOR EXPRESSIONS
* 9. WHILE EXPRESSIONS
* 10. CALL EXPRESSIONS
* 11. INDEX EXPRESSIONS
*
* (CTRL + F) IF NEEDED
 */

/*
* USEFUL HELPERS
 */

func checkParserErrors(t *testing.T, p *Parser) {
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

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, il ast.Expression, value float64) bool {
	float, ok := il.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("il not *ast.FloatLiteral. got=%T", il)
		return false
	}

	if float.Value != value {
		t.Errorf("float.Value not %f. got=%f", value, float.Value)
		return false
	}

	formattedValue := strconv.FormatFloat(value, 'f', -1, 64)
	if float.TokenLiteral() != formattedValue {
		t.Errorf("float.TokenLiteral not %s. got=%s", formattedValue, float.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

/*
* PRECEDENCE TESTS
 */

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b - c", "((a + b) - c)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

/*
* STATEMENT TESTS
 */

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return true;", true},
		{"return x;", "x"},
		{"return 1", 1},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt := program.Statements[0]
		val := stmt.(*ast.ReturnStatement).ReturnValue
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

/*
* LITERAL EXPRESSION TESTS
 */

func TestFunctionLiteralExpression(t *testing.T) {
	input := "fn(a: int, b: float, c: bool, d: string, e: array[array[int]], f: dict[float], g: fn(int, int, fn() -> none) -> int) -> int { x + y; }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 7 {
		t.Fatalf("function literal parameters wrong. want 7, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "a")
	testLiteralExpression(t, function.Parameters[1], "b")
	testLiteralExpression(t, function.Parameters[2], "c")
	testLiteralExpression(t, function.Parameters[3], "d")
	testLiteralExpression(t, function.Parameters[4], "e")
	testLiteralExpression(t, function.Parameters[5], "f")
	testLiteralExpression(t, function.Parameters[6], "g")

	if len(function.ParamTypes) != 7 {
		t.Fatalf("function literal ParamTypes wrong. want 7, got=%d\n", len(function.ParamTypes))
	}

	expectedTypes := []string{"int", "float", "bool", "string", "array[array[int]]", "dict[float]", "fn(int, int, fn() -> none) -> int"}
	for idx, ex := range expectedTypes {
		actual := function.ParamTypes[idx].String()
		if ex != actual {
			t.Fatalf("function param type wrong. want=%s, got=%s", ex, actual)
		}
	}

	if function.ReturnType.String() != "int" {
		t.Fatalf("function return type wrong. want=%s, got=%s", "int", function.ReturnType.String())
	}

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements does not contain %d statements. got=%d\n", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: "fn() -> int {};", expected: []string{}},
		{input: "fn(x: int) -> int {};", expected: []string{"x"}},
		{input: "fn(x: int, y: int, z: int) -> int {};", expected: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expected) {
			t.Errorf("length of parameters is wrong. want=%d, got=%d\n", len(tt.expected), len(function.Parameters))
		}

		for i, ident := range tt.expected {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestArrayLiteralParsing(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestDictLiteralParsing(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	dict, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("exp is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	if len(dict.Pairs) != 3 {
		t.Errorf("dict.Pairs has wrong length. got=%d", len(dict.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, val := range dict.Pairs {
		lit, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[lit.String()]
		testIntegerLiteral(t, val, expectedValue)
	}
}

func TestEmptyDictParsing(t *testing.T) {
	input := `{}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	dict, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("exp is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	if len(dict.Pairs) != 0 {
		t.Errorf("dict.Pairs has wrong length. got=%d", len(dict.Pairs))
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("ident.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", literal.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.3;", 5.3},
		{"0.234", 0.234},
		{"123.", 123.},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		testFloatLiteral(t, stmt.Expression, tt.expected)
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
	}
	if literal.Value != true {
		t.Errorf("ident.Value not %t. got=%t", true, literal.Value)
	}
	if literal.TokenLiteral() != "true" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", literal.TokenLiteral())
	}
}

/*
* PREFIX EXPRESSION TESTS
 */

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input      string
		operator   string
		rightValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!5.3;", "!", 5.3},
		{"-15.3;", "-", 15.3},
		{"!true;", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testPrefixExpression(t, stmt.Expression, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp is not ast.PrefixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

/*
* INFIX EXPRESSION TESTS
 */

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5.3 + 5.3", 5.3, "+", 5.3},
		{"true == true", true, "==", true},
		{"true == false", true, "==", false},
		{"false == false", false, "==", false},
		{"true && false", true, "&&", false},
		{"true || false", true, "||", false},
		{"myArg | myFunc", "myArg", "|", "myFunc"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

/*
* IF EXPRESSION TESTS
 */

func TestIfExpression(t *testing.T) {
	input := "if x < y { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition[0].(*ast.ExpressionStatement).Expression, "x", "<", "y") {
		return
	}
	if len(exp.TrueBranch.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d\n", 1, len(exp.TrueBranch.Statements))
	}

	consequence, ok := exp.TrueBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.TrueBranch.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.FalseBranch != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.FalseBranch)
	}
}

func TestMultiIfExpression(t *testing.T) {
	input := "if x = 5; (x < y) { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if len(exp.Condition) != 2 {
		t.Fatalf("exp.Condition is not %d statements. got=%d\n", 2, len(exp.Condition))
	}

	if !testInfixExpression(t, exp.Condition[1].(*ast.ExpressionStatement).Expression, "x", "<", "y") {
		return
	}
	if len(exp.TrueBranch.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d\n", 1, len(exp.TrueBranch.Statements))
	}

	consequence, ok := exp.TrueBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.TrueBranch.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.FalseBranch != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.FalseBranch)
	}
}

func TestMultiIfElifExpression(t *testing.T) {
	input := "if x = 5; (x < y) { x } else if y = 5; y > x { y } else if z = 5; z > y { z } else { w }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if len(exp.Condition) != 2 {
		t.Fatalf("exp.Condition is not %d statements. got=%d\n", 2, len(exp.Condition))
	}

	if !testInfixExpression(t, exp.Condition[1].(*ast.ExpressionStatement).Expression, "x", "<", "y") {
		return
	}
	if len(exp.TrueBranch.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d\n", 1, len(exp.TrueBranch.Statements))
	}

	consequence, ok := exp.TrueBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.TrueBranch.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.ElifConditions[0]) != 2 {
		t.Fatalf("exp.ElifConditions[0] is not %d statements. got=%d\n", 2, len(exp.ElifConditions[0]))
	}
	if !testInfixExpression(t, exp.ElifConditions[0][1].(*ast.ExpressionStatement).Expression, "y", ">", "x") {
		return
	}
	if len(exp.ElifBranches[0].Statements) != 1 {
		t.Errorf("exp.ElifBranches[0] is not %d statements. got=%d\n", 1, len(exp.ElifBranches[0].Statements))
	}
	if !testIdentifier(t, exp.ElifBranches[0].Statements[0].(*ast.ExpressionStatement).Expression, "y") {
		return
	}

	if len(exp.ElifConditions[1]) != 2 {
		t.Fatalf("exp.ElifConditions[1] is not %d statements. got=%d\n", 2, len(exp.ElifConditions[1]))
	}
	if !testInfixExpression(t, exp.ElifConditions[1][1].(*ast.ExpressionStatement).Expression, "z", ">", "y") {
		return
	}
	if len(exp.ElifBranches[1].Statements) != 1 {
		t.Errorf("exp.ElifBranches[1] is not %d statements. got=%d\n", 1, len(exp.ElifBranches[1].Statements))
	}
	if !testIdentifier(t, exp.ElifBranches[1].Statements[0].(*ast.ExpressionStatement).Expression, "z") {
		return
	}

	alternative, ok := exp.FalseBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.FalseBranch.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "w") {
		return
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition[0].(*ast.ExpressionStatement).Expression, "x", "<", "y") {
		return
	}
	if len(exp.TrueBranch.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d\n", 1, len(exp.TrueBranch.Statements))
	}

	consequence, ok := exp.TrueBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.TrueBranch.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.FalseBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.FalseBranch.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIfElifElseExpression(t *testing.T) {
	input := "if (x < y) { x } else if (x < 100) { x } else if y < 100 { y } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition[0].(*ast.ExpressionStatement).Expression, "x", "<", "y") {
		return
	}
	if len(exp.TrueBranch.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d\n", 1, len(exp.TrueBranch.Statements))
	}

	consequence, ok := exp.TrueBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.TrueBranch.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if !testInfixExpression(t, exp.ElifConditions[0][0].(*ast.ExpressionStatement).Expression, "x", "<", 100) {
		return
	}
	if len(exp.ElifBranches[0].Statements) != 1 {
		t.Errorf("elifBranch1 is not %d statements. got=%d\n", 1, len(exp.ElifBranches[0].Statements))
	}
	elifBranch1, ok := exp.ElifBranches[0].Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ElifBranches[0].Statements[0] is not an ast.ExpressionStatement. got=%T", exp.ElifBranches[0].Statements[0])
	}
	if !testIdentifier(t, elifBranch1.Expression, "x") {
		return
	}

	if !testInfixExpression(t, exp.ElifConditions[1][0].(*ast.ExpressionStatement).Expression, "y", "<", 100) {
		return
	}
	if len(exp.ElifBranches[1].Statements) != 1 {
		t.Errorf("elifBranch2 is not %d statements. got=%d\n", 1, len(exp.ElifBranches[1].Statements))
	}
	elifBranch2, ok := exp.ElifBranches[1].Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ElifBranches[1].Statements[0] is not an ast.ExpressionStatement. got=%T", exp.ElifBranches[1].Statements[0])
	}
	if !testIdentifier(t, elifBranch2.Expression, "y") {
		return
	}

	alternative, ok := exp.FalseBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.FalseBranch.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIfElifExpression(t *testing.T) {
	input := "if (x < y) { x } else if (x < 100) { x } else if y < 100 { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition[0].(*ast.ExpressionStatement).Expression, "x", "<", "y") {
		return
	}
	if len(exp.TrueBranch.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d\n", 1, len(exp.TrueBranch.Statements))
	}

	consequence, ok := exp.TrueBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not an ast.ExpressionStatement. got=%T", exp.TrueBranch.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if !testInfixExpression(t, exp.ElifConditions[0][0].(*ast.ExpressionStatement).Expression, "x", "<", 100) {
		return
	}
	if len(exp.ElifBranches[0].Statements) != 1 {
		t.Errorf("elifBranch1 is not %d statements. got=%d\n", 1, len(exp.ElifBranches[0].Statements))
	}
	elifBranch1, ok := exp.ElifBranches[0].Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ElifBranches[0].Statements[0] is not an ast.ExpressionStatement. got=%T", exp.ElifBranches[0].Statements[0])
	}
	if !testIdentifier(t, elifBranch1.Expression, "x") {
		return
	}

	if !testInfixExpression(t, exp.ElifConditions[1][0].(*ast.ExpressionStatement).Expression, "y", "<", 100) {
		return
	}
	if len(exp.ElifBranches[1].Statements) != 1 {
		t.Errorf("elifBranch2 is not %d statements. got=%d\n", 1, len(exp.ElifBranches[1].Statements))
	}
	elifBranch2, ok := exp.ElifBranches[1].Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ElifBranches[1].Statements[0] is not an ast.ExpressionStatement. got=%T", exp.ElifBranches[1].Statements[0])
	}
	if !testIdentifier(t, elifBranch2.Expression, "y") {
		return
	}

	if exp.FalseBranch != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.FalseBranch)
	}
}

/*
* FOR EXPRESSION TESTS
 */

func TestForExpression(t *testing.T) {
	input := "for x; y, z, w { u }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
	}

	if len(forExp.ForPrecondition) != 2 {
		t.Fatalf("forExp.ForPrecondition is not 2 statements. got=%d", len(forExp.ForPrecondition))
	}
	if len(forExp.ForCondition) != 1 {
		t.Fatalf("forExp.ForCondition is not 1 statement. got=%d", len(forExp.ForCondition))
	}
	if len(forExp.ForPostcondition) != 1 {
		t.Fatalf("forExp.ForPostcondition is not 1 statement. got=%d", len(forExp.ForPostcondition))
	}
	if !testIdentifier(t, forExp.ForPrecondition[0].(*ast.ExpressionStatement).Expression, "x") {
		return
	}
	if !testIdentifier(t, forExp.ForPrecondition[1].(*ast.ExpressionStatement).Expression, "y") {
		return
	}
	if !testIdentifier(t, forExp.ForCondition[0].(*ast.ExpressionStatement).Expression, "z") {
		return
	}
	if !testIdentifier(t, forExp.ForPostcondition[0].(*ast.ExpressionStatement).Expression, "w") {
		return
	}
	if !testIdentifier(t, forExp.Body.Statements[0].(*ast.ExpressionStatement).Expression, "u") {
		return
	}
}

func TestForExpressionOnlyCondition(t *testing.T) {
	input := "for x = true; false || x { u }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
	}

	if len(forExp.ForPrecondition) != 0 {
		t.Fatalf("forExp.ForPrecondition is not 0 statements. got=%d", len(forExp.ForPrecondition))
	}
	if len(forExp.ForCondition) != 2 {
		t.Fatalf("forExp.ForCondition is not 1 statement. got=%d", len(forExp.ForCondition))
	}
	if len(forExp.ForPostcondition) != 0 {
		t.Fatalf("forExp.ForPostcondition is not 0 statements. got=%d", len(forExp.ForPostcondition))
	}
	if !testIdentifier(t, forExp.Body.Statements[0].(*ast.ExpressionStatement).Expression, "u") {
		return
	}
}

func TestForExpressionOnlyConditionAndPostCondition(t *testing.T) {
	input := "for x < 2, x = x + 1 { u }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
	}

	if len(forExp.ForPrecondition) != 0 {
		t.Fatalf("forExp.ForPrecondition is not 0 statements. got=%d", len(forExp.ForPrecondition))
	}
	if len(forExp.ForCondition) != 1 {
		t.Fatalf("forExp.ForCondition is not 1 statement. got=%d", len(forExp.ForCondition))
	}
	if len(forExp.ForPostcondition) != 1 {
		t.Fatalf("forExp.ForPostcondition is not 0 statements. got=%d", len(forExp.ForPostcondition))
	}
	if !testIdentifier(t, forExp.Body.Statements[0].(*ast.ExpressionStatement).Expression, "u") {
		return
	}
}

func TestBreakAndContinue(t *testing.T) {
	input := "break; continue;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 2, len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.BreakStatement. got=%T", program.Statements[0])
	}
	_, ok = program.Statements[1].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ContinueStatement. got=%T", program.Statements[0])
	}
}

/*
* CALL EXPRESSION TESTS
 */

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: "add();", expected: []string{}},
		{input: "add(a);", expected: []string{"a"}},
		{input: "add(a, b, c);", expected: []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp := stmt.Expression.(*ast.CallExpression)

		if len(exp.Arguments) != len(tt.expected) {
			t.Errorf("length of parameters is wrong. want=%d, got=%d\n", len(tt.expected), len(exp.Arguments))
		}

		for i, ident := range tt.expected {
			testLiteralExpression(t, exp.Arguments[i], ident)
		}
	}
}

/*
* INDEX EXPRESSION TESTS
 */

func TestIndexExpressionParsing(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}
