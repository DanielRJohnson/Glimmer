package parser

import (
	"fmt"
	"glimmer/ast"
	"glimmer/token"
	"glimmer/types"
	"strconv"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAR) {
		return nil
	}

	lit.Parameters, lit.ParamTypes = p.parseFunctionParameters()

	if !p.expectPeek(token.ARROW) {
		return nil
	}

	p.nextToken()

	lit.ReturnType = p.parseTypeNode()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() ([]*ast.Identifier, []types.TypeNode) {
	ids := []*ast.Identifier{}
	parTypes := []types.TypeNode{}

	if p.peekTokenIs(token.RPAR) {
		p.nextToken()
		return ids, parTypes
	}
	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	ids = append(ids, ident)

	// p.nextToken() // curtok = ':'
	if !p.expectPeek(token.COLON) {
		return nil, nil
	}
	p.nextToken() // curtok = type
	parTypes = append(parTypes, p.parseTypeNode())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // curtok = comma
		p.nextToken() // curtok = id
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		ids = append(ids, ident)

		if !p.expectPeek(token.COLON) {
			return nil, nil
		}

		p.nextToken() // curtok = type

		parTypes = append(parTypes, p.parseTypeNode())
	}
	if !p.expectPeek(token.RPAR) {
		return nil, nil
	}
	return ids, parTypes
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		p.nextToken()
		array.ExplicitType = p.parseTypeNode()
		array.Elements = []ast.Expression{}
		return array
	}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseDictLiteral() ast.Expression {
	dict := &ast.DictLiteral{Token: p.curToken}
	dict.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		dict.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return dict
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as an integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as a float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseTypeNode() types.TypeNode {
	switch p.curToken.Type {
	case token.INTEGER_TYPE:
		return INT_T
	case token.FLOAT_TYPE:
		return FLOAT_T
	case token.BOOLEAN_TYPE:
		return BOOL_T
	case token.STRING_TYPE:
		return STRING_T
	case token.ARRAY_TYPE:
		typ := &types.ArrayType{}

		if !p.expectPeek(token.LBRACKET) {
			return nil
		}
		p.nextToken() // curtok = type

		innerType := p.parseTypeNode()
		if innerType == nil {
			return nil
		}

		typ.HeldType = innerType

		if !p.expectPeek(token.RBRACKET) {
			return nil
		}

		return typ
	case token.DICT_TYPE:
		typ := &types.DictType{}

		if !p.expectPeek(token.LBRACKET) {
			return nil
		}
		p.nextToken() // curtok = type

		innerType := p.parseTypeNode()
		if innerType == nil {
			return nil
		}

		typ.HeldType = innerType

		if !p.expectPeek(token.RBRACKET) {
			return nil
		}

		return typ
	case token.FUNCTION:
		typ := &types.FunctionType{}

		if !p.expectPeek(token.LPAR) {
			return nil
		}

		if p.peekTokenIs(token.RPAR) { // no params case
			p.nextToken() // curtok = )
			if !p.expectPeek(token.ARROW) {
				return nil
			}
			p.nextToken()
			typ.ReturnType = p.parseTypeNode()
			return typ
		}

		p.nextToken() // curtok = type

		typ.ParamTypes = append(typ.ParamTypes, p.parseTypeNode()) // first type
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // curtok = ','
		}
		for p.curTokenIs(token.COMMA) { // any more amount of types
			p.nextToken()
			typ.ParamTypes = append(typ.ParamTypes, p.parseTypeNode())
			if !p.peekTokenIs(token.RPAR) {
				p.nextToken()
			}
		}

		if !p.expectPeek(token.RPAR) {
			return nil
		}

		if !p.expectPeek(token.ARROW) {
			return nil
		}

		p.nextToken() // curtok = ret type

		typ.ReturnType = p.parseTypeNode()

		return typ
	case token.NONE_TYPE:
		return NONE_T
	default:
		p.typeNotRecognizedError(p.curToken.Type, p.curToken.Line, p.curToken.Col)
		return nil
	}
}
