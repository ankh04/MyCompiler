package parser

import (
	"MyCompiler/src/ast"
	"MyCompiler/src/lexer"
	"MyCompiler/src/parser"
	"MyCompiler/src/token"
	"fmt"
	"testing"
)

// region 帮助函数

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLetStatement(t *testing.T, s ast.Statement, name string, expectedVal interface{}) {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let'. got = %q", s.TokenLiteral())
		return
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("the statement is not LetStatement type. got = %T", s)
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got = '%s'", name, letStmt.Name.Value)
		return
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got = '%s'",
			name, letStmt.Name.TokenLiteral())
		return
	}

	if !testLiteralExpression(t, letStmt.Value, expectedVal) {
		return
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Error()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not IntegerLiteral. got %T", exp)
		return false
	}

	if integ.Value != value {
		t.Fatalf("integ.Value is not %d. got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("integ.TokenLiteral is not %d. got %s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testBoolLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	integ, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp is not BooleanLiteral. got %T", exp)
		return false
	}

	if integ.Value != value {
		t.Fatalf("integ.Value is not %t. got %t", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Fatalf("integ.TokenLiteral is not %t. got %s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not InfixExpression. got %T", exp)
		return false
	}
	// 测试左边的表达式
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	// 测试运算符
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s. got %s", operator, opExp.Operator)
		return false
	}
	// 测试右边的表达式
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

// endregion

// 构建"let myVar = anotherVar;"语句的ast，测试字符串输出
func TestString(t *testing.T) {
	program := &ast.Program{
		Statement: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got = %q", program.String())
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993;
`
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statement) != 3 {
		t.Fatalf("program.Statement does contain 3 statements. got = %d",
			len(program.Statement))
	}

	for _, stmt := range program.Statement {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not ReturnStatement. got = %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got = %q", returnStmt.TokenLiteral())
		}
	}
	expects := []struct {
		expectedVal int64
	}{
		{5},
		{10},
		{993},
	}
	for i, tt := range expects {
		stmt := program.Statement[i]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not ReturnStatement. got = %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got = %q", returnStmt.TokenLiteral())
		}
		testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedVal)
	}
}

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 78787878;
`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statement) != 3 {
		t.Fatalf("program.Statement does contain 3 statements. got = %d",
			len(program.Statement))
	}

	expects := []struct {
		expectedIdentifier string
		expectedVal        int64
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 78787878},
	}
	for i, tt := range expects {
		stmt := program.Statement[i]
		testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedVal)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statement) != 1 {
		t.Fatalf("IntegerLiteralExpression does contain 1 statements. got = %d",
			len(program.Statement))
	}
	stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ExpressionStatement. got = %T",
			program.Statement[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	testLiteralExpression(t, literal, 5)
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := `
true;
false;
`
	expects := []bool{true, false}

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statement) != 2 {
		t.Fatalf("BooleanLiteralExpression does contain 2 statments, got = %d", len(program.Statement))
	}
	for i, tt := range expects {
		stmt := program.Statement[i]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[i] should be ExpressionStatement, got %T",
				stmt)
		}
		testLiteralExpression(t, expStmt.Expression, tt)
	}

}

// 仅有两种前缀运算符： ！ 和 -
func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statement) != 1 {
			t.Fatalf("program.Statement should contain 1 statements. got = %d",
				len(program.Statement))
		}
		stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ExpressionStatement. got = %T",
				program.Statement[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not prefixExpression. got = %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %s. got %s", tt.operator, exp.Operator)
		}
		// 测试右值是否是整数字面量，并检查值是否正确
		testLiteralExpression(t, exp.Right, tt.value)

	}

}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  interface{}
		operator string
		rightVal interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true != true;", true, "!=", true},
		{"false != true;", false, "!=", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statement) != 1 {
			t.Fatalf("program.Statement should have 1 statement. got %d\n",
				len(program.Statement))
		}

		stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got %T\n",
				program.Statement[0])
		}
		testInfixExpression(t, stmt.Expression, tt.leftVal, tt.operator, tt.rightVal)
	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b",
			"((-a) * b)",
		}, {
			"!-a", "(!(-a))",
		}, {
			"a + b + c",
			"((a + b) + c)",
		}, {
			"a + b - c",
			"((a + b) - c)",
		}, {
			"a * b * c",
			"((a * b) * c)",
		}, {
			"a * b / c",
			"((a * b) / c)",
		}, {
			"a + b / c",
			"(a + (b / c))",
		}, {
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		}, {
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		}, {
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		}, {
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		}, {
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		}, {
			"3 < 5 == true",
			"((3 < 5) == true)",
		}, {
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		}, {
			"!(true == true)",
			"(!(true == true))",
		}, {
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foo"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statement) != 1 {
		t.Fatalf("IdentifierExpression does contain 1 statements. got = %d",
			len(program.Statement))
	}
	stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ExpressionStatement. got = %T",
			program.Statement[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not Indentifier. got = %T", stmt.Expression)
	}
	if ident.Value != "foo" {
		t.Errorf("ident.Value is not %s. got %s", input, ident.Value)
	}
	if ident.TokenLiteral() != "foo" {
		t.Errorf("ident.TokenLiteral is not %s. got %s", input, ident.TokenLiteral())
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statement) != 1 {
		t.Fatalf("program.Statements do not contain %d statements. got %d",
			1, len(program.Statement))
	}

	// 条件一定是表达式语句
	stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got %T",
			program.Statement[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ifExpression, got %T", stmt.Expression)
	}

	// 本示例的条件应该是一个中缀表达式
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statments. got %d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ExpressionStatement. got %T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("alternative is not 1 statments. got %d",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statement[0] is not ExpressionStatement. got %T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statement) != 1 {
		t.Fatalf("program.Statements do not contain %d statements. got %d",
			1, len(program.Statement))
	}

	// 条件一定是表达式语句
	stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got %T",
			program.Statement[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ifExpression, got %T", stmt.Expression)
	}

	// 本示例的条件应该是一个中缀表达式
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statments. got %d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ExpressionStatement. got %T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Fatalf("exp.Alternative should be 0. got %+v", exp.Alternative)
	}

}

func TestFnExpression(t *testing.T) {
	input := `
fn(x, y) {
	return x + y;
}
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statement) != 1 {
		t.Fatalf("program.Statement should have 1 Statement. got %d", len(program.Statement))
	}

	// 应该是表达式语句
	expStmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fnExpression should be ExpressionStatement. got %T", program.Statement[0])
	}

	// 应该是函数表达式
	fnExp, ok := expStmt.Expression.(*ast.FnExpression)
	if !ok {
		t.Fatalf("expStmt should be fnExpression. got %T", expStmt.Expression)
	}

	// 检查参数
	expectsParams := []string{"x", "y"}
	for i, param := range fnExp.Parameters {
		if !testIdentifier(t, param, expectsParams[i]) {
			return
		}
	}

	// 检查函数体
	if len(fnExp.Body.Statements) != 1 {
		t.Fatalf("fnExp should only have 1 statement. got %d", len(fnExp.Body.Statements))
	}

	// 应该是return语句
	returnStmt, ok := fnExp.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("the statement in fnExp should be RetunStatement. got %T", fnExp.Body.Statements[0])
	}

	// 检查中缀表达式
	testInfixExpression(t, returnStmt.ReturnValue, "x", "+", "y")

}

func TestFunctionParameterParsing(t *testing.T) {

	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statement[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FnExpression)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}

}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statement) != 1 {
		t.Fatalf("program.Statement should have contain 1 statement. got %d", len(program.Statement))
	}

	stmt, ok := program.Statement[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] should be ExpressionStatement. got %T", program.Statement[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression should have be CallExpression. got %T", stmt.Expression)
	}

	// 函数应该是个标识符
	fnIdent, ok := callExp.Function.(*ast.Identifier)
	if !ok {
		t.Fatalf("callExp.Function should be an Indentifer, but got %T", callExp.Function)
	}

	// 检查标识符
	if !testIdentifier(t, fnIdent, "add") {
		return
	}

	// 检查实参
	testLiteralExpression(t, callExp.Arguments[0], 1)
	testInfixExpression(t, callExp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExp.Arguments[2], 4, "+", 5)
}
