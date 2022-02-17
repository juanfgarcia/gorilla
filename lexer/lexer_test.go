package lexer

import (
	"github.com/juanfgarcia/gorilla/token"
	"testing"
)

func TestLexer(t *testing.T) {

	t.Run("lex a variable assignment", func(t *testing.T) {
		input := `let a := 3;`

		lexer := New(input)

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.LET, "let"},
			{token.IDENTIFIER, "a"},
			{token.ASSIGN, ":="},
			{token.INT, "3"},
			{token.SEMICOLON, ";"},
		}

		for i, tt := range tests {
			got := lexer.NextToken()

			if tt.expectedType != got.Typ {
				t.Errorf("[%d]Got %d but want %d", i, got.Typ, tt.expectedType)
			}
		}
	})
}
