package vm

import (
	"fmt"
	"github.com/looplanguage/compiler/compiler"
	"github.com/looplanguage/loop/lexer"
	"github.com/looplanguage/loop/models/ast"
	"github.com/looplanguage/loop/models/object"
	"github.com/looplanguage/loop/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"10 * 10", 100},
		{"10 / 10", 1},
		{"10 - 10", 0},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: "true", expected: true},
		{input: "false", expected: false},
		{input: "1 == 1", expected: true},
		{input: "true == false", expected: false},
		{input: "true == true", expected: true},
		{input: "1 == 10", expected: false},
	}

	runVmTests(t, tests)
}

func TestConditionalExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: "if (true) { 10 }", expected: 10},
		{input: "if (true) { 10 } else { 20 }", expected: 10},
		{input: "if (false) { 10 } else { 20 }", expected: 20},
		{input: "if (1 > 10) { 10 } else { 20 }", expected: 20},
		{input: "if (1 > 10) { 10 } else if(true) { 400 } else { 20 }", expected: 400},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	i := 0
	for _, tc := range tests {
		i++
		program := parse(tc.input)

		comp := compiler.Create()
		err := comp.Compile(program)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := Create(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElement := vm.LastPoppedStackElem()

		testExpectedObject(t, tc.expected, stackElement)
	}
}

func testExpectedObject(
	t *testing.T,
	expected interface{},
	actual object.Object,
) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBoolObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBoolObject failed: %s", err)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)

	if !ok {
		return fmt.Errorf("object is not integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d. expected=%d", result.Value, expected)
	}

	return nil
}

func testBoolObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)

	if !ok {
		return fmt.Errorf("object is not boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t. expected=%t", result.Value, expected)
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.Create(input)
	p := parser.Create(l)
	return p.Parse()
}
