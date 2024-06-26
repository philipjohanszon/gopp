package evaluator

import (
	"fmt"
	lex "go++/lexer"
	"go++/object"
	parse "go++/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"2 * 2 + 4 * 2", 12},
		{"5 + 5 + 5 + 5 + 5 - 10", 15},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 - 50", 0},
		{"5 * 2 + 10", 20},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEvaluation(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!!false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 > 1", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"0 == 0", true},
		{"1 != 0", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"(1 != 2) == true", true},
		{"(20 - 2 == 18) == true", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
		{"(1 > 2) == true", false},
		{"(1 < 2) == false", false},
	}

	for _, tt := range tests {
		evaluated := testEvaluation(tt.input)

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!0", true},
		{"!!0", false},
	}

	for _, tt := range tests {
		evaluated := testEvaluation(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 20 }", nil},
		{"if 1 { 10 } else { 5 }", 10},
		{"if 0 { 10 } else { 5 }", 5},
		{"if 1 > 2 { 10 } else { 5 }", 5},
		{"if 1 < 2 { 10 } else { 5 }", 10},
		{"if 1 > 2 { 10 }", nil},
	}

	for _, tt := range tests {
		evaluated := testEvaluation(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 5;", 5},
		{"return 5 * 2;", 10},
		{"9; return 5 * 3;", 15},
		{"9; return 5 * 2; 5;", 10},
		{
			`
			if 10 > 1 {
				if 10 > {
					return 10;
				}
				return 1;
			}`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEvaluation(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + true; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if 10 > 1 { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
		if 10 > 1 {
			if 10 > {
				return true + true;
			}
			return 1;
		}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEvaluation(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			println(tt.input)
			t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. got=%q, want=%q", errObj.Message, tt.expectedMessage)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 5; x;", 5},
		{"let x = 5 * 5; x;", 25},
		{"let x = 5 * 5; let y = x; y;", 25},
		{"let x = 5 * 5; let y = x + 5; y;", 30},
		{"let x = 5 * 5; let y = 5; let z = x - 4 * y; x * 2 + y * 3 / 5 + z", 58},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEvaluation(tt.input), tt.expected) {
			fmt.Println(tt.input)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEvaluation(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong number of parameters. got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x' . got=%q", fn.Parameters[0].String())
	}

	expectedBody := " {\n(x + 2)\n}"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { return x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, add(5, 5));", 15},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEvaluation(tt.input), tt.expected)
	}
}

func testEvaluation(input string) object.Object {
	lexer := lex.New(input)
	parser := parse.New(lexer)

	program := parser.ParseProgram()
	env := object.NewEnvironment()

	return Evaluate(program, env)
}
