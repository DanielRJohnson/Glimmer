package parser

import (
	"glimmer/ast"
	"glimmer/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type, p.curToken.Line, p.curToken.Col)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOL) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAR) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	for !p.peekTokenIs(token.LBRACE) {
		p.nextToken() // cur = IFE , peek = first of cond
		expression.Condition = append(expression.Condition, p.parseStatement())
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.TrueBranch = p.parseBlockStatement()

	// for peek token else, parse either elif or else
	hasEncounteredElse := false
	for p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IFE) { // elif
			p.nextToken() //curToken = If

			condStmts := []ast.Statement{}
			for !p.peekTokenIs(token.LBRACE) {
				p.nextToken() //curToken = first token of condition statement
				condStmts = append(condStmts, p.parseStatement())
			}
			p.nextToken()
			expression.ElifConditions = append(expression.ElifConditions, condStmts)
			expression.ElifBranches = append(expression.ElifBranches, p.parseBlockStatement())
		} else { // else
			if !p.expectPeek(token.LBRACE) || hasEncounteredElse {
				return nil
			}
			expression.FalseBranch = p.parseBlockStatement()
			hasEncounteredElse = true
		}
	}

	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAR)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}
