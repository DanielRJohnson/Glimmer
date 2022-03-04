package parser

import (
	"glimmer/ast"
	"glimmer/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
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

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.TrueBranch = p.parseBlockStatement()

	// for peek token else, parse either elif or else
	hasEncounteredElse := false
	for p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) { // elif
			p.nextToken() //curToken = If
			p.nextToken() //curToken = first token of expression
			expression.ElifConditions = append(expression.ElifConditions, p.parseExpression(LOWEST))

			if !p.expectPeek(token.LBRACE) {
				return nil
			}

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

// for {} -- infinite
// for true {} -- condition
// for i < 2, i = i + 1 -- condition and postcondition
// for let i = 2, i < 5, i = i + 1 -- precondition, condition, postcondition

func (p *Parser) parseForExpression() ast.Expression {
	fe := &ast.ForExpression{Token: p.curToken}

	if p.peekTokenIs(token.LBRACE) { // infinite loop
		p.nextToken()
		fe.Body = p.parseBlockStatement()
		return fe
	}

	hasOptionalParen := false
	if p.peekTokenIs(token.LPAR) {
		p.nextToken()
		hasOptionalParen = true
	}

	section1 := []ast.Statement{}
	section2 := []ast.Statement{}
	section3 := []ast.Statement{}
	commaCounter := 0

	for !p.peekTokenIs(token.LBRACE) && !(hasOptionalParen && p.peekTokenIs(token.RPAR)) {
		p.nextToken()
		switch commaCounter {
		case 0:
			section1 = append(section1, p.parseStatement())
		case 1:
			section2 = append(section2, p.parseStatement())
		case 2:
			section3 = append(section3, p.parseStatement())
		case 3:
			p.maxOccuranceError(token.COMMA, "ForExpression")
			return nil
		}
		if p.peekTokenIs(token.COMMA) {
			commaCounter++
			p.nextToken()
		}
	}

	if hasOptionalParen && !p.expectPeek(token.RPAR) {
		return nil
	}

	if len(section2) == 0 && len(section3) == 0 { // only condition
		fe.ForCondition = section1
	} else if len(section3) == 0 { // only condition and postcondition
		fe.ForCondition = section1
		fe.ForPostcondition = section2
	} else { // precondition, condition, and postcondition
		fe.ForPrecondition = section1
		fe.ForCondition = section2
		fe.ForPostcondition = section3
	}

	p.nextToken()
	fe.Body = p.parseBlockStatement()

	return fe
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
