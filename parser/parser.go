package parser

import (
	"fmt"
	"strconv"

	"github.com/naronA/monkey/ast"
	"github.com/naronA/monkey/lexer"
	"github.com/naronA/monkey/mtoken"
)

var precedenses = map[mtoken.TokenType]int{
	mtoken.EQ:       EQUALS,
	mtoken.NOTEQ:    EQUALS,
	mtoken.LT:       LESSGREATER,
	mtoken.GT:       LESSGREATER,
	mtoken.PLUS:     SUM,
	mtoken.MINUS:    SUM,
	mtoken.SLASH:    PRODUCT,
	mtoken.ASTERISK: PRODUCT,
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > または <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X または !X
	CALL        // myFunction(X)
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
	p.prefixParseFns = make(map[mtoken.TokenType]prefixParseFn)
	p.registerPrefix(mtoken.IDENT, p.parseIdentifier)
	p.registerPrefix(mtoken.INT, p.parseIntegerLiteral)
	p.registerPrefix(mtoken.BANNG, p.parsePrefixExpression)
	p.registerPrefix(mtoken.MINUS, p.parsePrefixExpression)
	p.registerPrefix(mtoken.TRUE, p.parseBoolean)
	p.registerPrefix(mtoken.FALSE, p.parseBoolean)
	p.registerPrefix(mtoken.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(mtoken.IF, p.parseIfExpression)

	p.infixParseFns = make(map[mtoken.TokenType]infixParseFn)
	p.registerInfix(mtoken.PLUS, p.parseInfixExpression)
	p.registerInfix(mtoken.MINUS, p.parseInfixExpression)
	p.registerInfix(mtoken.SLASH, p.parseInfixExpression)
	p.registerInfix(mtoken.ASTERISK, p.parseInfixExpression)
	p.registerInfix(mtoken.EQ, p.parseInfixExpression)
	p.registerInfix(mtoken.NOTEQ, p.parseInfixExpression)
	p.registerInfix(mtoken.LT, p.parseInfixExpression)
	p.registerInfix(mtoken.GT, p.parseInfixExpression)

	// 2つのトークンを読み込む. curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(mtoken.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(mtoken.RPAREN) {
		return nil
	}

	if !p.expectPeek(mtoken.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(mtoken.ELSE) {
		p.nextToken()

		if !p.expectPeek(mtoken.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(mtoken.RBRACE) && !p.curTokenIs(mtoken.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(mtoken.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(mtoken.TRUE)}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecendece()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
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

/*
単に*ast.Identifierを返すだけ
現在のトークンをTokenに格納
現在のトークンのリテラル値をValueに格納
*/
func (p *Parser) parseIdentifier() ast.Expression {
	/*
		全ての構文解析関数prefixParseFnやinfixParseFnは次の規約に従う
		構文解析関数に関連付けられたトークンがcurTokenにセットされている状態で動作を開始する
		そして、この関数の処理対象である式の一番最後のトークンがcurTokenにセットされた状態になるまですすんで終了する
	*/
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
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
		return p.parseExpressionStatement()
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

	// セミコロンに遭遇するまで式を読み飛ばしてしまっている.
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

	// セミコロンに遭遇するまで指揮を読み飛ばしている
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(mtoken.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(mtoken.SEMICOLON) && precedence < p.peekPrecendece() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)

		return nil
	}

	lit.Value = value

	return lit
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) noPrefixParseFnError(t mtoken.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecendece() int {
	if p, ok := precedenses[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecendece() int {
	if p, ok := precedenses[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
