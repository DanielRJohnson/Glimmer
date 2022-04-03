package ast

import (
	"bytes"
	"glimmer/token"
	"strings"
)

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " " + ls.Name.String() + " = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type AssignStatement struct {
	Token token.Token
	Type  token.TokenType
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	return as.Name.String() + " = " + as.Value.String() + ";"
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	return out.String()
}

type IfStatement struct {
	Token          token.Token
	Condition      []Statement
	TrueBranch     *BlockStatement
	ElifBranches   []*BlockStatement
	ElifConditions [][]Statement
	FalseBranch    *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	for _, condStmt := range is.Condition {
		out.WriteString(condStmt.String())
	}
	out.WriteString(") ")
	out.WriteString(is.TrueBranch.String())

	for index, branch := range is.ElifBranches {
		out.WriteString(" else if ")
		for _, condStmt := range is.ElifConditions[index] {
			out.WriteString(condStmt.String())
		}
		out.WriteString(branch.String())
	}

	if is.FalseBranch != nil {
		out.WriteString(" else ")
		out.WriteString(is.FalseBranch.String())
	}

	return out.String()
}

type ForStatement struct {
	Token      token.Token
	LoopVars   []*Identifier
	Collection Expression
	Body       *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.TokenLiteral() + " ")
	lvStrings := []string{}
	for _, lv := range fs.LoopVars {
		lvStrings = append(lvStrings, lv.Value)
	}
	out.WriteString(strings.Join(lvStrings, ", ") + " in ")
	out.WriteString(fs.Collection.String() + fs.Body.String())

	return out.String()
}

type WhileStatement struct {
	Token     token.Token
	Condition []Statement
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ws.TokenLiteral() + " (")
	for _, condStmt := range ws.Condition {
		out.WriteString(condStmt.String())
	}
	out.WriteString(" )" + ws.Body.String())

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")

	return out.String()
}

type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string {
	return bs.TokenLiteral() + ";"
}

type ContinueStatement struct {
	Token token.Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string {
	return cs.TokenLiteral() + ";"
}
