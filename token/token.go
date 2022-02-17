package token

type TokenType int

type Token struct {
	Typ     TokenType
	Literal string
}

const (
	ILLEGAL = iota
	EOF

	ASSIGN

	SEMICOLON

	IDENTIFIER
	INT
	LET
)

var keywords = map[string]TokenType{
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
