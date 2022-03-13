package parser

import (
	"glimmer/ast"
	"glimmer/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.RETURN:
		return p.parseReturnStatement()
	case token.ID:
		if isAssign(p.peekToken.Type) {
			return p.parseAssignStatement()
		} else {
			return p.parseExpressionStatement()
		}
	case token.BREAK:
		br := &ast.BreakStatement{Token: p.curToken}
		if p.peekTokenIs(token.SEMICOL) {
			p.nextToken()
		}
		return br
	case token.CONT:
		ct := &ast.ContinueStatement{Token: p.curToken}
		if p.peekTokenIs(token.SEMICOL) {
			p.nextToken()
		}
		return ct
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: p.peekToken, Type: p.peekToken.Type}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // curtok = assign
	p.nextToken() // curtok = value

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}
