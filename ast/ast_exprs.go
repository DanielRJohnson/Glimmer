package ast

import (
	"bytes"
	"glimmer/token"
	"strings"
)

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token          token.Token
	Condition      Expression
	TrueBranch     *BlockStatement
	ElifBranches   []*BlockStatement
	ElifConditions []Expression
	FalseBranch    *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(ie.Condition.String())
	out.WriteString(") ")
	out.WriteString(ie.TrueBranch.String())

	for index, branch := range ie.ElifBranches {
		out.WriteString(" else if ")
		out.WriteString(ie.ElifConditions[index].String())
		out.WriteString(branch.String())
	}

	if ie.FalseBranch != nil {
		out.WriteString(" else ")
		out.WriteString(ie.FalseBranch.String())
	}

	return out.String()
}

type ForExpression struct {
	Token            token.Token
	ForPrecondition  []Statement
	ForCondition     []Statement
	ForPostcondition []Statement
	Body             *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fe.TokenLiteral() + " ")
	for _, preStmt := range fe.ForPrecondition {
		out.WriteString(preStmt.String() + " ")
	}
	if len(fe.ForPrecondition) != 0 {
		out.WriteString(", ")
	}
	for _, condStmt := range fe.ForCondition {
		out.WriteString(condStmt.String() + " ")
	}
	if len(fe.ForCondition) != 0 {
		out.WriteString(", ")
	}
	for _, postStmt := range fe.ForPostcondition {
		out.WriteString(postStmt.String() + " ")
	}
	out.WriteString(" " + fe.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}
