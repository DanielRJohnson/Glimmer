package lexer

import (
	"glimmer/token"
	"strings"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	line         int  // current line number
	linePosition int  // current position on a given line (set to zero after each newline)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.linePosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.SkipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line, l.linePosition)
		}
	case ':':
		tok = newToken(token.COLON, l.ch, l.line, l.linePosition)
	case ';':
		tok = newToken(token.SEMICOL, l.ch, l.line, l.linePosition)
	case '(':
		tok = newToken(token.LPAR, l.ch, l.line, l.linePosition)
	case ')':
		tok = newToken(token.RPAR, l.ch, l.line, l.linePosition)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line, l.linePosition)
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUSEQ, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.PLUS, l.ch, l.line, l.linePosition)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MINUSEQ, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.MINUS, l.ch, l.line, l.linePosition)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NEQ, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.NOT, l.ch, l.line, l.linePosition)
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DIVEQ, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.DIV, l.ch, l.line, l.linePosition)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MULTEQ, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.MULT, l.ch, l.line, l.linePosition)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.GT, l.ch, l.line, l.linePosition)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.LT, l.ch, l.line, l.linePosition)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			literal := string(l.ch) + string(l.ch) // & gets parsed to &&, no need for error
			tok = token.Token{Type: token.AND, Literal: literal, Line: l.line, Col: l.linePosition}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: l.line, Col: l.linePosition}
		} else {
			tok = newToken(token.PIPE, l.ch, l.line, l.linePosition)
		}
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line, l.linePosition)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line, l.linePosition)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line, l.linePosition)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line, l.linePosition)
	case '#':
		ch := l.ch
		for ch != '\n' && ch != 0 {
			l.readChar()
			ch = l.ch
		}
		tok = l.NextToken()
		return tok // l pos has already been incremented. early exit.
	case '"':
		tok = token.Token{Type: token.STRING, Literal: l.readString(), Line: l.line, Col: l.linePosition}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.line, Col: l.linePosition}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Col = l.linePosition
			return tok // l pos has already been incremented. early exit.
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT
				if tok.Literal[len(tok.Literal)-1:] == "." { //if last character is .
					tok.Literal = tok.Literal[0 : len(tok.Literal)-1] //cut off .
				}
			} else {
				tok.Type = token.INT
			}
			tok.Line = l.line
			tok.Col = l.linePosition
			return tok // l pos has already been incremented. early exit.
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.linePosition)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) SkipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line += 1
			l.linePosition = 0
		}
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	start_position := l.position
	for isLetter(l.ch) || (l.position != start_position && isDigit(l.ch)) {
		l.readChar()
	}
	return l.input[start_position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' && l.input[l.position-1] != '\\' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	start_position := l.position
	//read some amount of numbers (int)
	for isDigit(l.ch) {
		l.readChar()
	}
	//optionally read a period (float)
	if l.ch == '.' {
		l.readChar()
		//read however many digits come after a period (could be zero)
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[start_position:l.position]
}

func newToken(tokenType token.TokenType, ch byte, line int, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Col: col}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
