package lexer

import (
	"testing"

	"glimmer/token"
)

func TestNextToken(t *testing.T) {
	input := "ife += -= *= /= for break continue : ==!==!abc+-,; # this is a line comment \n \t\r ()/*><{}100 123.456 123. fn -> $ \x00 = && & || <= >= | \"foobar\" \"foo\t\t\tbar\" [1, 2]; int float bool string array dict none"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IFE, "ife"},
		{token.PLUSEQ, "+="},
		{token.MINUSEQ, "-="},
		{token.MULTEQ, "*="},
		{token.DIVEQ, "/="},
		{token.FOR, "for"},
		{token.BREAK, "break"},
		{token.CONT, "continue"},
		{token.COLON, ":"},
		{token.EQ, "=="},
		{token.NEQ, "!="},
		{token.ASSIGN, "="},
		{token.NOT, "!"},
		{token.ID, "abc"},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.COMMA, ","},
		{token.SEMICOL, ";"},
		{token.LPAR, "("},
		{token.RPAR, ")"},
		{token.DIV, "/"},
		{token.MULT, "*"},
		{token.GT, ">"},
		{token.LT, "<"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.INT, "100"},
		{token.FLOAT, "123.456"},
		{token.FLOAT, "123"},
		{token.FUNCTION, "fn"},
		{token.ARROW, "->"},
		{token.ILLEGAL, "$"},
		{token.EOF, ""},
		{token.ASSIGN, "="},
		{token.AND, "&&"},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.PIPE, "|"},
		{token.STRING, "foobar"},
		{token.STRING, "foo\t\t\tbar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOL, ";"},
		{token.INTEGER_TYPE, "int"},
		{token.FLOAT_TYPE, "float"},
		{token.BOOLEAN_TYPE, "bool"},
		{token.STRING_TYPE, "string"},
		{token.ARRAY_TYPE, "array"},
		{token.DICT_TYPE, "dict"},
		{token.NONE_TYPE, "none"},
		{token.EOF, ""},
	}
	lex := New(input)

	for i, tt := range tests {
		tok := lex.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
