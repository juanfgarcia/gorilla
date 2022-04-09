package parser

import (
	"github.com/juanfgarcia/gorilla/ast"
	"github.com/juanfgarcia/gorilla/lexer"
	"github.com/juanfgarcia/gorilla/token"

	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQUALS:   EQUALS,
	token.NEQUALS:  EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(input string) *Parser {
	l := lexer.New(input)
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.IDENTIFIER] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.TRUE] = p.parseBoolean
	p.prefixParseFns[token.FALSE] = p.parseBoolean
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpressions
	p.prefixParseFns[token.IF] = p.parseIfExpression
	p.prefixParseFns[token.FUNCTION] = p.parseFunctionLiteral

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.SLASH] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[token.EQUALS] = p.parseInfixExpression
	p.infixParseFns[token.NEQUALS] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression

	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) peekPrecedence() int {
	if pred, ok := precedences[p.peekToken.Typ]; ok {
		return pred
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if pred, ok := precedences[p.curToken.Typ]; ok {
		return pred
	}
	return LOWEST
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Typ == token.SEMICOLON {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Typ]

	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", p.curToken.Literal))
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Typ != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Typ]

		if infix == nil {
			p.errors = append(p.errors, fmt.Sprintf("no infix parse function for %s found", p.curToken.Literal))
			return leftExp
		}

		p.NextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.NextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Typ == token.ELSE {
		p.NextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.NextToken()

	for p.curToken.Typ != token.RBRACE && p.curToken.Typ != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}
	return block
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.NextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.NextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseGroupedExpressions() ast.Expression {
	p.NextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value

	return lit
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekToken.Typ == token.RPAREN {
		p.NextToken()
		return identifiers
	}

	p.NextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekToken.Typ == token.COMMA {
		p.NextToken()
		p.NextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBoolean() ast.Expression {
	if p.curToken.Typ == token.TRUE {
		return &ast.Boolean{Token: p.curToken, Value: true}
	} else {
		return &ast.Boolean{Token: p.curToken, Value: false}
	}
}
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	//TODO: parseExpression, by the moment skip until semicolon

	for p.curToken.Typ != token.SEMICOLON {
		p.NextToken()
	}

	return stmt
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
