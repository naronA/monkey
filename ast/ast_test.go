package ast

import (
	"testing"

	"github.com/naronA/monkey/mtoken"
)

/*
このテストではASTを手で組み立てた。
デモンストレーションとして、このテストは構文解析器に大して文字列の比較を行うことで、
可読性の高いテストのレイヤーを追加する方法を教えてくれる
*/
func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: mtoken.Token{Type: mtoken.LET, Literal: "let"},
				Name: &Identifier{
					Token: mtoken.Token{Type: mtoken.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: mtoken.Token{Type: mtoken.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}

}
