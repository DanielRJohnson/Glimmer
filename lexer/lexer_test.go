package lexer

import (
	"testing"

	"glimmer/token"
)

func TestNextToken(t *testing.T) {
	input := "+= -= *= /= for break continue : ==!==!abc+-,; # this is a line comment \n \t\r ()/*><{}100 123.456 123. let fn $ \x00 = && & || <= >= | \"foobar\" \"foo\t\t\tbar\" [1, 2]; "

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
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
		{token.LET, "let"},
		{token.FUNCTION, "fn"},
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
