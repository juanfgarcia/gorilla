package parser

import (
	"github.com/juanfgarcia/gorilla/ast"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let a  = 23;`

	p := New(input)
	program := p.ParseProgram()

	AssertNoErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements should contain 2 elements, got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"a"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		AssertLetStmt(t, stmt, tt.expectedIdentifier)
	}
}

func AssertLetStmt(t testing.TB, stmt ast.Statement, name string) {
	t.Helper()

	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not 'x', got = %s", letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not 'x', got = %s", letStmt.Name.TokenLiteral())
	}
}

func AssertNoErrors(t testing.TB, p *Parser) {
	t.Helper()

	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser have %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parse error: %q", msg)
	}
	t.FailNow()
}
