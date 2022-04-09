package ast

import (
	"github.com/juanfgarcia/gorilla/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Typ: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Typ: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Typ: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%q", program.String())
	}
}
