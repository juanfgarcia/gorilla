package parser

import (
    "fmt"
	"github.com/juanfgarcia/gorilla/ast"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let a  = 23;`

	p := New(input)
	program := p.ParseProgram()

	AssertNoErrors(t, p)

    AssertNumberStatements(t,len(program.Statements), 2)

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

func TestReturnStatements(t *testing.T) {
    input := `return 5;
    return 123;
    return 889712;`

    p := New(input)
	program := p.ParseProgram()
	AssertNoErrors(t, p)


    AssertNumberStatements(t,len(program.Statements), 3)

    for _,stmt := range program.Statements {
        returnStmt, ok := stmt.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("stmt not *ast.ReturnStatement. got=%T",stmt)
            continue
        }

        if returnStmt.TokenLiteral() != "return" {
            t.Errorf("returnStmt.TokenLiteral not 'return', got=%q", returnStmt.TokenLiteral())
        }
    }
}


func TestPrefixExpressions(t *testing.T) {
    prefixTests := []struct {
        input string
        operator string
        right interface{}
    }{
        {"!5;", "!", 5},
        {"-15;", "-", 15},
        {"!false;", "!", false},
    }

    for _, tt := range prefixTests {
        p := New(tt.input)
        program := p.ParseProgram()
        AssertNoErrors(t, p)

        AssertNumberStatements(t, len(program.Statements), 1)

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", 
            program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("stmt is not ast.PrefixExpression, got=%T", stmt.Expression)
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
        }
        
        AssertLiteralExpression(t, exp.Right, tt.right)
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    p := New(input)
    program := p.ParseProgram()
    AssertNoErrors(t, p)


    AssertNumberStatements(t, len(program.Statements), 1)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    AssertIdentifier(t, stmt.Expression, "foobar")
}

func TestBooleanExpression(t *testing.T) {
    input := "true;"

    p := New(input)
    program := p.ParseProgram()
    AssertNoErrors(t, p)


    AssertNumberStatements(t, len(program.Statements), 1)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    AssertBoolean(t, stmt.Expression, true)
}

func TestInfixExpression(t *testing.T) {
    infixTests := []struct{
        input string
        leftValue interface{}
        operator string
        rightValue interface{}
    }{
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},
        {"true == true;",   true, "==", true},
        {"true != false;",  true, "!=", false},
        {"false == false;", false, "==", false},
    }
    
    for _, tt := range infixTests {
        p := New(tt.input)
        program := p.ParseProgram()

        AssertNoErrors(t,p)

        AssertNumberStatements(t, len(program.Statements), 1)

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Errorf("stmt not *ast.ExpressionStatement. got=%T",stmt)
        }

        AssertInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
    }

}

func AssertInfixExpression(t testing.TB, exp ast.Expression, left interface{}, operator string, right interface{}) {
    t.Helper()

    infExp, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Fatalf("exp not *ast.InfixExpression. got=%T", infExp)
    }

    AssertLiteralExpression(t, infExp.Left, left)
    
    if infExp.Operator != operator {
        t.Fatalf("exp.Operator is not '%s'. got=%s",  operator, infExp.Operator)
    }

    AssertLiteralExpression(t, infExp.Right, right)
}

func TestOperatorPrecedence(t *testing.T) {
    tests := []struct{
        input string
        expected string
    }{
        {"-a * b;", "((-a) * b)"},
        {"!-a;", "(!(-a))"},
        {"a + b + c;", "((a + b) + c)"},
        {"a + b / c;", "(a + (b / c))"},
        {" 3 > 5 == false;", "((3 > 5) == false)"},
        {"(2+3) * 3;", "((2 + 3) * 3)"},
    }

    for _, tt := range tests {
        p := New(tt.input)
        program := p.ParseProgram()
        AssertNoErrors(t, p)

        got := program.String()
        if got != tt.expected {
            t.Errorf("Want=%q, but got=%q", tt.expected, got)
        }
    }
}


func TestIfExpression(t *testing.T) {
    input := `if (x < y) { x }`

    p := New(input)
    program := p.ParseProgram()
    AssertNoErrors(t,p)

    AssertNumberStatements(t,len(program.Statements), 1)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("prorgam.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", program.Statements[0])
    }

    AssertInfixExpression(t, exp.Condition, "x", "<", "y")

    AssertNumberStatements(t, len(program.Statements), 1)

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
    }

    AssertIdentifier(t, consequence.Expression, "x")

    if exp.Alternative != nil {
        t.Errorf("exp.Alternative.Statements was not nil, got=%+v", exp.Alternative)
    }
}

func TestIfElseExpression(t *testing.T){
    input := `if (x < y) { x } else { y }`

    p := New(input)
    program := p.ParseProgram()
    AssertNoErrors(t,p)

    AssertNumberStatements(t,len(program.Statements), 1)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("prorgam.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", program.Statements[0])
    }

    AssertInfixExpression(t, exp.Condition, "x", "<", "y")

    AssertNumberStatements(t, len(exp.Consequence.Statements), 1)

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
    }

    AssertIdentifier(t, consequence.Expression, "x")

    if exp.Alternative == nil {
        t.Errorf("exp.Alternative.Statements was nil, got=%+v", exp.Alternative)
    }

    AssertNumberStatements(t, len(exp.Alternative.Statements), 1)

    alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
    }

    AssertIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteral(t *testing.T) {
    input := `fn(x, y) { x + y }`

    p := New(input)

    program := p.ParseProgram()
    AssertNoErrors(t,p)
    AssertNumberStatements(t, len(program.Statements), 1)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    function, ok := stmt.Expression.(*ast.FunctionLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
    }

    AssertNumberStatements(t, len(function.Parameters), 2)

    AssertLiteralExpression(t, function.Parameters[0], "x")
    AssertLiteralExpression(t, function.Parameters[1], "y")
    
    body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", function.Body.Statements[0])
    }

    AssertInfixExpression(t, body.Expression, "x", "+", "y")
    
}

func AssertLiteralExpression(t testing.TB, exp ast.Expression, expected interface{}) {
    t.Helper()

    switch v := expected.(type) {
    case int:
        AssertIntegerLiteral(t, exp, int64(v))
    case int64:
        AssertIntegerLiteral(t, exp, v)
    case string:
        AssertIdentifier(t, exp, v)
    case bool:
        AssertBoolean(t, exp, v)
    default:
        t.Errorf("not handler for %T", expected)
    }
}

func AssertIntegerLiteral(t testing.TB, il ast.Expression, value int64) {
    t.Helper()

    literal, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("exp not *ast.IntegerLiteral. got=%T", literal)
    }
    if literal.Value != value {
        t.Errorf("literal.Value not %d. got=%d", value, literal.Value)
    }
    if literal.TokenLiteral() != fmt.Sprintf("%d", value) {
        t.Errorf("literal.TokenLiteral not %d. got=%s", value, literal.TokenLiteral())
    }
}

func AssertIdentifier(t testing.TB, ie ast.Expression, value string) {
    t.Helper()

    ident, ok := ie.(*ast.Identifier)
    if !ok {
        t.Fatalf("exp not *ast.Identifier. got=%T", ident)
    }

    if ident.Value != value {
        t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
    }

    if ident.TokenLiteral() != value {
        t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
    }
}

func AssertBoolean(t testing.TB, b ast.Expression, value bool) {
    t.Helper()

    boolean, ok := b.(*ast.Boolean)
    if !ok {
        t.Fatalf("exp not *ast.Boolean. got=%T", boolean)
    }

    if boolean.Value != value {
        t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
    }

    if boolean.TokenLiteral() != fmt.Sprintf("%t",value) {
        t.Errorf("boolean.TokenLiteral not %s. got=%s", fmt.Sprintf("%t",value), boolean.TokenLiteral())
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


func AssertNumberStatements(t testing.TB, got, want int) {
    if got != want {
        t.Fatalf("program.Statements should contain %d elements, got=%d", want, got)
    }
}

