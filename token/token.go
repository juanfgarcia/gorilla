package token


type TokenType int

type Token struct {
	Typ     TokenType
	Literal string
}

const (
	ILLEGAL = iota
	EOF

	// Identifiers
	IDENTIFIER
	INT

	// Operators
	ASSIGN
	EQUALS
	PLUS
	MINUS
	ASTERISK
	RIGHTARROW
	SLASH
	 
	// Delimiters
	COMMA
	COLON
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// Keywords
	TYPE
	LET
	FUNCTION
	RETURN
)

func (t TokenType) String() string{
     return [...]string{
	"ILLEGAL",
	"EOF",
	"IDENTIFIER",
	"INT",
	"ASSIGN",
	"EQUALS",
	"PLUS",
	"MINUS",
	"ASTERISK",
	"RIGHTARROW",
	"SLASH",
	"COMMA",
	"COLON",
	"SEMICOLON",
	"LPAREN",
	"RPAREN",
	"LBRACE",
	"RBRACE",
    "TYPE",
	"LET",
	"FUNCTION",
	"RETURN"}[t]
}

var keywords = map[string]TokenType{
	"let": LET,
	"fn": FUNCTION,
	"return": RETURN,
	"Int" : TYPE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
