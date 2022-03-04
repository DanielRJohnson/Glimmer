package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	ID     = "ID"     // add, foobar, x, y, ...
	INT    = "INT"    // 123456
	FLOAT  = "FLOAT"  // 123.456
	STRING = "STRING" // "Hello, World!"

	// Operators
	ASSIGN  = "="
	PLUS    = "+"
	MINUS   = "-"
	NOT     = "!"
	MULT    = "*"
	DIV     = "/"
	PLUSEQ  = "+="
	MINUSEQ = "-="
	MULTEQ  = "*="
	DIVEQ   = "/="

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	EQ  = "=="
	NEQ = "!="
	AND = "||"
	OR  = "&&"

	// Special
	PIPE = "|"

	// Delimiters
	COMMA   = ","
	COLON   = ":"
	SEMICOL = ";"

	LPAR     = "("
	RPAR     = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	BREAK    = "break"
	CONT     = "continue"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONT,
	"return":   RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ID
}
