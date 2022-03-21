package lexer

import (
	"github.com/juanfgarcia/gorilla/token"
	"testing"
)

type tokenTest struct{
     	expectedType    token.TokenType
	expectedLiteral string
}

func LexAssert(t *testing.T, input string, want []tokenTest) {
    t.Helper()
    lexer := New(input)
     
    for i, tt := range want {
	got := lexer.NextToken()

	if tt.expectedType != got.Typ {
	    t.Errorf("[%d]Got %s but want %s", i, got.Typ, tt.expectedType)
	}

	if tt.expectedLiteral != got.Literal {
	    t.Errorf("[%d]Got %s but want %s", i, got.Literal, tt.expectedLiteral)
	}

    }
}

func TestLexer(t *testing.T) {

	t.Run("Variable assignment", func(t *testing.T) {
		input := `let a := 3;`

		tests := []tokenTest {
			{token.LET, "let"},
			{token.IDENTIFIER, "a"},
			{token.ASSIGN, ":="},
			{token.INT, "3"},
			{token.SEMICOLON, ";"},
            {token.EOF, ""},
		}

		LexAssert(t,input,tests)
		
	})

	t.Run("Function declaration", func(t *testing.T) {
		input := `fn add(x : Int, y: Int) -> Int {
		      return x + y;
		}`

		tests := []tokenTest {
			{token.FUNCTION, "fn"},
			{token.IDENTIFIER, "add"},
			{token.LPAREN, "("},
			{token.IDENTIFIER, "x"},
			{token.COLON, ":"},
			{token.TYPE, "Int"},
			{token.COMMA, ","},
			{token.IDENTIFIER, "y"},
			{token.COLON, ":"},
			{token.TYPE, "Int"},
			{token.RPAREN, ")"},
			{token.RIGHTARROW, "->"},
			{token.TYPE, "Int"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.IDENTIFIER,"x"},
			{token.PLUS, "+"},
			{token.IDENTIFIER, "y"},
			{token.SEMICOLON, ";"},
            {token.RBRACE, "}"},
            {token.EOF, ""},
		}

		LexAssert(t,input,tests)
		
	})

	t.Run("Arithmetic expressions", func(t *testing.T) {
		input := `2+3/5*4-;`

		want := []tokenTest {
		     {token.INT, "2"},
		     {token.PLUS, "+"},
		     {token.INT, "3"},
		     {token.SLASH, "/"},
		     {token.INT, "5"},
		     {token.ASTERISK, "*"},
		     {token.INT, "4"},
		     {token.MINUS, "-"},
             {token.SEMICOLON, ";"},
             {token.EOF, ""},
		}
		
		LexAssert(t, input, want)
	})
}
