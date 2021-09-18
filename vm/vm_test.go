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
		{input: "if (1 > 10) {}", expected: Null},
		{input: "if (false) {}", expected: Null},
	}

	runVmTests(t, tests)
}

func TestVM_VariableDeclarations(t *testing.T) {
	tests := []vmTestCase{
		{"var one = 1; one", 1},
		{"var _test = 10 + 10; _test", 20},
		{"var _test = 10; var _test2 = _test * 4; _test2", 40},
	}

	runVmTests(t, tests)
}

func TestVM_OpAdd(t *testing.T) {
	tests := []vmTestCase{
		{`"hello"`, "hello"},
		{`"hello " + "world"`, "hello world"},
		{`"hello world"`, "hello world"},
	}

	runVmTests(t, tests)
}

func TestVM_Arrays(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1]", []int{1}},
		{"[1 + 2]", []int{3}},
		{"[1 + 2, 1 * 1]", []int{3, 1}},
		{"[1 + 2, 1 * 1, 42 - 10]", []int{3, 1, 32}},
	}

	runVmTests(t, tests)
}

func TestVM_Hashes(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{"{1: 2, 3: 4}", map[object.HashKey]int64{
			(&object.Integer{Value: 1}).Hash(): 2,
			(&object.Integer{Value: 3}).Hash(): 4,
		}},
		{"{1: 2 + 2, 3: 4 * 2}", map[object.HashKey]int64{
			(&object.Integer{Value: 1}).Hash(): 4,
			(&object.Integer{Value: 3}).Hash(): 8,
		}},
		{"{1 + 1: 2 + 2, 3: 4 * 2}", map[object.HashKey]int64{
			(&object.Integer{Value: 2}).Hash(): 4,
			(&object.Integer{Value: 3}).Hash(): 8,
		}},
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
		err := testBoolObject(expected, actual)
		if err != nil {
			t.Errorf("testBoolObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not null: %T (%+v)", actual, actual)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.HashMap)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
			return
		}
		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong number of Pairs. want=%d, got=%d",
				len(expected), len(hash.Pairs))
			return
		}
		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}
			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object is not array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong number of elements. expected=%d. got=%d", len(expected), len(array.Elements))
			return
		}

		for i, expectedElement := range expected {
			err := testIntegerObject(int64(expectedElement), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	}
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)

	if !ok {
		return fmt.Errorf("object is not string. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%s. expected=%s", result.Value, expected)
	}

	return nil
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
