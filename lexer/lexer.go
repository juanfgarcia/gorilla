package lexer

import (
	"github.com/juanfgarcia/gorilla/token"
)

type Lexer struct {
	input    string
	position int
	start    int
	tokens   chan token.Token
}

// LexState is a function that represent a state in the lexer,
// the function lexes a token and return the next LexState.
type LexState func(*Lexer) LexState

const eof = 0

func New(input string) *Lexer {
	lex := &Lexer{
		input:    input,
		position: 0,
		start:    0,
		tokens:   make(chan token.Token),
	}
	go lex.run()
	return lex
}

// next returns the next char in the input
func (lex *Lexer) read() byte {
	if lex.position >= len(lex.input) {
		return eof
	}
	ch := lex.input[lex.position]
	lex.position++
	return ch
}

// backup steps back one position
func (lex *Lexer) backup() {
	lex.position--
}

// peek returns the next char but does not consume it
func (lex *Lexer) peek() byte {
	ch := lex.read()
	lex.backup()
	return ch
}

// emit passes a token to the client
func (lex *Lexer) emit(typ token.TokenType) {
	lex.tokens <- token.Token{typ, lex.input[lex.start:lex.position]}
	lex.start = lex.position
}

// NextToken is the public interface from the lexer
// to the client, it return the tokens concurrently
// as they are read.
func (lex *Lexer) NextToken() token.Token {
	return <-lex.tokens
}

func (lex *Lexer) run() {
	for state := startState; state != nil; {
		state = state(lex)
	}
	close(lex.tokens)
}

func (lex *Lexer) ignoreWhiteSpaces() {
	for ch := lex.read(); isSpace(ch); {
		ch = lex.read()
	}
	lex.backup()
	lex.start = lex.position
}

func startState(lex *Lexer) LexState {
	lex.ignoreWhiteSpaces()

	ch := lex.read()

	switch ch {
	case ';':
		{
			lex.emit(token.SEMICOLON)
			return startState(lex)
		}
	case ':':
		{
			if lex.peek() == '=' {
				lex.read()
				lex.emit(token.ASSIGN)
			}
			return startState(lex)
		}
	default:
		{
			if isLetter(ch) {
				lex.backup()
				return identifierState(lex)
			}
			if isNumber(ch) {
				lex.backup()
				return IntState(lex)
			}
		}
	}

	return nil
}

func identifierState(lex *Lexer) LexState {
	for ch := lex.read(); isLetter(ch); {
		ch = lex.read()
	}
	lex.backup()
	typ := token.LookupIdent(lex.input[lex.start:lex.position])
	lex.emit(typ)
	return startState
}

func IntState(lex *Lexer) LexState {
	for ch := lex.read(); isNumber(ch); {
		ch = lex.read()
	}
	lex.backup()
	lex.emit(token.INT)
	return startState
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}
