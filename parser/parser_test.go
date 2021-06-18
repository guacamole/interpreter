package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests:= []struct{
		input string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",

		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _,tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t,p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("Excepected %s. got %s",tt.expected,actual)
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct{
		input string
		leftValue int64
		operator string
		rightValue int64
	}{
		{"5 + 5;",5,"+",5},
		{"5 - 5;",5,"-",5},
		{"5 * 5;",5,"*",5},
		{"5 / 5;",5,"/",5},
		{"5 > 5;",5,">",5},
		{"5 < 5;",5,"<",5},
		{"5 == 5;",5,"==",5},
		{"5 != 5;",5,"!=",5},

	}

	for _,tt := range infixTests {

		l := lexer.New(tt.input)
		p := New(l)

		program:= p.ParseProgram()
		checkParseErrors(t,p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program statements not equal to %d got %d",1,len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.statement is not of ast.expression type, got =%T",program.Statements[0])
		}

		exp,ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.infix type. got= %T",stmt.Expression)
		}

		if !testIntegerLiteral(t,exp.Left,tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("operator is not %s. got=%s",tt.operator,exp.Operator)
		}

		if !testIntegerLiteral(t,exp.Right,tt.rightValue) {
			return
		}

	}
}

func TestParsingPrefixExpressions(t *testing.T) {

	prefixTests := []struct{
		input string
		operator string
		integerValue int64
	}{
		{"!5","!",5},
		{"-15","-",15},
	}

	for _,tt := range prefixTests{
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t,p)

		if len(program.Statements) != 1 {
			t.Fatalf("program doesn't have enough statements got= %q",len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not an ast.ExpressionStatement. got = %T",program.Statements[0])
		}

		exp,ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp is not ast.PrefixExpression type. got= %T",stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got %s",tt.operator,exp.Operator)
		}

		if !testIntegerLiteral(t,exp.Right,tt.integerValue) {
			return
		}
	}

}

func testIntegerLiteral(t *testing.T,il ast.Expression,value int64) bool {

	integ,ok := il.(*ast.IntegerLiteral)
	if !ok{
		t.Errorf("il is not ast.integerLiteral type. got= %T",il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ Value is not %d. got %d",value,integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d",value) {
		t.Errorf("integ token literal is not %d. got %s",value,integ.TokenLiteral())
		return false
	}
	return true

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t,p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have enough statements got= %q",len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not an ast.ExpressionStatement. got = %T",program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not ast.IntegerLiteral typr. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got %d",5,literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() not %s. got %s","5",literal.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	 l := lexer.New(input)
	 p := New(l)
	 program := p.ParseProgram()
	 checkParseErrors(t,p)

	 if len(program.Statements) != 1 {
	 	t.Fatalf("program doesn't have enough statements got= %q",len(program.Statements))
	 }

	 stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	 if !ok {
	 	t.Fatalf("program.Statement[0] is not an ast.ExpressionStatement. got = %T",program.Statements[0])
	 }

	 ident, ok := stmt.Expression.(*ast.Identifier)
	 if !ok {
		 t.Fatalf("stmt.Expression is not an ast.Identifier. got = %T",stmt.Expression)
	 }

	 if ident.Value != "foobar" {
		 t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	 }

	 if ident.TokenLiteral() != "foobar" {
		 t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
		 }


}

func TestReturnStatements(t *testing.T) {

	input := `
return  5;
return 10;
return 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t,p)

	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements does not cotain 3 statements got %d",len(program.Statements))
	}

	for _,stmt := range program.Statements {

		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok{
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("turnStmt.TokenLiteral not 'return', got = %q",returnStmt.TokenLiteral())
		}
	}


}

func TestLetStatements(t *testing.T) {

	input := `let x = 5;
	let y = 10;
	let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got= %d",len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i,tt := range tests {

		stmt := program.Statements[i]
		if !testLetStatement(t,stmt,tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool{

	if s.TokenLiteral() != "let" {
		t.Errorf("token literal is not let got= %q",s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement) // TODO
	if !ok {
		t.Errorf("s is not ast.LetStatement got %T",s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s. got= %s",name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %s. got= %s",name,letStmt.Name.TokenLiteral())
		return false
	}
	return true
}


func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0{
		return
	}
	t.Errorf("parser has %d errors", len(errors))

	for _,msg  := range errors {
		t.Errorf("parser error %q",msg)
	}
	t.FailNow()
}

