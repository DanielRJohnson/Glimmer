package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
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
	ARROW   = "->"

	LPAR     = "("
	RPAR     = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONT     = "CONTINUE"
	RETURN   = "RETURN"

	// Type Keywords
	INTEGER_TYPE = "INTEGER_TYPE"
	FLOAT_TYPE   = "FLOAT_TYPE"
	BOOLEAN_TYPE = "BOOLEAN_TYPE"
	STRING_TYPE  = "STRING_TYPE"
	ARRAY_TYPE   = "ARRAY_TYPE"
	DICT_TYPE    = "DICT_TYPE"
	// fn type is handled by fn
	//FUNCTION_TYPE = "FUNCTION_TYPE"
)

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONT,
	"return":   RETURN,
	"int":      INTEGER_TYPE,
	"float":    FLOAT_TYPE,
	"bool":     BOOLEAN_TYPE,
	"string":   STRING_TYPE,
	"array":    ARRAY_TYPE,
	"dict":     DICT_TYPE,
	// fn type is handled by fn
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ID
}
