package lexer

import (
	"github.com/naronA/monkey/mtoken"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
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
	l.readPosition++
}

func (l *Lexer) NextToken() mtoken.Token {
	var tok mtoken.Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(mtoken.ASSIGN, l.ch)
	case '+':
		tok = newToken(mtoken.PLUS, l.ch)
	case '-':
		tok = newToken(mtoken.MINUS, l.ch)
	case '!':
		tok = newToken(mtoken.BANNG, l.ch)
	case '/':
		tok = newToken(mtoken.SLASH, l.ch)
	case '*':
		tok = newToken(mtoken.ASTERISK, l.ch)
	case '<':
		tok = newToken(mtoken.LT, l.ch)
	case '>':
		tok = newToken(mtoken.GT, l.ch)
	case ';':
		tok = newToken(mtoken.SEMICOLON, l.ch)
	case '(':
		tok = newToken(mtoken.LPAREN, l.ch)
	case ')':
		tok = newToken(mtoken.RPAREN, l.ch)
	case ',':
		tok = newToken(mtoken.COMMA, l.ch)
	case '{':
		tok = newToken(mtoken.LBRACE, l.ch)
	case '}':
		tok = newToken(mtoken.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = mtoken.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = mtoken.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = mtoken.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(mtoken.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(mtokenType mtoken.TokenType, ch byte) mtoken.Token {
	return mtoken.Token{Type: mtokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition > len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
