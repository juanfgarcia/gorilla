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
    TRUE
    FALSE

	// Operators
	ASSIGN
	EQUALS
    NEQUALS
	PLUS
	MINUS
	ASTERISK
	RIGHTARROW
	SLASH
    BANG
    LT
    GT

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
	IF
	ELSE
)

func (t TokenType) String() string {
	return [...]string{
		"ILLEGAL",
		"EOF",
		"IDENTIFIER",
		"INT",
        "TRUE",
        "FALSE",
		"ASSIGN",
		"EQUALS",
        "NEQUALS",
		"PLUS",
		"MINUS",
		"ASTERISK",
		"RIGHTARROW",
		"SLASH",
        "BANG",
        "LT",
        "GT",
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
		"RETURN",
		"IF",
		"ELSE"}[t]
}

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"return": RETURN,
	"Int":    TYPE,
    "true":   TRUE,
    "false":  FALSE,
	"if":	  IF,
	"else":	  ELSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
