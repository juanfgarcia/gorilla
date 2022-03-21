package parser

import (
    "github.com/juanfgarcia/gorilla/ast"
    "github.com/juanfgarcia/gorilla/lexer"
    "github.com/juanfgarcia/gorilla/token"

    "fmt"
)

type Parser struct {
    l *lexer.Lexer
    
    curToken token.Token
    peekToken token.Token

    errors []string
}

func New(input string) *Parser {
    l := lexer.New(input)
    p := &Parser{
        l : l,
        errors: []string{},
    }

    p.NextToken()
    p.NextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken)

    p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Typ != token.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.NextToken()
    }
    return program
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Typ {
    case token.LET:
        return p.parseLetStatement()
    default:
        return nil
    }
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}

    if !p.expectPeek(token.IDENTIFIER) {
        return nil
    }

    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    if !p.expectPeek(token.ASSIGN) {
        return nil
    }

    //TODO: parseExpression, by the moment skip until semicolon

    for p.curToken.Typ != token.SEMICOLON {
        p.NextToken()
    }

    return stmt
}


func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekToken.Typ == t {
        p.NextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
}
