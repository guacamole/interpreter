package evaluator

import (
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"testing"
)

func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 10 - 5", 15},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 - 50", 0},
		{"5 * 2 + 10", 20},
		{"50 / 2 * 2 + 10", 60},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObj(t, evaluated, tt.expected)
	}

}

func testIntegerObj(t *testing.T, obj object.Object, expected int64) bool {

	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf(" object is not of integer type. got= %T (%v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value.got=%d,want= %d", result.Value, expected)
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {

	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T", obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value.got=%t, want=%t", result.Value, expected)
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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 < 2) {10}", 10},
		{"if (1 > 2) {10} else {20}", 20},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObj(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"2;return 10;", 10},
		{"return 2*5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObj(t, evaluated, tt.expected)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {

	if obj != NULL {
		t.Errorf("object is not nill. got=%T(%v)", obj, obj)
		return false
	}
	return true
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) {true + false}", "unknown operator: BOOLEAN + BOOLEAN"},
		{
			`if (10 > 1) {
					if (10 > 1) {
						return true + false;
					}
				  return 1;
				 }
				`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got= %T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong message. expected= %q,got= %q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a= 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c= a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObj(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not function. got= %T (%+v)",evaluated,evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("funtion has wrong parameters. Parameters= %+v",fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x. got= %q",fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody{
		t.Fatalf("body is not %q. got= %q",expectedBody,fn.Body.String())
	}
}

