package lexer

import (
	"testing"

	"github.com/naronA/monkey/mtoken"
)

func TestNextToken2(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};
let result = add(five, ten);
	`
	tests := []struct {
		expectedType    mtoken.TokenType
		expectedLiteral string
	}{
		// let five = 5
		{mtoken.LET, "let"},
		{mtoken.IDENT, "five"},
		{mtoken.ASSIGN, "="},
		{mtoken.INT, "5"},
		{mtoken.SEMICOLON, ";"},
		// let ten = 10;
		{mtoken.LET, "let"},
		{mtoken.IDENT, "ten"},
		{mtoken.ASSIGN, "="},
		{mtoken.INT, "10"},
		{mtoken.SEMICOLON, ";"},
		// let add = fn(x, y) {
		// 	x + y;
		// };
		{mtoken.LET, "let"},
		{mtoken.IDENT, "add"},
		{mtoken.ASSIGN, "="},
		{mtoken.FUNCTION, "fn"},
		{mtoken.LPAREN, "("},
		{mtoken.IDENT, "x"},
		{mtoken.COMMA, ","},
		{mtoken.IDENT, "y"},
		{mtoken.RPAREN, ")"},
		{mtoken.LBRACE, "{"},
		{mtoken.IDENT, "x"},
		{mtoken.PLUS, "+"},
		{mtoken.IDENT, "y"},
		{mtoken.SEMICOLON, ";"},
		{mtoken.RBRACE, "}"},
		{mtoken.SEMICOLON, ";"},
		// let result = add(five, ten);
		{mtoken.LET, "let"},
		{mtoken.IDENT, "result"},
		{mtoken.ASSIGN, "="},
		{mtoken.IDENT, "add"},
		{mtoken.LPAREN, "("},
		{mtoken.IDENT, "five"},
		{mtoken.COMMA, ","},
		{mtoken.IDENT, "ten"},
		{mtoken.RPAREN, ")"},
		{mtoken.SEMICOLON, ";"},
		{mtoken.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken1(t *testing.T) {
	input := `=+(){},;`
	tests := []struct {
		expectedType    mtoken.TokenType
		expectedLiteral string
	}{
		{mtoken.ASSIGN, "="},
		{mtoken.PLUS, "+"},
		{mtoken.LPAREN, "("},
		{mtoken.RPAREN, ")"},
		{mtoken.LBRACE, "{"},
		{mtoken.RBRACE, "}"},
		{mtoken.COMMA, ","},
		{mtoken.SEMICOLON, ";"},
		{mtoken.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

	}
}
