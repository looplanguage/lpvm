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
		{input: "if (true) { return 10 }", expected: 10},
		{input: "if (true) { return 10 } else { 20 }", expected: 10},
		{input: "if (false) { return 10 } else { return 20 }", expected: 20},
		{input: "if (1 > 10) { return 10 } else { return 20 }", expected: 20},
		{input: "if (1 > 10) { return 10 } else if(true) { return 400 } else { return 20 }", expected: 400},
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

func TestVM_IndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1][0]", 1},
		{"[1, 2, 3][2]", 3},
		{"[1, 2, 3][3]", Null},
		{"[1 + 1, 2, 3][3 - 1]", 3},
		{"[1 + 1, 2, 3][0 * 0]", 2},
		{"{}[0]", Null},
		{"{0: 100}[0]", 100},
		{"{0: 100}[1]", Null},
		{"{0 + 1: 100}[0]", Null},
		{"{0 + 1: 100 * 2}[1]", 200},
		{"{0: 100, 3: 400}[0]", 100},
		{"{0: 100, 1: 50 * 1}[1]", 50},
		{"{0 + 1: 100, 2: 2, 3: 3}[3]", 3},
		{"{0 + 1: 100 * 2}[2]", Null},
	}

	runVmTests(t, tests)
}

func TestVM_CallExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"fun() { return 20 }()", 20},
		{"fun() { return 20; return 40; }()", 20},
		{"var early = fun() { return 20; return 40; }; early()", 20},
		{"fun() { 20 }()", Null},
		{"fun() { }()", Null},
		{"fun() { return 20 * 2 }()", 40},
		{"var test = fun() { return 20 * 2 }; test()", 40},
		{"var test = fun() { 20 * 2 + 100 }; test()", Null},
		{"var test = fun() { return fun() { return 20 } }; test()()", 20},
		{"var test = fun() { return 10 }; var testTwo = fun() { return 4 }; test() * testTwo()", 40},
	}

	runVmTests(t, tests)
}

func TestVM_CallExpressionsWithArguments(t *testing.T) {
	tests := []vmTestCase{
		{"fun(a) { return a }(20)", 20},
		{"fun(a, b) { return a * b }(20, 2)", 40},
		{"var double = fun(x) { return x * 2 }; double(500)", 1000},
		{"var double = fun(x) { return x * 2 }; double(500) + double(500)", 2000},
		{"var double = fun(x) { return x * 2 }; var test = fun() { return double(2) + double(2) }; test()", 8},
	}

	runVmTests(t, tests)
}

func TestVM_CallExpressionsWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{"fun() { return 100 }(2)", "wrong number of arguments. expected=0. got=1"},
		{"fun() { return 100 }(2, 2)", "wrong number of arguments. expected=0. got=2"},
		{"fun(x) { return x }(2, 2)", "wrong number of arguments. expected=1. got=2"},
		{"fun(x, y) { return x }(2)", "wrong number of arguments. expected=2. got=1"},
		{"fun(x, y) { return x }()", "wrong number of arguments. expected=2. got=0"},
	}

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.Create()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := Create(comp.Bytecode())
		err = vm.Run(nil)
		if err == nil {
			t.Fatalf("expected VM error but resulted in none.")
		}
		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}
}

func TestVM_FunctionBindings(t *testing.T) {
	tests := []vmTestCase{
		{"var test = fun() { var one = 1; return one }; test()", 1},
		{"var test = fun() { var one = 1; var two = 2; return one + two }; test()", 3},
		{"var global = 100; var test = fun() { return global - 1 }; test()", 99},
		{
			input: `
        var oneAndTwo = fun() { var one = 1; var two = 2; return one + two; };
        oneAndTwo();
        `,
			expected: 3,
		},
		{
			input: `
        var oneAndTwo = fun() { var one = 1; var two = 2; return one + two; };
        var threeAndFour = fun() { var three = 3; var four = 4; return three + four; };
        oneAndTwo() + threeAndFour();
        `,
			expected: 10,
		},
		{
			input: `
        var firstFoobar = fun() { var foobar = 50; return foobar; };
        var secondFoobar = fun() { var foobar = 100; return foobar; };
        firstFoobar() + secondFoobar();
        `,
			expected: 150,
		},
		{
			input: `
        var globalSeed = 50;
        var minusOne = fun() {
            var num = 1;
            return globalSeed - num;
        }
        var minusTwo = fun() {
            var num = 2;
            return globalSeed - num;
        }
        minusOne() + minusTwo();
        `,
			expected: 97,
		},
	}

	runVmTests(t, tests)
}

func TestVM_BuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len([])`, 0},
		{`len([1])`, 1},
		{`len([1, 2, 3])`, 3},
		{`len("hello")`, 5},
		{`len(1)`, &object.Error{Message: `incorrect argument type, can not iterate. got="INTEGER"`}},
		{`len({})`, &object.Error{Message: `incorrect argument type, can not iterate. got="HASHMAP"`}},
	}

	runVmTests(t, tests)
}

func TestVM_Closures(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			var newClosure = fun(a) {
				return fun() { return a; };
			};
			var closure = newClosure(99);
			closure();
			`,
			99,
		},
		{
			`
        var newAdderOuter = fun(a, b) {
            var c = a + b;
            return fun(d) {
                var e = d + c;
                return fun(f) { return e + f; };
            };
        };
        var newAdderInner = newAdderOuter(1, 2)
        var adder = newAdderInner(3);
        adder(8);
        `,
			14,
		},
		{
			`
        var a = 1;
        var newAdderOuter = fun(b) {
            return fun(c) {
                return fun(d) { return a + b + c + d };
            };
        };
        var newAdderInner = newAdderOuter(2)
        var adder = newAdderInner(3);
        adder(8);
        `,
			14,
		},
		{
			`
        var newClosure = fun(a, b) {
            var one = fun() { return a; };
            var two = fun() { return b; };
            return fun() { return one() + two(); };
        };
        var closure = newClosure(9, 90);
        closure();
        `,
			99,
		},
	}

	runVmTests(t, tests)
}

func TestVM_Recursive(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			var fibonacci = fun(x) {
				return if (x == 0) {
					return 0;
				} else {
					if (x == 1) {
						return 1;
					} else {
						return fibonacci(x - 1) + fibonacci(x - 2);
					}
				}
			};

			fibonacci(15);
`,
			610,
		},
	}

	runVmTests(t, tests)
}

func TestVM_EarlyReturns(t *testing.T) {
	tests := []vmTestCase{
		{"if(true) { return 50; return 20 }", 50},
		{"if(true) { if(true) { return 10; } return 50; return 20 }", 10},
		{"if(true) { if(true) { if(true) { return 2000; } return 10; } return 50; return 20 }", 2000},
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
		err = vm.Run(nil)
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
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not error. got=%T (%+v)", actual, actual)
		}

		if errObj.Message != expected.Message {
			t.Errorf("wrong error. expected=%q. got=%q", expected.Message, errObj.Message)
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
