package parser

import (
	"fmt"

	"github.com/naronA/monkey/ast"
	"github.com/naronA/monkey/lexer"
	"github.com/naronA/monkey/mtoken"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string // errorsフィールドができた。単に文字列のスライスだ。

	curToken  mtoken.Token
	peekToken mtoken.Token

	prefixParseFns map[mtoken.TokenType]prefixParseFn
	infixParseFns  map[mtoken.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 2つのトークンを読み込む. curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t mtoken.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t, p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != mtoken.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// このswitch分岐がどんどん増えていく
	switch p.curToken.Type {
	case mtoken.LET:
		return p.parseLetStatement()
	case mtoken.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(mtoken.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(mtoken.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(mtoken.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t mtoken.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t mtoken.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t mtoken.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// 46 page
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()

	// TODO: セミコロンに遭遇するまで指揮を読み飛ばしている
	for !p.curTokenIs(mtoken.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) registerPrefix(tokenType mtoken.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType mtoken.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)
