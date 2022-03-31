package parser

import (
	"glimmer/token"
)

const (
	_ int = iota
	LOWEST
	PIPE
	EQUALS
	BOOLEANOP
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.PIPE:     PIPE,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.AND:      BOOLEANOP,
	token.OR:       BOOLEANOP,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.MULT:     PRODUCT,
	token.DIV:      PRODUCT,
	token.LPAR:     CALL,
	token.LBRACKET: INDEX,
}
