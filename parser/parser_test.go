package parser

import (
	"fmt"
	"go++/ast"
	lex "go++/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5", "x", 5},
		{"let y = true", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)

	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T", stmt.TokenLiteral())
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not %s. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		returnValue interface{}
	}{
		{"return 5", 5},
		{"return 10944", 10944},
		{"return true", true},
		{"return x", "x"},
	}

	for _, tt := range tests {

		lexer := lex.New(tt.input)

		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt.Statement not *ast.ReturnStatement. got=%T", returnStmt.TokenLiteral())
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got=%q", returnStmt.TokenLiteral())
		}

		testLiteralExpression(t, returnStmt.ReturnValue, tt.returnValue)
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	lexer := lex.New(input)

	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	lexer := lex.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	input := `true;`

	lexer := lex.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.Boolean)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.Boolean. got=%T", stmt.Expression)
	}

	if literal.Value != true {
		t.Fatalf("literal.Value not %t. got=%t", true, literal.Value)
	}

	if literal.TokenLiteral() != "true" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!false;", "!", false},
		{"!true;", "!", true},
	}

	for _, tt := range prefixTests {
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"false != true;", false, "!=", true},
		{"false == false;", false, "==", false},
		{"true == true;", true, "==", true},
	}

	for _, tt := range tests {
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integerLiteral, ok := exp.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral. got=%T", exp)
		return false
	}

	if integerLiteral.Value != value {
		t.Errorf("integerLiteral.Value not %d. got=%d", value, integerLiteral.Value)
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integerLiteral.TokenLiteral not %d. got=%s", value, integerLiteral.TokenLiteral())
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
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
			"a + b / c",
			"(a + (b / c))",
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
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"false != true",
			"(false != true)",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a + b.value",
			"(a + (b.value))",
		},
		{
			"a + b.toInt()",
			"(a + (b.toInt)())",
		},
		{
			"add(a + b.toInt() + c * d / f + g)",
			"add((((a + (b.toInt)()) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if x < y { x }`

	lexer := lex.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative is not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if x < y { x } else { y }`

	lexer := lex.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Fatalf("exp.Alternative is not nil. got=%+v", exp.Alternative)
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `fn (x, y) {
	x + y;
}`

	lexer := lex.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if !testParameterList(t, exp.Parameters, []string{"x", "y"}) {
		return
	}

	if exp.Body == nil {
		t.Fatalf("exp.Body was nil")
	}

	if len(exp.Body.Statements) != 1 {
		t.Fatalf("exp.Body does not contain 1 statements, Got=%d", len(exp.Body.Statements))
	}

	bodyStmt, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body is not ast.ExpressionStatement. got=%T", exp.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterList(t *testing.T) {
	tests := []struct {
		input              string
		expectedParameters []string
	}{
		{
			"fn (x, y) {}",
			[]string{"x", "y"},
		},
		{
			"fn () {}",
			[]string{},
		},
		{
			"fn (a, b, c, d) {}",
			[]string{"a", "b", "c", "d"},
		},
	}

	for _, tt := range tests {
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		testParameterList(t, function.Parameters, tt.expectedParameters)
	}
}

func testParameterList(t *testing.T, parameters []*ast.Identifier, expected []string) bool {
	if len(parameters) != len(expected) {
		t.Errorf("len(parameters)=%d, want=%d", len(parameters), len(expected))
		return false
	}

	for i := 0; i < len(parameters); i++ {
		if !testIdentifier(t, parameters[i], expected[i]) {
			return false
		}
	}

	return true
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	lexer := lex.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Errorf("len(exp.Arguments)=%d, want=3", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"helloWorld"`, "helloWorld"},
		{`"whats up"`, "whats up"},
		{`"Howdy World"`, "Howdy World"},
		{`"H"`, "H"},
	}

	for _, tt := range tests {
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		literal, ok := stmt.Expression.(*ast.StringLiteral)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.StringLiteral. got=%T", stmt.Expression)
		}

		if literal.Value != tt.expected {
			t.Errorf("literal.Value = %s, want=%s", literal.Value, tt.input)
		}
	}
}

func TestAssignExpressionParsing(t *testing.T) {
	tests := []struct {
		input    string
		assignee string
		value    string
	}{
		{"a = 5", "a", "5"},
		{"a = 5 * 5", "a", "(5 * 5)"},
		{`a = "hey"`, "a", "hey"},
	}

	for _, tt := range tests {
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", stmt.Expression)
		}

		exp, ok := stmt.Expression.(*ast.AssignExpression)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.AssignExpression. got=%T", stmt.Expression)
		}

		if exp.Value.String() != tt.value {
			t.Errorf("exp.Value = %s, want=%s", exp.Value, tt.value)
		}

		if exp.Assignee.String() != tt.assignee {
			t.Errorf("exp.Assignee = %s, want=%s", exp.Assignee, tt.assignee)
		}
	}
}

func TestForStatementParsing(t *testing.T) {
	input := `
	for x != 0 {
		"hello"
	}
`

	lexer := lex.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	exp, ok := stmt.Expression.(*ast.ForLoopLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForLoopLiteral. got=%T", stmt.Expression)
	}

	if exp.Condition.String() != `(x != 0)` {
		t.Fatalf("exp.Condition is wrong, want=%s got=%s", "(x != 0)", exp.Condition.String())
	}

	if exp.Body.String() != " {\nhello\n}" {
		t.Fatalf("exp.Body is wrong, want=%s got=%s", " {\nhello\n}", exp.Body.String())
	}
}

func TestMemberAccessExpressionParsing(t *testing.T) {
	tests := []struct {
		input string
		valFn func(t *testing.T, exp *ast.ExpressionStatement) bool
	}{
		{
			`"string".length`,
			func(t *testing.T, stmt *ast.ExpressionStatement) bool {
				exp, ok := stmt.Expression.(*ast.MemberAccessExpression)

				if !ok {
					t.Errorf("stmt.Expression is not ast.MemberAccessExpression. got=%T", stmt.Expression)
				}

				str, ok := exp.Expression.(*ast.StringLiteral)

				if !ok {
					t.Errorf("exp.Expression is not ast.StringLiteral. got=%T", exp)
					return false
				}

				if str.Value != "string" {
					t.Errorf("str.Value = %s, want=%s", str.Value, "string")
					return false
				}

				if exp.AccessedMember.Value != "length" {
					t.Errorf("Accessed member is not length. got=%s", exp.AccessedMember.Value)
					return false
				}

				return true
			},
		},
		{
			`6.isNumber`,
			func(t *testing.T, stmt *ast.ExpressionStatement) bool {
				exp, ok := stmt.Expression.(*ast.MemberAccessExpression)

				if !ok {
					t.Errorf("stmt.Expression is not ast.MemberAccessExpression. got=%T", stmt.Expression)
				}

				integer, ok := exp.Expression.(*ast.IntegerLiteral)

				if !ok {
					t.Errorf("exp.Expression is not ast.IntegerLiteral. got=%T", exp)
					return false
				}

				if integer.Value != 6 {
					t.Errorf("str.Value = %d, want=%d", integer.Value, 6)
					return false
				}

				if exp.AccessedMember.Value != "isNumber" {
					t.Errorf("Accessed member is not isNumber. got=%s", exp.AccessedMember.Value)
					return false
				}

				return true
			},
		},
		{
			`5 + 5.isNumber`,
			func(t *testing.T, stmt *ast.ExpressionStatement) bool {
				exp, ok := stmt.Expression.(*ast.InfixExpression)

				if !ok {
					t.Errorf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
				}

				member, ok := exp.Right.(*ast.MemberAccessExpression)

				if !ok {
					t.Errorf("exp.Right is not ast.MemberAccessExpression. got=%T", exp)
					return false
				}

				integer, ok := member.Expression.(*ast.IntegerLiteral)

				if !ok {
					t.Errorf("exp.Expression is not ast.IntegerLiteral. got=%T", exp)
					return false
				}

				if integer.Value != 5 {
					t.Errorf("str.Value = %d, want=%d", integer.Value, 5)
					return false
				}

				if member.AccessedMember.Value != "isNumber" {
					t.Errorf("Accessed member is not isNumber. got=%s", member.AccessedMember.Value)
					return false
				}

				return true
			},
		},
		{
			`foobar.length`,
			func(t *testing.T, stmt *ast.ExpressionStatement) bool {
				exp, ok := stmt.Expression.(*ast.MemberAccessExpression)

				if !ok {
					t.Errorf("stmt.Expression is not ast.MemberAccessExpression. got=%T", stmt.Expression)
				}

				ident, ok := exp.Expression.(*ast.Identifier)

				if !ok {
					t.Errorf("exp.Expression is not ast.Ident. got=%T", exp)
					return false
				}

				if ident.Value != "foobar" {
					t.Errorf("ident.value = %s, want=%s", ident.Value, "foobar")
					return false
				}

				if exp.AccessedMember.Value != "length" {
					t.Errorf("Accessed member is not length. got=%s", exp.AccessedMember.Value)
					return false
				}

				return true
			},
		},

		{
			`foobar.length()`,
			func(t *testing.T, stmt *ast.ExpressionStatement) bool {
				call, ok := stmt.Expression.(*ast.CallExpression)

				if !ok {
					t.Errorf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
				}

				exp, ok := call.Function.(*ast.MemberAccessExpression)

				if !ok {
					t.Errorf("stmt.Function is not ast.MemberAccessExpression. got=%T", stmt.Expression)
				}

				ident, ok := exp.Expression.(*ast.Identifier)

				if !ok {
					t.Errorf("exp.Expression is not ast.Ident. got=%T", exp)
					return false
				}

				if ident.Value != "foobar" {
					t.Errorf("ident.value = %s, want=%s", ident.Value, "foobar")
					return false
				}

				if exp.AccessedMember.Value != "length" {
					t.Errorf("Accessed member is not length. got=%s", exp.AccessedMember.Value)
					return false
				}

				return true
			},
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.input)
		lexer := lex.New(tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.ExpressionStatement. got=%T", stmt.Expression)
		}

		if !tt.valFn(t, stmt) {
			return
		}
	}
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
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("unexpected expected type: %T, expression: %s", expected, exp.String())
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	booleanLiteral, ok := exp.(*ast.Boolean)

	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if booleanLiteral.Value != value {
		t.Errorf("booleanLiteral.Value not %t. got=%t", value, booleanLiteral.Value)
		return false
	}

	if booleanLiteral.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("booleanLiteral.TokenLiteral not %t. got=%s", value, booleanLiteral.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	indentifier, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.Identifier. Should be=%s got=%T", value, exp)
		return false
	}

	if indentifier.Value != value {
		t.Errorf("indentifier.Value not %s. got=%s", value, indentifier.Value)
		return false
	}

	if indentifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral not %s. got=%s", value, indentifier.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err)
	}

	t.FailNow()
}
