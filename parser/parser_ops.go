package parser

import (
	"fmt"
	"glimmer/ast"
	"glimmer/token"
)

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t, p.curToken.Line, p.curToken.Col)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) typeNotRecognizedError(t token.TokenType, line int, col int) {
	msg := fmt.Sprintf("[%d,%d]: type not recognized: %s", line, col, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType, line int, col int) {
	msg := fmt.Sprintf("[%d,%d]: expected next token to be %s, got %s instead",
		line, col, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType, line int, col int) {
	msg := fmt.Sprintf("[%d,%d]: no prefix parse function for %s found", line, col, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func isAssign(tok token.TokenType) bool {
	return tok == token.ASSIGN || tok == token.PLUSEQ ||
		tok == token.MINUSEQ || tok == token.MULTEQ || tok == token.DIVEQ
}
