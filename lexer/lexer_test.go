package lexer

import (
	"testing"

	"github.com/naronA/monkey/mtoken"
)

func TestNextToken4(t *testing.T) {
	input := `
10 == 10;
10 != 9;
`
	tests := []struct {
		expectedType    mtoken.TokenType
		expectedLiteral string
	}{
		{mtoken.INT, "10"},
		{mtoken.EQ, "=="},
		{mtoken.INT, "10"},
		{mtoken.SEMICOLON, ";"},

		{mtoken.INT, "10"},
		{mtoken.NOTEQ, "!="},
		{mtoken.INT, "9"},
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

func TestNextToken3(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};
let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if ( 5 < 10 ) {
	return true;
} else {
	return false;
}
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
		//!-/*5;
		{mtoken.BANNG, "!"},
		{mtoken.MINUS, "-"},
		{mtoken.SLASH, "/"},
		{mtoken.ASTERISK, "*"},
		{mtoken.INT, "5"},
		{mtoken.SEMICOLON, ";"},
		//5 < 10 > 5;
		{mtoken.INT, "5"},
		{mtoken.LT, "<"},
		{mtoken.INT, "10"},
		{mtoken.GT, ">"},
		{mtoken.INT, "5"},
		{mtoken.SEMICOLON, ";"},

		// if ( 5 < 10 ) {
		// 	return true;
		// } else {
		// 	return false;
		// }
		{mtoken.IF, "if"},
		{mtoken.LPAREN, "("},
		{mtoken.INT, "5"},
		{mtoken.LT, "<"},
		{mtoken.INT, "10"},
		{mtoken.RPAREN, ")"},
		{mtoken.LBRACE, "{"},
		{mtoken.RETURN, "return"},
		{mtoken.TRUE, "true"},
		{mtoken.SEMICOLON, ";"},
		{mtoken.RBRACE, "}"},
		{mtoken.ELSE, "else"},
		{mtoken.LBRACE, "{"},
		{mtoken.RETURN, "return"},
		{mtoken.FALSE, "false"},
		{mtoken.SEMICOLON, ";"},
		{mtoken.RBRACE, "}"},

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
