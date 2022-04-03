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
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.WHILE:
		return p.parseWhileStatement()
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

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	for !p.peekTokenIs(token.LBRACE) {
		p.nextToken() // cur = IF , peek = first of cond
		stmt.Condition = append(stmt.Condition, p.parseStatement())
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.TrueBranch = p.parseBlockStatement()

	// for peek token else, parse either elif or else
	hasEncounteredElse := false
	for p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) { // elif
			p.nextToken() //curToken = If

			condStmts := []ast.Statement{}
			for !p.peekTokenIs(token.LBRACE) {
				p.nextToken() //curToken = first token of condition statement
				condStmts = append(condStmts, p.parseStatement())
			}
			p.nextToken()
			stmt.ElifConditions = append(stmt.ElifConditions, condStmts)
			stmt.ElifBranches = append(stmt.ElifBranches, p.parseBlockStatement())
		} else { // else
			if !p.expectPeek(token.LBRACE) || hasEncounteredElse {
				return nil
			}
			stmt.FalseBranch = p.parseBlockStatement()
			hasEncounteredElse = true
		}
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(token.ID) {
		return nil
	} // curtok = loopvar

	firstLVar := p.parseIdentifier().(*ast.Identifier)
	stmt.LoopVars = append(stmt.LoopVars, firstLVar)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if !p.expectPeek(token.ID) {
			return nil
		} //curtok = loopvar
		lv := p.parseIdentifier().(*ast.Identifier)
		stmt.LoopVars = append(stmt.LoopVars, lv)
	}

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken() // curtok = collection
	stmt.Collection = p.parseExpression(LOWEST)

	p.nextToken() // curtok = body
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	for !p.peekTokenIs(token.LBRACE) {
		p.nextToken() // cur = IF , peek = first of cond
		stmt.Condition = append(stmt.Condition, p.parseStatement())
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

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
